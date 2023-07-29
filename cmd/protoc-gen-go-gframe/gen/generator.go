package gen

import (
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
	"log"
	"sync"
)

type ProjectGenerator struct {
	plugin  *protogen.Plugin
	once    sync.Once
	subGens []Component
	config  helper.GenerateConfig
}

func NewProjectGenerator(plugin *protogen.Plugin) *ProjectGenerator {
	return &ProjectGenerator{
		plugin: plugin,
	}
}

func (g *ProjectGenerator) Generate(opts ...Option) {
	g.initOnce(opts...)

	for _, sub := range g.subGens {
		sub.Generate(g.config)
	}
}

func (g *ProjectGenerator) initOnce(opts ...Option) {
	g.once.Do(func() {
		options := genOptions{}
		for _, opt := range opts {
			opt(&options)
		}
		for _, name := range options.components {
			if com, ok := GetComponent(name); ok {
				com.Setup(g.plugin)
				g.subGens = append(g.subGens, com)
			} else {
				log.Printf("component %s not found\n", name)
			}
		}
		g.config = helper.GenerateConfig{
			Project: options.project,
			PbGoDir: options.pbGoDir,
			ModPath: options.modPath,
		}
	})
}
