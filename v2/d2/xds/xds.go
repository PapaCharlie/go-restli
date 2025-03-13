package xds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"log"
	"maps"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PapaCharlie/go-restli/v2/d2"
	"github.com/PapaCharlie/go-restli/v2/d2/internal"
	indisapi "github.com/PapaCharlie/go-restli/v2/d2/xds/indis"
	"github.com/PapaCharlie/go-restli/v2/restli"
	"github.com/linkedin/diderot"
	"github.com/linkedin/diderot/ads"
)

var Logger = log.New(d2.Logger.Writer(), "[xDS] ", d2.Logger.Flags())

var _ restli.HostnameResolver = (*Client)(nil)

// The Client wraps a [diderot.ADSClient] and uses it to implement the [restli.HostnameResolver]
// interface by first subscribing to the requested D2 Service, and then the set of URIs within that
// service.
type Client struct {
	adsClient *diderot.ADSClient

	fetchTimeout time.Duration

	lock     sync.Mutex
	services map[string]*pendingResource[*internal.Service]
	uris     map[string]*pendingResource[internal.URIs]
}

type Option interface {
	apply(*Client)
}

type optionFunc func(*Client)

func (o optionFunc) apply(c *Client) {
	o(c)
}

func NewClient(adsClient *diderot.ADSClient, opts ...Option) (c *Client, err error) {
	c = &Client{
		adsClient:    adsClient,
		fetchTimeout: 10 * time.Second,
		services:     make(map[string]*pendingResource[*internal.Service]),
		uris:         make(map[string]*pendingResource[internal.URIs]),
	}

	for _, o := range opts {
		o.apply(c)
	}

	return c, nil
}

// WithFetchTimeout updates the default fetch timeout of 10 seconds to the given duration. If the xDS
// backend does not respond to the subscription within this timeout, an error will be returned by
// [Client.ResolveHostnameAndContextForQuery].
func WithFetchTimeout(timeout time.Duration) Option {
	return optionFunc(func(client *Client) {
		client.fetchTimeout = timeout
	})
}

// pendingResource is used to store a resource that will eventually be set. When signal is closed,
// the resource has been received from the xDS backend.
type pendingResource[T any] struct {
	signal   chan struct{}
	notify   func()
	resource atomic.Pointer[T]
}

func newPendingResource[T any]() *pendingResource[T] {
	ch := make(chan struct{})
	return &pendingResource[T]{
		signal: ch,
		notify: sync.OnceFunc(func() {
			close(ch)
		}),
	}
}

// await blocks and returns only once the resource is set by setAndNotify, or when context is
// canceled in which case it returns an error.
func (r *pendingResource[T]) await(ctx context.Context) (T, error) {
	select {
	case <-r.signal:
		return *r.resource.Load(), nil
	case <-ctx.Done():
		var t T
		return t, ctx.Err()
	}
}

// setAndNotify sets the resource and notifies any goroutines waiting on it.
func (r *pendingResource[T]) setAndNotify(t *T) {
	r.resource.Swap(t)
	r.notify()
}

// serviceWatcher implements the [diderot.Watcher] interface and waits for updates to the
// [indisapi.Node] that contains the definition of the D2 Service.
type serviceWatcher struct {
	serviceName string
	c           *Client
}

func (w *serviceWatcher) Notify(resources iter.Seq2[string, *ads.Resource[*indisapi.Node]]) error {
	// There should only be one element in the seq, and it will either be the resource itself, or nil to
	// indicate that it is deleted.
	var resource *ads.Resource[*indisapi.Node]
	for _, resource = range resources {
		break
	}

	var service *internal.Service
	if resource != nil {
		service = new(internal.Service)
		err := json.Unmarshal(resource.Resource.Data, service)
		if err != nil {
			return err
		}
	}

	w.c.lock.Lock()
	defer w.c.lock.Unlock()

	w.c.services[w.serviceName].setAndNotify(&service)

	return nil
}

// clusterWatcher implements the [diderot.Watcher] interface and waits for updates to the glob
// collection of [indisapi.D2URI] that represents the set of URIs in the D2 Cluster.
type clusterWatcher struct {
	clusterName string
	c           *Client
}

