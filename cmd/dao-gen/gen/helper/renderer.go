package helper

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type TmplRenderer struct {
	name string
	text string
	data interface{}
}

func NewTmplRenderer(name string) *TmplRenderer {
	return &TmplRenderer{name: name}
}

func (r *TmplRenderer) Text(text string) *TmplRenderer {
	r.text = text
	return r
}

func (r *TmplRenderer) Data(data interface{}) *TmplRenderer {
	r.data = data
	return r
}

func (r *TmplRenderer) RenderTmpl() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New(r.name).Parse(strings.TrimSpace(r.text))
	if err != nil {
		return fmt.Errorf("render template failed %v", err).Error()
	}
	if err := tmpl.ExecuteTemplate(buf, r.name, r.data); err != nil {
		return fmt.Errorf("render template failed %v", err).Error()
	}
	return buf.String()
}
