include .bingo/Variables.mk

.DEFAULT_GOAL := help

PROJECT := $(shell go list -m)
VERSION_NS := $(PROJECT)/command

GIT_TAGCMD := $(shell git --no-pager tag --points-at HEAD)
GIT_VERSION := $(if $(GIT_TAGCMD),$(GIT_TAGCMD),$(shell git rev-parse --short HEAD))
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_PRERELEASE := $(if $(shell git status --porcelain),dev)

GO_LDFLAGS := -X $(VERSION_NS).Version=$(GIT_VERSION) \
-X $(VERSION_NS).GitCommit=$(GIT_COMMIT) \
-X $(VERSION_NS).VersionPreRelease=$(GIT_PRERELEASE)
GO_BUILD := go build -trimpath -ldflags "$(GO_LDFLAGS)"
GO_PKG := ./cmd

PLATFORMS := darwin/amd64 darwin/arm64 \
linux/386 linux/amd64 linux/arm linux/arm64 linux/riscv64 \
windows/386 windows/amd64 windows/arm

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: setup
## setup: initialize the project environment
setup:
	@.tools/bin/asdf setup
	@.tools/bin/pre-commit setup

.PHONY: dev
## dev: quickstart
dev: setup build

.PHONY: clean
## clean: clean the working directory
clean:
	@.tools/bin/go clean

.PHONY: test
## test: run the test suite
test:
	@.tools/bin/go test

.PHONY: test-coverage
## test-report: run the test suite and generate a coverage report
test-coverage: $(GOVERALLS)
	@.tools/bin/goveralls $(GOVERALLS)

.PHONY: lint
## lint: run linters
lint: $(GOLANGCI_LINT)
	@.tools/bin/golangci $(GOLANGCI_LINT)

.PHONY: pre-commit
## pre-commit: run all pre-commit hooks
pre-commit:
	@.tools/bin/pre-commit run

.PHONY: build
## build: build the application for the local target only
build: clean build/$(shell go env GOOS)/$(shell go env GOARCH)

## build/%: build for a specific os/arch, format like `go tool dist list`
build/%: GO_OUT ?= $@
build/%: export GOOS = $(firstword $(subst /, ,$*))
build/%: export GOARCH = $(lastword $(subst /, ,$*))
build/%: export CGO_ENABLED = 0
build/%:
	$(info *** building release for target $(GOOS)/$(GOARCH))
	@$(GO_BUILD) -o $(GO_OUT) $(GO_PKG)

.PHONY: release
## release: build binaries for all targets
release: clean $(foreach p, $(PLATFORMS), build/$(p))
