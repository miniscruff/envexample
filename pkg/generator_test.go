package pkg

import (
	"testing"

	"github.com/miniscruff/envexample/then"
)

func Test_NewPkgTypes_UnableToGetGoPackage(t *testing.T) {
	_, err := NewPackageTypes("badgopkg")
	then.Err(t, errUnableToFindPackage, err)
}

func Test_NewGen_UnableToGetGoPackage(t *testing.T) {
	_, err := NewGenerator(&Config{
		Directory: "badpkgdir",
	})
	then.Err(t, errUnableToFindPackage, err)
}

func Test_Run_UnableToWrite(t *testing.T) {
	gen := &Generator{}
	badWriter := then.NewErrWriter()
	err := gen.Run(badWriter)
	badWriter.Raised(t, err)
}
