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
	Directory    string
	Version      string
	DryRun       bool
	ShowVersion  bool
	ShowHelp     bool
}

// NewConfig will load a config from CLI flags.
func NewConfig(args []string) (*Config, error) {
	cfg := Config{}
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	flagSet.StringVar(&cfg.ExportFile, "export", ".env.example", "`filepath` to export generated example to")
	flagSet.StringVar(&cfg.ConfigStruct, "type", "", "`struct` to build example from")
	flagSet.StringVar(&cfg.Directory, "dir", ".", "`directory` our config struct is located in")
	flagSet.BoolVar(&cfg.DryRun, "dry", false, "output to stdout instead of writing to file")
	flagSet.BoolVar(&cfg.ShowVersion, "v", false, "show version")
	flagSet.BoolVar(&cfg.ShowHelp, "h", false, "show help")

	if err := flagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("%w: %w", errUnableToParseFlags, err)
	}

	return &cfg, nil
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
