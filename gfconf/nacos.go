package gfconf

import (
	"bytes"
	"context"
	"github.com/brunowang/gframe/gfconf/internal/nacos"
	"github.com/brunowang/gframe/gfcontainer"
	"github.com/brunowang/gframe/gflog"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

type NacosConfigure struct {
	basicAuth nacos.BasicAuth
	paths     []nacos.PathParsed
	endpoint  string
	viperMap  gfcontainer.SafeMap[nacos.PathParsed, *viper.Viper]
	loadOnce  sync.Once
	loadErr   error
}

func NewNacosConfigure(paths []string, endpoint string, auth ...BasicAuth) *NacosConfigure {
	parsedPaths := make([]nacos.PathParsed, 0, len(paths))
	for _, p := range paths {
		parsedPaths = append(parsedPaths, nacos.Path(p).Parse())
	}
	ret := &NacosConfigure{
		paths:    parsedPaths,
		endpoint: endpoint,
	}
	if len(auth) > 0 {
		ret.basicAuth = auth[0]
	}
	return ret
}

func (c *NacosConfigure) GetConfig(path string) (*viper.Viper, error) {
	parsedPath := nacos.Path(path).Parse()
	ret, has := c.viperMap.Load(parsedPath)
	if !has {
		return nil, ErrorConfigNotLoaded
	}
	return ret, nil
}

func (c *NacosConfigure) Init() error {
	if len(c.paths) == 0 {
		return ErrorConfigInitFailed
	}
	loadFn := func(namespaces []string) error {
		if err := nacos.Init(namespaces, []string{c.endpoint}, c.basicAuth); err != nil {
			return err
		}
		if err := nacos.SupportViper(); err != nil {
			return err
		}
		for _, path := range c.paths {
			vip := viper.New()
			if err := vip.AddRemoteProvider("nacos", c.endpoint, path.String()); err != nil {
				return err
			}
			vip.SetConfigType(path.Format)
			if err := vip.ReadRemoteConfig(); err != nil {
				return err
			}
			c.viperMap.Store(path, vip)
		}
		return nil
	}
	c.loadOnce.Do(func() {
		namespaceSet := gfcontainer.NewSortedSet[string]()
		for _, p := range c.paths {
			namespaceSet.Add(p.Namespace)
		}
		c.loadErr = loadFn(namespaceSet.List())
	})
	return c.loadErr
}

func (c *NacosConfigure) WatchConfig(path string, watchFn func(vip *viper.Viper)) error {
	ctx := context.Background()
	rp := nacos.NewNacosRemoteProvider(c.endpoint, path)
	rspCh, _ := viper.RemoteConfig.WatchChannel(rp)
	vip, err := c.GetConfig(path)
	if err != nil {
		return err
	}
	go func() {
		for rsp := range rspCh {
			if rsp.Error != nil {
				gflog.Warn(ctx, "[NacosConfigure] config OnChange failed", zap.Error(rsp.Error))
				continue
			}
			if err := vip.ReadConfig(bytes.NewReader(rsp.Value)); err != nil {
				gflog.Warn(ctx, "[NacosConfigure] config OnChange failed", zap.Error(err))
				continue
			}
			watchFn(vip)
		}
	}()
	return nil
}
