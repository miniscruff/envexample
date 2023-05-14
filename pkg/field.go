package pkg

import (
	"fmt"
	"go/ast"
	"io"
	"reflect"
	"strings"
)

type EnvTagOptions struct {
	Key string

	LoadFile bool
	Unset    bool
	NotEmpty bool
	Required bool
}

func NewEnvTagOptions(name string, structTag *ast.BasicLit) *EnvTagOptions {
	opts := &EnvTagOptions{
		Key: name,
	}

	if structTag == nil {
		return opts
	}

	structTags := reflect.StructTag(strings.Trim(structTag.Value, "`"))

	optValues := strings.Split(structTags.Get("env"), ",")
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

	return opts
}

func WriteField(writer io.Writer, field *ast.Field, prefix string) error {
	// just checking the first write
	if _, err := writer.Write([]byte("# ")); err != nil {
		return fmt.Errorf("%w: %w", errUnableToWrite, err)
	}

	_, _ = writer.Write([]byte(strings.Replace(strings.Trim(field.Doc.Text(), "\n"), "\n", "\n# ", -1)))
	_, _ = writer.Write([]byte("\n"))

	linePrefix := ""
	defaultValue := ""
	opts := NewEnvTagOptions(field.Names[0].Name, field.Tag)

	// part of #6
	// if def, ok := structTags.Lookup("envDefault"); ok {
	// defaultValue = def
	// }

	if opts.NotEmpty && defaultValue == "" {
		defaultValue = "can not be empty"
	}

	if !opts.Required && !opts.NotEmpty {
		linePrefix = "#"
	}

	// wrap values that have a space or are empty with quotes
	if strings.Contains(defaultValue, " ") || defaultValue == "" {
		defaultValue = `"` + defaultValue + `"`
	}

	if opts.LoadFile {
		_, _ = writer.Write([]byte(helpFile))
	}

	if opts.Unset {
		_, _ = writer.Write([]byte(helpUnset))
	}

	_, _ = writer.Write([]byte(fmt.Sprintf(
		"%s%s%s=%s",
		linePrefix,
		prefix,
		opts.Key,
		defaultValue,
	)))

	// will need this later for slice logic I believe
	/*
		switch fieldType := field.Type.(type) {
		case *ast.Ident:
			// props.TypeName = fieldType.Name
		case *ast.ArrayType:
			// rootType := fieldType.Elt.(*ast.Ident)
			// props.TypeName = rootType.Name
			// props.Slice = true
		case *ast.StarExpr:
			// rootType := fieldType.X.(*ast.Ident)
			// props.TypeName = rootType.Name
		case *ast.SelectorExpr:
			// rootType := fieldType.Sel
			// props.TypeName = rootType.Name
		default:
			return fmt.Errorf(
				"%w: %T for field '%v': %w",
				errUnknownFieldType,
				fieldType,
				field.Names[0],
				errUnableToWriteField,
			)
		}
	*/

	// if envPrefix, ok := tags.Lookup("envPrefix"); ok {
	// *queue = append(*queue, QueueEntry{
	// TypeName: props.Name,
	// Prefix:   prefix + envPrefix,
	// })
	// }

	_, _ = writer.Write([]byte("\n"))

	return nil
}
