package suite

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
	"testing"

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

func TestGoRestli(rootT *testing.T) {
	manifest := ReadManifest()

	s := new(TestServer)
	s.oLock = new(sync.Mutex)
	s.server = httptest.NewServer(s)
	serverUrl, _ := url.Parse(s.server.URL)
	s.client = &protocol.RestLiClient{
		Client:                        &http.Client{},
		HostnameResolver:              &protocol.SimpleHostnameResolver{Hostname: serverUrl},
		StrictResponseDeserialization: true,
	}

	operations := make(map[string]Operation)
	for _, testData := range manifest.WireProtocolTestData {
		if testData.GetClient(s) == nil {
			rootT.Run(testData.Name, func(t *testing.T) {
				for _, o := range testData.Operations {
					t.Run(o.Name, func(t *testing.T) { t.SkipNow() })
				}
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
