package srv

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

type GrpcServer struct {
	plugin *protogen.Plugin
	goPkg  string
}

func (a *GrpcServer) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
	a.goPkg = "frontend"
}

func (a *GrpcServer) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		importDomain := strings.Split(string(file.GoImportPath), "/")[0]
		projName := string(file.GoPackageName)
		fdir := fmt.Sprintf("%s/projects/%s", importDomain, projName)
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).
			Import(string(file.GoImportPath)).Import("context").
			Import("net").Import("go.uber.org/zap").
			Import("github.com/brunowang/gframe/gflog").
			Import(fdir + "/service").Import("google.golang.org/grpc").
			Import("github.com/grpc-ecosystem/go-grpc-middleware").
			Import("google.golang.org/grpc/reflection")

		for _, svc := range file.Services {
			fpath := fmt.Sprintf("%s/%s/%s.go",
				fdir, a.goPkg, "grpcserver")
			g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
			g.P(fhead)
			tmpl := GrpcServerTmpl{
				ProjName: projName,
				SvcName:  svc.GoName,
			}
			g.P(tmpl.Render())
		}
	}
}
