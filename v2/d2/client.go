package d2

import (
	"encoding/json"
	"math"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PapaCharlie/go-restli/v2/d2/lazymap"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
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

	services lazymap.LazySyncMap
	uris     lazymap.LazySyncMap
	caches   sync.Map
}

const DefaultInitialUriWatchTimeout = 10 * time.Second

// getServiceUris returns a snapshot of the current Service and all announced URIs. Both the Service and serviceUris
// objects are read-only, as updates to those objects will be overwritten by the next update from ZK. This also means
// they are thread safe.
func (c *Client) getServiceUris(serviceName string) (*Service, *serviceUris, error) {
	timeout := c.newTimeout()

	var serviceEvents chan TreeCacheEvent
	s := c.services.LoadOrStore(serviceName, func() interface{} {
		path := ServicesPath(serviceName)

		exists, _, err := c.Conn.Exists(path)
		if err != nil {
			return err
		}
		if !exists {
			return errors.Errorf("No service node found at %q", path)
		}

		serviceEvents = make(chan TreeCacheEvent)
		c.caches.Store(path, NewTreeCache(c.Conn, path, serviceEvents))

		for {
			select {
			case e := <-serviceEvents:
				s := c.handleServiceUpdate(serviceName, e)
				if s != nil {
					return s
				}
			case <-timeout:
				return errors.Errorf("Failed to get service definition for %q within timeout", serviceName)
			}
		}
	})
	if err, ok := s.(error); ok {
		return nil, nil, err
	}

	if serviceEvents != nil {
		go c.waitForServiceUpdates(serviceName, serviceEvents)
	}

	service := s.(*Service)

	var uriEvents chan TreeCacheEvent
	uris := c.uris.LoadOrStore(service.ClusterName, func() interface{} {
		Logger.Printf("Creating new URI watcher for %q", service.ClusterName)
		watcher := &serviceUris{
			zkPath: UrisPath(service.ClusterName),
			uris:   make(map[string]*Uri),
		}
		uriEvents = make(chan TreeCacheEvent)

		exists, _, err := c.Conn.Exists(watcher.zkPath)
		if err != nil {
			return err
		}
		if !exists {
			return errors.Errorf("No URIs node found at %q", watcher.zkPath)
		}

		c.caches.Store(watcher.zkPath, NewTreeCache(c.Conn, watcher.zkPath, uriEvents))

		for {
			select {
			case e := <-uriEvents:
				watcher = c.handleUriUpdate(watcher, e)
				// Only return once at least one host is found (respecting the prioritized schemes)
				if watcher.chooseHost(service.PrioritizedSchemes) != nil {
					Logger.Println(watcher)
					return watcher
				}
			case <-timeout:
				return errors.Errorf("Failed to find a valid URI for %q within timeout", service.ClusterName)
			}
		}
	})
	if err, ok := uris.(error); ok {
		return nil, nil, err
	}

	if uriEvents != nil {
		go c.waitForUriUpdates(service.ClusterName, uriEvents)
	}

	return service, uris.(*serviceUris), nil
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

func (c *Client) waitForServiceUpdates(serviceName string, events chan TreeCacheEvent) {
	for e := range events {
		s := c.handleServiceUpdate(serviceName, e)
		if s != nil {
			c.services.Store(serviceName, s)
		}
	}
}

func (c *Client) handleServiceUpdate(serviceName string, event TreeCacheEvent) *Service {
	path := ServicesPath(serviceName)
	if event.Path != path || event.Data == nil {
		return nil
	}

	s := new(Service)
	err := json.Unmarshal(*event.Data, s)
	if err != nil {
		Logger.Printf("Ignoring update to %q (contents: %q) due to error: %v", event.Path, string(*event.Data), err)
	}
	Logger.Printf("Got service definition for %q: %+v", serviceName, s)
	return s
}

func (c *Client) waitForUriUpdates(clusterName string, events chan TreeCacheEvent) {
	for e := range events {
		uri, _ := c.uris.Load(clusterName)
		c.uris.Store(clusterName, c.handleUriUpdate(uri.(*serviceUris), e))
	}
}

func (c *Client) handleUriUpdate(watcher *serviceUris, event TreeCacheEvent) *serviceUris {
	path := strings.TrimPrefix(event.Path, watcher.zkPath)

	if path == "" {
		return watcher
	}

	if event.Data == nil {
		watcher = watcher.copy()
		delete(watcher.uris, path)
		return watcher
	}

	uri := new(Uri)
	err := json.Unmarshal(*event.Data, uri)
	if err != nil {
		Logger.Printf("Ignoring update to %s (contents: %s) due to error: %v", event.Path, string(*event.Data), err)
		return watcher
	}

	if len(uri.Weights) == 0 {
		Logger.Printf("Ignoring partitioned URI at %q: %+v", event.Path, uri)
		return watcher
	}

	Logger.Printf("Received update for %q: %+v", event.Path, uri)

	watcher = watcher.copy()
	watcher.uris[path] = uri
	return watcher
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
