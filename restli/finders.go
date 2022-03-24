package restli

import (
	"context"
	"log"
	"net/http"

	"github.com/PapaCharlie/go-restli/restlicodec"
	"github.com/PapaCharlie/go-restli/restlidata"
)

// Find executes a rest.li find request
func Find[V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
) (results *restlidata.Elements[V], err error) {
	req, err := NewGetRequest(c, ctx, rp, query, Method_finder)
	if err != nil {
		return nil, err
	}

	results, _, err = DoAndUnmarshal(c, req, restlicodec.UnmarshalRestLi[*restlidata.Elements[V]])
	return results, nil
}

// FindWithMetadata executes a rest.li find request for finders that declare metadata
func FindWithMetadata[V, M restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
) (results *restlidata.ElementsWithMetadata[V, M], err error) {
	req, err := NewGetRequest(c, ctx, rp, query, Method_finder)
	if err != nil {
		return nil, err
	}

	results, _, err = DoAndUnmarshal(c, req, restlicodec.UnmarshalRestLi[*restlidata.ElementsWithMetadata[V, M]])
	return results, nil
}

func RegisterFinder[RP ResourcePathUnmarshaler[RP], QP restlicodec.QueryParamsDecoder[QP], V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	name string,
	find func(*RequestContext, RP, QP) (*restlidata.Elements[V], error),
) {
	registerFinder(s, segments, name,
		func(ctx *RequestContext, rp RP, qp QP) (restlicodec.Marshaler, error) {
			return find(ctx, rp, qp)
		})
}

func RegisterFinderWithMetadata[RP ResourcePathUnmarshaler[RP], QP restlicodec.QueryParamsDecoder[QP], V, M restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	name string,
	find func(*RequestContext, RP, QP) (*restlidata.ElementsWithMetadata[V, M], error),
) {
	registerFinder(s, segments, name,
		func(ctx *RequestContext, rp RP, p QP) (restlicodec.Marshaler, error) {
			return find(ctx, rp, p)
		})
}

func registerFinder[RP ResourcePathUnmarshaler[RP], QP restlicodec.QueryParamsDecoder[QP]](
	s Server,
	segments []ResourcePathSegment,
	name string,
	h func(*RequestContext, RP, QP) (restlicodec.Marshaler, error),
) {
	p := s.subNode(segments)
	if _, ok := p.finders[name]; ok {
		log.Panicf("go-restli: Cannot register finder %q twice for %v", name, segments)
	}

	p.finders[name] = func(
		ctx *RequestContext,
		segments []restlicodec.Reader,
		body []byte,
	) (responseBody restlicodec.Marshaler, err error) {
		rp, err := UnmarshalResourcePath[RP](segments)
		if err != nil {
			return newErrorResponsef(err, http.StatusBadRequest, "Invalid path for finder %q: %s", name)
		}

		queryParams, err := restlicodec.UnmarshalQueryParamsDecoder[QP](ctx.Request.URL.RawQuery)
		if err != nil {
			return newErrorResponsef(err, http.StatusBadRequest, "Invalid query params for finder %q: %s", name)
		}

		if len(body) != 0 {
			return newErrorResponsef(nil, http.StatusBadRequest, "Finders do not accept request bodies")
		}

		results, err := h(ctx, rp, queryParams)
		if err != nil {
			return newErrorResponsef(err, http.StatusInternalServerError, "Finder %q failed: %s", name)
		}
		return results, nil
	}
}
