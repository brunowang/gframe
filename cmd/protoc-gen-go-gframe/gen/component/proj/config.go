package proj

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
)

type Config struct {
	plugin *protogen.Plugin
	goPkg  string
	name   string
}

func (a *Config) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
	a.goPkg = "conf"
	a.name = "configure"
}

func (a *Config) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		fdir := helper.GetFileBaseDir(file, config)
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).
			Import("github.com/BurntSushi/toml")

		fpath := fmt.Sprintf("%s/%s/%s.go", fdir, a.goPkg, a.name)
		g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
		g.P(fhead)

		tmpl := ConfigTmpl{}
		g.P(tmpl.Render())
	}
}
