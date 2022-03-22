package restli

import "net/url"

type SimpleHostnameResolver struct {
	Hostname *url.URL
}

func (s *SimpleHostnameResolver) ResolveHostnameAndContextForQuery(string, *url.URL) (*url.URL, error) {
	return s.Hostname, nil
}

type HostnameResolver interface {
	// ResolveHostnameAndContextForQuery takes in the name of the service for which to resolve the hostname, along with
	// the URL for the query that is about to be sent. The service name is often the top-level parent resource's name,
	// but can be any unique identifier for a D2 endpoint. Some HostnameResolver implementations will choose to ignore
	// this parameter and resolve hostnames using a different strategy. By default, the generated code will always pass
	// in the top-level parent resource's name.
	ResolveHostnameAndContextForQuery(serviceName string, query *url.URL) (*url.URL, error)
}
