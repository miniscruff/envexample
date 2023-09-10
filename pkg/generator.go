package pkg

import (
	"fmt"
	"io"
	"sort"

	"k8s.io/gengo/namer"
	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

type StructQueueEntry struct {
	TypeName  types.Name
	FieldDocs []string
	Prefix    string
}

type Generator struct {
	queue    []*StructQueueEntry
	version  string
	builder  *parser.Builder
	universe types.Universe

	RequiredIfNoDef       bool
	UseFieldNameByDefault bool
	TagName               string
}

func NewGenerator(cfg *Config) (*Generator, error) {
	builder := parser.New()

	err := builder.AddDirRecursive(".")
	if err != nil {
		return nil, fmt.Errorf("finding package directory: %w", err)
	}

	uni, err := builder.FindTypes()
	if err != nil {
		return nil, fmt.Errorf("finding package types: %w", err)
	}

	return &Generator{
		version:  cfg.Version,
		builder:  builder,
		universe: uni,
		queue: []*StructQueueEntry{
			{
				TypeName: types.Name{
					Package: cfg.Directory,
					Name:    cfg.ConfigStruct,
				},
				Prefix: cfg.Prefix,
			},
		},
		RequiredIfNoDef:       cfg.RequiredIfNoDef,
		UseFieldNameByDefault: cfg.UseFieldNameByDefault,
		TagName:               cfg.TagName,
	}, nil
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

	writeLines(writer, field.CommentLines, "# ", "")

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
