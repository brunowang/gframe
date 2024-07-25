package gfconf

import (
	"github.com/brunowang/gframe/gfcontainer"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type LocalConfigure struct {
	paths    []string
	viperMap gfcontainer.SafeMap[string, *viper.Viper]
	loadOnce sync.Once
	loadErr  error
}

func NewLocalConfigure(paths []string) *LocalConfigure {
	return &LocalConfigure{paths: paths}
}

func (c *LocalConfigure) GetConfig(path string) (*viper.Viper, error) {
	path = toAbsolutePath(path)
	ret, has := c.viperMap.Load(path)
	if !has {
		return nil, ErrorConfigNotLoaded
	}
	return ret, nil
}

func (c *LocalConfigure) Init() error {
	if len(c.paths) == 0 {
		return ErrorConfigInitFailed
	}
	loadFn := func() error {
		for _, path := range c.paths {
			vip := viper.New()
			path = toAbsolutePath(path)
			dir, file := filepath.Split(path)
			ext := filepath.Ext(file)
			vip.SetConfigName(strings.TrimSuffix(file, ext))
			vip.SetConfigType(strings.TrimPrefix(ext, "."))
			vip.AddConfigPath(dir)
			if err := vip.ReadInConfig(); err != nil {
				return err
			}
			c.viperMap.Store(path, vip)
		}
		return nil
	}
	c.loadOnce.Do(func() {
		c.loadErr = loadFn()
	})
	return c.loadErr
}

func (c *LocalConfigure) WatchConfig(path string, watchFn func(vip *viper.Viper)) error {
	path = toAbsolutePath(path)
	vip, err := c.GetConfig(path)
	if err != nil {
		return err
	}
	go func() {
		vip.WatchConfig()
		vip.OnConfigChange(func(e fsnotify.Event) {
			watchFn(vip)
		})
	}()
	return nil
}

func toAbsolutePath(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}
	workDir, err := os.Getwd()
	if err != nil {
		// TODO log
		return path
	}
	return workDir + "/" + path
}
