package internal

import (
	"encoding/json"
	"fmt"
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
	Path                        string
	ServiceName                 string   `json:"serviceName"`
	ClusterName                 string   `json:"clusterName"`
	PrioritizedSchemes          []string `json:"prioritizedSchemes"`
	SslSessionValidationStrings []string `json:"sslSessionValidationStrings"`
}

type UriProperty struct {
	AppName    string `json:"com.linkedin.app.name"`
	AppVersion string `json:"com.linkedin.app.version"`
}

type URI struct {
	URI           *url.URL
	PartitionDesc map[int32]float64
}

func (u *URI) UnmarshalJSON(data []byte) error {
	raw := &struct {
		Weights       map[*url.URL]float64 `json:"weights"`
		PartitionDesc map[*url.URL]map[int32]struct {
			Weight float64 `json:"weight"`
		} `json:"partitionDesc"`
	}{}

	err := json.Unmarshal(data, raw)
	if err != nil {
		return err
	}

	if len(raw.PartitionDesc) > 0 {
		uri, partitionDesc := getFirstEntry(raw.PartitionDesc)
		u.URI = uri
		u.PartitionDesc = make(map[int32]float64, len(partitionDesc))
		for partition, weight := range partitionDesc {
			u.PartitionDesc[partition] = weight.Weight
		}
	} else if len(raw.Weights) > 0 {
		uri, weight := getFirstEntry(raw.Weights)
		u.URI = uri
		u.PartitionDesc = map[int32]float64{0: weight}
	} else {
		return fmt.Errorf("no weights or partitionDesc in URI: %s", string(data))
	}

	return nil
}

// getFirstEntry returns the first entry in the given map. If the map is empty, it returns the
// 0-value for the key and value. Primarily used to unwrap the various maps present in a ZK D2 URI,
// which are always represented as a singleton map where the key is the URI itself, and the value
// contains the relevant information for that field (see the documentation on indisapi.D2URI for an
// example of a full ZK URI)
func getFirstEntry[K comparable, V any](m map[K]V) (k K, v V) {
	for k, v = range m {
		break
	}
	return k, v
}
