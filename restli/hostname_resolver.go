package restli

import "net/url"

type SimpleHostnameResolver struct {
	Hostname *url.URL
}

func (s *SimpleHostnameResolver) ResolveHostnameAndContextForQuery(string, *url.URL) (*url.URL, error) {
	return s.Hostname, nil
}

type HostnameResolver interface {
	// ResolveHostnameAndContextForQuery takes in the name of the resource for which to resolve the hostname, along with
	// the URL for the query that is about to be sent. The root resource will be the top-level parent resource's name.
	// Some HostnameResolver implementations will choose to ignore this parameter and resolve hostnames using a
	// different strategy.
	ResolveHostnameAndContextForQuery(rootResource string, query *url.URL) (*url.URL, error)
}
