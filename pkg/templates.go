package pkg

import (
	"errors"
	"unicode"
)

const (
	header    = "# Generated by envexample %s.\n\n"
	helpFile  = "# Info: Value is loaded from a file at the path defined.\n"
	helpUnset = "# Info: After loading environment variable is unset.\n"
	helpSlice = "# Info: Values are separated by '%s'.\n"
)

var (
	errInvalidConfigNoStruct = errors.New("config struct not defined")
	errInvalidConfigNoExport = errors.New("export file empty")
	errUnableToParseFlags    = errors.New("unable to parse cli flags")
	errUnableToCreateWriter  = errors.New("unable to create writer")
	errUnableToWrite         = errors.New("unable to write")
)

const underscore rune = '_'

func toEnvName(input string) string {
	var output []rune
	for i, c := range input {
		if i > 0 && output[i-1] != underscore && c != underscore && unicode.ToUpper(c) == c {
			output = append(output, underscore)
		}

		output = append(output, unicode.ToUpper(c))
	}

	return string(output)
}
