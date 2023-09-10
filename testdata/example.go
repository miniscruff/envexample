package testdata

// ExampleConfig handles our applications configuration
type ExampleConfig struct {
	// Home is our users home directory
	Home string `env:"HOME"`
	// Port is our http listeners port address
	Port int `env:"PORT" envDefault:"3000"`
	// IsProduction is whether or not we are running in production
	IsProduction bool `env:"PRODUCTION"`
}
