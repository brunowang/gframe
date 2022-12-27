package template

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/helper"
)

type CacheDefTmpl struct{}

//go:embed cachedef.tmpl
var cachedeftmpl string

func (t *CacheDefTmpl) Render() string {
	return helper.NewTmplRenderer("cachedef.tmpl").
		Text(cachedeftmpl).Data(t).RenderTmpl()
}
