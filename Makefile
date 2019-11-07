SHELL := /bin/bash

PACKAGE_PREFIX := github.com/PapaCharlie/go-restli/generated

SNAPSHOTS ?= $(sort $(shell find . -name '*.snapshot.json'))
.PHONY: $(SNAPSHOTS)
PACKAGES := ./codegen ./d2 ./protocol

build: generate test integration-test

generate:
	go generate $(PACKAGES)

test: generate imports
	go test $(PACKAGES)

imports:
	goimports -w main.go $(PACKAGES)

integration-test: # clean
	go generate ./tests
	go test -count=1 ./tests/...

$(SNAPSHOTS): clean
	mkdir -p tmp
	go run main.go \
		--package-prefix $(PACKAGE_PREFIX) \
		--output-dir tmp \
		$(@)
	mv tmp/$(PACKAGE_PREFIX) .
	rm -r tmp
	go test -count=1 $(PACKAGE_PREFIX)

clean:
	git -C tests/rest.li-test-suite reset --hard origin/master
	rm -rf generated
