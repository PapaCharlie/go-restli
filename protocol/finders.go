package protocol

import (
	"context"
	"log"
	"net/http"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/PapaCharlie/go-restli/protocol/stdtypes"
)

type Elements[V restlicodec.Marshaler] struct {
	Elements []V
	Paging   *stdtypes.CollectionMedata
}

func (f *Elements[V]) marshalElements(keyWriter func(key string) restlicodec.Writer) (err error) {
	return restlicodec.WriteArray(keyWriter(elementsField), f.Elements, V.MarshalRestLi)
}

func (f *Elements[V]) marshalPaging(keyWriter func(key string) restlicodec.Writer) error {
	if f.Paging != nil {
		return f.Paging.MarshalRestLi(keyWriter(pagingField))
	} else {
		return nil
	}
}

func (f *Elements[V]) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		err = f.marshalElements(keyWriter)
		if err != nil {
			return err
		}

		err = f.marshalPaging(keyWriter)
		if err != nil {
			return err
		}
		return nil
	})
}

func (f *Elements[V]) unmarshalRestLi(reader restlicodec.Reader, field string) (err error) {
	switch field {
	case elementsField:
		f.Elements, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[V])
		return err
	case pagingField:
		f.Paging = new(stdtypes.CollectionMedata)
		return f.Paging.UnmarshalRestLi(reader)
	default:
		return reader.Skip()
	}
}

func (f *Elements[V]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(elementsRequiredResponseFields, f.unmarshalRestLi)
}

type ElementsWithMetadata[V, M restlicodec.Marshaler] struct {
	Elements[V]
	Metadata M
}

func (f *ElementsWithMetadata[V, M]) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		err = f.marshalElements(keyWriter)
		if err != nil {
			return err
		}

		err = f.Metadata.MarshalRestLi(keyWriter(metadataField))
		if err != nil {
			return err
		}

		err = f.marshalPaging(keyWriter)
		if err != nil {
			return err
		}

		return nil
	})
}

func (f *ElementsWithMetadata[V, M]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(elementsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == metadataField {
			f.Metadata, err = restlicodec.UnmarshalRestLi[M](reader)
			return err
		} else {
			return f.unmarshalRestLi(reader, field)
		}
	})
}

// Find executes a rest.li find request
func Find[V restlicodec.Marshaler](
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
) (results *Elements[V], err error) {
	req, err := NewGetRequest(c, ctx, rp, query, Method_finder)
	if err != nil {
		return nil, err
	}

	results, _, err = DoAndUnmarshal(c, req, restlicodec.UnmarshalRestLi[*Elements[V]])
	return results, nil
}

// FindWithMetadata executes a rest.li find request for finders that declare metadata
func FindWithMetadata[V, M restlicodec.Marshaler](
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
) (results *ElementsWithMetadata[V, M], err error) {
	req, err := NewGetRequest(c, ctx, rp, query, Method_finder)
	if err != nil {
		return nil, err
	}

	results, _, err = DoAndUnmarshal(c, req, restlicodec.UnmarshalRestLi[*ElementsWithMetadata[V, M]])
	return results, nil
}

func RegisterFinder[RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder, V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	name string,
	find func(*RequestContext, RP, QP) (*Elements[V], error),
) {
	registerFinder(s, segments, name,
		func(ctx *RequestContext, rp RP, qp QP) (restlicodec.Marshaler, error) {
			return find(ctx, rp, qp)
		})
}

func RegisterFinderWithMetadata[RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder, V, M restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	name string,
	find func(*RequestContext, RP, QP) (*ElementsWithMetadata[V, M], error),
) {
	registerFinder(s, segments, name,
		func(ctx *RequestContext, rp RP, p QP) (restlicodec.Marshaler, error) {
			return find(ctx, rp, p)
		})
}

func registerFinder[RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder](
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
