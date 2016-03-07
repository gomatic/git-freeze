.PHONY: build

export GOPATH := $(shell pwd)/build:$(GOPATH)

build: build/src/github.com/nicerobot/git-freeze
	@cd build/src/github.com/nicerobot/git-freeze; go build

build/src/github.com/nicerobot/git-freeze: build/src/github.com/nicerobot
	@ln -s ../../../.. $@

build/src/github.com/nicerobot:
	@mkdir -p $@
