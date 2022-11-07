package restli

import (
	"bufio"
	"context"
	"crypto/rand"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"

	"github.com/PapaCharlie/go-restli/v2/restlidata/generated/com/linkedin/restli/common"
	"github.com/stretchr/testify/require"
)

func TestQueryTunnellingWithBody(t *testing.T) {
	testQueryTunnelling(t, true)
}

func TestQueryTunnellingWithoutBody(t *testing.T) {
	testQueryTunnelling(t, false)
}

func testQueryTunnelling(t *testing.T, withBody bool) {
	t.Run("encode", func(t *testing.T) {
		// This reader always returns 0, so the boundary will be consistent, though to prevent breaking other tests,
		// reset it back to the original value once the test finishes
		defer func(oldReader io.Reader) {
			rand.Reader = oldReader
		}(rand.Reader)
		rand.Reader = zeroReader{}

		// Read the expected tunnelled request
		expectedReq := readTestTunnellingRequest(t, withBody)
		// Generate a new tunnelled request (setting the threshold to 1 effectively enables tunnelling for all requests)
		actualReq := newTestTunnellingRequest(t, 1, withBody)

		compareRequests(t, expectedReq, actualReq)
	})

	t.Run("decode", func(t *testing.T) {
		// Setting the threshold to 0 disables tunnelling, so we should get an untunnelled request to compare against
		expectedReq := newTestTunnellingRequest(t, 0, withBody)

		actualReq := readTestTunnellingRequest(t, withBody)
		// Untunnel the request after reading the tunnelled request from disk
		require.NoError(t, DecodeTunnelledQuery(actualReq))

		compareRequests(t, expectedReq, actualReq)
	})
}

type zeroReader struct{}

func (z zeroReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

func newTestTunnellingRequest(t *testing.T, tunnellingThreshold int, withBody bool) (req *http.Request) {
	c := &Client{
		HostnameResolver: &SimpleHostnameResolver{Hostname: &url.URL{
			Scheme: "https",
			Host:   "localhost",
		}},
		QueryTunnellingThreshold: tunnellingThreshold,
	}
	ctx := context.TODO()
	rp := ResourcePathString("/testTunnelling")
	qp := QueryParamsString("param=bar")

	var err error

	if withBody {
		req, err = NewJsonRequest(c, ctx, rp, qp, http.MethodPut, Method_partial_update, new(common.EmptyRecord), nil)
	} else {
		req, err = NewGetRequest(c, ctx, rp, qp, Method_get)
	}

	require.NoError(t, err)
	return req
}

func readTestTunnellingRequest(t *testing.T, withBody bool) *http.Request {
	var request string
	if withBody {
		request = strings.Join([]string{
			"POST /testTunnelling HTTP/1.1\r",
			"Host: localhost\r",
			"Accept: application/json\r",
			"Content-Type: multipart/mixed; boundary=000000000000000000000000000000000000000000000000000000000000\r",
			"X-Http-Method-Override: PUT\r",
			"X-Restli-Method: partial_update\r",
			"X-Restli-Protocol-Version: 2.0.0\r",
			"\r",
			"--000000000000000000000000000000000000000000000000000000000000\r",
			"Content-Type: application/x-www-form-urlencoded\r",
			"\r",
			"param=bar\r",
			"--000000000000000000000000000000000000000000000000000000000000\r",
			"Content-Type: application/json\r",
			"\r",
			"{}\r",
			"--000000000000000000000000000000000000000000000000000000000000--\r",
			"",
		}, "\n")
	} else {
		request = strings.Join([]string{
			"POST /testTunnelling HTTP/1.1\r",
			"Host: localhost\r",
			"Accept: application/json\r",
			"Content-Type: application/x-www-form-urlencoded\r",
			"X-Http-Method-Override: GET\r",
			"X-Restli-Method: get\r",
			"X-Restli-Protocol-Version: 2.0.0\r",
			"\r",
			"param=bar",
		}, "\n")
	}

	r := bufio.NewReader(strings.NewReader(request))
	expectedReq, err := http.ReadRequest(r)
	expectedReq.Body = io.NopCloser(r)
	require.NoError(t, err)

	return expectedReq
}

func compareRequests(t *testing.T, expectedReq, actualReq *http.Request) {
	dump := func(req *http.Request) string {
		data, err := httputil.DumpRequest(req, true)
		require.NoError(t, err)
		return string(data)
	}

	require.Equal(t, dump(expectedReq), dump(actualReq))
}
