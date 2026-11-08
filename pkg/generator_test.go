package pkg

import (
	"bytes"
	"os"
	"testing"

	"github.com/miniscruff/envexample/then"
)

func generate(t *testing.T, cfg *Config) string {
	t.Helper()

	cfg.Version = "dev"
	cfg.TagName = defaultTag

	gen, err := NewGenerator(cfg)
	then.Nil(t, err)

	var writer bytes.Buffer

	then.Nil(t, gen.Run(&writer))

	return writer.String()
}

func TestGeneratorRelativePackage(t *testing.T) {
	then.RunFromDir(t, "../testdata")

	golden, err := os.ReadFile("example.golden.env")
	then.Nil(t, err)

	for _, pkg := range []string{".", "./"} {
		t.Run(pkg, func(t *testing.T) {
			output := generate(t, &Config{Package: pkg, ConfigStruct: "ExampleConfig"})
			then.Equals(t, string(golden), output)
		})
	}
}

func TestGeneratorRelativeSubPackage(t *testing.T) {
	then.RunFromDir(t, "..")

	output := generate(t, &Config{Package: "./testdata", ConfigStruct: "ExampleConfig"})
	then.Contains(t, "#PORT=3000", output)
}

func TestGeneratorStructNotFound(t *testing.T) {
	_, err := NewGenerator(&Config{
		Package:      "github.com/miniscruff/envexample/testdata",
		ConfigStruct: "MissingConfig",
		TagName:      defaultTag,
	})
	then.Err(t, errStructNotFound, err)
}

func TestGeneratorPatternMatchesManyPackages(t *testing.T) {
	_, err := NewGenerator(&Config{
		Package:      "github.com/miniscruff/envexample/testdata/...",
		ConfigStruct: "ExampleConfig",
		TagName:      defaultTag,
	})
	then.Err(t, errNotOnePackage, err)
}
