package d2

import (
	"encoding/json"
	"net/url"
)

type Cluster struct {
	ClusterName         string `json:"clusterName"`
	PartitionProperties struct {
		HashAlgorithm     string `json:"hashAlgorithm"`
		PartitionCount    int    `json:"partitionCount"`
		PartitionKeyRegex string `json:"partitionKeyRegex"`
		PartitionType     string `json:"partitionType"`
	} `json:"partitionProperties"`
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
