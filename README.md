# envexample
Generate a .env.example from an env struct

## Installation
TODO

## Usage

Using go generate, specify the `envexample` arguments
```
//go:generate envexample -struct Config
```

## Tasks
Below is a list of common development tasks, these can easily be run using [xc](https://xcfile.dev/).
For example `xc test` will run the test suite.

### golden
Run unit test suite with code coverage enabled
```
go generate ./testdata/
```

### lint
```
goimports -w -local github.com/miniscruff/envexample .
golangci-lint run ./...
```

## License
Distributed under the [MIT License](LICENSE).
