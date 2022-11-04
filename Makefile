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
FAT_JAR := go-restli-spec-parser.jar
GRADLEW := cd spec-parser && ./gradlew

TESTDATA := internal/tests/testdata
TEST_SUITE := $(TESTDATA)/rest.li-test-suite/client-testsuite
EXTRA_TEST_SUITE := $(TESTDATA)/extra-test-suite
COVERPKG = -coverpkg=github.com/PapaCharlie/go-restli/v2/...
COVERPROFILE = -coverprofile=internal/tests/coverage
GENERATOR_TEST_ARGS = $(COVERPKG) -count=1 -v -args --
PACKAGE_PREFIX := github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated
PACKAGES := ./codegen ./d2 $(wildcard ./restli*)
TOTALCOV := ./internal/tests/coverage/total.cov

build: generate test integration-test
	rm -rf bin
	$(MAKE) $(foreach goos,linux darwin,$(foreach goarch,amd64,bin/go-restli_$(goos)-$(goarch)))

bin/go-restli_%:
	export GOOS=$(word 1,$(subst -, ,$(*F))) ; \
	export GOARCH=$(word 2,$(subst -, ,$(*F))) ; \
	go build -ldflags "-s -w -X github.com/PapaCharlie/go-restli/v2/cmd.Version=$(VERSION).$(*F)" -o "$(@)" ./

generate: $(FAT_JAR)
	go generate ./...
	go test . $(COVERPROFILE)/core_generator.cov $(GENERATOR_TEST_ARGS) \
		--output-dir restlidata/generated \
		--package-root github.com/PapaCharlie/go-restli/v2/restlidata/generated \
		--namespace-allow-list com.linkedin.restli.common \
		$(FAT_JAR)

test: generate imports
	go test -count=1 $(COVERPKG) $(COVERPROFILE)/protocol.cov $(foreach p,$(PACKAGES),$p/...)

imports:
	goimports -w $$(git ls-files | grep '.go$$' | grep -v '.gr.go$$')

integration-test: generate-extras generate-restli run-testsuite

generate-extras: $(FAT_JAR)
	go test . $(COVERPROFILE)/extras_generator.cov $(GENERATOR_TEST_ARGS) \
		--output-dir $(TESTDATA)/generated_extras \
		--package-root $(PACKAGE_PREFIX)_extras \
		--raw-records extras.Any \
		--dependencies $(FAT_JAR) \
		--manifest-dependencies ./restlidata \
		$(EXTRA_TEST_SUITE)/schemas $(EXTRA_TEST_SUITE)/restspecs

generate-restli: clean $(FAT_JAR)
	go test . $(COVERPROFILE)/generator.cov $(GENERATOR_TEST_ARGS) \
		--output-dir $(TESTDATA)/generated \
		--dependencies $(FAT_JAR) \
		--package-root $(PACKAGE_PREFIX) \
		$(TEST_SUITE)/schemas $(TEST_SUITE)/restspecs
	go generate ./...

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
	git submodule update --init --recursive
	git -C $(TEST_SUITE) fetch --all
	git -C $(TEST_SUITE) reset --hard origin/master

fat-jar: $(FAT_JAR)
$(FAT_JAR): $(shell git ls-files spec-parser)
	$(GRADLEW) build fatJar
	ln -f spec-parser/build/libs/go-restli-spec-parser.jar $(FAT_JAR)
