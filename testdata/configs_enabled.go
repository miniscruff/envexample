package testdata

// EnvConfigsEnabled is generated with env configs turned on.
// cenv is our custom tag name.
type EnvConfigsEnabled struct {
	// NoEnvStructTag at all
	NoEnvStructTag string
	// NoEnvName where env struct tag is empty
	NoEnvName string `cenv:""`
	// WithEnvName sets our env var key
	WithEnvName string `cenv:"WITH_CUSTOM_NAME"`

	// RequiredWithDef is required because of the env tag
	RequiredWithDef string `cenv:",required" envDefault:"secret_value"`
	// RequiredNoDef is required because RequiredIfNoDef is on
	RequiredNoDef string
}
