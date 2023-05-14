package pkg

import (
	"bytes"
	"os"
	"path/filepath"
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
	cfg.PackageName = "testdata"

	var writer bytes.Buffer
    err := Run(&writer, "dev", cfg)
	then.Nil(t, err)

	t.Log(writer.String())

	NormalizedFileContents(t, writer.Bytes(), "testdata", golden.GoldenFile)
}

func TestGoldens(t *testing.T) {
    then.RunFromDir(t, "..")

	for _, golden := range []GoldenTest{
		{
			GoldenFile: "starting.golden",
			Config: &Config{
				ConfigStruct: "StartingConfig",
			},
		},
	} {
		t.Run(golden.GoldenFile, golden.Run)
	}
}

// FileContents will check the contents of a file.
func NormalizedFileContents(t *testing.T, contents []byte, paths ...string) {
	t.Helper()

	fullPath := filepath.Join(paths...)

	bs, err := os.ReadFile(fullPath)
	if err != nil {
		t.Errorf("reading file: '%v'", fullPath)
	}

    if len(bs) == 0 {
        t.Error("file is empty", fullPath)
    }

	normalized := bytes.ReplaceAll(bs, []byte("\r\n"), []byte("\n"))

	then.SliceEquals(t, contents, normalized)
}
