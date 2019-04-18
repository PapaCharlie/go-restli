package d2

import (
	"math"
	"net/url"
	"sync"
	"testing"
)

func TestClient_addUri(t *testing.T) {
	c := &Client{
		lock:        new(sync.RWMutex),
		hostWeights: make(map[url.URL]float64),
	}
	addUri(c, "a", 1.0)
	addUri(c, "b", 5.0)
	if c.totalWeight != 6.0 {
		t.Fatalf("totalWeight was %g instead of 6", c.totalWeight)
	}

	addUri(c, "a", 5.0)
	if c.totalWeight != 10.0 {
		t.Fatalf("totalWeight was %g instead of 10", c.totalWeight)
	}
}

func TestClient_GetHostname(t *testing.T) {
	c := &Client{
		lock:        new(sync.RWMutex),
		hostWeights: make(map[url.URL]float64),
	}
	a := addUri(c, "a", 1.0)
	b := addUri(c, "b", 3.0)

	hits := make(map[url.URL]int)
	for i := 0; i < 10000000; i++ {
		h, err := c.GetHostname()
		if err != nil {
			t.Fatal(err)
		}
		hits[*h] += 1
	}

	ratio := math.Round(float64(hits[b])/float64(hits[a])*100) / 100
	if ratio != 3.0 {
		t.Fatalf("ratio of hits to a vs b was not 3 (was %g)", ratio)
	}
}

func addUri(c *Client, u string, w float64) url.URL {
	h, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	c.addUrl(*h, w)
	return *h
}