func (w *clusterWatcher) Notify(resources iter.Seq2[string, *ads.Resource[*indisapi.D2URI]]) error {
	update := make(internal.URIs)

	var errs []error
	for name, resource := range resources {
		if resource != nil {
			if _, ok := resource.Resource.PartitionDesc[0]; !ok {
				Logger.Printf("Ignoring partitioned URI: %+v", resource.Resource)
				continue
			}

			uri, err := url.Parse(resource.Resource.Uri)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid URI %q: %w", resource.Resource.Uri, err))
				continue
			}

			update[name] = &internal.URI{
				URI:           uri,
				PartitionDesc: resource.Resource.PartitionDesc,
			}
		} else {
			update[name] = nil
		}
	}

	w.c.lock.Lock()
	defer w.c.lock.Unlock()

	r := w.c.uris[w.clusterName]
	urisRef := r.resource.Load()
	var uris internal.URIs
	if urisRef != nil {
		uris = maps.Clone(*urisRef)
	} else {
		uris = make(internal.URIs)
	}

	for name, uri := range update {
		if uri == nil {
			delete(uris, name)
		} else {
			uris[name] = uri
		}
	}

	r.setAndNotify(&uris)

	return errors.Join(errs...)
}

// getPendingResource gets or creates a [pendingResource] for the given type and resource name. The
// returned boolean is only true when this function *creates* the [pendingResource], otherwise every
// subsequent invocation it is false.
func getPendingResource[K comparable, T any](
	c *Client,
	m map[K]*pendingResource[T],
	k K,
) (r *pendingResource[T], created bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	r, ok := m[k]
	if !ok {
		r = newPendingResource[T]()
		m[k] = r
	}

	return r, !ok
}

func (c *Client) getService(serviceName string) (s *pendingResource[*internal.Service], created bool) {
	return getPendingResource(c, c.services, serviceName)
}

func (c *Client) getURIs(clusterName string) (u *pendingResource[internal.URIs], created bool) {
	return getPendingResource(c, c.uris, clusterName)
}

// getServiceAndURIs subscribes to the given service name, then subscribes to the returned cluster's
// glob collection of URIs.
func (c *Client) getServiceAndURIs(serviceName string) (*internal.Service, internal.URIs, error) {
	ctx := context.Background()
	if c.fetchTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.fetchTimeout)
		defer cancel()
	}

	pendingService, created := c.getService(serviceName)
	if created {
		// If created is true, then this is the first time the service is being queried, so add the watch on
		// the ADSClient.
		Logger.Printf("Subscribing to service %q", serviceName)
		diderot.Watch(
			c.adsClient,
			d2.ServicesPath(serviceName),
			&serviceWatcher{serviceName: serviceName, c: c},
		)
	}

	service, err := pendingService.await(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting service %q: %w", serviceName, err)
	}

	if service == nil {
		return nil, nil, fmt.Errorf("service %q does not exist", serviceName)
	}

	pendingURIs, created := c.getURIs(service.ClusterName)
	if created {
		// Same as the service, this is the first time the cluster is being queried.
		Logger.Printf("Subscribing to cluster %q", service.ClusterName)
		diderot.Watch(
			c.adsClient,
			ads.NewGlobCollectionURL[*indisapi.D2URI]("", service.ClusterName, nil).String(),
			&clusterWatcher{c: c, clusterName: service.ClusterName},
		)
	}

	uris, err := pendingURIs.await(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting cluster URIs for %q: %w", service.ClusterName, err)
	}

	if len(uris) == 0 {
		return nil, nil, fmt.Errorf("cluster %q does not exist", service.ClusterName)
	}

	return service, uris, nil
}

func (c *Client) ResolveHostnameAndContextForQuery(rootResource string, _ *url.URL) (*url.URL, error) {
	service, uris, err := c.getServiceAndURIs(rootResource)
	if err != nil {
		return nil, err
	}

	chosenHost := uris.ChooseHost(service.PrioritizedSchemes)
	if chosenHost == nil {
		return nil, fmt.Errorf("could not find a host for service %q", rootResource)
	}
	return chosenHost, nil
}

type SingleServiceClient struct {
	c           *Client
	serviceName string
}

func (s *SingleServiceClient) ResolveHostnameAndContextForQuery(_ string, query *url.URL) (*url.URL, error) {
	return s.c.ResolveHostnameAndContextForQuery(s.serviceName, query)
}
