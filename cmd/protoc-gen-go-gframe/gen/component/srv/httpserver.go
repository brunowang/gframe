package srv

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
)

type HttpServer struct {
	plugin *protogen.Plugin
	goPkg  string
	name   string
}

func (a *HttpServer) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
	a.goPkg = "frontend"
	a.name = "httpserver"
}

func (a *HttpServer) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		projName := string(file.GoPackageName)
		fdir := helper.GetFileBaseDir(file, config)
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).
			Import("net").Import("go.uber.org/zap").
			Import("github.com/gin-gonic/gin").
			Import("github.com/brunowang/gframe/gflog").
			Import(fdir + "/service")

		fpath := fmt.Sprintf("%s/%s/%s.go", fdir, a.goPkg, a.name)
		g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
		g.P(fhead)
		for _, svc := range file.Services {
			tmpl := HttpServerTmpl{
				ProjName: projName,
				SvcName:  svc.GoName,
			}
			g.P(tmpl.Render())
		}
	}
}
