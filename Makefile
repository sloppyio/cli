.PHONY: default cross test test-local fmt  build shell bundle deploy clean
# sloppy-cli makefile
VERSION=$(shell cat ./version.txt)
GIT_COMMIT=$(shell git rev-parse HEAD)

# Default docker image
DOCKER_IMAGE := sloppy/go-cross:latest

# env vars passed through directly to Docker's build scripts
DOCKER_ENVS := \
	-e GIT_COMMIT=$(GIT_COMMIT)


DOCKER_MOUNT := -v "$(CURDIR):/go/src/github.com/sloppyio/cli" \
								-w "/go/src/github.com/sloppyio/cli"

DOCKER_FLAGS := docker run --rm -i -e SLOPPY_APITOKEN=$(SLOPPY_APITOKEN) $(DOCKER_ENVS) $(DOCKER_MOUNT)

# if this session isn't interactive, then we don't want to allocate a
# TTY, which would fail, but if it is interactive, we do want to attach
# so that the user can send e.g. ^C through.
INTERACTIVE := $(shell [ -t 0 ] && echo 1 || echo 0)
ifeq ($(INTERACTIVE), 1)
	DOCKER_FLAGS += -t
endif

DOCKER_RUN_DOCKER := $(DOCKER_FLAGS) "$(DOCKER_IMAGE)"

default: cross

local: scripts/make.sh build

cross: bundle
	$(DOCKER_RUN_DOCKER) scripts/make.sh cross release

beta: bundle
	$(DOCKER_RUN_DOCKER) scripts/make.sh cross beta

test: bundle
	$(DOCKER_RUN_DOCKER) scripts/test.sh

update-vendor:
	$(DOCKER_RUN_DOCKER) dep ensure -update github.com/sloppyio/sloppose

test-local: bundle
	scripts/make.sh test

fmt:
	$(DOCKER_RUN_DOCKER) gofmt -w .

deploy: bundle
	scripts/make.sh deploy release

build:
	docker build -t sloppy/go-cross:latest .

shell:
	$(DOCKER_RUN_DOCKER) bash

bundle:
	mkdir -p bundles

coverage-show:
	go tool cover -html=coverage.txt

coverage-report:
	goveralls -coverprofile=coverage.txt -service=travis-ci

clean:
	rm -rf bundles
