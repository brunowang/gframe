package template

import (
	_ "embed"
	"fmt"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/helper"
)

type Field struct {
	Name    string
	Type    string
	ColName string
	Comment string
}

type Fields []Field

func (a *Fields) Find(name string) (Field, error) {
	for _, f := range *a {
		if f.Name == name {
			return f, nil
		}
	}
	return Field{}, fmt.Errorf("field not found")
}

type Index struct {
	Uniq bool
	Cols Fields
}

type StructDefTmpl struct {
	Name    string
	TabName string
	Fields  Fields
	Indexes []Index
}

//go:embed structdef.tmpl
var structdeftmpl string

func (t *StructDefTmpl) Render() string {
	return helper.NewTmplRenderer("structdef.tmpl").
		Text(structdeftmpl).Data(t).RenderTmpl()
}
