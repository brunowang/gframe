package srv

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
)

type GrpcServer struct {
	plugin *protogen.Plugin
	goPkg  string
	name   string
}

func (a *GrpcServer) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
	a.goPkg = "frontend"
	a.name = "grpcserver"
}

func (a *GrpcServer) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		projName := string(file.GoPackageName)
		fdir := helper.GetFileBaseDir(file, config)
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).
			Import(string(file.GoImportPath)).
			Import("net").Import("go.uber.org/zap").
			Import("github.com/brunowang/gframe/gflog").
			Import(fdir + "/service").Import("google.golang.org/grpc").
			Import("github.com/grpc-ecosystem/go-grpc-middleware").
			Import("google.golang.org/grpc/reflection")

		fpath := fmt.Sprintf("%s/%s/%s.go", fdir, a.goPkg, a.name)
		g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
		g.P(fhead)
		for _, svc := range file.Services {
			tmpl := GrpcServerTmpl{
				ProjName: projName,
				SvcName:  svc.GoName,
			}
			g.P(tmpl.Render())
		}
	}
}
