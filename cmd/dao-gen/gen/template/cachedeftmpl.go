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

//go:embed cachelocal.tmpl
var cachelocaltmpl string

func (t *CacheDefTmpl) RenderLocal() string {
	return helper.NewTmplRenderer("cachelocal.tmpl").
		Text(cachelocaltmpl).Data(t).RenderTmpl()
}

//go:embed cacheredis.tmpl
var cacheredistmpl string

func (t *CacheDefTmpl) RenderRedis() string {
	return helper.NewTmplRenderer("cacheredis.tmpl").
		Text(cacheredistmpl).Data(t).RenderTmpl()
}

//go:embed cacheserial.tmpl
var cacheserialtmpl string

func (t *CacheDefTmpl) RenderSerial() string {
	return helper.NewTmplRenderer("cacheserial.tmpl").
		Text(cacheserialtmpl).Data(t).RenderTmpl()
}
