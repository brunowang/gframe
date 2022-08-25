package srv

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type GrpcServerTmpl struct {
	ProjName string
	SvcName  string
}

//go:embed grpcserver.tmpl
var grpcservertmpl string

func (t *GrpcServerTmpl) Render() string {
	return helper.NewTmplRenderer("grpcserver.tmpl").Text(grpcservertmpl).Data(t).RenderTmpl()
}
