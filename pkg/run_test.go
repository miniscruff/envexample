package pkg

import (
	"bytes"
	"testing"

	"github.com/miniscruff/envexample/then"
)

func Test_NewPkgTypes_UnableToGetGoPackage(t *testing.T) {
	_, err := NewPackageTypes("badgopkg")
	then.Err(t, errUnableToFindPackage, err)
}

func Test_Run_UnableToWrite(t *testing.T) {
	badWriter := then.NewErrWriter()
	err := Run(badWriter, "ver", &Config{})
	badWriter.Raised(t, err)
}

func Test_Run_UnableToBuildPackages(t *testing.T) {
	var writer bytes.Buffer
	err := Run(&writer, "ver", &Config{
		PackageName: "badgopkgagain",
	})
	then.Err(t, errUnableToBuildPackages, err)
}

func Test_Run_UnableToWriteStruct(t *testing.T) {
    then.RunFromDir(t, "..")

	badWriter := then.NewCountWriter(1)
	err := Run(badWriter, "ver", &Config{
		PackageName: "pkg",
	})
	badWriter.Raised(t, err)
}
