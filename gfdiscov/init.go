package gfdiscov

import "google.golang.org/grpc/resolver"

func init() {
	resolver.Register(NewEtcdResolver())
}
