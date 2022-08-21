package dto

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

type Params struct {
	plugin *protogen.Plugin
	goPkg  string
}

func (a *Params) Setup(plugin *protogen.Plugin) {
	a.plugin = plugin
	a.goPkg = "dto"
}

func (a *Params) Generate(config helper.GenerateConfig) {
	for _, file := range a.plugin.Files {
		if !file.Generate {
			continue
		}
		fhead := helper.NewCodeHeader().Pkg(a.goPkg).
			Import(string(file.GoImportPath)).
			Import("github.com/golang/protobuf/jsonpb")
		projName := string(file.GoPackageName)

		importDomain := strings.Split(string(file.GoImportPath), "/")[0]
		for _, svc := range file.Services {
			fpath := fmt.Sprintf("%s/projects/%s/%s/%s.go",
				importDomain, projName, a.goPkg, strings.ToLower(svc.GoName))
			g := a.plugin.NewGeneratedFile(fpath, file.GoImportPath)
			g.P(fhead)
			tmpl := ParamsTpl{
				ProjName: projName,
			}
			reqMap := make(map[protoreflect.Name]struct{})
			rspMap := make(map[protoreflect.Name]struct{})
			for _, s := range file.Services {
				for _, method := range s.Methods {
					reqMap[method.Input.Desc.Name()] = struct{}{}
					rspMap[method.Output.Desc.Name()] = struct{}{}
				}
			}
			for _, pbmsg := range file.Messages {
				isReq := false
				if _, ok := reqMap[pbmsg.Desc.Name()]; ok {
					isReq = true
				}
				isRsp := false
				if _, ok := rspMap[pbmsg.Desc.Name()]; ok {
					isRsp = true
				}
				msg := Message{
					Name:  helper.ToCamelCase(string(pbmsg.Desc.Name())),
					IsReq: isReq,
					IsRsp: isRsp,
				}
				for _, field := range pbmsg.Fields {
					defaultVal := "0"
					switch field.Desc.Kind() {
					case protoreflect.StringKind:
						defaultVal = "\"\""
					case protoreflect.MessageKind:
						defaultVal = "nil"
					}
					checkLen := false
					if field.Desc.IsList() ||
						field.Desc.IsMap() ||
						field.Desc.Kind() == protoreflect.BytesKind {
						checkLen = true
					}
					msg.Fields = append(msg.Fields, Field{
						Name:       field.GoName,
						CheckLen:   checkLen,
						DefaultVal: defaultVal,
					})
				}
				tmpl.Messages = append(tmpl.Messages, msg)
			}
			g.P(tmpl.Render())
		}
	}
}
