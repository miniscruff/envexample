package pkg

import (
	"fmt"
	godoc "go/doc"
	"go/parser"
	"go/token"
	"io"
)

type PackageTypes map[string]*godoc.Type

func NewPackageTypes(packageName string) (PackageTypes, error) {
	fset := token.NewFileSet()
	packages, _ := parser.ParseDir(fset, "./", nil, parser.ParseComments)

	projectPackage, ok := packages[packageName]
	if !ok {
		return nil, fmt.Errorf("%w: '%v'", errUnableToFindPackage, packageName)
	}

	p := godoc.New(projectPackage, "./", 0)

	projectPackages := make(PackageTypes)
	for _, t := range p.Types {
		projectPackages[t.Name] = t
	}

	return projectPackages, nil
}

func Run(writer io.Writer, version string, cfg *Config) error {
	_, err := writer.Write([]byte(fmt.Sprintf(header, version)))
	if err != nil {
		return fmt.Errorf("%w: %w", errUnableToWrite, err)
	}

	pkgTypes, err := NewPackageTypes(cfg.PackageName)
	if err != nil {
		return fmt.Errorf("%w: %w", errUnableToBuildPackages, err)
	}

	err = WriteStruct(writer, pkgTypes[cfg.ConfigStruct], "")
	if err != nil {
		return fmt.Errorf("%w: %w", errUnableToWriteStruct, err)
	}

	return nil
}
