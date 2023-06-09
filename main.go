package main

import (
	"fmt"
	"log"
	"os"

	"github.com/miniscruff/envexample/pkg"
)

var version = "dev" // injected by goreleaser

func main() {
	cfg, err := pkg.NewConfig(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ShowVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if cfg.ShowHelp {
		cfg.PrintDefaults()
		os.Exit(0)
	}

	if err = cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	writer, err := cfg.Writer()
	if err != nil {
		log.Fatal(err)
	}

	defer writer.Close()

	cfg.Version = version

	gen, err := pkg.NewGenerator(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = gen.Run(writer)
	if err != nil {
		log.Fatal(err)
	}
}
