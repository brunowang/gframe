package proj

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
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
		projName := string(file.GoPackageName)

		importDomain := strings.Split(string(file.GoImportPath), "/")[0]
		fdir := fmt.Sprintf("%s/projects/%s", importDomain, projName)
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
		g.P("\tgithub.com/golang/protobuf v1.5.2")
		g.P("\tgithub.com/grpc-ecosystem/go-grpc-middleware v1.3.0")
		g.P("\tgithub.com/soheilhy/cmux v0.1.5")
		g.P("\tgo.uber.org/zap v1.22.0")
		g.P("\tgopkg.in/alecthomas/kingpin.v2 v2.2.6")
		g.P(")")
		g.P("require (")
		g.P("\tgoogle.golang.org/grpc v1.48.0 // indirect")
		g.P(")")
		break
	}
}
