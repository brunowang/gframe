package gfdiscov

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"time"
)

var (
	errClosed = errors.New("etcd monitor chan has been closed")
)

// EtcdClient interface represents an etcd client.
type EtcdClient interface {
	ActiveConnection() *grpc.ClientConn
	Close() error
	Ctx() context.Context
	Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error)
	KeepAlive(ctx context.Context, id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error)
	Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
	Revoke(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error)
	Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan
}

func NewEtcdCli(endpoints []string) (EtcdClient, error) {
	cfg := clientv3.Config{
		Endpoints:           endpoints,
		AutoSyncInterval:    time.Minute,
		DialTimeout:         5 * time.Second,
		RejectOldCluster:    true,
		PermitWithoutStream: true,
	}
	return clientv3.New(cfg)
}
