package proj

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
)

type ConfigTmpl struct{}

//go:embed config.tmpl
var configtmpl string

func (t *ConfigTmpl) Render() string {
	return helper.NewTmplRenderer("config.tmpl").Text(configtmpl).Data(t).RenderTmpl()
}
