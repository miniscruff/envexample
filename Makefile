.ONESHELL:
.DEFAULT: help

.PHONY: help
help:
	@grep -E '^[a-z-]+:.*#' Makefile | \
		sort | \
		while read -r l; do printf "\033[1;32m$$(echo $$l | \
		cut -d':' -f1)\033[00m:$$(echo $$l | cut -d'#' -f2-)\n"; \
	done

.PHONY: test
test: # Run unit test suite
	go test -race -coverprofile=c.out ./...
	go tool cover -html=c.out -o=coverage.html

.PHONY: format
format: # Run linter and formatters
	golangci-lint run --fix ./...

.PHONY: release
release: # Generate release PR
	go run main.go -h &> DOCS.md
	changie batch auto
	changie merge
	git checkout --branch release-$$(changie latest)
	gh pr create
