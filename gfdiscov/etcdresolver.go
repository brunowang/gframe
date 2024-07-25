package gfdiscov

import (
	"context"
	"errors"
	"fmt"
	"github.com/brunowang/gframe/gferr"
	"github.com/brunowang/gframe/gflog"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"strings"
	"sync"
	"time"
)

type EtcdResolver struct {
	cli     EtcdClient
	once    sync.Once
	stat    map[string]*resolver.State
	addrSet map[string]map[string]struct{}
	lock    sync.RWMutex
	cc      resolver.ClientConn
}

func NewEtcdResolver() *EtcdResolver {
	return &EtcdResolver{
		stat:    make(map[string]*resolver.State),
		addrSet: make(map[string]map[string]struct{}),
	}
}

func (e *EtcdResolver) dialOnce(endpoints []string) error {
	var initErr error
	e.once.Do(func() {
		cli, err := NewEtcdCli(endpoints)
		if err != nil {
			initErr = err
			return
		}
		e.cli = cli
	})
	if initErr != nil {
		return initErr
	}
	return nil
}

func (e *EtcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	hosts := strings.Split(target.URL.Host, ",")
	if err := e.dialOnce(hosts); err != nil {
		return nil, err
	}

	e.cc = cc
	key := target.Endpoint()

	revision, err := e.load(key)
	if err != nil {
		return nil, err
	}
	go func() {
		defer gferr.Recovery()
		e.watch(key, revision)
	}()
	return nil, nil
}

func (e *EtcdResolver) Scheme() string {
	return etcdScheme
}

func (e *EtcdResolver) load(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(e.cli.Ctx(), 3*time.Second)
	rsp, err := e.cli.Get(ctx, makeKeyPrefix(key), clientv3.WithPrefix())
	cancel()
	if err != nil {
		return 0, err
	}

	addrs := make([]resolver.Address, 0, len(rsp.Kvs))
	addrSet := make(map[string]struct{}, len(rsp.Kvs))
	for _, kv := range rsp.Kvs {
		val := string(kv.Value)
		if _, has := addrSet[val]; has {
			continue
		}
		addrSet[val] = struct{}{}

		addrs = append(addrs, resolver.Address{Addr: val})
	}

	e.lock.Lock()
	e.stat[key] = &resolver.State{Addresses: addrs}
	e.addrSet[key] = addrSet
	e.lock.Unlock()

	if err := e.cc.UpdateState(e.getResolverState(key)); err != nil {
		gflog.Error(ctx, "EtcdResolver.load.UpdateState.failed", zap.Error(err))
	}

	return rsp.Header.Revision, nil
}

func (e *EtcdResolver) watch(key string, revision int64) {
	ctx := context.Background()
	for {
		err := e.watchStream(key, revision)
		if err == nil {
			return
		}
		gflog.Error(ctx, "EtcdResolver.watch.watchStream.failed", zap.Error(err))

		if revision != 0 && errors.Is(err, rpctypes.ErrCompacted) {
			gflog.Error(ctx, "EtcdResolver.watch.watchStream has been compacted, try to reload", zap.Int64("revision", revision))
			rev, err := e.load(key)
			if err != nil {
				gflog.Error(ctx, "EtcdResolver.watch.load.failed", zap.Error(err))
				continue
			}
			revision = rev
		}
	}
}

func (e *EtcdResolver) watchStream(key string, revision int64) error {
	opts := []clientv3.OpOption{clientv3.WithPrefix()}
	if revision != 0 {
		opts = append(opts, clientv3.WithRev(revision+1))
	}

	rspCh := e.cli.Watch(clientv3.WithRequireLeader(e.cli.Ctx()), makeKeyPrefix(key), opts...)
	for {
		select {
		case rsp, ok := <-rspCh:
			if !ok {
				return errClosed
			}
			if rsp.Canceled {
				return fmt.Errorf("etcd monitor chan has been canceled, error: %w", rsp.Err())
			}
			if rsp.Err() != nil {
				return fmt.Errorf("etcd monitor chan error: %w", rsp.Err())
			}

			e.handleWatchEvents(key, rsp.Events)
		}
	}
}

func (e *EtcdResolver) handleWatchEvents(key string, events []*clientv3.Event) {
	ctx := context.Background()
	for _, ev := range events {
		if !strings.HasPrefix(string(ev.Kv.Key), key) {
			continue
		}
		switch ev.Type {
		case clientv3.EventTypePut:
			e.lock.Lock()
			e.ensureMap(key)

			if _, has := e.addrSet[key][string(ev.Kv.Value)]; has {
				continue
			}
			e.stat[key].Addresses = append(e.stat[key].Addresses, resolver.Address{Addr: string(ev.Kv.Value)})
			e.addrSet[key][string(ev.Kv.Value)] = struct{}{}

			e.lock.Unlock()

			if err := e.cc.UpdateState(e.getResolverState(key)); err != nil {
				gflog.Error(ctx, "EtcdResolver.handlePutEvent.UpdateState.failed", zap.Error(err))
			}
		case clientv3.EventTypeDelete:
			e.lock.Lock()
			e.ensureMap(key)

			if _, has := e.addrSet[key][string(ev.Kv.Value)]; !has {
				continue
			}
			for idx, addr := range e.stat[key].Addresses {
				if addr.Addr == string(ev.Kv.Value) {
					e.stat[key].Addresses = append(e.stat[key].Addresses[:idx], e.stat[key].Addresses[idx+1:]...)
					break
				}
			}
			delete(e.addrSet[key], string(ev.Kv.Value))

			e.lock.Unlock()

			if err := e.cc.UpdateState(e.getResolverState(key)); err != nil {
				gflog.Error(ctx, "EtcdResolver.handleDeleteEvent.UpdateState.failed", zap.Error(err))
			}
		}
	}
}

func (e *EtcdResolver) getResolverState(key string) resolver.State {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return *e.stat[key]
}

func (e *EtcdResolver) ensureMap(key string) {
	if _, has := e.stat[key]; !has {
		e.stat[key] = &resolver.State{}
	}
	if _, has := e.addrSet[key]; !has {
		e.addrSet[key] = make(map[string]struct{})
	}
}
