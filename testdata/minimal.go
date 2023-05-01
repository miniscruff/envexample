package testdata

//go:generate go run ../main.go -struct MinimalConfig -export minimal.golden

type MinimalConfig struct {
	// Listen address for our http server
	Address string `env:"ADDRESS" envDefault:"localhost"`
	// Listen port for our http server
	Port int `env:"PORT" envDefault:"8080"`
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
