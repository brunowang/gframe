package gfconf

import (
	"sync"
)

type NacosConf struct {
	Enabled  bool     `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Host     string   `json:"host" yaml:"host" mapstructure:"host"`
	Port     uint64   `json:"port" yaml:"port" mapstructure:"port"`
	Username string   `json:"username" yaml:"username" mapstructure:"username"`
	Password string   `json:"password" yaml:"password" mapstructure:"password"`
	Paths    []string `json:"paths" yaml:"paths" mapstructure:"paths"`
}

type LoaderRegister interface {
	Register() error
}

var (
	once sync.Once
)

func MustInitNacos(nacos NacosConf, register LoaderRegister) {
	if err := initOnce(nacos, register); err != nil {
		panic(err)
	}
}

func initOnce(nacos NacosConf, register LoaderRegister) error {
	var initErr error
	once.Do(func() {
		initErr = initConfig(nacos, register)
	})
	if initErr != nil {
		return initErr
	}
	return nil
}

func initConfig(nacos NacosConf, register LoaderRegister) error {
	if err := register.Register(); err != nil {
		return err
	}
	remoteCfg := RemoteConfig{
		Enabled:  nacos.Enabled,
		Host:     nacos.Host,
		Port:     nacos.Port,
		Username: nacos.Username,
		Password: nacos.Password,
	}
	return LoadConfig(nacos.Paths, &remoteCfg)
}
