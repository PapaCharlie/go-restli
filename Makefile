SHELL := zsh

define get-version
ref=HEAD; tag=""; while true ; do \
  tag=$$(git tag -l "v*" --contains "$$ref" | tail -n1) ; \
  [[ -n "$$tag" ]] && break ; \
  ref="$${ref}^" ; \
done ; \
echo "$${tag#v}"$$([[ HEAD == "$$ref" ]] || echo "-SNAPSHOT")
endef

VERSION := $(shell $(get-version))
JARGO := ./cmd/classpath_jar.go
FAT_JAR := spec-parser/build/libs/go-restli-spec-parser-$(VERSION).jar
GRADLEW := cd spec-parser && ./gradlew -Pversion=$(VERSION)

TESTDATA := internal/tests/testdata
TEST_SUITE := $(TESTDATA)/rest.li-test-suite/client-testsuite
EXTRA_TEST_SUITE := $(TESTDATA)/extra-test-suite
COVERPKG = -coverpkg=github.com/PapaCharlie/go-restli/...
COVERPROFILE = -coverprofile=internal/tests/coverage
GENERATOR_TEST_ARGS = $(COVERPKG) -count=1 -v -tags=jar -args --
PACKAGE_PREFIX := github.com/PapaCharlie/go-restli/internal/tests/testdata/generated
PACKAGES := ./codegen ./d2 $(wildcard ./restli*)
TOTALCOV := ./internal/tests/coverage/total.cov

build: generate test integration-test
	rm -rf bin
	$(MAKE) $(foreach goos,linux darwin,$(foreach goarch,amd64,bin/go-restli_$(goos)-$(goarch)))

bin/go-restli_%:
	export GOOS=$(word 1,$(subst -, ,$(*F))) ; \
	export GOARCH=$(word 2,$(subst -, ,$(*F))) ; \
	go build -tags=jar -ldflags "-s -w -X github.com/PapaCharlie/go-restli/cmd.Version=$(VERSION).$(*F)" -o "$(@)" ./

generate: $(JARGO)
	go generate ./...
	go test . $(COVERPROFILE)/core_generator.cov $(GENERATOR_TEST_ARGS) \
		--output-dir restlidata/generated \
		--package-root github.com/PapaCharlie/go-restli/restlidata/generated \
		--namespace-allow-list com.linkedin.restli.common \
		"$(FAT_JAR)"

test: generate imports
	go test -count=1 $(COVERPKG) $(COVERPROFILE)/protocol.cov $(foreach p,$(PACKAGES),$p/...)

imports:
	goimports -w $$(git ls-files | grep '.go$$' | grep -v '.gr.go$$')

integration-test: generate-restli run-testsuite


generate-restli: clean $(JARGO)
	go test . $(COVERPROFILE)/extras_generator.cov $(GENERATOR_TEST_ARGS) \
		--output-dir $(TESTDATA)/generated_extras \
		--package-root $(PACKAGE_PREFIX)_extras \
		--raw-records extras.Any \
		--dependencies $(FAT_JAR) \
		--manifest-dependencies ./restlidata \
		$(EXTRA_TEST_SUITE)/schemas $(EXTRA_TEST_SUITE)/restspecs
	go test . $(COVERPROFILE)/generator.cov $(GENERATOR_TEST_ARGS) \
		--output-dir $(TESTDATA)/generated \
		--dependencies $(FAT_JAR) \
		--package-root $(PACKAGE_PREFIX) \
		$(TEST_SUITE)/schemas $(TEST_SUITE)/restspecs

run-testsuite:
	go test $(COVERPKG) $(COVERPROFILE)/suite.cov -count=1 ./internal/tests/...
	go test -json ./internal/tests/suite | go run ./internal/tests/coverage
	rm -f $(TOTALCOV)
	gocovmerge ./internal/tests/coverage/*.cov | grep -v .gr.go > $(TOTALCOV)

bench:
	go test -bench=. ./restlicodec


coverage:
	go tool cover -html ./internal/tests/coverage/total.cov

install-gocovmerge:
	go install github.com/wadey/gocovmerge@latest

generate-tests:
	cd internal/tests/suite && go run ./generator

clean:
	git -C $(TEST_SUITE) reset --hard origin/master

fat-jar: $(FAT_JAR)
$(FAT_JAR): $(shell git ls-files spec-parser)
	$(GRADLEW) build fatJar
	touch $(FAT_JAR) # touch the jar after the build to inform make that the file is fresh

jargo: $(JARGO)
$(JARGO): $(FAT_JAR)
	echo -e '//go:build jar\n// +build jar\n\npackage cmd\n\nimport "encoding/base64"\n\nvar Jar, _ = base64.StdEncoding.DecodeString(`' > $(JARGO)
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

codegen/parser/antlr4-4.11.1-complete.jar:
	wget "https://www.antlr.org/download/antlr-4.11.1-complete.jar" -O $@

codegen/parser/jar/go-restli-manifest.gr.json: $(JARGO)
	go test . $(COVERPROFILE)/jar.cov $(GENERATOR_TEST_ARGS) \
		--output-dir $(@D) \
		--package-root garbage \
		--raw-records extras.Any \
		$(FAT_JAR) \
		internal/tests/testdata/extra-test-suite/schemas \
		internal/tests/testdata/rest.li-test-suite/client-testsuite/schemas
	find $(@D) \( -mindepth 1 -maxdepth 1 -type d -or -name all_imports_test.gr.go \) -exec rm -rf {} \;

codegen/parser/gopdl/go-restli-manifest.gr.json: $(JARGO)
	go run ./$(@D) \
		$(FAT_JAR) \
		internal/tests/testdata/extra-test-suite/schemas \
		internal/tests/testdata/rest.li-test-suite/client-testsuite/schemas
	find $(@D) \( -mindepth 1 -maxdepth 1 -type d -or -name all_imports_test.gr.go \) -exec rm -rf {} \;
