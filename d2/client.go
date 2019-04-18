package d2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

var Logger = log.New(ioutil.Discard, "[D2] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.LUTC)

func EnableD2Logging() {
	Logger.SetOutput(os.Stderr)
}

func ClustersPath(cluster string) string {
	return filepath.Join("/", "d2", "clusters", cluster)
}

func ServicesPath(service string) string {
	return filepath.Join("/", "d2", "services", service)
}

func UrisPath(uri string) string {
	return filepath.Join("/", "d2", "uris", uri)
}

type FixedClient struct {
	*zk.Conn
	service *watchedService
}

func (c *FixedClient) GetHostnameForQuery(query string) (*url.URL, error) {
	return c.service.GetHostnameForQuery(query)
}

type watchedService struct {
	name                  string
	uris                  map[string]Uri
	urisWatcher           *ChildWatcher
	hostWeights           map[url.URL]float64
	totalWeight           float64
	partitionHostWeights  map[int]map[url.URL]float64
	partitionTotalWeights map[int]float64

	lock sync.RWMutex
}

func newService(name string, conn *zk.Conn) (*watchedService, error) {
	w := &watchedService{
		name: name,
		uris: make(map[string]Uri),

		hostWeights: make(map[url.URL]float64),
		totalWeight: 0,

		partitionHostWeights:  make(map[int]map[url.URL]float64),
		partitionTotalWeights: make(map[int]float64),
	}

	path := ServicesPath(name)
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

	path = UrisPath(s.ClusterName)
	w.urisWatcher, err = NewChildWatcher(conn, path, w.handleUpdate)

	if err != nil {
		err = errors.Wrapf(err, "failed to read", path)
		return nil, err
	}

	return w, nil
}

func (c *watchedService) getWeightedHostname(hostWeights map[url.URL]float64, totalWeight float64) (*url.URL, error) {
	randomWeight := rand.Float64() * totalWeight
	for h, w := range hostWeights {
		randomWeight -= w
		if randomWeight <= 0 {
			return &h, nil
		}
	}
	return nil, errors.Errorf("Could not find a host for %s", c.name)
}

func (c *watchedService) GetHostnameForQuery(query string) (*url.URL, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.getWeightedHostname(c.hostWeights, c.totalWeight)
}

func (c *watchedService) addUrl(h url.URL, w float64) {
	if oldW, ok := c.hostWeights[h]; ok {
		c.totalWeight -= oldW
		Logger.Println(h, "NEW_WEIGHT", w)
	} else {
		Logger.Println(h, "UP")
	}
	c.totalWeight += w
	c.hostWeights[h] = w
}

func (c *watchedService) addUrlToPartition(p int, h url.URL, w float64) {
	partition := c.partitionHostWeights[p]
	if partition == nil {
		c.partitionHostWeights[p] = make(map[url.URL]float64)
	}

	if oldW, ok := partition[h]; ok {
		c.partitionTotalWeights[p] -= oldW
		Logger.Println(h, p, "NEW_WEIGHT", p, w)
	} else {
		Logger.Println(h, p, "UP")
	}
	c.partitionTotalWeights[p] += w
	c.partitionHostWeights[p][h] = w
}

func (c *watchedService) handleUpdate(child string, data []byte, err error) {
	if err != nil {
		fmt.Println(err)
		log.Panicln(err)
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	if data == nil {
		if oldUri, ok := c.uris[child]; ok {
			for h := range oldUri.Weights {
				Logger.Println(h, "DOWN")
				delete(c.hostWeights, h)
			}
			for h, partitions := range oldUri.PartitionDesc {
				for p := range partitions {
					Logger.Println(h, p, "DOWN")
					delete(c.partitionHostWeights[p], h)
				}
			}
			delete(c.uris, child)
		}
		return
	}

	var uri Uri
	err = json.Unmarshal(data, &uri)
	if err != nil {
		log.Panicln(child, string(data), err)
	}
	c.uris[child] = uri

	for h, w := range uri.Weights {
		c.addUrl(h, w)
	}

	for h, partitions := range uri.PartitionDesc {
		for p, w := range partitions {
			c.addUrlToPartition(p, h, w)
		}
	}
}

func NewFixedClient(name string, conn *zk.Conn) (c *FixedClient, err error) {
	c = &FixedClient{Conn: conn}
	c.service, err = newService(name, conn)
	return c, err
}
