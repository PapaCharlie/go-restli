package d2

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

type watchedUri struct {
	zkPath                string
	name                  string
	uris                  map[string]Uri
	urisWatcher           *TreeCache
	hostWeights           map[url.URL]float64
	totalWeight           float64
	partitionHostWeights  map[int]map[url.URL]float64
	partitionTotalWeights map[int]float64

	lock sync.RWMutex
}

var (
	services  = make(map[string]*Service)
	urisCache = make(map[string]*watchedUri)
	d2Lock    = new(sync.Mutex)
)

func getOrCreateService(serviceName string, conn *zk.Conn) (*watchedUri, error) {
	d2Lock.Lock()
	defer d2Lock.Unlock()

	service, ok := services[serviceName]
	if !ok {
		path := ServicesPath(serviceName)
		data, _, err := conn.Get(path)
		if err != nil {
			err = errors.Wrapf(err, "failed to read %s", path)
			return nil, err
		}

		var s Service
		err = json.Unmarshal(data, &s)
		if err != nil {
			err = errors.Wrapf(err, "could not unmarshal data from %s: %s", path, string(data))
			return nil, err
		}

		services[serviceName] = &s
		service = &s
		Logger.Printf(`Got service definition for "%s": %+v`, serviceName, s)
	} else {
		Logger.Printf(`Using known service mapping from "%s" -> "%s"`, serviceName, service.ClusterName)
	}

	path := UrisPath(service.ClusterName)
	var uri *watchedUri
	if uri = urisCache[service.ClusterName]; uri == nil {
		uri = &watchedUri{
			zkPath: path,
			name:   service.ClusterName,
			uris:   make(map[string]Uri),

			hostWeights: make(map[url.URL]float64),
			totalWeight: 0,

			partitionHostWeights:  make(map[int]map[url.URL]float64),
			partitionTotalWeights: make(map[int]float64),
		}
		events := make(chan TreeCacheEvent)
		uri.urisWatcher = NewTreeCache(conn, path, events)
		uri.start(events)
		Logger.Printf(`Created new URI watcher for "%s"`, service.ClusterName)
		urisCache[service.ClusterName] = uri
	} else {
		Logger.Printf(`Reusing cached URI watcher for "%s"`, service.ClusterName)
	}

	return uri, nil
}

func (u *watchedUri) start(events chan TreeCacheEvent) {
	for e := range events {
		u.handleUpdate(e)
		if len(u.uris) > 0 { // wait for the first host update before returning from start
			break
		}
	}
	go func() {
		for e := range events {
			u.handleUpdate(e)
		}
	}()
}

func (u *watchedUri) addUrl(h url.URL, w float64) {
	if oldW, ok := u.hostWeights[h]; ok {
		u.totalWeight -= oldW
		Logger.Println(h, "NEW_WEIGHT", w)
	} else {
		Logger.Println(h, "UP")
	}
	u.totalWeight += w
	u.hostWeights[h] = w
}

func (u *watchedUri) addUrlToPartition(p int, h url.URL, w float64) {
	partition := u.partitionHostWeights[p]
	if partition == nil {
		u.partitionHostWeights[p] = make(map[url.URL]float64)
	}

	if oldW, ok := partition[h]; ok {
		u.partitionTotalWeights[p] -= oldW
		Logger.Println(h, p, "NEW_WEIGHT", p, w)
	} else {
		Logger.Println(h, p, "UP")
	}
	u.partitionTotalWeights[p] += w
	u.partitionHostWeights[p][h] = w
}

func (u *watchedUri) handleUpdate(event TreeCacheEvent) {
	u.lock.Lock()
	defer u.lock.Unlock()

	path := strings.TrimPrefix(event.Path, u.zkPath)

	if path != "" {
		if event.Data == nil {
			if oldUri, ok := u.uris[path]; ok {
				for h := range oldUri.Weights {
					Logger.Println(h, "DOWN")
					delete(u.hostWeights, h)
				}
				for h, partitions := range oldUri.PartitionDesc {
					for p := range partitions {
						Logger.Println(h, p, "DOWN")
						delete(u.partitionHostWeights[p], h)
					}
				}
				delete(u.uris, path)
			}
			return
		}

		var uri Uri
		err := json.Unmarshal(*event.Data, &uri)
		if err != nil {
			Logger.Printf("Ignoring update to %s (contents: %s) due to error: %v", event.Path, string(*event.Data), err)
			return
		}
		u.uris[path] = uri

		for h, w := range uri.Weights {
			u.addUrl(h, w)
		}

		for h, partitions := range uri.PartitionDesc {
			for p, w := range partitions {
				u.addUrlToPartition(p, h, w)
			}
		}
	}
}

func NewSingleServiceClient(name string, conn *zk.Conn) (c *SingleServiceClient, err error) {
	c = &SingleServiceClient{}
	c.uri, err = getOrCreateService(name, conn)
	return c, err
}

type SingleServiceClient struct {
	uri *watchedUri
}

func (c *SingleServiceClient) GetHostnameForQuery(query string) (*url.URL, error) {
	return c.uri.getHostnameForQuery()
}

func (u *watchedUri) getHostnameForQuery() (*url.URL, error) {
	u.lock.RLock()
	defer u.lock.RUnlock()
	randomWeight := rand.Float64() * u.totalWeight
	for h, w := range u.hostWeights {
		randomWeight -= w
		if randomWeight <= 0 {
			return &h, nil
		}
	}
	return nil, errors.Errorf("Could not find a host for %s", u.name)
}

func NewR2D2Client(conn *zk.Conn) *R2D2Client {
	return &R2D2Client{conn: conn}
}

type R2D2Client struct {
	conn *zk.Conn
}

func (c *R2D2Client) GetHostnameForQuery(query string) (*url.URL, error) {
	idx := strings.Index(query[1:], "/")
	if idx == -1 {
		idx = len(query[1:])
	}
	serviceName := query[1 : idx+1]

	uri, err := getOrCreateService(serviceName, c.conn)
	if err != nil {
		return nil, err
	}
	return uri.getHostnameForQuery()
}
