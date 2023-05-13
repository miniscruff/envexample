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

	tags := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

	if envName, ok := tags.Lookup("env"); ok {
		_, _ = writer.Write([]byte(prefix + envName))
	} // will need to find what the name is without the env struct tag

	_, _ = writer.Write([]byte("="))
	if def, ok := tags.Lookup("envDefault"); ok {
		_, _ = writer.Write([]byte(def))
	} // need to write something else as the value

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

	// if envPrefix, ok := tags.Lookup("envPrefix"); ok {
	// *queue = append(*queue, QueueEntry{
	// TypeName: props.Name,
	// Prefix:   prefix + envPrefix,
	// })
	// }

	_, _ = writer.Write([]byte("\n"))

	return nil
}
