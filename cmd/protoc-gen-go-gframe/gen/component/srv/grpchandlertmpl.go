package srv

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type GrpcHandlerTmpl struct {
	ProjName      string
	ProjNameUpper string
	SvcName       string
	Handlers      []Handler
}

//go:embed grpchandler.tmpl
var grpchandlertmpl string

func (t *GrpcHandlerTmpl) Render() string {
	return helper.NewTmplRenderer("grpchandler.tmpl").Text(grpchandlertmpl).Data(t).RenderTmpl()
}
