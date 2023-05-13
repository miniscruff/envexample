package pkg

import (
	"fmt"
	"go/ast"
	godoc "go/doc"
	"io"
	"sort"
	"strings"
)

func WriteStruct(writer io.Writer, docType *godoc.Type, prefix string) error {
	// just checking the first write
	if _, err := writer.Write([]byte("# ")); err != nil {
		return fmt.Errorf("%w: %w", errUnableToWrite, err)
	}

	_, _ = writer.Write([]byte(strings.Replace(strings.Trim(docType.Doc, "\n"), "\n", "\n# ", -1)))
	_, _ = writer.Write([]byte("\n#\n"))

	fields := make([]*ast.Field, 0)

	specs := docType.Decl.Specs
	for _, spec := range specs {
		typeSpec, _ := spec.(*ast.TypeSpec)

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		fields = append(fields, structType.Fields.List...)
	}

	sort.SliceStable(fields, func(i, j int) bool {
		return fields[i].Names[0].Name < fields[j].Names[0].Name
	})

	for _, field := range fields {
		if err := WriteField(writer, field, prefix); err != nil {
			return fmt.Errorf("%w '%v': %w", errUnableToWriteField, field.Names[0].Name, err)
		}
	}

	return nil
}
