# envexample
Generate a '.env.example' from an [env configured](https://github.com/caarlos0/env) struct

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

### test
Run unit test suite with code coverage enabled.
```
go test ./... -coverprofile=c.out
```

### coverage
Run unit tests and preview the html coverage results.
requires: test
```
go tool cover -html=c.out
```

### lint
```
goimports -w -local github.com/miniscruff/envexample .
golangci-lint run ./...
```

## License
Distributed under the [MIT License](LICENSE).
