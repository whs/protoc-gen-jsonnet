package main

import (
	"embed"
	"github.com/flosch/pongo2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"io/ioutil"
	"os"
	"regexp"
)

//go:embed assets/*
var assets embed.FS
var loader = &FsLoader{Fs: assets}

var fileRegex = regexp.MustCompile("\\.proto$")

func protoDescriptorsMap(in []*descriptorpb.FileDescriptorProto) map[string]*descriptorpb.FileDescriptorProto {
	out := make(map[string]*descriptorpb.FileDescriptorProto)

	for _, inFile := range in {
		out[inFile.GetName()] = inFile
	}

	return out
}

func generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	pongo2.SetAutoescape(false)

	set := pongo2.NewSet("assets", loader)
	template, err := set.FromFile("assets/template.jsonnet")
	if err != nil {
		return nil, err
	}

	fdSet := descriptorpb.FileDescriptorSet{
		File: req.ProtoFile,
	}
	protofiles, err := protodesc.NewFiles(&fdSet)
	if err != nil {
		return nil, err
	}

	out := make([]*pluginpb.CodeGeneratorResponse_File, len(req.FileToGenerate))
	for i, filename := range req.FileToGenerate {
		inputFile, err := protofiles.FindFileByPath(filename)
		if err != nil {
			return nil, err
		}
		outFilename := fileRegex.ReplaceAllString(filename, ".libsonnet")
		jsonnetOut, err := template.Execute(pongo2.Context{
			"request":  req,
			"filename": filename,
			"input":    inputFile,
			"range":    rangeFunc,
			"docUtils": "github.com/jsonnet-libs/docsonnet/doc-util",
			// force lazy evaluate on recursive template
			"messageFile": "message.jsonnet",
		})
		if err != nil {
			return nil, err
		}
		out[i] = &pluginpb.CodeGeneratorResponse_File{
			Name:    &outFilename,
			Content: &jsonnetOut,
		}
	}

	return &pluginpb.CodeGeneratorResponse{
		File: out,
	}, nil
}

func errorResponse(err error) {
	outBytes, err := proto.Marshal(&pluginpb.CodeGeneratorResponse{
		Error: proto.String(err.Error()),
	})
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(outBytes)
	os.Exit(0)
}

func main() {
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		errorResponse(err)
	}

	var req pluginpb.CodeGeneratorRequest
	err = proto.Unmarshal(body, &req)
	if err != nil {
		errorResponse(err)
	}

	resp, err := generate(&req)
	if err != nil {
		errorResponse(err)
	}

	outBytes, err := proto.Marshal(resp)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(outBytes)
}
