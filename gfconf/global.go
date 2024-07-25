package gfconf

import (
	"fmt"
	"github.com/brunowang/gframe/gfconf/internal/nacos"
	"github.com/brunowang/gframe/gfcontainer"
	"github.com/spf13/viper"
	"path/filepath"
)

var (
	ErrorConfigNotLoaded  = fmt.Errorf("config not loaded")
	ErrorConfigInitFailed = fmt.Errorf("config init failed")
)

type Configure interface {
	Init() error
	GetConfig(path string) (*viper.Viper, error)
	WatchConfig(path string, watchFn func(vip *viper.Viper)) error
}

type BasicAuth = nacos.BasicAuth

var loaderMap = gfcontainer.SafeMap[string, func(*viper.Viper)]{}

func RegisterLoader(filename string, loader func(*viper.Viper)) error {
	if _, has := loaderMap.Load(filename); has {
		return fmt.Errorf("got repeat loader filename: %s", filename)
	}
	loaderMap.Store(filename, loader)
	return nil
}

type RemoteConfig struct {
	Enabled  bool
	Host     string
	Port     uint64
	Username string
	Password string
}

func MustLoadConfig(paths []string, remoteCfg *RemoteConfig) {
	if err := LoadConfig(paths, remoteCfg); err != nil {
		panic(err)
	}
}

func LoadConfig(paths []string, remoteCfg *RemoteConfig) error {
	if len(paths) == 0 {
		return fmt.Errorf("got empty config paths")
	}
	var cfg Configure
	if remoteCfg != nil && remoteCfg.Enabled {
		endpoint := fmt.Sprintf("%s:%d", remoteCfg.Host, remoteCfg.Port)
		auth := BasicAuth{
			Username: remoteCfg.Username,
			Password: remoteCfg.Password,
		}
		cfg = NewNacosConfigure(paths, endpoint, auth)
	} else {
		cfg = NewLocalConfigure(paths)
	}
	if err := cfg.Init(); err != nil {
		return fmt.Errorf("[InitConfigure] init config got error: %v, paths: %v\n", err, paths)
	}

	for _, path := range paths {
		vip, err := cfg.GetConfig(path)
		if err != nil {
			return fmt.Errorf("[InitConfigure] get config got error: %v, path: %s\n", err, path)
		}
		filename := filepath.Base(path)
		loadFunc, has := loaderMap.Load(filename)
		if !has {
			return fmt.Errorf("[InitConfigure] loader filename not found: %s\n", filename)
		}

		loadFunc(vip)
		if err := cfg.WatchConfig(path, func(vip *viper.Viper) {
			loadFunc(vip)
		}); err != nil {
			return fmt.Errorf("[InitConfigure] watch config got error: %v, path: %s\n", err, path)
		}
	}
	return nil
}
