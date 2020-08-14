package tests

import (
	"testing"
	"time"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite"
	. "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/params"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) ParamsGetWithQueryparams(t *testing.T, c Client) {
	long := int64(100)
	apple := conflictresolution.Fruits_APPLE
	params := &GetParams{
		Int:         3,
		String:      "string",
		Long:        9223372036854775807,
		StringArray: []string{"string one", "string two"},
		MessageArray: []*conflictresolution.Message{
			{
				Message: "test message",
			},
			{
				Message: "another message",
			},
		},
		StringMap: map[string]string{
			"one": "string one",
			"two": "string two",
		},
		PrimitiveUnion: testsuite.UnionOfPrimitives{
			PrimitivesUnion: testsuite.UnionOfPrimitives_PrimitivesUnion{
				Long: &long,
			},
		},
		ComplexTypesUnion: testsuite.UnionOfComplexTypes{
			ComplexTypeUnion: testsuite.UnionOfComplexTypes_ComplexTypeUnion{
				Fruits: &apple,
			},
		},
		UrlTyperef: "http://rest.li",
	}

	// Because the map iteration order is undetermined, it's possible the encoding will not match the expected order and
	// fail the test. To mitigate this, just retry the query a few times until it succeeds. If it doesn't succeed within
	// 100ms, fail the test with the most recent error
	complete := make(chan error)
	go func() {
		for {
			_, err := c.Get(100, params)
			if err != nil {
				complete <- err
			} else {
				close(complete)
				return
			}
		}
	}()

	var err error
	var channelOpen bool
	timeout := time.After(100 * time.Millisecond)
	for {
		select {
		case err, channelOpen = <-complete:
			if !channelOpen {
				return
			}
		case <-timeout:
			require.FailNowf(t, "Failed to encode GetParams", "Most recent error: %+v", err)
		}
	}
}
