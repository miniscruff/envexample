package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/miniscruff/envexample/pkg"
)

var version = "dev" // injected by goreleaser

// buildVersion returns the version injected by goreleaser, or for builds
// using go install or go tool, the module version recorded by Go.
func buildVersion() string {
	if version != "dev" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return version
	}

	// match the goreleaser version, which has no v prefix
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
