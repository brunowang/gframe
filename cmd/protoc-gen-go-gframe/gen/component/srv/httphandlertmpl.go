package srv

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type HttpHandlerTmpl struct {
	ProjName string
	SvcName  string
}

//go:embed httphandler.tmpl
var httphandlertmpl string

func (t *HttpHandlerTmpl) Render() string {
	return helper.NewTmplRenderer("httphandler.tmpl").Text(httphandlertmpl).Data(t).RenderTmpl()
}
