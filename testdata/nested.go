package testdata

import "github.com/miniscruff/envexample/testdata/a"

// Nested checks whether we can import packages outside our local.
type Config struct {
	// A is data related to A.
	A a.Data `envPrefix:"A_"`
	// B is data related to B.
	B b.Data `envPrefix:"B_"`
}
