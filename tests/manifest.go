package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen/cli"
	"github.com/PapaCharlie/go-restli/codegen/schema"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/pkg/errors"
)

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
	Name       string
	Schema     *schema.Resource
	Operations []Operation
}

type Operation struct {
	Name          string
	Method        *protocol.RestLiMethod
	Finder        *string
	Action        *string
	Request       *http.Request
	RequestBytes  []byte
	Response      *http.Response
	ResponseBytes []byte
	Status        int
}

func (d *WireProtocolTestData) UnmarshalJSON(data []byte) error {
	testData := &struct {
		Name       string      `json:"name"`
		Snapshot   string      `json:"snapshot"`
		Operations []Operation `json:"operations"`
	}{}
	err := json.Unmarshal(data, testData)
	if err != nil {
		return errors.WithStack(err)
	}

	d.Name = testData.Name
	d.Operations = testData.Operations

	snapshot, err := cli.LoadSnapshotFromFile(filepath.Join(restLiClientTestSuite, testData.Snapshot))
	if err != nil {
		return errors.WithStack(err)
	}
	d.Schema = snapshot.Schema

	return nil
}

func (o *Operation) UnmarshalJSON(data []byte) error {
	operation := &struct {
		Name     string `json:"name"`
		Method   string `json:"method"`
		Request  string `json:"request"`
		Response string `json:"response"`
		Status   int    `json:"status"`
	}{}

	err := json.Unmarshal(data, operation)
	if err != nil {
		return errors.WithStack(err)
	}

	o.Name = operation.Name

	const (
		actionMethodPrefix = "action:"
		finderMethodPrefix = "finder:"
	)

	switch {
	case strings.HasPrefix(operation.Method, actionMethodPrefix):
		action := strings.TrimPrefix(operation.Method, actionMethodPrefix)
		o.Action = &action
	case strings.HasPrefix(operation.Method, finderMethodPrefix):
		finder := strings.TrimPrefix(operation.Method, finderMethodPrefix)
		o.Finder = &finder
	default:
		if m, ok := protocol.RestLiMethodNameMapping[operation.Method]; ok {
			o.Method = &m
		} else {
			return errors.Errorf("No such method: %s", operation.Method)
		}
	}

	o.Status = operation.Status

	o.Request, o.RequestBytes, err = ReadRequestFromFile(changeToV2Path(operation.Request))
	if err != nil {
		return errors.WithStack(err)
	}

	o.Response, o.ResponseBytes, err = ReadResponseFromFile(changeToV2Path(operation.Response), o.Request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Manifest) GetOperations(selector func(testData WireProtocolTestData, operation Operation) bool) []Operation {
	var operations []Operation
	for _, testData := range m.WireProtocolTestData {
		for _, o := range testData.Operations {
			if selector(testData, o) {
				operations = append(operations, o)
			}
		}
	}
	return operations
}

func ReadManifest() *Manifest {
	f := filepath.Join(restLiClientTestSuite, "manifest.json")
	m := new(Manifest)
	err := cli.ReadJSONFromFile(f, m)
	if err != nil {
		log.Panicln(err)
	}
	return m
}

// This client only supports Rest.li 2.0, so we need to look at the -v2 requests/responses
func changeToV2Path(path string) string {
	return filepath.Join(restLiClientTestSuite, filepath.Dir(path)+"-v2", filepath.Base(path))
}
