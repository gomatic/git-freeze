.PHONY: build

export GOPATH := $(shell pwd)/build:$(GOPATH)

build: build/src/github.com/gomatic/git-freeze
	@cd build/src/github.com/gomatic/git-freeze; go build

build/src/github.com/gomatic/git-freeze: build/src/github.com/gomatic
	@ln -s ../../../.. $@

build/src/github.com/gomatic:
	@mkdir -p $@
