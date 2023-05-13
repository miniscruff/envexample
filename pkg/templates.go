package pkg

import "errors"

const (
	header = "# Generated by envexample %s.\n\n"
)

var (
	errInvalidConfigNoStruct = errors.New("config struct not defined")
	errInvalidConfigNoExport = errors.New("export file empty")
	errUnknownFieldType      = errors.New("unknown field type")
	errUnableToCreateWriter  = errors.New("unable to create writer")
	errUnableToWriteField    = errors.New("unable to write field")
	errUnableToWriteStruct   = errors.New("unable to write struct")
	errUnableToWrite         = errors.New("unable to write")
	errUnableToFindPackage   = errors.New("unable to find go package")
	errUnableToBuildPackages = errors.New("unable to build go packages")
)
