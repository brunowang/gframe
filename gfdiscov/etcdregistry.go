package gfdiscov

import (
	"github.com/brunowang/gframe/gferr"
	"github.com/brunowang/gframe/gflog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

var LeaseTTL int64 = 10

type EtcdRegistry struct {
	cfg      EtcdConf
	cli      EtcdClient
	once     sync.Once
	listenOn string
}

func NewEtcdRegistry(etcdConf EtcdConf, listenOn string) *EtcdRegistry {
	return &EtcdRegistry{
		cfg:      etcdConf,
		listenOn: listenOn,
	}
}

func (e *EtcdRegistry) dialOnce() error {
	var initErr error
	e.once.Do(func() {
		cli, err := NewEtcdCli(e.cfg.Hosts)
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

func (e *EtcdRegistry) Register() error {
	if err := e.dialOnce(); err != nil {
		return err
	}

	lease, err := e.register()
	if err != nil {
		return err
	}
	return e.keepalive(lease)
}

func (e *EtcdRegistry) register() (clientv3.LeaseID, error) {
	rsp, err := e.cli.Grant(e.cli.Ctx(), LeaseTTL)
	if err != nil {
		return clientv3.NoLease, err
	}
	lease := rsp.ID
	key := makeEtcdKey(e.cfg.Key, lease)
	_, err = e.cli.Put(e.cli.Ctx(), key, figureOutListenOn(e.listenOn), clientv3.WithLease(lease))
	if err != nil {
		return clientv3.NoLease, err
	}
	return lease, nil
}

func (e *EtcdRegistry) keepalive(lease clientv3.LeaseID) error {
	ch, err := e.cli.KeepAlive(e.cli.Ctx(), lease)
	if err != nil {
		return err
	}

	go func() {
		defer gferr.Recovery()
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					e.revoke(lease)
					e.loopRegister()
					return
				}
			}
		}
	}()

	return nil
}

func (e *EtcdRegistry) revoke(lease clientv3.LeaseID) {
	if _, err := e.cli.Revoke(e.cli.Ctx(), lease); err != nil {
		gflog.Error(e.cli.Ctx(), "EtcdRegistry.revoke.failed", zap.Error(err))
	}
}

func (e *EtcdRegistry) loopRegister() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		lease, err := e.register()
		if err != nil {
			gflog.Error(e.cli.Ctx(), "EtcdRegistry.loopRegister.failed", zap.Error(err))
			continue
		}
		if err := e.keepalive(lease); err != nil {
			gflog.Error(e.cli.Ctx(), "EtcdRegistry.loopKeepalive.failed", zap.Error(err))
			continue
		}
		break
	}
}
