package proj

import (
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
)

type GoMod struct {
	plugin *protogen.Plugin
}

func (a *GoMod) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
}

func (a *GoMod) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		fdir := helper.GetFileBaseDir(file, config)
		fpath := fdir + "/go.mod"
		pbGoDir := string(file.GoImportPath)
		if config.PbGoDir != "" {
			pbGoDir = config.PbGoDir
		}
		g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)

		g.P("module " + fdir)
		g.P()
		g.P("go 1.18")
		g.P()
		g.P("require (")
		g.P("\t" + pbGoDir + " latest")
		g.P("\tgithub.com/brunowang/gframe latest")
		g.P("\tgithub.com/gin-gonic/gin v1.9.1")
		g.P("\tgithub.com/golang/protobuf v1.5.2")
		g.P("\tgithub.com/gorilla/websocket v1.4.2")
		g.P("\tgithub.com/grpc-ecosystem/go-grpc-middleware v1.3.0")
		g.P("\tgithub.com/soheilhy/cmux v0.1.5")
		g.P("\tgo.uber.org/zap v1.22.0")
		g.P("\tgoogle.golang.org/grpc v1.49.0")
		g.P("\tgopkg.in/alecthomas/kingpin.v2 v2.2.6")
		g.P(")")
		break
	}
}
