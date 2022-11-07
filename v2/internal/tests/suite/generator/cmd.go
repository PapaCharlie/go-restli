package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PapaCharlie/go-restli/v2/internal/tests/suite"
	"github.com/pkg/errors"
)

const (
	defaultPrefix = "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated"
	extrasPrefix  = defaultPrefix + "_extras"
)

func main() {
	for _, wd := range suite.ReadTestManifest().WireProtocolTestData {
		var testFileContents string

		testFilename := wd.Name + "_test.go"
		f, err := os.Open(testFilename)
		if err != nil {
			if !os.IsNotExist(err) {
				panicIfErrf(err, "Failed to open test file for %s: %s", wd.Name, testFilename)
			}
			f, err = os.Create(testFilename)
			panicIfErrf(err, "Could not create %s", testFilename)
			var importPath string
			if strings.HasPrefix(wd.PackagePath, "extras") {
				importPath = filepath.Join(extrasPrefix, wd.PackagePath)
			} else {
				importPath = filepath.Join(defaultPrefix, wd.PackagePath)
			}
			_, err = fmt.Fprintf(f, `package suite

import (
	"testing"

	. "%s"
)
`, importPath)
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
func (o *Operation) %s(t *testing.T, c Client) {
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
