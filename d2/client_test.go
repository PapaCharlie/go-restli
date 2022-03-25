package d2

import (
	"math"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testServiceName = "testService"
	testClusterName = "TestCluster"
)

func init() {
	rng.Seed(0) // get consistent tests
	EnableD2Logging()
}

var (
	serviceDefinitionHttpsAndHttp = []byte(`{
  "serviceName": "` + testServiceName + `",
  "clusterName": "` + testClusterName + `",
  "prioritizedSchemes": [
    "https",
    "http"
  ]
}`)

	serviceDefinitionHttpOnly = []byte(`{
  "serviceName": "` + testServiceName + `",
  "clusterName": "` + testClusterName + `",
  "prioritizedSchemes": [
    "http"
  ]
}`)

	serviceDefinitionNoPrioritizedSchemes = []byte(`{
  "serviceName": "` + testServiceName + `",
  "clusterName": "` + testClusterName + `",
  "prioritizedSchemes": [ ]
}`)
)

type host struct {
	url  url.URL
	data []byte
}

func newHost(u, data string) (h host) {
	parsed, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return host{
		url:  *parsed,
		data: []byte(data),
	}
}

const (
	httpsAndHttp = "httpAndHttp"
	httpsOnly    = "httpsOnly"
	httpOnly     = "httpOnly"
	heavyweight  = "heavyweight"
)

var (
	httpsAndHttpHost = newHost("https://"+httpsAndHttp+":443", `{
  "weights": {
    "https://`+httpsAndHttp+`:443": 1,
    "http://`+httpsAndHttp+`:80": 1
  }
}`)

	httpsOnlyHost = newHost("https://"+httpsOnly+":443", `{
  "weights": {
    "https://`+httpsOnly+`:443": 1
  }
}`)

	httpOnlyHost = newHost("http://"+httpOnly+":80", `{
  "weights": {
    "http://`+httpOnly+`:80": 1
  }
}`)

	heavyweightHost = newHost("http://"+heavyweight+":443", `{
  "weights": {
    "https://`+heavyweight+`:443": 99
  }
}`)
)

func (c *Client) spoofUpdate(event TreeCacheEvent, f func(chan TreeCacheEvent)) {
	events := make(chan TreeCacheEvent)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		f(events)
		wg.Done()
	}()

	events <- event
	close(events)

	wg.Wait()
}

func (c *Client) spoofServiceUpdate(data *[]byte) {
	c.spoofUpdate(TreeCacheEvent{
		Path: ServicesPath(testServiceName),
		Data: data,
	}, func(events chan TreeCacheEvent) {
		c.waitForServiceUpdates(testServiceName, events)
	})
	// initialize an empty serviceUris object to make sure the code does not try to call ZK
	c.uris.LoadOrStore(testClusterName, func() interface{} {
		return &serviceUris{zkPath: UrisPath(testClusterName)}
	})
}

func (c *Client) spoofUriUpdate(h host) {
	c.spoofUpdate(TreeCacheEvent{
		Path: h.url.Hostname(),
		Data: &h.data,
	}, func(events chan TreeCacheEvent) {
		c.waitForUriUpdates(testClusterName, events)
	})
}

func (c *Client) spoofUriDelete(h host) {
	c.spoofUpdate(TreeCacheEvent{
		Path: h.url.Hostname(),
		Data: nil,
	}, func(events chan TreeCacheEvent) {
		c.waitForUriUpdates(testClusterName, events)
	})
}

func TestR2D2Client_BasicTest(t *testing.T) {
	c := new(Client)

	c.spoofServiceUpdate(&serviceDefinitionHttpsAndHttp)

	_, err := c.ResolveHostnameAndContextForQuery(testServiceName, nil)
	require.Error(t, err)

	c.spoofUriUpdate(httpsAndHttpHost)
	_, err = c.ResolveHostnameAndContextForQuery(testServiceName, nil)
	require.NoError(t, err)
}

func TestR2D2Client_PrioritizedSchemes(t *testing.T) {
	c := new(Client)

	c.spoofServiceUpdate(&serviceDefinitionHttpsAndHttp)

	c.spoofUriUpdate(httpsOnlyHost)
	c.spoofUriUpdate(httpOnlyHost)

	ratios := hostRatios(t, c)
	for h := range ratios {
		// since the top priority scheme is https, no host should have http as a scheme
		require.Equal(t, "https", h.Scheme)
	}

	c.spoofServiceUpdate(&serviceDefinitionHttpOnly)

	ratios = hostRatios(t, c)
	for h := range ratios {
		// now the the only priority scheme is http, so no host should have https as a scheme
		require.Equal(t, "http", h.Scheme)
	}

	// spoof the http-only host going offline. Since the current scheme is http-only, we should not be able to get new
	// hosts
	c.spoofUriDelete(httpOnlyHost)

	_, err := c.ResolveHostnameAndContextForQuery(testServiceName, nil)
	require.Error(t, err)
}

func TestR2D2Client_CallDistribution(t *testing.T) {
	c := new(Client)

	c.spoofServiceUpdate(&serviceDefinitionNoPrioritizedSchemes)

	c.spoofUriUpdate(httpOnlyHost)
	c.spoofUriUpdate(httpsOnlyHost)

	ratios := hostRatios(t, c)
	require.Less(t, math.Abs(ratios[httpOnlyHost.url]-ratios[httpsOnlyHost.url]), 0.01)

	// remove the httpOnly host and add the heavyweight host
	c.spoofUriDelete(httpOnlyHost)
	c.spoofUriUpdate(heavyweightHost)

	ratios = hostRatios(t, c)
	t.Log(ratios)
	require.Less(t, ratios[httpsOnlyHost.url]-0.01, 0.01)
	require.Less(t, ratios[heavyweightHost.url]-0.99, 0.01)
}

func hostRatios(t *testing.T, c *Client) map[url.URL]float64 {
	const samples = 100_000
	ratios := make(map[url.URL]float64)
	for range [samples]struct{}{} {
		h, err := c.ResolveHostnameAndContextForQuery(testServiceName, nil)
		require.NoError(t, err)
		ratios[*h] += 1
	}

	for k, v := range ratios {
		ratios[k] = v / samples
	}

	return ratios
}
