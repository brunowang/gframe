package internal

import (
	"github.com/brunowang/gframe/gfconf"
	"github.com/spf13/viper"
	"sync/atomic"
)

type HotConf struct {
	Application Application
	Whitelist   Whitelist
}

var (
	hotConf = atomic.Pointer[HotConf]{}
)

func GetHotConf() *HotConf {
	conf := hotConf.Load()
	if conf == nil {
		return &HotConf{}
	}
	return conf
}

func CopyHotConf() *HotConf {
	copyConf := *GetHotConf()
	return &copyConf
}

type HotConfRegister struct{}

func (HotConfRegister) Register() error {
	if err := gfconf.RegisterLoader("app.yaml", func(vip *viper.Viper) {
		copyConf := CopyHotConf()
		vip.Unmarshal(&copyConf.Application)
		hotConf.Store(copyConf)
	}); err != nil {
		return err
	}
	if err := gfconf.RegisterLoader("whitelist.yaml", func(vip *viper.Viper) {
		copyConf := CopyHotConf()
		vip.UnmarshalKey("whitelist", &copyConf.Whitelist)
		hotConf.Store(copyConf)
	}); err != nil {
		return err
	}
	// TODO 可以增加更多 gfconf.RegisterLoader
	return nil
}
