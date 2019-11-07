package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime/debug"
	"sync"
	"testing"

	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/iancoleman/strcase"
)

func TestGoRestli(rootT *testing.T) {
	manifest := ReadManifest()

	s := new(TestServer)
	s.oLock = new(sync.Mutex)
	s.server = httptest.NewServer(s)
	serverUrl, _ := url.Parse(s.server.URL)
	s.client = protocol.RestLiClient{
		Client:           &http.Client{},
		HostnameResolver: &protocol.SimpleHostnameSupplier{Hostname: serverUrl},
	}

	operations := make(map[string]Operation)
	for _, testData := range manifest.WireProtocolTestData {
		rootT.Run(testData.Name, func(t *testing.T) {
			skippedTests := false
			for _, o := range testData.Operations {
				funcName := strcase.ToCamel(o.Name)

				if dup, ok := operations[o.Name]; ok {
					rootT.Fatalf("Multiple operations named %s: %v and %v", o.Name, o, dup)
				} else {
					operations[o.Name] = o
				}

				testMethod := reflect.ValueOf(s).MethodByName(funcName)
				if testMethod != (reflect.Value{}) {
					s.oLock.Lock()
					s.o = o
					s.oLock.Unlock()
					t.Run(o.Name, func(t *testing.T) {
						testMethod.Call([]reflect.Value{reflect.ValueOf(t)})
					})
				} else {
					skippedTests = true
					t.Run(o.Name, func(t *testing.T) {
						t.Skipf("Skipping undefined test \"%s\"", o.Name)
					})
				}
			}
			if skippedTests {
				t.Skip("All tests were skipped!")
			}
		})
	}
}

type TestServer struct {
	o      Operation
	oLock  *sync.Mutex
	server *httptest.Server
	client protocol.RestLiClient
}

const UnexpectedRequestStatus = 666

func (s *TestServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// memory barrier for reading the Operation field since requests are served from different goroutines
	s.oLock.Lock()
	defer s.oLock.Unlock()

	defer func() {
		if r := recover(); r != nil {
			writeErrorResponse(res, "%+v", r)
		}
	}()

	if req.RequestURI != s.o.Request.RequestURI {
		writeErrorResponse(res, `RequestURIs did not match! Expected "%s", got "%s".`,
			s.o.Request.RequestURI, req.RequestURI)
		return
	}

	for h := range s.o.Request.Header {
		if req.Header.Get(h) != s.o.Request.Header.Get(h) {
			writeErrorResponse(res, `%s did not match! Expected "%+v", got "%+v".`,
				h, s.o.Request.Header.Get(h), req.Header.Get(h))
			return
		}
	}

	if len(s.o.RequestBytes) > 0 {
		reqBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writeErrorResponse(res, "Failed to read request: %+v", err)
			return
		}

		expectedMap := make(map[string]interface{})
		_ = json.Unmarshal(s.o.RequestBytes, &expectedMap)

		actualMap := make(map[string]interface{})
		_ = json.Unmarshal(reqBytes, &actualMap)

		if !reflect.DeepEqual(expectedMap, actualMap) {
			writeErrorResponse(res, "Request does not match! Expected\n\n%s\n\nGot\n\n%s",
				expectedMap, actualMap)
			return
		}
	}

	res.WriteHeader(s.o.Response.StatusCode)
	for h := range s.o.Response.Header {
		res.Header().Set(h, s.o.Response.Header.Get(h))
	}
	_, _ = res.Write(s.o.ResponseBytes)
}

func writeErrorResponse(res http.ResponseWriter, format string, args ...interface{}) {
	err := &protocol.RestLiError{
		Status:     UnexpectedRequestStatus,
		StackTrace: string(debug.Stack()),
		Message:    fmt.Sprintf(format, args...),
	}
	response, _ := json.Marshal(err)
	res.Header().Add(protocol.RestLiHeader_ErrorResponse, fmt.Sprint(true))
	http.Error(res, string(response), err.Status)
}
