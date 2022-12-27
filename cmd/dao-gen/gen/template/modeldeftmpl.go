package template

import (
	_ "embed"
	"fmt"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/helper"
)

type Field struct {
	Name    string
	Type    string
	ZeroVal string
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

type ModelDefTmpl struct {
	Name   string
	Fields Fields
}

//go:embed modeldef.tmpl
var modeldeftmpl string

func (t *ModelDefTmpl) Render() string {
	return helper.NewTmplRenderer("modeldef.tmpl").
		Text(modeldeftmpl).Data(t).RenderTmpl()
}
