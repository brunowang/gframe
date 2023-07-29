package srv

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
)

type HttpHandler struct {
	plugin *protogen.Plugin
	goPkg  string
	name   string
}

func (a *HttpHandler) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
	a.goPkg = "frontend"
	a.name = "httphandler"
}

func (a *HttpHandler) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		projName := string(file.GoPackageName)
		fdir := helper.GetFileBaseDir(file, config)
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).
			Import("github.com/gin-gonic/gin").
			Import(fdir+"/dto").Import("go.uber.org/zap").
			Import("github.com/brunowang/gframe/gflog").
			Import("github.com/brunowang/gframe/gfhttp").
			Import(fdir+"/service").Import("time").
			ImportWithAlias("github.com/gorilla/websocket", "ws").
			Import("net/http").
			Import("sync").Import("encoding/json")

		fpath := fmt.Sprintf("%s/%s/%s.go", fdir, a.goPkg, a.name)
		g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
		fnMap := make(map[string]struct{})
		if helper.Exists(fpath) {
			var gosrc []byte
			gosrc, fnMap = helper.ParseGoFile(fpath)
			g.P(string(gosrc))
		} else {
			g.P(fhead)
			for _, svc := range file.Services {
				tmpl := HttpHandlerTmpl{
					ProjName: projName,
					SvcName:  svc.GoName,
				}
				g.P(tmpl.Render())
			}
		}
		for _, svc := range file.Services {
			for _, method := range svc.Methods {
				if _, ok := fnMap[method.GoName]; ok {
					continue
				}
				hand := Handler{
					Method:            method.GoName,
					Request:           helper.ToCamelCase(string(method.Input.Desc.Name())),
					Response:          helper.ToCamelCase(string(method.Output.Desc.Name())),
					IsStreamingClient: method.Desc.IsStreamingClient(),
					IsStreamingServer: method.Desc.IsStreamingServer(),
				}
				tmpl := HttpMethodTmpl{
					Handler: hand,
				}
				g.P(tmpl.Render())
			}
		}
	}
}
