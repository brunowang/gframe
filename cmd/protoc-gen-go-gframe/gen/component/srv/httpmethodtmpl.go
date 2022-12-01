package srv

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type HttpMethodTmpl struct {
	Handler
}

//go:embed httpmethod.tmpl
var httpmethodtmpl string

func (t *HttpMethodTmpl) Render() string {
	return helper.NewTmplRenderer("httpmethod.tmpl").Text(httpmethodtmpl).Data(t).RenderTmpl()
}
