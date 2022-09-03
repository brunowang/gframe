package sdk

import (
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
)

type GrpcSDK struct {
	plugin *protogen.Plugin
}

func (g *GrpcSDK) Setup(plugin *protogen.Plugin) {
	g.plugin = plugin
}

func (g *GrpcSDK) Generate(config helper.GenerateConfig) {

}
