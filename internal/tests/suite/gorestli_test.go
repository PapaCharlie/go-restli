package suite

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"testing"

	actionset "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/actionSet"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection"
	colletionSubCollection "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection/subcollection"
	colletionSubSimple "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection/subsimple"
	collectionreturnentity "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collectionReturnEntity"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/complexkey"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/keywithunion/keywithunion"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/params"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/simple"
	collectiontyperef "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/typerefs/collectionTyperef"
	collectionwithannotations "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithAnnotations"
	collectionwithtyperefkey "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey"
	simplecomplexkey "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleComplexKey"
	simplewithpartialupdate "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleWithPartialUpdate"
	"github.com/PapaCharlie/go-restli/restli"
	"github.com/PapaCharlie/go-restli/restli/batchkeyset"
	"github.com/PapaCharlie/go-restli/restlicodec"
	"github.com/PapaCharlie/go-restli/restlidata"
	"github.com/stretchr/testify/require"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func deliberateSkip(t *testing.T, message string) {
	t.Log("GORESTLI_SKIPPED: " + message)
}

func TestGoRestli(t *testing.T) {
	manifest := ReadManifest()

	operations := make(map[string]*Operation)
	for _, testData := range manifest.WireProtocolTestData {
		for _, o := range testData.Operations {
			if dup, ok := operations[o.Name]; ok {
				t.Fatalf("Multiple operations named %s: %v and %v", o.Name, o, dup)
			} else {
				operations[o.Name] = o
			}

			prefix := testData.Name + "/" + o.Name

			testMethod := o.testMethod()
			var out []reflect.Value

			t.Run(prefix+"/client", func(t *testing.T) {
				if testMethod == nil {
					t.SkipNow()
				}
				out = testMethod.Func.Call([]reflect.Value{
					reflect.ValueOf(o),
					reflect.ValueOf(t),
					o.newClient(t, true),
				})
			})

			t.Run(prefix+"/server", func(t *testing.T) {
				if testMethod == nil || len(out) == 0 {
					t.SkipNow()
				}
				resource := out[0].Call([]reflect.Value{reflect.ValueOf(t)})[0]
				if resource.IsNil() {
					t.SkipNow()
				}
				server := restli.NewServer()
				o.getResource().register.Call([]reflect.Value{reflect.ValueOf(server), resource})

				res := httptest.NewRecorder()
				o.Request.Body = ioutil.NopCloser(bytes.NewReader(o.RequestBytes))
				server.(http.Handler).ServeHTTP(res, o.Request)

				compareResponses(t, o, res.Result())
			})
		}
	}
}

func requireMapEquals[K comparable, V any](t *testing.T, left, right map[K]V) {
	if len(left) == len(right) && left == nil {
		return
	}
	require.Equal(t, left, right)
}

func requiredBatchResponseEquals[K comparable, V restlicodec.Marshaler](t *testing.T, left, right *restlidata.BatchResponse[K, V]) {
	requireMapEquals(t, left.Statuses, right.Statuses)
	requireMapEquals(t, left.Results, right.Results)
	requireMapEquals(t, left.Errors, right.Errors)
}

type ckey[K any] interface {
	comparable
	batchkeyset.ComplexKey[K]
}

func requireComplexKeyMapEquals[K ckey[K], V any](t *testing.T, left, right map[K]V) {
	require.Equal(t, len(left), len(right))

	source := batchkeyset.NewBatchKeySet[K]()
	for k := range left {
		require.NoError(t, source.AddKey(k))
	}

	for k, v := range right {
		leftK, found := source.LocateOriginalKey(k)
		require.True(t, found, "Could not find %+v in left map", k)
		require.Equal(t, left[leftK], v)
	}
}

type supportedResource struct {
	client   reflect.Value
	register reflect.Value
}

type roundTripper func(*http.Request) (*http.Response, error)

func (r roundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return r(request)
}

func (o *Operation) testMethod() *reflect.Method {
	if m, ok := reflect.TypeOf(&Operation{}).MethodByName(o.TestMethodName()); ok {
		return &m
	} else {
		return nil
	}
}

func (o *Operation) getResource() *supportedResource {
	return resourceRegistry[o.testMethod().Type.In(2)]
}

