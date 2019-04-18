package protocol

import (
	"net/url"
	"testing"
)

var (
	seasBroker       = mustParse("/seas-broker")
	seasBrokerSearch = mustParse("/seas-broker/search")
	emptyContext     = mustParse("")
	slashContext     = mustParse("/")
)

const query = "/search?action=search"

func mustParse(u string) *url.URL {
	ur, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return ur
}

func TestRestLiClient_FormatQuery(t *testing.T) {
	supplier := &SimpleHostnameSupplier{}
	c := &RestLiClient{
		Client:           nil,
		HostnameResolver: supplier,
	}

	expected := "/seas-broker/search?action=search"

	supplier.Hostname = seasBroker
	if expected != c.formatQuery(t, query) {
		t.Errorf("Host: %s, Expected: %s, Got: %s", supplier.Hostname, expected, c.formatQuery(t, query))
	}

	supplier.Hostname = seasBrokerSearch
	if expected != c.formatQuery(t, query) {
		t.Errorf("Host: %s, Expected: %s, Got: %s", supplier.Hostname, expected, c.formatQuery(t, query))
	}

	expected = "/search?action=search"

	supplier.Hostname = emptyContext
	if expected != c.formatQuery(t, query) {
		t.Errorf("Host: %s, Expected: %s, Got: %s", supplier.Hostname, expected, c.formatQuery(t, query))
	}

	supplier.Hostname = slashContext
	if expected != c.formatQuery(t, query) {
		t.Errorf("Host: %s, Expected: %s, Got: %s", supplier.Hostname, expected, c.formatQuery(t, query))
	}
}

func (c *RestLiClient) formatQuery(t *testing.T, query string) string {
	u, err := c.FormatQueryUrl(query)
	if err != nil {
		t.Fatal(err)
	}
	return u.String()
}
