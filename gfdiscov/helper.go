package gfdiscov

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
	"os"
	"strings"
)

const (
	allEths    = "0.0.0.0"
	envPodIp   = "POD_IP"
	etcdScheme = "etcd"
)

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = internalIP()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}

// internalIP returns an internal ip.
func internalIP() string {
	infs, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, inf := range infs {
		if isEthDown(inf.Flags) || isLoopback(inf.Flags) {
			continue
		}

		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return ""
}

func isEthDown(f net.Flags) bool {
	return f&net.FlagUp != net.FlagUp
}

func isLoopback(f net.Flags) bool {
	return f&net.FlagLoopback == net.FlagLoopback
}

func makeEtcdKey(key string, lease clientv3.LeaseID) string {
	return fmt.Sprintf("%s/%d", key, lease)
}

func makeKeyPrefix(key string) string {
	if strings.HasSuffix(key, "/") {
		return key
	}
	return key + "/"
}

// BuildEtcdTarget returns a string that represents the given endpoints with etcd schema.
func BuildEtcdTarget(endpoints []string, key string) string {
	return fmt.Sprintf("%s://%s/%s", etcdScheme, strings.Join(endpoints, ","), key)
}
