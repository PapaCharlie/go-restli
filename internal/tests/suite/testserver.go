package suite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/actionSet"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collectionReturnEntity"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/complexkey"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/keywithunion/keywithunion"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/params"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/simple"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/typerefs/collectionTyperef"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithAnnotations"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleComplexKey"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleWithPartialUpdate"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/PapaCharlie/go-restli/protocol/stdtypes"
)

type TestServer struct {
	o      Operation
	oLock  *sync.Mutex
	server *httptest.Server
	client *protocol.RestLiClient
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
	case "simpleWithPartialUpdate":
		*v = reflect.ValueOf(simplewithpartialupdate.NewClient(s.client))
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
	case "collectionWithTyperefKey":
		*v = reflect.ValueOf(collectionwithtyperefkey.NewClient(s.client))
		return v
	case "collectionWithAnnotations":
		*v = reflect.ValueOf(collectionwithannotations.NewClient(s.client))
		return v
	case "simpleComplexKey":
		*v = reflect.ValueOf(simplecomplexkey.NewClient(s.client))
		return v
	}
	return nil
}

const UnexpectedRequestStatus = 666
const MismatchedQueriesStatus = 667

func (s *TestServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// memory barrier for reading the Operation field since requests are served from different goroutines
	s.oLock.Lock()
	s.oLock.Unlock()

	defer func() {
		if r := recover(); r != nil {
			writeErrorResponse(res, "%+v: %s", r, string(debug.Stack()))
		}
	}()

	if expected, got := s.o.Request.Method, req.Method; expected != got {
		writeErrorResponse(res, "Methods did not match! Expected %q, got %q.", expected, got)
		return
	}

	if expected, got := s.o.Request.URL.Path, req.URL.Path; expected != got {
		writeErrorResponse(res, "Request paths did not match! Expected %q, got %q.", expected, got)
		return
	}

	if expected, got := s.o.Request.URL.RawQuery, req.URL.RawQuery; expected != got {
		writeErrorResponseWithStatus(res, MismatchedQueriesStatus, "Request queries did not match!\nExpected: %q\nGot:      %q.", expected, got)
		return
	}

	for h := range s.o.Request.Header {
		if req.Header.Get(h) != s.o.Request.Header.Get(h) {
			writeErrorResponse(res, "%s did not match! Expected %q, got %q.",
				h, s.o.Request.Header.Get(h), req.Header.Get(h))
			return
		}
	}

	if len(s.o.RequestBytes) > 0 {
		reqBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writeErrorResponse(res, "Failed to read request: %q", err)
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
	writeErrorResponseWithStatus(res, UnexpectedRequestStatus, format, args...)
}

func writeErrorResponseWithStatus(res http.ResponseWriter, status int, format string, args ...interface{}) {
	err := &protocol.RestLiError{
		ErrorResponse: stdtypes.ErrorResponse{
			StackTrace: protocol.StringPointer(string(debug.Stack())),
			Message:    protocol.StringPointer(fmt.Sprintf(format, args...)),
			Status:     protocol.Int32Pointer(int32(status)),
		},
	}
	response, _ := json.Marshal(err)
	res.Header().Add(protocol.RestLiHeader_ErrorResponse, fmt.Sprint(true))
	http.Error(res, string(response), status)
}
