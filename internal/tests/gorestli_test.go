package tests

import (
	"encoding/json"
	"fmt"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/complexkey"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/keywithunion/keywithunion"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime/debug"
	"sync"
	"testing"

	actionset "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/actionSet"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/collection"
	collectionreturnentity "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/collectionReturnEntity"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/params"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/simple"
	collectiontyperef "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/typerefs/collectionTyperef"
	"github.com/PapaCharlie/go-restli/protocol"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func (o *Operation) TestMethod() *reflect.Method {
	if m, ok := reflect.TypeOf(&TestServer{}).MethodByName(o.TestMethodName()); ok {
		return &m
	} else {
		return nil
	}
}

func (d *WireProtocolTestData) GetClient(s *TestServer) *reflect.Value {
	v := new(reflect.Value)
	switch d.Name {
	case "collectionReturnEntity":
		*v = reflect.ValueOf(collectionreturnentity.NewClient(s.client))
		return v
	case "collection":
		*v = reflect.ValueOf(collection.NewClient(s.client))
		return v
	case "simple":
		*v = reflect.ValueOf(simple.NewClient(s.client))
		return v
	case "actionSet":
		*v = reflect.ValueOf(actionset.NewClient(s.client))
		return v
	case "params":
		*v = reflect.ValueOf(params.NewClient(s.client))
		return v
	case "collectionTyperef":
		*v = reflect.ValueOf(collectiontyperef.NewClient(s.client))
		return v
	case "complexkey":
		*v = reflect.ValueOf(complexkey.NewClient(s.client))
		return v
	case "keywithunion":
		*v = reflect.ValueOf(keywithunion.NewClient(s.client))
		return v
	}
	return nil
}

func TestGoRestli(rootT *testing.T) {
	manifest := ReadManifest()

	s := new(TestServer)
	s.oLock = new(sync.Mutex)
	s.server = httptest.NewServer(s)
	serverUrl, _ := url.Parse(s.server.URL)
	s.client = &protocol.RestLiClient{
		Client:           &http.Client{},
		HostnameResolver: &protocol.SimpleHostnameSupplier{Hostname: serverUrl},
	}

	operations := make(map[string]Operation)
	for _, testData := range manifest.WireProtocolTestData {
		if testData.GetClient(s) == nil {
			rootT.Run(testData.Name, func(t *testing.T) {
				t.Skipf("Skipping tests for unsupported resource: \"%s\"", testData.Name)
			})
			continue
		}
		rootT.Run(testData.Name, func(t *testing.T) {
			skippedTests := false
			for _, o := range testData.Operations {
				if dup, ok := operations[o.Name]; ok {
					rootT.Fatalf("Multiple operations named %s: %v and %v", o.Name, o, dup)
				} else {
					operations[o.Name] = o
				}

				client := testData.GetClient(s)
				testMethod := o.TestMethod()
				if testMethod != nil {
					s.oLock.Lock()
					s.o = o
					s.oLock.Unlock()
					t.Run(o.Name, func(t *testing.T) {
						testMethod.Func.Call([]reflect.Value{reflect.ValueOf(s), reflect.ValueOf(t), *client})
						if t.Skipped() {
							skippedTests = true
						}
					})
				} else {
					skippedTests = true
					t.Run(o.Name, func(t *testing.T) {
						t.Skipf("Skipping undefined test \"%s\"", o.Name)
					})
				}
			}
			if skippedTests {
				t.Skip("Some tests were skipped!")
			}
		})
	}
}

type TestServer struct {
	o      Operation
	oLock  *sync.Mutex
	server *httptest.Server
	client *protocol.RestLiClient
}

const UnexpectedRequestStatus = 666

func (s *TestServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// memory barrier for reading the Operation field since requests are served from different goroutines
	s.oLock.Lock()
	s.oLock.Unlock()

	defer func() {
		if r := recover(); r != nil {
			writeErrorResponse(res, "%+v", r)
		}
	}()

	if req.Method != s.o.Request.Method {
		writeErrorResponse(res, `Methods did not match! Expected "%s", got "%s".`,
			s.o.Request.Method, req.Method)
		return
	}

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

	for h := range s.o.Response.Header {
		res.Header().Set(h, s.o.Response.Header.Get(h))
	}
	res.WriteHeader(s.o.Response.StatusCode)
	_, err := res.Write(s.o.ResponseBytes)
	if err != nil {
		log.Panicln(err)
	}
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
