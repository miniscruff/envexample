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
	// NotEmpty specifies a field is required and can not be empty
	NotEmpty string `env:",notEmpty,,"`
	// FromFile loads our value from a file at path
	FromFile string `env:",file"`
	// UnsetVar will unset our environment variable after it is loaded
	UnsetVar string `env:",unset"`

	// Not part of our .env.example
	unexportedIgnored string `env:"ignored"`
}

/*


	Home         string        `env:"HOME"`
	Port         int           `env:"PORT" envDefault:"3000"`
	Password     string        `env:"PASSWORD,unset"`
	IsProduction bool          `env:"PRODUCTION"`
	Hosts        []string      `env:"HOSTS" envSeparator:":"`
	Duration     time.Duration `env:"DURATION"`
	TempFolder   string        `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`
*/
