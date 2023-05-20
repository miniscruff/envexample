package pkg

import (
	"go/ast"
	"reflect"
	"strings"
)

type EnvTagOptions struct {
	Key string

	LoadFile bool
	Unset    bool
	NotEmpty bool
	Required bool

	IsSlice         bool
	SliceSeperator  string
	DefaultValueTag string

	TypeName string
	Prefix   string
}

func NewEnvTagOptions(name string, gen *Generator, field *ast.Field) *EnvTagOptions {
	opts := &EnvTagOptions{
		Key:            name,
		SliceSeperator: ",",
		Required:       gen.RequiredIfNoDef,
	}

	if gen.UseFieldNameByDefault {
		opts.Key = toEnvName(name)
	}

	switch fieldType := field.Type.(type) {
	case *ast.ArrayType:
		opts.IsSlice = true
		rootType := fieldType.Elt.(*ast.Ident)
		opts.TypeName = rootType.Name
	case *ast.Ident:
		opts.TypeName = fieldType.Name
	}

	structTag := field.Tag
	if structTag == nil {
		return opts
	}

	structTags := reflect.StructTag(strings.Trim(structTag.Value, "`"))

	optValues := strings.Split(structTags.Get(gen.TagName), ",")
	if optValues[0] != "" {
		opts.Key = optValues[0]
	}

	for _, tag := range optValues[1:] {
		switch tag {
		case "":
			continue
		case "file":
			opts.LoadFile = true
		case "required":
			opts.Required = true
		case "unset":
			opts.Unset = true
		case "notEmpty":
			opts.NotEmpty = true
		}
	}

	if sep, ok := structTags.Lookup("envSeperator"); ok {
		opts.SliceSeperator = sep
	}

	if def, ok := structTags.Lookup("envDefault"); ok {
		opts.DefaultValueTag = def
		opts.Required = false
	}

	if prefix, ok := structTags.Lookup("envPrefix"); ok {
		opts.Prefix = prefix
	}

	return opts
}

func (opt *EnvTagOptions) LinePrefix() string {
	if opt.Required || opt.NotEmpty {
		return ""
	}

	return "#"
}

func (opt *EnvTagOptions) DefaultValue() string {
	defValue := opt.DefaultValueTag
	if opt.NotEmpty && defValue == "" {
		defValue = "can not be empty"
	}

	// wrap values that have a space or are empty with quotes
	if strings.Contains(defValue, " ") || defValue == "" {
		return `"` + defValue + `"`
	}

	return defValue
}
