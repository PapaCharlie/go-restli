SHELL := zsh

VERSION ?= $(shell make -s get-latest-version)
JARGO := internal/codegen/cmd/classpath_jar.go
FAT_JAR := spec-parser/build/libs/go-restli-spec-parser.jar

PACKAGE_PREFIX := github.com/PapaCharlie/go-restli/generated
PACKAGES := ./internal/codegen ./d2 ./protocol

build: generate test integration-test
	rm -rf bin
	make $(foreach goos,linux darwin,$(foreach goarch,amd64,bin/go-restli_$(goos)-$(goarch)))

bin/go-restli_%: $(shell git ls-files | grep "\.go")
	export GOOS=$(word 1,$(subst -, ,$(*F))) ; \
	export GOARCH=$(word 2,$(subst -, ,$(*F))) ; \
	go build -tags=jar -ldflags "-s -w -X github.com/PapaCharlie/go-restli/internal/codegen/cmd.Version=$(VERSION).$(*F)" -o "$(@)" ./

generate:
	go generate $(PACKAGES)

test: generate imports
	go test $(PACKAGES)

imports:
	goimports -w main.go $(PACKAGES)

integration-test: clean $(JARGO)
	cd internal/tests && go run -tags=jar ./test_generator
	go test -tags=jar -count=1 ./internal/tests/...

clean:
	git -C internal/tests/rest.li-test-suite reset --hard origin/master
	rm -rf internal/tests/generated

$(FAT_JAR): $(shell git ls-files spec-parser)
	cd spec-parser && ./gradlew build fatJar
	touch $(FAT_JAR) # touch the jar after the build to inform make that the file is fresh

$(JARGO): $(FAT_JAR)
	echo -e '// +build jar\n\npackage cmd\n\nimport "encoding/base64"\n\nvar Jar, _ = base64.StdEncoding.DecodeString(`' > $(JARGO)
	gzip -9 -c $(FAT_JAR) | base64 -w 120 >> $(JARGO)
	echo '`)' >> $(JARGO)

get-latest-version:
	@ref=HEAD; tag=""; while [[ -z "$$tag" ]] ; do tag=$$(git tag -l "v*" --contains "$$ref") ; ref="$${ref}^" ; done && echo "$${tag#v}"$$([[ HEAD == "$$ref" ]] || echo "-SNAPSHOT")

release:
	rm -rf ~/.m2/repository/io/papacharlie/
	cd spec-parser && ./gradlew -Pversion=$(VERSION) publishToMavenLocal
