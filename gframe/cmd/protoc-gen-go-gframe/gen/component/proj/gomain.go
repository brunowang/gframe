package proj

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
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
		importDomain := strings.Split(string(file.GoImportPath), "/")[0]
		projName := string(file.GoPackageName)
		fdir := fmt.Sprintf("%s/projects/%s", importDomain, projName)
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).
			Import("gopkg.in/alecthomas/kingpin.v2").
			Import("github.com/brunowang/gframe/gflog").
			Import("os").Import("os/signal").
			Import("syscall").Import("runtime").
			Import(fdir + "/frontend").Import(fdir + "/conf")

		fpath := fmt.Sprintf("%s/cmd/%s/%s.go", fdir, projName, a.name)
		g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
		g.P(fhead)

		tmpl := GoMainTmpl{
			ProjName: projName,
		}
		g.P(tmpl.Render())
	}
}
