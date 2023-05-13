package pkg

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Config handles all the envexample configurations.
type Config struct {
	ExportFile   string
	ConfigStruct string
	PackageName  string
	DryRun       bool
	ShowVersion  bool
	ShowHelp     bool
}

// NewConfig will load a config from CLI flags.
func NewConfig() *Config {
	cfg := Config{}

	flag.StringVar(&cfg.ExportFile, "export", ".env.example", "`filepath` to export generated example to")
	flag.StringVar(&cfg.ConfigStruct, "struct", "", "`struct` to build example from")
	flag.StringVar(&cfg.PackageName, "package", "main", "`package` our config struct is located in")
	flag.BoolVar(&cfg.DryRun, "dry", false, "output to stdout instead of writing to file")
	flag.BoolVar(&cfg.ShowVersion, "v", false, "show version")
	flag.BoolVar(&cfg.ShowHelp, "h", false, "show help")
	flag.Parse()

	return &cfg
}

// Validate will check if the config is an actual valid set of configurations.
func (c *Config) Validate() error {
	if c.ConfigStruct == "" {
		flag.PrintDefaults()
		return errInvalidConfigNoStruct
	}

	if !c.DryRun && c.ExportFile == "" {
		flag.PrintDefaults()
		return errInvalidConfigNoExport
	}

	return nil
}

// Writer will create the proper io.Writer based on the configuration.
// Be sure to close the writer when done.
func (c *Config) Writer() (io.WriteCloser, error) {
	if c.DryRun {
		return os.Stdout, nil
	}

	file, err := os.Create(c.ExportFile)
	if err != nil {
		return nil, fmt.Errorf("%w: '%s': %w", errUnableToCreateWriter, c.ExportFile, err)
	}

	return file, nil
}
