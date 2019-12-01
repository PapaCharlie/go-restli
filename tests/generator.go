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
	if err != nil {
		log.Panicln("Could list read", restspecsDir, err)
	}
	if len(restspecs) == 0 {
		log.Panicln("No restspecs found in", restspecsDir)
	}
	sort.Strings(restspecs)

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Panicln("Failed to create temp directory", err)
	}

	err = cli.Run(restspecs, tmpDir)
	if err != nil {
		log.Panicln("Could not generate code from restspecs", err)
	}

	_ = os.RemoveAll(generatedPackageSuffix)
	err = os.Rename(filepath.Join(tmpDir, codegen.PackagePrefix), generatedPackageSuffix)
	if err != nil {
		log.Panicln("Failed to move the generated code", err)
	}

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
			if os.IsNotExist(err) {
				f, err = os.Create(testFilename)
				if err != nil {
					log.Panicln(err)
				}
				_, err = fmt.Fprintf(f, `package main

import (
	"testing"

	. "%s"
)
`, filepath.Join(codegen.PackagePrefix, wd.PackagePath))
				if err != nil {
					log.Panicln(err)
				}
				err = f.Close()
				if err != nil {
					log.Panicln(err)
				}
			} else {
				log.Panicf("Failed to open test file for %s: %s (%+v)", wd.Name, testFilename, err)
			}
		} else {
			var testFileBytes []byte
			testFileBytes, err = ioutil.ReadAll(f)
			if err != nil {
				log.Panicln(err)
			}
			testFileContents = string(testFileBytes)
		}

		f, err = os.OpenFile(testFilename, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Panicln(err)
		}

		for _, o := range wd.Operations {
			if !strings.Contains(testFileContents, o.TestMethodName()) {
				_, err = fmt.Fprintf(f, `
func (s *TestServer) %s(t *testing.T, c *Client) {
	t.SkipNow()
}
`, o.TestMethodName())
				if err != nil {
					log.Panicln(err)
				}
			}
		}
	}
}
