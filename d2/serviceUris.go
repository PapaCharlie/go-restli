package d2

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/go-zookeeper/zk"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type serviceUris map[string]zk.NodeData[*Uri]

func (uris serviceUris) iterateHostWeights(receiver func(host *url.URL, weight float64) bool) {
	for _, uri := range uris {
		for host, weight := range uri.T.Weights {
			if !receiver(&host, weight) {
				return
			}
		}
	}
}

func (uris serviceUris) filterAndChooseHost(hostFilter func(*url.URL) bool) *url.URL {
	var totalWeight float64
	uris.iterateHostWeights(func(host *url.URL, weight float64) bool {
		if hostFilter(host) {
			totalWeight += weight
		}
		return true
	})

	randomWeight := rng.Float64() * totalWeight

	var chosenHost *url.URL
	uris.iterateHostWeights(func(host *url.URL, weight float64) bool {
		if hostFilter(host) {
			randomWeight -= weight
			if randomWeight <= 0 {
				chosenHost = host
				return false
			}
		}
		return true
	})

	return chosenHost
}

func (uris serviceUris) chooseHost(prioritizedSchemes []string) *url.URL {
	if len(prioritizedSchemes) == 0 {
		return uris.filterAndChooseHost(func(*url.URL) bool { return true })
	}

	for _, scheme := range prioritizedSchemes {
		chosenHost := uris.filterAndChooseHost(func(u *url.URL) bool {
			return u.Scheme == scheme
		})
		if chosenHost != nil {
			return chosenHost
		}
	}

	return nil
}
