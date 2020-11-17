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

	actionset "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/actionSet"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection"
	collectionreturnentity "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collectionReturnEntity"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/complexkey"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/keywithunion/keywithunion"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/params"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/simple"
	collectiontyperef "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/typerefs/collectionTyperef"
	collectionwithannotations "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithAnnotations"
	collectionwithtyperefkey "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey"
	simplecomplexkey "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleComplexKey"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
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

	if err := queriesEqual(s.o.Request.URL.RawQuery, req.URL.RawQuery); err != nil {
		writeErrorResponse(res, err.Error())
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
	err := &protocol.RestLiError{
		StackTrace: string(debug.Stack()),
		Message:    fmt.Sprintf(format, args...),
	}
	response, _ := json.Marshal(err)
	res.Header().Add(protocol.RestLiHeader_ErrorResponse, fmt.Sprint(true))
	http.Error(res, string(response), UnexpectedRequestStatus)
}

func queriesEqual(expectedStr, actualStr string) error {
	if expectedStr == actualStr {
		return nil
	}
	read := func(rawQuery string) (map[string]interface{}, error) {
		reader := restlicodec.NewRestLiQueryParamsReader(rawQuery)
		query := make(map[string]interface{})
		err := reader.ReadParams(func(reader restlicodec.Reader, field string) error {
			raw, err := reader.ReadInterface()
			if err != nil {
				return err
			}
			query[field] = raw
			return nil
		})
		if err != nil {
			return nil, err
		}
		return query, nil
	}

	expected, err := read(expectedStr)
	if err != nil {
		return err
	}
	actual, err := read(actualStr)
	if err != nil {
		return err
	}

	if len(expected) != len(actual) {
		return fmt.Errorf("query parameter count differs (expected: %d, actual: %d)", len(expected), len(actual))
	}

	for k, v := range expected {
		actualV, ok := actual[k]
		if !ok {
			return fmt.Errorf("query parameter missing in actual: %q", k)
		}
		err = deepEquals(nil, v, actualV)
		if err != nil {
			return err
		}
	}

	return nil
}

func deepEquals(scope []string, expected, actual interface{}) error {
	switch expected := expected.(type) {
	case map[string]interface{}:
		actual := actual.(map[string]interface{})
		if len(expected) != len(actual) {
			return fmt.Errorf("%v: unequal map sizes (expected: %d, actual: %d)", scope, len(expected), len(actual))
		}
		for k, v := range expected {
			actualV, ok := actual[k]
			if !ok {
				return fmt.Errorf("%v: actual map does not contain key %q", scope, k)
			}
			err := deepEquals(append(append([]string(nil), scope...)), v, actualV)
			if err != nil {
				return err
			}
		}
		return nil
	case []interface{}:
		actual := actual.([]interface{})
		if len(expected) != len(actual) {
			return fmt.Errorf("%v: unequal array sizes (expected: %d, actual: %d)", scope, len(expected), len(actual))
		}
		for index, v := range expected {
			err := deepEquals(append(append([]string(nil), scope...)), v, actual[index])
			if err != nil {
				return err
			}
		}
		return nil
	default:
		if expected.(string) != actual.(string) {
			return fmt.Errorf("%v: unequal primitives (expected: %s, actual: %s)", scope, expected, actual)
		}
		return nil
	}
}
