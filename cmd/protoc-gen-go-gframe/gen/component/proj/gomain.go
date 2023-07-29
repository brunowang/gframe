package proj

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
)

type GoMain struct {
	plugin *protogen.Plugin
	goPkg  string
	name   string
}

func (a *GoMain) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
	a.goPkg = "main"
	a.name = "main"
}

func (a *GoMain) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		projName := string(file.GoPackageName)
		fdir := helper.GetFileBaseDir(file, config)
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).
			Import("gopkg.in/alecthomas/kingpin.v2").
			Import("github.com/brunowang/gframe/gflog").
			Import("os").Import("os/signal").
			Import("syscall").Import("runtime").
			Import(fdir + "/frontend").Import(fdir + "/conf").
			Import("go.uber.org/zap")

		fpath := fmt.Sprintf("%s/%s.go", fdir, a.name)
		g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
		g.P(fhead)

		tmpl := GoMainTmpl{
			ProjName: projName,
		}
		g.P(tmpl.Render())
	}
}
