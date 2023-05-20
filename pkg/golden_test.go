package pkg

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/miniscruff/envexample/then"
)

type GoldenTest struct {
	Config     *Config
	GoldenFile string
}

func (golden GoldenTest) Run(t *testing.T) {
	// override config values used by all golden tests
	cfg := golden.Config
	cfg.Directory = "testdata"
	cfg.Version = "dev"

	var writer bytes.Buffer

	gen, err := NewGenerator(cfg)
	then.Nil(t, err)

	err = gen.Run(&writer)
	then.Nil(t, err)

	fullPath := filepath.Join("testdata", golden.GoldenFile+".golden.env")
	overrideKey := "GOLDEN_" + strings.ToUpper(golden.GoldenFile)

	shouldOverride := os.Getenv(overrideKey)
	if shouldOverride == "true" {
		var goldenFile io.WriteCloser

		goldenFile, err = os.Create(fullPath)
		if err != nil {
			t.Errorf("opening file: '%v'", fullPath)
		}

		_, err = io.Copy(goldenFile, &writer)
		if err != nil {
			t.Errorf("writing golden file: '%v'", fullPath)
		}

		defer goldenFile.Close()

		return
	}

	t.Logf("Config: %+v", cfg)
	t.Logf(`Run "%s=true xc test" if the output below matches the updated value`, overrideKey)
	t.Log(writer.String())

	bs, err := os.ReadFile(fullPath)
	if err != nil {
		t.Errorf("reading file: '%v'", fullPath)
	}

	normalized := bytes.ReplaceAll(bs, []byte("\r\n"), []byte("\n"))
	generatedBytes := writer.Bytes()

	then.SliceEquals(t, generatedBytes, normalized)
}

func TestGoldens(t *testing.T) {
	then.RunFromDir(t, "..")

	for _, golden := range []GoldenTest{
		{
			GoldenFile: "starting",
			Config: &Config{
				ConfigStruct: "StartingConfig",
				TagName:      "env",
				Prefix:       "",
			},
		},
		{
			GoldenFile: "configs_enabled",
			Config: &Config{
				ConfigStruct:          "EnvConfigsEnabled",
				RequiredIfNoDef:       true,
				UseFieldNameByDefault: true,
				TagName:               "cenv",
				Prefix:                "PROJ_",
			},
		},
	} {
		t.Run(golden.GoldenFile, golden.Run)
	}
}
