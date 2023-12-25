SHELL := /bin/bash

# constant variables
PROJECT_NAME	= quic-checker
BINARY_NAM	= quic-checker
GIT_COMMIT	= $(shell git rev-parse HEAD)
BINARY_TAR_DIR	= $(BINARY_NAME)-$(GIT_COMMIT)
BINARY_TAR_FILE	= $(BINARY_TAR_DIR).tar.gz
BUILD_VERSION	= $(shell cat VERSION.txt)
BUILD_DATE	= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# golangci-lint config
golangci_lint_version=latest
vols=-v `pwd`:/app -w /app
run_lint=docker run --rm $(vols) golangci/golangci-lint:$(golangci_lint_version)

.PHONY: lint fmt test

fmt:
	@gofmt -l -w $(SRC)

lint:
	@printf "$(OK_COLOR)==> Running golang-ci-linter via Docker$(NO_COLOR)\n"
	@$(run_lint) golangci-lint run --timeout=5m --verbose
#
# Build: TODO
#
## test: run tests with coverage
test:
	@printf "$(OK_COLOR)==> Running tests$(NO_COLOR)\n"
	@go test -v -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage.txt ./...
	@go tool cover -func coverage.txt
