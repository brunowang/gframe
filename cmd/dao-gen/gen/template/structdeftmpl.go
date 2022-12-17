package template

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/helper"
)

type Field struct {
	Name    string
	Type    string
	ColName string
	Comment string
}

type StructDefTmpl struct {
	Name    string
	Fields  []Field
	TabName string
}

//go:embed structdef.tmpl
var structdeftmpl string

func (t *StructDefTmpl) Render() string {
	return helper.NewTmplRenderer("structdef.tmpl").
		Text(structdeftmpl).Data(t).RenderTmpl()
}
