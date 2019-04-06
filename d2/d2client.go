package d2

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"math/rand"
	"net/url"
	"path/filepath"
	"sync"
)

func ServicesPath(service string) string {
	return filepath.Join("/", "d2", "services", service)
}

func UrisPath(uri string) string {
	return filepath.Join("/", "d2", "uris", uri)
}

type Service struct {
	Path        string
	ServiceName string `json:"serviceName"`
	ClusterName string `json:"clusterName"`
}

type UriProperty struct {
	AppName    string `json:"com.linkedin.app.name"`
	AppVersion string `json:"com.linkedin.app.version"`
}

type Uri struct {
	Weights       map[url.URL]float64
	ClusterName   string
	Properties    map[url.URL]UriProperty
	PartitionDesc map[url.URL]map[int]float64
}

func (u *Uri) UnmarshalJSON(data []byte) error {
	uri := &struct {
		Weights       map[string]float64
		ClusterName   string                 `json:"clusterName"`
		Properties    map[string]UriProperty `json:"uriSpecificProperties"`
		PartitionDesc map[string]map[int]struct {
			Weight float64
		} `json:"partitionDesc"`
	}{}

	err := json.Unmarshal(data, uri)
	if err != nil {
		return err
	}

	u.ClusterName = uri.ClusterName
	u.Weights = make(map[url.URL]float64)
	u.Properties = make(map[url.URL]UriProperty)
	u.PartitionDesc = make(map[url.URL]map[int]float64)

	var hostUrl *url.URL
	for host, w := range uri.Weights {
		hostUrl, err = url.Parse(host)
		if err != nil {
			return err
		}
		u.Weights[*hostUrl] = w
	}

	for host, p := range uri.Properties {
		hostUrl, err = url.Parse(host)
		if err != nil {
			return err
		}
		u.Properties[*hostUrl] = p
	}

	for host, desc := range uri.PartitionDesc {
		hostUrl, err = url.Parse(host)
		if err != nil {
			return err
		}
		u.PartitionDesc[*hostUrl] = make(map[int]float64)
		for p, w := range desc {
			u.PartitionDesc[*hostUrl][p] = w.Weight
		}
	}

	return nil
}

type Client struct {
	zk.Conn
	Service               string
	lock                  *sync.RWMutex
	uris                  map[string]Uri
	watcher               *ChildWatcher
	hostWeights           map[url.URL]float64
	totalWeight           float64
	partitionHostWeights  map[int]map[url.URL]float64
	partitionTotalWeights map[int]float64
}

func (c *Client) getWeightedHostname(hostWeights map[url.URL]float64, totalWeight float64) (*url.URL, error) {
	randomWeight := rand.Float64() * totalWeight
	for h, w := range hostWeights {
		randomWeight -= w
		if randomWeight <= 0 {
			return &h, nil
		}
	}
	return nil, errors.Errorf("Could not find a host for %s", c.Service)
}

func (c *Client) GetHostname() (*url.URL, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.getWeightedHostname(c.hostWeights, c.totalWeight)
}

func (c *Client) GetHostnameForPartition(partition int) (*url.URL, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.getWeightedHostname(c.partitionHostWeights[partition], c.partitionTotalWeights[partition])
}

func (c *Client) addUrl(h url.URL, w float64) {
	if oldW, ok := c.hostWeights[h]; ok {
		c.totalWeight -= oldW
		log.Println(h, "NEW_WEIGHT", w)
	} else {
		log.Println(h, "UP")
	}
	c.totalWeight += w
	c.hostWeights[h] = w
}

func (c *Client) addUrlToPartition(p int, h url.URL, w float64) {
	partition := c.partitionHostWeights[p]
	if partition == nil {
		c.partitionHostWeights[p] = make(map[url.URL]float64)
	}

	if oldW, ok := partition[h]; ok {
		c.partitionTotalWeights[p] -= oldW
		log.Println(h, p, "NEW_WEIGHT", p, w)
	} else {
		log.Println(h, p, "UP")
	}
	c.partitionTotalWeights[p] += w
	c.partitionHostWeights[p][h] = w
}

func (c *Client) handleUpdate(child string, data []byte, err error) {
	if err != nil {
		fmt.Println(err)
		log.Panicln(err)
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	if data == nil {
		if oldUri, ok := c.uris[child]; ok {
			for h := range oldUri.Weights {
				log.Println(h, "DOWN")
				delete(c.hostWeights, h)
			}
			for h, partitions := range oldUri.PartitionDesc {
				for p := range partitions {
					log.Println(h, p, "DOWN")
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

func NewClient(name string, conn *zk.Conn) (c *Client, err error) {
	data, _, err := conn.Get(ServicesPath(name))
	if err != nil {
		return
	}

	c = &Client{
		lock: new(sync.RWMutex),
		uris: make(map[string]Uri),

		hostWeights: make(map[url.URL]float64),
		totalWeight: 0,

		partitionHostWeights:  make(map[int]map[url.URL]float64),
		partitionTotalWeights: make(map[int]float64),
	}

	var s Service
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	c.watcher, err = NewChildWatcher(conn, UrisPath(s.ClusterName), c.handleUpdate)

	if err != nil {
		return
	}

	return
}
