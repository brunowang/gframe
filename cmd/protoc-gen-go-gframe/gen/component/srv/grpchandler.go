package srv

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

type GrpcHandler struct {
	plugin *protogen.Plugin
	goPkg  string
	name   string
}

func (a *GrpcHandler) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
	a.goPkg = "frontend"
	a.name = "grpchandler"
}

func (a *GrpcHandler) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		importDomain := strings.Split(string(file.GoImportPath), "/")[0]
		projName := string(file.GoPackageName)
		fdir := fmt.Sprintf("%s/projects/%s", importDomain, projName)
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).
			Import(string(file.GoImportPath)).
			Import(fdir + "/dto").Import("go.uber.org/zap").
			Import("github.com/brunowang/gframe/gflog").
			Import(fdir + "/service").Import("io").
			Import("context").Import("time")

		fpath := fmt.Sprintf("%s/%s/%s.go", fdir, a.goPkg, a.name)
		g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
		g.P(fhead)
		for _, svc := range file.Services {
			tmpl := GrpcHandlerTmpl{
				ProjName:      projName,
				ProjNameUpper: strings.ToUpper(projName),
				SvcName:       svc.GoName,
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
