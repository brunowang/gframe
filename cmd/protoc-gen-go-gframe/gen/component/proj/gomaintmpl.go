package proj

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type GoMainTmpl struct {
	ProjName string
}

//go:embed gomain.tmpl
var gomaintmpl string

func (t *GoMainTmpl) Render() string {
	return helper.NewTmplRenderer("gomain.tmpl").Text(gomaintmpl).Data(t).RenderTmpl()
}
