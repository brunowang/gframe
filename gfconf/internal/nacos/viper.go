package nacos

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io"
)

var Manager *NacosMgr

func Init(namespaces, endpoints []string, auth ...BasicAuth) error {
	mgr, err := NewNacosMgr(namespaces, endpoints, auth...)
	if err != nil {
		return err
	}
	Manager = mgr
	return nil
}

func SupportViper() error {
	if Manager == nil {
		return fmt.Errorf("nacos manager not initialized")
	}
	viper.SupportedRemoteProviders = []string{"nacos"}
	viper.RemoteConfig = &RemoteConfigProvider{ConfigManager: Manager}
	return nil
}

type RemoteConfigProvider struct {
	ConfigManager *NacosMgr
}

func (r *RemoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	if err := check(rp); err != nil {
		return nil, err
	}
	bs, err := r.ConfigManager.GetData(rp.Path())
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bs), nil
}

func (r *RemoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	return r.Get(rp)
}

func (r *RemoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	if err := check(rp); err != nil {
		return nil, nil
	}
	quit := make(chan bool)
	rspCh, _ := r.ConfigManager.watchRemote(rp.Path(), quit)
	return rspCh, quit
}

func check(rp viper.RemoteProvider) error {
	if rp.Provider() != "nacos" {
		return fmt.Errorf("unexpected remote provider %s of nacos manager", rp.Provider())
	}
	return nil
}

type NacosRemoteProvider struct {
	endpoint string
	path     string
}

func NewNacosRemoteProvider(endpoint, path string) NacosRemoteProvider {
	return NacosRemoteProvider{endpoint: endpoint, path: path}
}

func (rp NacosRemoteProvider) Provider() string {
	return "nacos"
}

func (rp NacosRemoteProvider) Endpoint() string {
	return rp.endpoint
}

func (rp NacosRemoteProvider) Path() string {
	return rp.path
}

func (rp NacosRemoteProvider) SecretKeyring() string {
	return ""
}
