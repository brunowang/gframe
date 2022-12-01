package srv

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type HttpServerTmpl struct {
	ProjName string
	SvcName  string
	Handlers []Handler
}

//go:embed httpserver.tmpl
var httpservertmpl string

func (t *HttpServerTmpl) Render() string {
	return helper.NewTmplRenderer("httpserver.tmpl").Text(httpservertmpl).Data(t).RenderTmpl()
}
