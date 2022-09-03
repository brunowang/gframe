package srv

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type ServerTmpl struct {
	SvcName string
}

//go:embed server.tmpl
var servertmpl string

func (t *ServerTmpl) Render() string {
	return helper.NewTmplRenderer("server.tmpl").Text(servertmpl).Data(t).RenderTmpl()
}
