package internal

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalURI(t *testing.T) {
	raw := `{
		"weights": {
			"https://foo.stg.linkedin.com:18792/Toki/resources": 1.0
		},
		"partitionDesc": {
			"https://foo.stg.linkedin.com:18792/Toki/resources": {
				"0": {
					"weight": 1.0
				}
				"1": {
					"weight": 2.0
				}
			}
		},
		"uriSpecificProperties": {
			"https://foo.stg.linkedin.com:18792/Toki/resources": {
				"com.linkedin.app.version": "0.1.76"
			}
		},
		"clusterName": "Toki"
	}`
	var actual URI
	require.Error(t, json.Unmarshal([]byte(raw), &actual))

	u, err := url.Parse("https://foo.stg.linkedin.com:18792/Toki/resources")
	require.NoError(t, err)

	require.Equal(t, URI{
		URI:           u,
		PartitionDesc: map[int32]float64{0: 1, 1: 2},
	}, actual)
}
