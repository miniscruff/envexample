package pkg

import (
	"fmt"
	"io"
	"slices"
	"sort"

	"k8s.io/gengo/v2/namer"
	"k8s.io/gengo/v2/parser"
	"k8s.io/gengo/v2/types"
)

type StructQueueEntry struct {
	TypeName  types.Name
	FieldDocs []string
	Prefix    string
}

type Generator struct {
	queue    []*StructQueueEntry
	version  string
	universe types.Universe

	RequiredIfNoDef       bool
	UseFieldNameByDefault bool
	TagName               string
}

func NewGenerator(cfg *Config) (*Generator, error) {
	root := types.Name{
		Package: cfg.Package,
		Name:    cfg.ConfigStruct,
	}

	uni, err := loadUniverse(root)
	if err != nil {
		return nil, err
	}

	return &Generator{
		version:  cfg.Version,
		universe: uni,
		queue: []*StructQueueEntry{
			{
				TypeName: root,
				Prefix:   cfg.Prefix,
			},
		},
		RequiredIfNoDef:       cfg.RequiredIfNoDef,
		UseFieldNameByDefault: cfg.UseFieldNameByDefault,
		TagName:               cfg.TagName,
	}, nil
}

// loadUniverse loads the types of the root struct and every struct nested in
// it, which can come from any package or module the root package imports.
// gengo only records the doc comments of types in packages we request, so
// the packages of nested structs are found first and then requested as well.
func loadUniverse(root types.Name) (types.Universe, error) {
	requested := []string{root.Package}

	for {
		builder := parser.New()

		err := builder.LoadPackages(requested...)
		if err != nil {
			return nil, fmt.Errorf("loading packages: %w", err)
		}

		uni, err := builder.NewUniverse()
		if err != nil {
			return nil, fmt.Errorf("finding package types: %w", err)
		}

		missing := missingPackages(uni.Type(root), requested)
		if len(missing) == 0 {
			return uni, nil
		}

		requested = append(requested, missing...)
	}
}

// missingPackages returns the packages of structs nested in t that are not
// in requested.
func missingPackages(t *types.Type, requested []string) []string {
	var missing []string

	seen := map[*types.Type]bool{}
	queue := []*types.Type{t}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if seen[current] {
			continue
		}

		seen[current] = true

		for _, member := range current.Members {
			if len(member.Type.Members) == 0 {
				continue
			}

			pkg := member.Type.Name.Package
			if pkg != "" && !slices.Contains(requested, pkg) && !slices.Contains(missing, pkg) {
				missing = append(missing, pkg)
			}

			queue = append(queue, member.Type)
		}
	}

	return missing
}

func (g *Generator) WriteStruct(writer io.Writer, entry *StructQueueEntry) {
	docType := g.universe.Type(entry.TypeName)

	writeLines(writer, entry.FieldDocs, "# ", "#\n")
	writeLines(writer, docType.CommentLines, "# ", "#\n")

	membs := docType.Members
	sort.Slice(membs, func(i, j int) bool {
		return membs[i].Name < membs[j].Name
	})

	for _, field := range membs {
		g.WriteField(writer, field, entry.Prefix)
	}

	_, _ = writer.Write([]byte("\n"))
}

func (g *Generator) WriteField(writer io.Writer, field types.Member, prefix string) {
	if namer.IsPrivateGoName(field.Name) {
		return
	}

	opts := NewEnvTagOptions(field.Name, g, field.Type.Kind == types.Slice, field.Tags)

	if len(field.Type.Members) > 0 {
		g.queue = append(g.queue, &StructQueueEntry{
			TypeName:  field.Type.Name,
			FieldDocs: field.CommentLines,
			Prefix:    prefix + opts.Prefix,
		})

		// write our nested struct later
		return
	}

	// ignore fields that are defaulted to a dash
	if opts.Key == "-" {
		return
	}

	writeLines(writer, field.CommentLines, "# ", "")

	if opts.LoadFile {
		_, _ = writer.Write([]byte(helpFile))
	}

	if opts.Unset {
		_, _ = writer.Write([]byte(helpUnset))
	}

	if opts.IsSlice {
		_, _ = fmt.Fprintf(writer, helpSlice, opts.SliceSeperator)
	}

	_, _ = fmt.Fprintf(writer,
		"%s%s%s=%s",
		opts.LinePrefix(),
		prefix,
		opts.Key,
		opts.DefaultValue())

	_, _ = writer.Write([]byte("\n"))
}

func (g *Generator) Run(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, header, g.version)
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

func writeLines(w io.Writer, lines []string, prefix, ender string) {
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return
	}

	for _, line := range lines {
		if line == "" {
			continue
		}

		_, _ = w.Write([]byte(prefix + line + "\n"))
	}

	_, _ = w.Write([]byte(ender))
}
