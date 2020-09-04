package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	"github.com/PapaCharlie/go-restli/internal/tests/suite"
	"github.com/pkg/errors"
)

func init() {
	utils.PackagePrefix = "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated"
}

func main() {
	server := new(suite.TestServer)
	for _, wd := range suite.ReadManifest().WireProtocolTestData {
		if wd.GetClient(server) == nil {
			// don't generate tests for unsupported resource types
			continue
		}
		var testFileContents string

		testFilename := wd.Name + "_test.go"
		f, err := os.Open(testFilename)
		if err != nil {
			if !os.IsNotExist(err) {
				panicIfErrf(err, "Failed to open test file for %s: %s", wd.Name, testFilename)
			}
			f, err = os.Create(testFilename)
			panicIfErrf(err, "Could not create %s", testFilename)
			_, err = fmt.Fprintf(f, `package tests

import (
	"testing"

	. "%s"
)
`, filepath.Join(utils.PackagePrefix, wd.PackagePath))
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
func (s *TestServer) %s(t *testing.T, c Client) {
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
