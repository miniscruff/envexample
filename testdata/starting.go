package testdata

// StartingConfig handles a configuration for a project with default env options.
type StartingConfig struct {
	// NoEnvStructTag at all
	NoEnvStructTag string
	// NoEnvName where env struct tag is empty
	NoEnvName string `env:""`
	// WithEnvName sets our env var key
	WithEnvName string `env:"WITH_ENV_NAME"`

	// RequiredNoDefault specifies a required value but no default
	RequiredNoDefault string `env:",required"`
	// RequiredWithDef preloads the default value for our field
	RequiredWithDef string `env:"VALUE_WITH_DEF,required" envDefault:"secret_value"`
	// OptionalWithDef preloads the default value for our field as a comment
	OptionalWithDef string `env:"VALUE_WITH_DEF" envDefault:"secret_value"`

	// NotEmpty specifies a field is required and can not be empty
	NotEmpty string `env:",notEmpty,,"`
	// FromFile loads our value from a file at path
	FromFile string `env:",file"`
	// UnsetVar will unset our environment variable after it is loaded
	UnsetVar string `env:",unset"`

	// SliceWithDefaultSep loads a slice of values separated by a comma
	SliceWithDefaultSep []string `env:"SLICE_WITH_DEFAULT_SEP,required"`
	// SliceCustomSep loads a slice of values separated by a slash
	SliceCustomSep []string `env:"SLICE_CUSTOM_SEP,required" envSeperator:"/"`

	// Server loads a separate struct with a prefix
	Server ServerConfig `envPrefix:"SERVER_"`

	// Admin loads a separate struct without a custom prefix
	Admin AdminConfig

	// Not part of our .env.example
	unexportedIgnored string `env:"ignored"`
}

// ServerConfig handles configurations for an HTTP server
type ServerConfig struct {
	Host   string             `env:"HOST" envDefault:"localhost"`
	Port   int                `env:"PORT" envDefault:"8080"`
	Nested NestedServerConfig `envPrefix:"NESTED_"`
}

// NestedServerConfig is a nested config type under server config.
type NestedServerConfig struct {
	Debug bool `env:"DEBUG" envDefault:"true"`
}

type AdminConfig struct {
	AdminUsername string `env:"ADMIN_USER" envDefault:"ADMIN_USER"`
	Adminpassword string `env:"ADMIN_PASSWORD" envDefault:"ADMIN_PASS"`
}
