package srv

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

type Service struct {
	plugin *protogen.Plugin
	goPkg  string
}

func (a *Service) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
	a.goPkg = "service"
}

func (a *Service) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		importDomain := strings.Split(string(file.GoImportPath), "/")[0]
		projName := string(file.GoPackageName)
		fdir := fmt.Sprintf("%s/projects/%s", importDomain, projName)
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).
			Import(fdir + "/dto").Import("context")

		for _, svc := range file.Services {
			fpath := fmt.Sprintf("%s/%s/%s.go",
				fdir, a.goPkg, strings.ToLower(svc.GoName))
			g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
			g.P(fhead)
			tmpl := ServiceTmpl{
				SvcName: svc.GoName,
			}
			for _, method := range svc.Methods {
				tmpl.Handlers = append(tmpl.Handlers, Handler{
					Method:            method.GoName,
					Request:           helper.ToCamelCase(string(method.Input.Desc.Name())),
					Response:          helper.ToCamelCase(string(method.Output.Desc.Name())),
					IsStreamingClient: method.Desc.IsStreamingClient(),
					IsStreamingServer: method.Desc.IsStreamingServer(),
				})
			}
			g.P(tmpl.Render())
		}
	}
}
