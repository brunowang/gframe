package srv

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

type Server struct {
	plugin *protogen.Plugin
	goPkg  string
	name   string
}

func (a *Server) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
	a.goPkg = "frontend"
	a.name = "server"
}

func (a *Server) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		importDomain := strings.Split(string(file.GoImportPath), "/")[0]
		projName := string(file.GoPackageName)
		fdir := fmt.Sprintf("%s/projects/%s", importDomain, projName)
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).Import("fmt").
			Import("net").Import("go.uber.org/zap").
			Import("github.com/brunowang/gframe/gflog").
			Import(fdir + "/service").Import("github.com/soheilhy/cmux")

		fpath := fmt.Sprintf("%s/%s/%s.go", fdir, a.goPkg, a.name)
		g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
		g.P(fhead)
		for _, svc := range file.Services {
			tmpl := ServerTmpl{
				SvcName: svc.GoName,
			}
			g.P(tmpl.Render())
		}
	}
}
