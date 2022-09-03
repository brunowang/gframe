package dto

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type ParamsTpl struct {
	ProjName string
	Messages []Message
}

type Message struct {
	Name   string
	IsReq  bool
	IsRsp  bool
	Fields []Field
}

type Field struct {
	Name       string
	CheckLen   bool
	DefaultVal string
}

//go:embed params.tmpl
var paramstmpl string

func (t *ParamsTpl) Render() string {
	return helper.NewTmplRenderer("params.tmpl").Text(paramstmpl).Data(t).RenderTmpl()
}
