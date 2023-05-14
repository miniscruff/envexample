package pkg

import (
	"go/ast"
	godoc "go/doc"
	"testing"

	"github.com/miniscruff/envexample/then"
)

func Test_Struct_BadField(t *testing.T) {
	badWriter := then.NewCountWriter(2)

	// creating this struct is kinda insane...
	structType := godoc.Type{
		Doc: "some docs",
		Decl: &ast.GenDecl{
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Type: &ast.StructType{
						Fields: &ast.FieldList{
							Opening: 0,
							Closing: 100,
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{
											Name: "FieldName",
										},
									},
									Doc: &ast.CommentGroup{
										List: []*ast.Comment{
											{
												Slash: 0,
												Text:  "field doc text",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	err := WriteStruct(badWriter, &structType, "")
	then.Err(t, errUnableToWriteField, err)
}
