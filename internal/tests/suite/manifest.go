package suite

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PapaCharlie/go-restli/internal/tests"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

var testSuite = "rest.li-test-suite/client-testsuite"

type Manifest struct {
	JsonTestData []struct {
		Data string `json:"data"`
	} `json:"jsonTestData"`
	SchemaTestData []struct {
		Schema string `json:"schema"`
		Data   string `json:"data"`
	} `json:"schemaTestData"`
	WireProtocolTestData []WireProtocolTestData `json:"wireProtocolTestData"`
}

type WireProtocolTestData struct {
	Name        string
	PackagePath string
	Operations  []*Operation
}

type Operation struct {
	Name          string
	Request       *http.Request
	RequestBytes  []byte
	Response      *http.Response
	ResponseBytes []byte
}

func (d *WireProtocolTestData) UnmarshalJSON(data []byte) error {
	testData := &struct {
		Name       string       `json:"name"`
		Restspec   string       `json:"restspec"`
		Operations []*Operation `json:"operations"`
	}{}
	err := json.Unmarshal(data, testData)
	if err != nil {
		return errors.WithStack(err)
	}

	d.Name = testData.Name
	d.PackagePath = strings.Replace(strings.TrimSuffix(strings.TrimPrefix(testData.Restspec, "restspecs/"), ".restspec.json"), ".", "/", -1)
	d.Operations = testData.Operations

	return nil
}

func (o *Operation) UnmarshalJSON(data []byte) error {
	var operation struct {
		Name string `json:"name"`
	}

	err := json.Unmarshal(data, &operation)
	if err != nil {
		return errors.WithStack(err)
	}

	o.Name = operation.Name

	o.Request, o.RequestBytes, err = tests.ReadRequestFromFile(filepath.Join(testSuite, "requests-v2", o.Name+".req"))
	if err != nil {
		return errors.WithStack(err)
	}

	o.Response, o.ResponseBytes, err = tests.ReadResponseFromFile(filepath.Join(testSuite, "responses-v2", o.Name+".res"), o.Request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (o *Operation) TestMethodName() string {
	return strcase.ToCamel(o.Name)
}

func ReadManifest() *Manifest {
	var aggregateManifest Manifest
	for _, testSuite = range []string{"../testdata/rest.li-test-suite/client-testsuite", "../testdata/extra-test-suite"} {
		f, err := os.Open(filepath.Join(testSuite, "manifest.json"))
		if err != nil {
			log.Panicln(err)
		}

		m := new(Manifest)
		err = json.NewDecoder(f).Decode(m)
		if err != nil {
			log.Panicln(err)
		}
		aggregateManifest.JsonTestData = append(aggregateManifest.JsonTestData, m.JsonTestData...)
		aggregateManifest.SchemaTestData = append(aggregateManifest.SchemaTestData, m.SchemaTestData...)
		for _, wd := range m.WireProtocolTestData {
			// Association resources will likely never be supported so skip loading them altogether
			if wd.Name == "association" || wd.Name == "associationTyperef" {
				continue
			}
			aggregateManifest.WireProtocolTestData = append(aggregateManifest.WireProtocolTestData, wd)
		}
	}
	return &aggregateManifest
}
