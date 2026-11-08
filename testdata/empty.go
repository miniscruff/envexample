package testdata

import "github.com/miniscruff/envexample/testdata/a"

// EmptyConfig has nested structs without any env fields of their own.
type EmptyConfig struct {
	// Name is loaded from the environment.
	Name string `env:"NAME"`
	// Group only groups other config structs.
	Group GroupConfig `envPrefix:"GROUP_"`
	// Ignored only has ignored fields.
	Ignored IgnoredFieldsConfig
	// Private only has unexported fields.
	Private PrivateFieldsConfig
}

// GroupConfig only groups other config structs.
type GroupConfig struct {
	// Data is data related to A.
	Data a.Data `envPrefix:"DATA_"`
}

// IgnoredFieldsConfig only has ignored fields.
type IgnoredFieldsConfig struct {
	// Skipped is ignored.
	Skipped string `env:"-"`
}

// PrivateFieldsConfig only has unexported fields.
type PrivateFieldsConfig struct {
	client string
}
