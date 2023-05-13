package main

/* todo list
env tag:
no env tag at all ``
name defined `env:"NAME"`
name not defined `env:""`
required `env:"NAME,required"`
not empty `env:"NAME,notempty"`
unset `env:"NAME,unset"`
from file `env:"NAME,fromfile"`


tags:
env prefix `envPrefix:"DB_"`
env expand `envExpand:"true"`
env separator ( for slices ) `envSeperator:":"`
env default `envDefault:"DEFAULT_VALUE"`

options:
tag name
use field name by default
required if no def
*/

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/miniscruff/envexample/pkg"
)

var version = "dev" // injected by goreleaser

func main() {
	cfg := pkg.NewConfig()
	if cfg.ShowVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if cfg.ShowHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	writer, err := cfg.Writer()
	if err != nil {
		log.Fatal(err)
	}

	defer writer.Close()

	err = pkg.Run(writer, version, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
