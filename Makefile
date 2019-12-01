SHELL := /bin/bash

PACKAGE_PREFIX := github.com/PapaCharlie/go-restli/generated

PACKAGES := ./codegen ./d2 ./protocol

build: generate test integration-test

generate:
	go generate $(PACKAGES)

test: generate imports
	go test $(PACKAGES)

imports:
	goimports -w main.go $(PACKAGES)

integration-test: clean
	go generate ./tests
	go test -count=1 ./tests/...

clean:
	git -C tests/rest.li-test-suite reset --hard origin/master
	rm -rf tests/generated
