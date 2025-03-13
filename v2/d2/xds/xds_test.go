package xds

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/PapaCharlie/go-restli/v2/d2"
	"github.com/PapaCharlie/go-restli/v2/d2/internal"
	indisapi "github.com/PapaCharlie/go-restli/v2/d2/xds/indis"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/linkedin/diderot"
	"github.com/linkedin/diderot/ads"
	"github.com/linkedin/diderot/testutils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestClient(t *testing.T) {
	testServer := testutils.NewTestGRPCServer(t)

	sl := &simpleLocator{
		t:         t,
		nodeCache: diderot.NewCache[*indisapi.Node](),
		urisCache: diderot.NewCache[*indisapi.D2URI](),
	}
	discoveryv3.RegisterAggregatedDiscoveryServiceServer(testServer.Server, diderot.NewADSServer(sl))
	testServer.Start()

	expectedService := &internal.Service{
		Path:               "/foo",
		ServiceName:        "fooService",
		ClusterName:        "fooCluster",
		PrioritizedSchemes: []string{"https"},
	}
	serviceData := must(t, expectedService, func(in *internal.Service) ([]byte, error) { return json.Marshal(in) })

	sl.nodeCache.Set(d2.ServicesPath("fooService"), "0", &indisapi.Node{Data: serviceData}, time.Time{})

	uri1 := &internal.URI{
		URI:           must(t, "https://host1:8080/foo", url.Parse),
		PartitionDesc: map[int32]float64{0: 1},
	}
	uri2 := &internal.URI{
		URI:           must(t, "https://host2:8080/foo", url.Parse),
		PartitionDesc: map[int32]float64{1: 1},
	}

	uri1Resource := toD2URI(expectedService.ClusterName, uri1)
	sl.urisCache.SetResource(uri1Resource, time.Time{})
	sl.urisCache.SetResource(toD2URI(expectedService.ClusterName, uri2), time.Time{})

	c, err := NewClient(diderot.NewADSClient(testServer.Dial(), new(ads.Node)))
	require.NoError(t, err)

	service, uris, err := c.getServiceAndURIs("fooService")
	require.NoError(t, err)
	require.Equal(t, expectedService, service)
	require.Equal(t, internal.URIs{uri1Resource.Name: uri1}, uris)
}

func toD2URI(clusterName string, u *internal.URI) *ads.Resource[*indisapi.D2URI] {
	escaped := url.QueryEscape(u.URI.String())
	return ads.NewResource(
		ads.NewGlobCollectionURL[*indisapi.D2URI]("", clusterName, nil).MemberURN(escaped),
		"0",
		&indisapi.D2URI{
			Version:       0,
			ModifiedTime:  timestamppb.Now(),
			ClusterName:   clusterName,
			Uri:           u.URI.String(),
			PartitionDesc: u.PartitionDesc,
		},
	)
}

func must[IN, OUT any](t *testing.T, in IN, f func(IN) (OUT, error)) OUT {
	out, err := f(in)
	require.NoError(t, err)
	return out
}

type simpleLocator struct {
	t         *testing.T
	nodeCache diderot.Cache[*indisapi.Node]
	urisCache diderot.Cache[*indisapi.D2URI]
}

func (s *simpleLocator) Subscribe(
	_ context.Context,
	typeURL, resourceName string,
	handler ads.RawSubscriptionHandler,
) (unsubscribe func()) {
	var c diderot.RawCache
	switch typeURL {
	case s.nodeCache.Type().URL():
		c = s.nodeCache
	case s.urisCache.Type().URL():
		c = s.urisCache
	default:
		s.t.Errorf("Unknown type URL: %q", typeURL)
		return func() {}
	}
	diderot.Subscribe(c, resourceName, handler)
	return func() {
		diderot.Unsubscribe(c, resourceName, handler)
	}
}
