package suite

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/PapaCharlie/go-restli/v2/restli"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var testSuite = "rest.li-test-suite/client-testsuite"

type TestManifest struct {
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
	Name     string
	Request  func(*testing.T) *http.Request
	Response func(*testing.T) *http.Response
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

	testSuite := testSuite
	o.Request = func(t *testing.T) *http.Request {
		filename := filepath.Join(testSuite, "requests-v2", o.Name+".req")
		r := readFile(t, filename)
		req, err := http.ReadRequest(r)
		require.NoError(t, err, "Could not read request from: %q", filename)
		req.Body = adjustContentLength(t, filename, r, req.Header)

		require.NoError(t, restli.DecodeTunnelledQuery(req))

		return req
	}

	o.Response = func(t *testing.T) *http.Response {
		filename := filepath.Join(testSuite, "responses-v2", o.Name+".res")
		r := readFile(t, filename)
		res, err := http.ReadResponse(r, o.Request(t))
		require.NoError(t, err, "Could not read response from: %q", filename)
		res.Body = adjustContentLength(t, filename, r, res.Header)
		return res
	}

	return nil
}

func (o *Operation) TestMethodName() string {
	return strcase.ToCamel(o.Name)
}

func ReadTestManifest() *TestManifest {
	var aggregateManifest TestManifest
	for _, testSuite = range []string{"../testdata/rest.li-test-suite/client-testsuite", "../testdata/extra-test-suite"} {
		f, err := os.Open(filepath.Join(testSuite, "manifest.json"))
		if err != nil {
			log.Panicln(err)
		}

		m := new(TestManifest)
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

func readFile(t *testing.T, filename string) *bufio.Reader {
	reqBytes, err := ioutil.ReadFile(filename)
	require.NoError(t, err, "Could not read %s", filename)
	return bufio.NewReader(bytes.NewBuffer(reqBytes))
}

func adjustContentLength(t *testing.T, filename string, r *bufio.Reader, h http.Header) io.ReadCloser {
	const contentLength = "Content-Length"
	// ReadRequest and ReadResponse only read the leading HTTP protocol bytes (e.g. GET /foo HTTP/1.1) and the headers.
	// What remains of the buffer is the body of the request
	b, err := io.ReadAll(r)
	require.NoError(t, err)

	b = bytes.Trim(b, "\r\n")
	cl := h.Get(contentLength)
	if cl != "" {
		cli, _ := strconv.Atoi(cl)
		if len(b) != cli {
			require.FailNowf(t,
				"Content-Length header in %s indicates %d bytes, but body was %d bytes", filename, cli, len(b))
		}
	}
	return io.NopCloser(bytes.NewReader(b))
}
