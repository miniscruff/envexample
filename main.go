package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/miniscruff/envexample/pkg"
)

// buildVersion returns the module version Go records in the binary when
// installed with go install or go tool, or dev for local builds.
func buildVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "dev"
	}

	return strings.TrimPrefix(info.Main.Version, "v")
}

func main() {
	cfg, err := pkg.NewConfig(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ShowVersion {
		fmt.Println(buildVersion())
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

	cfg.Version = buildVersion()

	gen, err := pkg.NewGenerator(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = gen.Run(writer)
	if err != nil {
		log.Fatal(err)
	}
}