func (o *Operation) newClient(t *testing.T, strictResponseDeserialization bool) reflect.Value {
	c := &restli.Client{
		Client: &http.Client{
			Transport: roundTripper(func(req *http.Request) (*http.Response, error) {
				compareRequests(t, o, req)
				o.Response.Body = ioutil.NopCloser(bytes.NewReader(o.ResponseBytes))
				return o.Response, nil
			}),
		},
		HostnameResolver:              &restli.SimpleHostnameResolver{Hostname: &url.URL{}},
		StrictResponseDeserialization: strictResponseDeserialization,
	}

	return o.getResource().client.Call([]reflect.Value{reflect.ValueOf(c)})[0]
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf(new(T)).Elem()
}

var resourceRegistry = map[reflect.Type]*supportedResource{
	typeOf[collectionreturnentity.Client](): {
		client:   reflect.ValueOf(collectionreturnentity.NewClient),
		register: reflect.ValueOf(collectionreturnentity.RegisterResource),
	},
	typeOf[collection.Client](): {
		client:   reflect.ValueOf(collection.NewClient),
		register: reflect.ValueOf(collection.RegisterResource),
	},
	typeOf[colletionSubCollection.Client](): {
		client:   reflect.ValueOf(colletionSubCollection.NewClient),
		register: reflect.ValueOf(colletionSubCollection.RegisterResource),
	},
	typeOf[colletionSubSimple.Client](): {
		client:   reflect.ValueOf(colletionSubSimple.NewClient),
		register: reflect.ValueOf(colletionSubSimple.RegisterResource),
	},
	typeOf[simple.Client](): {
		client:   reflect.ValueOf(simple.NewClient),
		register: reflect.ValueOf(simple.RegisterResource),
	},
	typeOf[simplewithpartialupdate.Client](): {
		client:   reflect.ValueOf(simplewithpartialupdate.NewClient),
		register: reflect.ValueOf(simplewithpartialupdate.RegisterResource),
	},
	typeOf[actionset.Client](): {
		client:   reflect.ValueOf(actionset.NewClient),
		register: reflect.ValueOf(actionset.RegisterResource),
	},
	typeOf[params.Client](): {
		client:   reflect.ValueOf(params.NewClient),
		register: reflect.ValueOf(params.RegisterResource),
	},
	typeOf[collectiontyperef.Client](): {
		client:   reflect.ValueOf(collectiontyperef.NewClient),
		register: reflect.ValueOf(collectiontyperef.RegisterResource),
	},
	typeOf[complexkey.Client](): {
		client:   reflect.ValueOf(complexkey.NewClient),
		register: reflect.ValueOf(complexkey.RegisterResource),
	},
	typeOf[keywithunion.Client](): {
		client:   reflect.ValueOf(keywithunion.NewClient),
		register: reflect.ValueOf(keywithunion.RegisterResource),
	},
	typeOf[collectionwithtyperefkey.Client](): {
		client:   reflect.ValueOf(collectionwithtyperefkey.NewClient),
		register: reflect.ValueOf(collectionwithtyperefkey.RegisterResource),
	},
	typeOf[collectionwithannotations.Client](): {
		client:   reflect.ValueOf(collectionwithannotations.NewClient),
		register: reflect.ValueOf(collectionwithannotations.RegisterResource),
	},
	typeOf[simplecomplexkey.Client](): {
		client:   reflect.ValueOf(simplecomplexkey.NewClient),
		register: reflect.ValueOf(simplecomplexkey.RegisterResource),
	},
}

func niceHeaders(h http.Header) string {
	buf := new(strings.Builder)
	_ = h.Write(buf)
	return buf.String()
}

func dumpRequest(t *testing.T, req *http.Request, body []byte) string {
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	data, err := httputil.DumpRequest(req, true)
	require.NoError(t, err)
	return string(data)
}

