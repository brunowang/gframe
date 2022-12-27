package template

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/helper"
)

type CurdDefTmpl struct {
	Name    string
	TabName string
	Fields  Fields
	Indexes []Index
}

//go:embed curddef.tmpl
var curddeftmpl string

func (t *CurdDefTmpl) Render() string {
	return helper.NewTmplRenderer("curddef.tmpl").
		Text(curddeftmpl).Data(t).RenderTmpl()
}
