package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/cli"
	"github.com/pkg/errors"
)

const (
	restLiClientTestSuite  = "rest.li-test-suite/client-testsuite"
	generatedPackageSuffix = "generated"
)

func init() {
	codegen.PdscDirectory = filepath.Join(restLiClientTestSuite, "schemas")
	codegen.PackagePrefix = "github.com/PapaCharlie/go-restli/tests/" + generatedPackageSuffix
}

//go:generate go run .
func main() {
	log.SetFlags(log.Lshortfile)

	restspecsDir := filepath.Join(restLiClientTestSuite, "restspecs")
	restspecs, err := filepath.Glob(restspecsDir + "/*.restspec.json")
	panicIfErrf(err, "No restspecs found in %s", restspecsDir)
	if len(restspecs) == 0 {
		log.Panicln("No restspecs found in", restspecsDir)
	}
	sort.Strings(restspecs)

	tmpDir, err := ioutil.TempDir("", "")
	panicIfErrf(err, "Failed to create temp directory")

	err = cli.Run(restspecs, tmpDir)
	panicIfErrf(err, "Could not generate code from restspecs")

	_ = os.RemoveAll(generatedPackageSuffix)
	err = os.Rename(filepath.Join(tmpDir, codegen.PackagePrefix), generatedPackageSuffix)
	panicIfErrf(err, "Failed to move the generated code")

	_ = os.RemoveAll(tmpDir)
	//generateClientTests()
}

// generateClientTests is ignored by default, but it can be used to bootstrap the test framework by generating empty
// tests for all the tests that need to be implemented
func generateClientTests() {
	for _, wd := range ReadManifest().WireProtocolTestData {
		var testFileContents string

		testFilename := wd.Name + "_test.go"
		f, err := os.Open(testFilename)
		if err != nil {
			if ! os.IsNotExist(err) {
				panicIfErrf(err, "Failed to open test file for %s: %s", wd.Name, testFilename)
			}
			f, err = os.Create(testFilename)
			panicIfErrf(err, "Could not create %s", testFilename)
			_, err = fmt.Fprintf(f, `package main

import (
	"testing"

	. "%s"
)
`, filepath.Join(codegen.PackagePrefix, wd.PackagePath))
			panicIfErr(err)
			panicIfErr(f.Close())
		} else {
			var testFileBytes []byte
			testFileBytes, err = ioutil.ReadAll(f)
			panicIfErr(err)
			testFileContents = string(testFileBytes)
		}

		f, err = os.OpenFile(testFilename, os.O_WRONLY|os.O_APPEND, 0666)
		panicIfErr(err)

		for _, o := range wd.Operations {
			if !strings.Contains(testFileContents, o.TestMethodName()) {
				_, err = fmt.Fprintf(f, `
func (s *TestServer) %s(t *testing.T, c *Client) {
	t.SkipNow()
}
`, o.TestMethodName())
				panicIfErr(err)
			}
		}
	}
}

func panicIfErr(err error) {
	if err = errors.WithStack(err); err != nil {
		log.Fatalf("%+v", err)
	}
}

func panicIfErrf(err error, fmt string, args ...interface{}) {
	if err = errors.Wrapf(err, fmt, args...); err != nil {
		log.Fatalf("%+v", err)
	}
}
