package testdata

import "github.com/miniscruff/envexample/testdata/a"

// IgnoredConfig has fields that are not loaded from the environment.
type IgnoredConfig struct {
	// Name is loaded from the environment.
	Name string `env:"NAME"`
	// Internal is ignored cause it has a dash as the env key.
	Internal a.Data `env:"-"`
	// Prefixed is ignored even with a prefix.
	Prefixed a.Data `env:"-" envPrefix:"PREFIXED_"`
}
