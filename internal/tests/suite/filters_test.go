package suite

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	actionset "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/actionSet"
	actionsettest "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/actionSet_test"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection"
	collectiontest "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection_test"
	"github.com/PapaCharlie/go-restli/restli"
	"github.com/PapaCharlie/go-restli/restlicodec"
	"github.com/stretchr/testify/require"
)

type (
	filter1   []int
	filter2   []int
	customKey struct{}
)

const (
	ctxValue   = int(5)
	testHeader = "TestHeader"
)

func (f *filter1) PreRequest(req *http.Request) (context.Context, error) {
	*f = append(*f, 1)
	return context.WithValue(req.Context(), customKey{}, ctxValue), nil
}

func (f *filter1) PostRequest(ctx context.Context, responseHeaders http.Header) error {
	*f = append(*f, 1)
	responseHeaders.Add(testHeader, testHeader)
	return nil
}

func (f *filter2) PreRequest(*http.Request) (context.Context, error) {
	*f = append(*f, 2)
	return nil, nil
}

func (f *filter2) PostRequest(context.Context, http.Header) error {
	*f = append(*f, 2)
	return nil
}

func TestContextAndFilters(t *testing.T) {
	filterValues := make([]int, 0, 4)
	server := restli.NewServer((*filter1)(&filterValues), (*filter2)(&filterValues))

	actionset.RegisterResource(server, &actionsettest.MockResource{
		MockReturnBoolAction: func(ctx *restli.RequestContext) (actionResult bool, err error) {
			require.Equal(t, restli.Method_action, restli.GetMethodFromContext(ctx.Request.Context()))
			require.Equal(t, "returnBool", restli.GetActionNameFromContext(ctx.Request.Context()))

			require.Equal(t, ctx.Request.Context().Value(customKey{}).(int), ctxValue)

			return true, nil
		},
	})
	id := int64(42)
	message := &conflictresolution.Message{Message: "test"}
	elements := &collection.FindBySearchElements{
		Elements: []*conflictresolution.Message{message},
		Metadata: new(testsuite.Optionals),
	}
	collection.RegisterResource(server, &collectiontest.MockResource{
		MockGet: func(ctx *restli.RequestContext, collectionId int64) (entity *conflictresolution.Message, err error) {
			require.Equal(t, []restli.ResourcePathSegment{
				restli.NewResourcePathSegment("collection", true),
			}, restli.GetResourcePathSegmentsFromContext(ctx.Request.Context()))

			entitySegments := restli.GetEntitySegmentsFromContext(ctx.Request.Context())
			require.Equal(t, 1, len(entitySegments))
			actualId, err := restlicodec.UnmarshalRestLi[int64](entitySegments[0])
			require.NoError(t, err)
			require.Equal(t, id, actualId)

			require.Equal(t, restli.Method_get, restli.GetMethodFromContext(ctx.Request.Context()))

			return message, nil
		},
		MockFindBySearch: func(ctx *restli.RequestContext, queryParams *collection.FindBySearchParams) (results *collection.FindBySearchElements, err error) {
			require.Equal(t, restli.Method_finder, restli.GetMethodFromContext(ctx.Request.Context()))
			require.Equal(t, "search", restli.GetFinderNameFromContext(ctx.Request.Context()))
			return elements, nil
		},
	})

	c := newClient(server, func(req *http.Request, res *http.Response) {
		require.Equal(t, res.Header.Get(testHeader), testHeader)
	})

	resBool, err := actionset.NewClient(c).ReturnBoolAction()
	require.NoError(t, err)
	require.True(t, resBool)

	resMessage, err := collection.NewClient(c).Get(id)
	require.NoError(t, err)
	require.Equal(t, message, resMessage)

	resElements, err := collection.NewClient(c).FindBySearch(&collection.FindBySearchParams{})
	require.NoError(t, err)
	require.Equal(t, elements, resElements)

	require.Equal(t, []int{1, 2, 2, 1, 1, 2, 2, 1, 1, 2, 2, 1}, filterValues)
}

func TestCustomHeaders(t *testing.T) {
	server := restli.NewServer()

	const (
		h1, v1 = "h1", "v1"
		h2, v2 = "h2", "v2"
	)

	actionset.RegisterResource(server, &actionsettest.MockResource{
		MockReturnBoolAction: func(ctx *restli.RequestContext) (actionResult bool, err error) {
			require.Equal(t, v1, ctx.Request.Header.Get(h1))
			ctx.ResponseHeaders.Set(h2, v2)
			return true, nil
		},
	})

	c := newClient(server, nil)

	ctx := context.Background()
	ctx = restli.ExtraRequestHeaders(ctx, http.Header{http.CanonicalHeaderKey(h1): {v1}})
	ctx, resHeaders := restli.AddResponseHeadersCaptor(ctx)

	_, err := actionset.NewClient(c).ReturnBoolActionWithContext(ctx)
	require.NoError(t, err)
	require.Equal(t, v2, resHeaders.Get(h2))
}

func newClient(server restli.Server, f func(req *http.Request, res *http.Response)) *restli.Client {
	return &restli.Client{
		Client: &http.Client{
			Transport: roundTripper(func(req *http.Request) (*http.Response, error) {
				record := httptest.NewRecorder()
				server.(http.Handler).ServeHTTP(record, req)
				res := record.Result()
				if f != nil {
					f(req, res)
				}
				return res, nil
			}),
		},
		HostnameResolver: &restli.SimpleHostnameResolver{Hostname: &url.URL{}},
	}
}
