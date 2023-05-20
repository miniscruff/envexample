package pkg

import (
	"fmt"
	"go/ast"
	godoc "go/doc"
	"go/parser"
	"go/token"
	"io"
	"sort"
	"strings"
)

type PackageTypes map[string]*godoc.Type

func NewPackageTypes(directory string) (PackageTypes, error) {
	fset := token.NewFileSet()
	packages, _ := parser.ParseDir(fset, "./"+directory, nil, parser.ParseComments)

	projectPackage, ok := packages[directory]
	if !ok {
		return nil, fmt.Errorf("%w: '%v'", errUnableToFindPackage, directory)
	}

	p := godoc.New(projectPackage, "./", 0)

	projectPackages := make(PackageTypes)
	for _, t := range p.Types {
		projectPackages[t.Name] = t
	}

	return projectPackages, nil
}

type StructQueueEntry struct {
	TypeName  string
	FieldDocs string
	Prefix    string
}

type Generator struct {
	packageTypes PackageTypes
	queue        []*StructQueueEntry
	version      string
}

func NewGenerator(cfg *Config) (*Generator, error) {
	pkgTypes, err := NewPackageTypes(cfg.Directory)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errUnableToBuildPackages, err)
	}

	return &Generator{
		version:      cfg.Version,
		packageTypes: pkgTypes,
		queue: []*StructQueueEntry{
			{
				TypeName: cfg.ConfigStruct,
			},
		},
	}, nil
}

func (g *Generator) WriteStruct(writer io.Writer, entry *StructQueueEntry) {
	docType := g.packageTypes[entry.TypeName]

	if entry.FieldDocs != "" {
		_, _ = writer.Write([]byte("# "))
		_, _ = writer.Write([]byte(strings.Replace(strings.Trim(entry.FieldDocs, "\n"), "\n", "\n# ", -1)))
		_, _ = writer.Write([]byte("\n#\n"))
	}

	if docType.Doc != "" {
		_, _ = writer.Write([]byte("# "))
		_, _ = writer.Write([]byte(strings.Replace(strings.Trim(docType.Doc, "\n"), "\n", "\n# ", -1)))
		_, _ = writer.Write([]byte("\n#\n"))
	}

	fields := make([]*ast.Field, 0)

	specs := docType.Decl.Specs
	for _, spec := range specs {
		// can't remember when either of these will fail
		typeSpec, _ := spec.(*ast.TypeSpec)
		structType, _ := typeSpec.Type.(*ast.StructType)

		fields = append(fields, structType.Fields.List...)
	}

	sort.SliceStable(fields, func(i, j int) bool {
		return fields[i].Names[0].Name < fields[j].Names[0].Name
	})

	for _, field := range fields {
		g.WriteField(writer, field, entry.Prefix)
	}

	_, _ = writer.Write([]byte("\n"))
}

func (g *Generator) WriteField(writer io.Writer, field *ast.Field, prefix string) {
	opts := NewEnvTagOptions(field.Names[0].Name, field)
	if _, found := g.packageTypes[opts.TypeName]; found {
		g.queue = append(g.queue, &StructQueueEntry{
			TypeName:  opts.TypeName,
			FieldDocs: field.Doc.Text(),
			Prefix:    opts.Prefix,
		})

		// write our nested struct later
		return
	}

	if field.Doc.Text() != "" {
		_, _ = writer.Write([]byte("# "))
		_, _ = writer.Write([]byte(strings.Replace(strings.Trim(field.Doc.Text(), "\n"), "\n", "\n# ", -1)))
		_, _ = writer.Write([]byte("\n"))
	}

	if opts.LoadFile {
		_, _ = writer.Write([]byte(helpFile))
	}

	if opts.Unset {
		_, _ = writer.Write([]byte(helpUnset))
	}

	if opts.IsSlice {
		_, _ = writer.Write([]byte(fmt.Sprintf(helpSlice, opts.SliceSeperator)))
	}

	_, _ = writer.Write([]byte(fmt.Sprintf(
		"%s%s%s=%s",
		opts.LinePrefix(),
		prefix,
		opts.Key,
		opts.DefaultValue(),
	)))

	_, _ = writer.Write([]byte("\n"))
}

func (g *Generator) Run(writer io.Writer) error {
	_, err := writer.Write([]byte(fmt.Sprintf(header, g.version)))
	if err != nil {
		return fmt.Errorf("%w: %w", errUnableToWrite, err)
	}

	for {
		entry := g.queue[0]
		g.queue = g.queue[1:]

		g.WriteStruct(writer, entry)

		if len(g.queue) == 0 {
			break
		}
	}

	return nil
}
