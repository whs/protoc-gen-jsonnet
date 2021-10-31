package main

import (
	"encoding/json"
	"github.com/flosch/pongo2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

var reservedKeywords = []string{
	"assert", "else", "error", "false", "for", "function", "if",
	"import", "importstr", "in", "local", "null", "tailstrict",
	"then", "self", "super", "true",
}

func init() {
	pongo2.RegisterFilter("json", jsonEncode)
	pongo2.RegisterFilter("identifier", identifier)
	pongo2.RegisterFilter("indent", indent)
	pongo2.RegisterFilter("strip", strip)
	pongo2.RegisterFilter("ucfirst", ucFirst)
}

func jsonEncode(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var message proto.Message
	switch v := in.Interface().(type) {
	case protoreflect.MessageDescriptor:
		message = protodesc.ToDescriptorProto(v)
	case protoreflect.EnumDescriptor:
		message = protodesc.ToEnumDescriptorProto(v)
	case protoreflect.EnumValueDescriptor:
		message = protodesc.ToEnumValueDescriptorProto(v)
	case protoreflect.FieldDescriptor:
		message = protodesc.ToFieldDescriptorProto(v)
	case protoreflect.FileDescriptor:
		message = protodesc.ToFileDescriptorProto(v)
	case protoreflect.MethodDescriptor:
		message = protodesc.ToMethodDescriptorProto(v)
	case protoreflect.OneofDescriptor:
		message = protodesc.ToOneofDescriptorProto(v)
	case protoreflect.ServiceDescriptor:
		message = protodesc.ToServiceDescriptorProto(v)
	case protoreflect.Value:
		return jsonEncode(pongo2.AsValue(v.Interface()), nil)
	default:
		b, err := json.Marshal(in.Interface())
		if err != nil {
			return nil, &pongo2.Error{
				Sender:    "filter:json",
				OrigError: err,
			}
		}
		return pongo2.AsValue(string(b)), nil
	}
	return pongo2.AsValue(protojson.Format(message)), nil
}

func identifier(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := in.String()
	for _, r := range reservedKeywords {
		if val == r {
			return pongo2.AsValue(val + "_"), nil
		}
	}
	return in, nil
}

func indent(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := in.String()
	indentSize := param.Integer()
	out := strings.ReplaceAll(val, "\n", "\n"+strings.Repeat("\t", indentSize))
	return pongo2.AsValue(out), nil
}

func rangeFunc(in int) []int {
	out := make([]int, in)
	for i := 0; i < in; i++ {
		out[i] = i
	}
	return out
}

func strip(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(strings.TrimSpace(in.String())), nil
}

func ucFirst(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	s = strings.ToUpper(s[0:1]) + s[1:]
	return pongo2.AsValue(s), nil
}
