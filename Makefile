SHELL := zsh

define get-version
ref=HEAD; tag=""; while true ; do \
  tag=$$(git tag -l "v*" --contains "$$ref") ; \
  [[ -n "$$tag" ]] && break ; \
  ref="$${ref}^" ; \
done ; \
echo "$${tag#v}"$$([[ HEAD == "$$ref" ]] || echo "-SNAPSHOT")
endef

VERSION := $(shell $(get-version))
JARGO := internal/codegen/cmd/classpath_jar.go
FAT_JAR := spec-parser/build/libs/go-restli-spec-parser-$(VERSION).jar
GRADLEW := cd spec-parser && ./gradlew -Pversion=$(VERSION)

TEST_SUITE := internal/tests/rest.li-test-suite/client-testsuite
PACKAGE_PREFIX := github.com/PapaCharlie/go-restli/internal/tests/generated
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
	rm -rf internal/tests/generated
	go run -tags=jar . \
		--output-dir internal/tests/generated \
		--resolver-path $(TEST_SUITE)/schemas \
		--package-prefix $(PACKAGE_PREFIX) \
		--named-schemas-to-generate testsuite.Primitives \
		--named-schemas-to-generate testsuite.ComplexTypes \
		--named-schemas-to-generate testsuite.Include \
		$(TEST_SUITE)/restspecs/* internal/tests/extra-test-suite/restspecs/*
	go test -tags=jar -count=1 ./internal/tests/...

generate-tests:
	cd internal/tests && go run ./test_generator

clean:
	git -C internal/tests/rest.li-test-suite reset --hard origin/master

fat-jar: $(FAT_JAR)
$(FAT_JAR): $(shell git ls-files spec-parser)
	$(GRADLEW) build fatJar
	touch $(FAT_JAR) # touch the jar after the build to inform make that the file is fresh

$(JARGO): $(FAT_JAR)
	echo -e '// +build jar\n\npackage cmd\n\nimport "encoding/base64"\n\nvar Jar, _ = base64.StdEncoding.DecodeString(`' > $(JARGO)
	gzip -9 -c $(FAT_JAR) | base64 | fold -w 120 >> $(JARGO)
	echo '`)' >> $(JARGO)

release: build
	rm -rf ~/.m2/repository/io/papacharlie/
	$(GRADLEW) publishToMavenLocal
	mkdir -p releases/$(VERSION)
	cp \
		bin/go-restli_darwin-amd64 \
		bin/go-restli_linux-amd64 \
		spec-parser/build/libs/go-restli-spec-parser-$(VERSION).jar \
		spec-parser/build/libs/spec-parser-$(VERSION)-javadoc.jar \
		spec-parser/build/libs/spec-parser-$(VERSION)-sources.jar \
		spec-parser/build/libs/spec-parser-$(VERSION).jar \
		spec-parser/build/publications/mavenJava/pom-default.xml \
		releases/$(VERSION)
