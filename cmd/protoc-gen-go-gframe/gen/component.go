package gen

import (
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/component/dto"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/component/proj"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/component/sdk"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/component/srv"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	ComponentGoMod       = "go_mod"
	ComponentParams      = "params"
	ComponentService     = "service"
	ComponentGrpcHandler = "grpc_handler"
	ComponentGrpcSDK     = "grpc_sdk"
)

func init() {
	RegisterComponent(ComponentGoMod, new(proj.GoMod))
	RegisterComponent(ComponentParams, new(dto.Params))
	RegisterComponent(ComponentService, new(srv.Service))
	RegisterComponent(ComponentGrpcHandler, new(srv.GrpcHandler))
	RegisterComponent(ComponentGrpcSDK, new(sdk.GrpcSDK))
}

type Component interface {
	Setup(plugin *protogen.Plugin)
	Generate(config helper.GenerateConfig)
}

var components = make(map[string]Component, 10)

func RegisterComponent(name string, com Component) {
	components[name] = com
}

func GetComponent(name string) (Component, bool) {
	com, ok := components[name]
	return com, ok
}
