package nacos

import (
	"context"
	"fmt"
	"github.com/brunowang/gframe/gfcontainer"
	"github.com/brunowang/gframe/gflog"
	"github.com/nacos-group/nacos-sdk-go/clients"
	nacosCli "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"strconv"
	"strings"
)

type NacosMgr struct {
	clients gfcontainer.SafeMap[string, nacosCli.IConfigClient]
	cache   gfcontainer.SafeMap[string, []byte]
}

type BasicAuth struct {
	Username string
	Password string
}

type PathParsed struct {
	IsRelative bool
	Namespace  string
	Group      string
	DataID     string
	Format     string
}

func (p PathParsed) String() string {
	ret := strings.Join([]string{p.Namespace, p.Group, p.DataID}, "/")
	if p.IsRelative {
		return ret
	}
	return "/" + ret
}

type Path string

func (p Path) Parse() PathParsed {
	ret := PathParsed{
		IsRelative: true,
		Namespace:  "public",
		Group:      "DEFAULT_GROUP",
		DataID:     "DEFAULT_DATA_ID",
		Format:     "json",
	}
	path := string(p)
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
		ret.IsRelative = false
	}
	fileExt := filepath.Ext(path)
	if fileExt != "" {
		ret.Format = strings.TrimPrefix(fileExt, ".")
	}
	arr := strings.Split(path, "/")
	if len(arr) == 1 {
		ret.DataID = arr[0]
	} else if len(arr) == 2 {
		ret.Group = arr[0]
		ret.DataID = arr[1]
	} else if len(arr) >= 3 {
		ret.Namespace = arr[0]
		ret.Group = arr[1]
		ret.DataID = strings.Join(arr[2:], "/")
	}
	return ret
}

func NewNacosMgr(namespaces, endpoints []string, auth ...BasicAuth) (*NacosMgr, error) {
	ctx := context.Background()
	if len(namespaces) == 0 {
		namespaces = []string{"public"}
	}
	if len(endpoints) == 0 {
		endpoints = []string{"127.0.0.1:8848"}
	}
	ret := new(NacosMgr)
	srvConfigs := make([]constant.ServerConfig, 0, len(endpoints))
	for _, ep := range endpoints {
		hostPort := strings.Split(ep, ":")
		host, port := hostPort[0], uint64(80)
		if len(hostPort) == 2 {
			port, _ = strconv.ParseUint(hostPort[1], 10, 64)
		}
		srvConfigs = append(srvConfigs, constant.ServerConfig{
			ContextPath: "/nacos",
			IpAddr:      host,
			Port:        port,
		})
	}
	for _, ns := range namespaces {
		cliConfig := constant.ClientConfig{
			NamespaceId:         ns,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogLevel:            "info",
			CacheDir:            ".nacos/cache",
			LogDir:              ".nacos/log",
			LogRollingConfig: &lumberjack.Logger{
				MaxAge: 3,
			},
		}
		if len(auth) > 0 {
			cliConfig.Username = auth[0].Username
			cliConfig.Password = auth[0].Password
		}
		client, err := clients.NewConfigClient(vo.NacosClientParam{
			ServerConfigs: srvConfigs,
			ClientConfig:  &cliConfig,
		})
		if err != nil {
			gflog.Warn(ctx, "[NacosMgr] CreateConfigClient failed", zap.Error(err),
				zap.Any("srvConfigs", srvConfigs), zap.Any("cliConfig", cliConfig))
			return nil, err
		}
		ret.clients.Store(ns, client)
	}
	return ret, nil
}

func (m *NacosMgr) GetData(path string) ([]byte, error) {
	pathParsed := Path(path).Parse()
	path = pathParsed.String()

	if bs, ok := m.getFromCache(pathParsed); ok {
		return bs, nil
	}
	bs, err := m.getFromRemote(pathParsed)
	if err != nil {
		return nil, err
	}
	m.cache.Store(path, bs)
	return bs, nil
}

func (m *NacosMgr) getFromCache(path PathParsed) ([]byte, bool) {
	return m.cache.Load(path.String())
}

func (m *NacosMgr) getFromRemote(path PathParsed) ([]byte, error) {
	client, ok := m.clients.Load(path.Namespace)
	if !ok {
		return nil, fmt.Errorf("namespace %s not found", path.Namespace)
	}
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: path.DataID,
		Group:  path.Group,
	})
	return []byte(content), err
}

func (m *NacosMgr) watchRemote(path string, stop <-chan bool) (<-chan *viper.RemoteResponse, error) {
	ctx := context.Background()
	pathParsed := Path(path).Parse()
	path = pathParsed.String()
	rsp := make(chan *viper.RemoteResponse)

	configParams := vo.ConfigParam{
		DataId: pathParsed.DataID,
		Group:  pathParsed.Group,
		OnChange: func(namespace, group, dataID, data string) {
			gflog.Info(ctx, "[NacosMgr] config OnChange", zap.String("dataID", dataID), zap.String("data", data))
			m.cache.Store(dataID, []byte(data))
			rsp <- &viper.RemoteResponse{Value: []byte(data)}
		},
	}

	client, ok := m.clients.Load(pathParsed.Namespace)
	if !ok {
		return nil, fmt.Errorf("namespace %s not found", pathParsed.Namespace)
	}
	err := client.ListenConfig(configParams)
	if err != nil {
		gflog.Warn(ctx, "[NacosMgr] ListenConfig failed", zap.Any("params", configParams))
		return nil, err
	}
	if stop != nil {
		go func() {
			for {
				select {
				case <-stop:
					_ = client.CancelListenConfig(configParams)
					return
				}
			}
		}()
	}
	return rsp, nil
}
