package internal

import (
	"math/rand/v2"
	"net/url"
)

type URIs map[string]*URI

func (uris URIs) IterateHostWeights(receiver func(host *url.URL, weight float64) bool) {
	for _, uri := range uris {
		_, weight := getFirstEntry(uri.PartitionDesc)
		if !receiver(uri.URI, weight) {
			return
		}
	}
}

func (uris URIs) FilterAndChooseHost(hostFilter func(*url.URL) bool) *url.URL {
	var totalWeight float64
	uris.IterateHostWeights(func(host *url.URL, weight float64) bool {
		if hostFilter(host) {
			totalWeight += weight
		}
		return true
	})

	randomWeight := rand.Float64() * totalWeight

	var chosenHost *url.URL
	uris.IterateHostWeights(func(host *url.URL, weight float64) bool {
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

func (uris URIs) ChooseHost(prioritizedSchemes []string) *url.URL {
	if len(prioritizedSchemes) == 0 {
		return uris.FilterAndChooseHost(func(*url.URL) bool { return true })
	}

	for _, scheme := range prioritizedSchemes {
		chosenHost := uris.FilterAndChooseHost(func(u *url.URL) bool {
			return u.Scheme == scheme
		})
		if chosenHost != nil {
			return chosenHost
		}
	}

	return nil
}
