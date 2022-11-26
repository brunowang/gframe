package srv

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type GrpcMethodTmpl struct {
	Handler
	ProjName string
	SvcName  string
}

//go:embed grpcmethod.tmpl
var grpcmethodtmpl string

func (t *GrpcMethodTmpl) Render() string {
	return helper.NewTmplRenderer("grpcmethod.tmpl").Text(grpcmethodtmpl).Data(t).RenderTmpl()
}
