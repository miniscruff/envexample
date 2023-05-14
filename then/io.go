package then

import (
	"os"
	"testing"
)

// RunFromDir will change our directory temporarily for a test
func RunFromDir(t *testing.T, relativeDirChange string) {
	oldDir, err := os.Getwd()
	Nil(t, err)
	t.Cleanup(func() {
		Nil(t, os.Chdir(oldDir))
	})

	Nil(t, os.Chdir(relativeDirChange))
}
