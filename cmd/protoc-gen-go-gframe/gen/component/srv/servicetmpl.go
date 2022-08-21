package srv

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type ServiceTmpl struct {
	SvcName  string
	Handlers []Handler
}

type Handler struct {
	Path              string
	Method            string
	Request           string
	Response          string
	IsStreamingClient bool
	IsStreamingServer bool
}

//go:embed service.tmpl
var servicetmpl string

func (t *ServiceTmpl) Render() string {
	return helper.NewTmplRenderer("service.tmpl").Text(servicetmpl).Data(t).RenderTmpl()
}
