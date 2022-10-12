package d2

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"sync"
	"time"

	"github.com/PapaCharlie/go-restli/d2/lazymap"
	"github.com/go-zookeeper/zk"
	"github.com/pkg/errors"
)

type Client struct {
	Conn *zk.Conn
	// During the initial listing of the /d2/uris node for a new cluster, this duration specifies how long to for the
	// first host to show up. If the /d2/services node exists for a service, it is impossible to know whether or not a
	// host will ever show up in the /d2/uris node for that service, which is why this timeout is provided.
	// A positive value will be used as-is, a 0 value will default to DefaultInitialUriWatchTimeout and a negative value
	// disables the timeout altogether
	InitialZkWatchTimeout time.Duration

	// TODO: Support sslSessionValidationStrings. This can be done by opening a connection to the host, pulling its cert
	//  and calling this function. This should happen when the update is received from D2 so that hosts are only added
	//  to a service after getting validated.
	// SSLSessionValidator func(service string, host *url.URL, validationStrings []string, rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error
	// SSLSessionValidatorTimeout time.Duration

	services lazymap.LazySyncMap[string, *zk.NodeCache[*Service]]
	uris     lazymap.LazySyncMap[string, *zk.TreeCache[*Uri]]
	caches   sync.Map
}

const DefaultInitialUriWatchTimeout = 10 * time.Second

func unmarshalJson[T any](data []byte) (*T, error) {
	t := new(T)
	err := json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// getServiceUris returns a snapshot of the current Service and all announced URIs. Both the Service and serviceUris
// objects are read-only, as updates to those objects will be overwritten by the next update from ZK. This also means
// they are thread safe.
func (c *Client) getServiceUris(serviceName string) (service *Service, uris serviceUris, err error) {
	done := make(chan struct{})

	go func() {
		defer close(done)

		var serviceCache *zk.NodeCache[*Service]
		serviceCache, err = c.services.LazyLoad(serviceName, func() (*zk.NodeCache[*Service], error) {
			Logger.Printf("Creating new service watcher for %q", serviceName)
			path := ServicesPath(serviceName)
			return zk.NewNodeCache(c.Conn, path, func(data []byte) (s *Service, err error) {
				return unmarshalJson[Service](data)
			})
		})
		if err != nil {
			return
		}

		service, err = serviceCache.Get()
		if err != nil {
			return
		}

		var urisCache *zk.TreeCache[*Uri]
		urisCache, err = c.uris.LazyLoad(service.ClusterName, func() (*zk.TreeCache[*Uri], error) {
			Logger.Printf("Creating new URI watcher for %q", service.ClusterName)

			return zk.NewTreeCacheWithOpts(
				c.Conn,
				UrisPath(service.ClusterName),
				func(_ string, data []byte) (u *Uri, err error) {
					return unmarshalJson[Uri](data)
				},
				zk.TreeCacheOpts{
					MinRelativeDepth: 2,
					MaxRelativeDepth: 2,
				},
			)
		})
		if err != nil {
			return
		}

		uris = urisCache.Children(urisCache.RootPath)
		for _, v := range uris {
			if err = v.Err; err != nil {
				return
			}
		}
	}()

	select {
	case <-done:
	case <-c.newTimeout():
		return nil, nil, fmt.Errorf("go-restli: Failed to find a valid URI for %q within timeout", service.ClusterName)
	}

	return service, uris, err
}

func (c *Client) newTimeout() <-chan time.Time {
	var d time.Duration
	switch {
	case c.InitialZkWatchTimeout > 0:
		d = c.InitialZkWatchTimeout
	case c.InitialZkWatchTimeout == 0:
		d = DefaultInitialUriWatchTimeout
	default:
		d = math.MaxInt64
	}
	return time.After(d)
}

// SingleServiceClient will always return URIs from the given serviceName instead of getting them from the service
// specified in resourceBaseName in ResolveHostnameAndContextForQuery. Some services announce to special D2 cluster
// instead of the resource name (i.e. by default, the service name is equal to the resource name).
func (c *Client) SingleServiceClient(serviceName string) *SingleServiceClient {
	return &SingleServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

type SingleServiceClient struct {
	c           *Client
	serviceName string
}

func (c *SingleServiceClient) ResolveHostnameAndContextForQuery(_ string, query *url.URL) (*url.URL, error) {
	return c.c.ResolveHostnameAndContextForQuery(c.serviceName, query)
}

func (c *Client) ResolveHostnameAndContextForQuery(rootResource string, _ *url.URL) (*url.URL, error) {
	service, uris, err := c.getServiceUris(rootResource)
	if err != nil {
		return nil, err
	}

	chosenHost := uris.chooseHost(service.PrioritizedSchemes)
	if chosenHost == nil {
		return nil, errors.Errorf("Could not find a host for %q", rootResource)
	}
	return chosenHost, nil
}
