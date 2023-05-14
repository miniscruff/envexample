package pkg

import (
	"fmt"
	"go/ast"
	"io"
	"reflect"
	"strings"
)

func WriteField(writer io.Writer, field *ast.Field, prefix string) error {
	// just checking the first write
	if _, err := writer.Write([]byte("# ")); err != nil {
		return fmt.Errorf("%w: %w", errUnableToWrite, err)
	}

	_, _ = writer.Write([]byte(strings.Replace(strings.Trim(field.Doc.Text(), "\n"), "\n", "\n# ", -1)))
	_, _ = writer.Write([]byte("\n"))

	var (
		structTags reflect.StructTag
		tags       []string
		linePrefix string
		loadFile   bool
		unset      bool
		notEmpty   bool
		required   bool
	)

	defaultValue := ""
	ownKey := field.Names[0].Name

	if field.Tag != nil {
		structTags = reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

		opts := strings.Split(structTags.Get("env"), ",")
		if opts[0] != "" {
			ownKey = opts[0]
		}

		tags = opts[1:]
	}

	for _, tag := range tags {
		switch tag {
		case "":
			continue
		case "file":
			loadFile = true
		case "required":
			required = true
		case "unset":
			unset = true
		case "notEmpty":
			notEmpty = true
		}
	}

	// part of #6
	// if def, ok := structTags.Lookup("envDefault"); ok {
	// defaultValue = def
	// }

	if notEmpty && defaultValue == "" {
		defaultValue = "can not be empty"
	}

	if !required && !notEmpty {
		linePrefix = "#"
	}

	// wrap values that have a space or are empty with quotes
	if strings.Contains(defaultValue, " ") || defaultValue == "" {
		defaultValue = `"` + defaultValue + `"`
	}

	if loadFile {
		_, _ = writer.Write([]byte(helpFile))
	}

	if unset {
		_, _ = writer.Write([]byte(helpUnset))
	}

	_, _ = writer.Write([]byte(fmt.Sprintf(
		"%s%s%s=%s",
		linePrefix,
		prefix,
		ownKey,
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
