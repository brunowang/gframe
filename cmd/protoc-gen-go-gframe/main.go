package main

import (
	"flag"
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
	"strings"
)

func main() {
	var flags flag.FlagSet
	test := flags.Bool("test", false, "log and exit")
	components := flags.String("components", "", "component name list, split by +")
	project := flags.String("project", "unnamed", "specify project name, will use as go mod name")
	pbGoDir := flags.String("pbGoDir", "", "specify proto go code dir, will use for go mod require")
	options := protogen.Options{
		ParamFunc: flags.Set,
	}

	baseComponents := []string{
		gen.ComponentGoMod,
		gen.ComponentGoMain,
		gen.ComponentParams,
		gen.ComponentService,
		gen.ComponentServer,
		gen.ComponentGrpcServer,
		gen.ComponentGrpcHandler,
	}

	options.Run(func(plugin *protogen.Plugin) error {
		if *test {
			return fmt.Errorf("protoc-gen-go-gframe")
		}
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		components := append(baseComponents, strings.Split(*components, "+")...)
		gen.NewProjectGenerator(plugin).Generate(
			gen.WithComponents(components),
			gen.WithProject(*project),
			gen.WithPbGoDir(*pbGoDir),
		)
		return nil
	})
}
