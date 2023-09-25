package gfcfg

import (
	"context"
	"fmt"
	"sync"
)

type Configure interface {
	Init(ctx context.Context) error
	Read(ctx context.Context) ([]byte, error)
	Watch(ctx context.Context) (<-chan []byte, error)
}

type ConfigMgr interface {
	ReadAndWatch(ctx context.Context, key string) ([]byte, <-chan []byte, error)
	Read(ctx context.Context, key string) ([]byte, error)
	Watch(ctx context.Context, key string) (<-chan []byte, error)
}

type ConfigMgrImpl struct {
	configs syncMap[Configure]
	once    sync.Once
	iniErr  error
	caches  syncMap[[]byte]
}

func NewConfigMgr() ConfigMgr {
	return &ConfigMgrImpl{}
}

func (m *ConfigMgrImpl) initOnce(ctx context.Context) error {
	m.once.Do(func() {
		m.configs.Range(func(key, val any) bool {
			if err := val.(Configure).Init(ctx); err != nil {
				m.iniErr = err
				return false
			}
			return true
		})
	})
	if m.iniErr != nil {
		return m.iniErr
	}
	return nil
}

func (m *ConfigMgrImpl) ReadAndWatch(ctx context.Context, key string) ([]byte, <-chan []byte, error) {
	bs, err := m.Read(ctx, key)
	if err != nil {
		return nil, nil, err
	}
	ch, err := m.Watch(ctx, key)
	if err != nil {
		return nil, nil, err
	}
	return bs, ch, nil
}

func (m *ConfigMgrImpl) Read(ctx context.Context, key string) ([]byte, error) {
	if err := m.initOnce(ctx); err != nil {
		return nil, err
	}
	if val, has := m.caches.Load(key); has {
		return val, nil
	}
	cfg, has := m.configs.Load(key)
	if !has {
		return nil, fmt.Errorf("read config got unknown config key %s", key)
	}
	val, err := cfg.Read(ctx)
	if err != nil {
		return nil, err
	}
	m.caches.Store(key, val)
	return val, nil
}

func (m *ConfigMgrImpl) Watch(ctx context.Context, key string) (<-chan []byte, error) {
	if err := m.initOnce(ctx); err != nil {
		return nil, err
	}
	cfg, has := m.configs.Load(key)
	if !has {
		return nil, fmt.Errorf("watch config got unknown config key %s", key)
	}
	ch, err := cfg.Watch(ctx)
	if err != nil {
		return nil, err
	}
	valCh := make(chan []byte, 1)
	go func() {
		for val := range ch {
			m.caches.Store(key, val)
			valCh <- val
		}
	}()
	return valCh, nil
}

/*
make sync.Map support generic type
*/
type syncMap[T any] struct {
	sync.Map
}

func (m *syncMap[T]) Load(key string) (T, bool) {
	return m.Map.Load(key)
}

func (m *syncMap[T]) Store(key string, val T) {
	m.Map.Store(key, val)
}
