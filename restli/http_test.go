package restli

import (
	"net/url"
	"testing"
)

func mustParse(u string) *url.URL {
	ur, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return ur
}

func TestRestLiClient_FormatQuery(t *testing.T) {
	supplier := &SimpleHostnameResolver{}
	c := &Client{
		Client:           nil,
		HostnameResolver: supplier,
	}

	tests := []struct {
		Name     string
		Expected string
		Values   []*url.URL
	}{
		{
			Name:     "One segment",
			Expected: "/seas-broker/search?action=search",
			Values: []*url.URL{
				mustParse("/seas-broker"),
				mustParse("/seas-broker/search"),
			},
		},
		{
			Name:     "Multiple segments",
			Expected: "/api/v2/search?action=search",
			Values: []*url.URL{
				mustParse("/api/v2"),
				mustParse("/api/v2/"),
				mustParse("/api/v2/search"),
			},
		},
		{
			Name:     "Simple contexts",
			Expected: "/search?action=search",
			Values: []*url.URL{
				mustParse(""),
				mustParse("/"),
			},
		},
		{
			Name:     "Segment with partial match of basename",
			Expected: "/seas-searcher/search?action=search",
			Values: []*url.URL{
				mustParse("/seas-searcher"),
			},
		},
	}

	const query = QueryParamsString("action=search")
	const rp = ResourcePathString("/search")

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			for _, v := range test.Values {
				supplier.Hostname = v

				u, err := c.FormatQueryUrl(rp, query)
				if err != nil {
					t.Fatal(err)
				}
				q := u.String()

				if test.Expected != q {
					t.Errorf("Host: %s, Expected: %s, Got: %s", supplier.Hostname, test.Expected, q)
				}
			}
		})
	}
}