func compareRequests(t *testing.T, left *Operation, right *http.Request) {
	reqBytes, err := ioutil.ReadAll(right.Body)
	require.NoError(t, err)

	equal := func(l, r any, msg string, params ...any) {
		require.Equalf(t, l, r, msg+"\n\nExpected:\n\n%s\n\nGot:\n\n%s", append(params,
			dumpRequest(t, left.Request, left.RequestBytes),
			dumpRequest(t, right, reqBytes),
		)...)
	}

	equal(left.Request.Method, right.Method, "methods did not match")
	equal(left.Request.URL.Path, right.URL.Path, "paths did not match")
	equal(left.Request.URL.RawQuery, right.URL.RawQuery, "queries did not match")

	rightHeaders := right.Header.Clone()
	for k := range left.Request.Header {
		equal(left.Request.Header.Get(k), right.Header.Get(k), "%q header did not match", k)
		rightHeaders.Del(k)
	}
	// go-restli always sends the method header so ignore it if the expected response doesn't have it
	rightHeaders.Del(restli.MethodHeader)
	if len(rightHeaders) != 0 {
		t.Fatalf("Unexpected headers:\n%s", niceHeaders(rightHeaders))
	}

	if len(left.RequestBytes) > 0 {
		if expectedMap, actualMap, match := compareJson(left.RequestBytes, reqBytes); !match {
			require.FailNow(t, "Request does not match! Expected\n\n%s\n\nGot\n\n%s",
				expectedMap, actualMap)
		}
	}

}

func compareAndRemoveMapValue(key string, left, right map[string]any) bool {
	l, r := left[key], right[key]
	isEmptyMap := func(v any) bool {
		value := reflect.ValueOf(v)
		return !value.IsValid() || value.Len() == 0
	}

	if (isEmptyMap(l) && isEmptyMap(r)) || reflect.DeepEqual(l, r) {
		delete(left, key)
		delete(right, key)
		return true
	} else {
		return false
	}
}

func shallowCopy(source map[string]any) map[string]any {
	if source == nil {
		return nil
	}
	dest := make(map[string]any, len(source))
	for k, v := range source {
		dest[k] = v
	}
	return dest
}

func dumpResponse(t *testing.T, res *http.Response, body []byte) string {
	res.Body = ioutil.NopCloser(bytes.NewReader(body))
	data, err := httputil.DumpResponse(res, true)
	require.NoError(t, err)
	return string(data)
}

func compareResponses(t *testing.T, left *Operation, right *http.Response) {
	rightBytes, err := ioutil.ReadAll(right.Body)
	require.NoError(t, err)

	equal := func(l, r any, msg string, params ...any) {
		require.Equalf(t, l, r, msg+"\n\nExpected:\n\n%s\n\nGot:\n\n%s", append(params,
			dumpResponse(t, left.Response, left.ResponseBytes),
			dumpResponse(t, right, rightBytes),
		)...)
	}

	equal(left.Response.StatusCode, right.StatusCode, "Statuses did not match")

	equal(
		dropContentLength(left.Response.Header),
		dropContentLength(right.Header),
		"Headers did not match",
	)

	// It's not super important that the error responses match
	if strings.ToLower(left.Response.Header.Get(restli.ErrorResponseHeader)) == "true" {
		return
	}

	expectedMap, actualMap, _ := compareJson(left.ResponseBytes, rightBytes)

	if reflect.DeepEqual(expectedMap, actualMap) {
		return
	}

	t.Log(string(left.ResponseBytes))
	t.Log(string(rightBytes))

	if strings.Contains(left.Name, "batch") {
		// there's a number of top-level fields that make life difficult, especially since reflect.DeepEqual doesn't
		// consider the nil map and the empty map equal, so these maps need to be checked individually
		expectedCopy, actualCopy := shallowCopy(expectedMap), shallowCopy(actualMap)
		compareAndRemoveMapValue("statuses", expectedCopy, actualCopy)
		compareAndRemoveMapValue("errors", expectedCopy, actualCopy)
		if reflect.DeepEqual(expectedCopy, actualCopy) {
			return
		}
	}

	t.Fatalf("Implementation response did not match expected response! Expected\n\n%s\n\nGot\n\n%s",
		expectedMap, actualMap)
}

func compareJson(expected, actual []byte) (expectedMap, actualMap map[string]any, match bool) {
	expectedMap = make(map[string]any)
	_ = json.Unmarshal(expected, &expectedMap)

	actualMap = make(map[string]any)
	_ = json.Unmarshal(actual, &actualMap)

	return expectedMap, actualMap, reflect.DeepEqual(expectedMap, actualMap)
}

func dropContentLength(h http.Header) http.Header {
	h = h.Clone()
	delete(h, "Content-Length")
	return h
}
