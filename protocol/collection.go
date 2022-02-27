package protocol

import (
	"context"
	"net/http"

	"github.com/PapaCharlie/go-restli/protocol/batchkeyset"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/PapaCharlie/go-restli/protocol/stdtypes"
)

type CollectionClient[K comparable, V, PartialV restlicodec.Marshaler] struct {
	SimpleClient[V, PartialV]
	KeyUnmarshaler      restlicodec.GenericUnmarshaler[K]
	BatchKeySetProvider func() batchkeyset.BatchKeySet[K]

	ReadOnlyFields   restlicodec.PathSpec
	CreateOnlyFields restlicodec.PathSpec
}

func (c *CollectionClient[K, V, PV]) NewCreateRequest(
	ctx context.Context,
	rp ResourcePath,
	query QueryParams,
	method RestLiMethod,
	create restlicodec.Marshaler,
) (*http.Request, error) {
	return c.NewJsonRequest(ctx, rp, query, http.MethodPost, method, create, c.ReadOnlyFields)
}

// Create executes a rest.li create request with the given object. The X-RestLi-Id header field will be parsed into id
// (though a CreateResponseHasNoEntityHeaderError will be returned if the header is not set). The response body will
// always be ignored.
func (c *CollectionClient[K, V, PV]) Create(
	ctx context.Context,
	rp ResourcePath,
	create V,
	query QueryParams,
) (*CreatedEntity[K], error) {
	req, err := c.NewCreateRequest(ctx, rp, query, Method_create, create)
	if err != nil {
		return nil, err
	}

	res, err := c.RestLiClient.DoAndIgnore(req)
	if err != nil {
		return nil, err
	}

	return c.unmarshalReturnEntityKey(res)
}

// CreateWithReturnEntity is like CollectionClient.Create, except it parses the returned entity from the response.
func (c *CollectionClient[K, V, PV]) CreateWithReturnEntity(
	ctx context.Context,
	rp ResourcePath,
	create V,
	query QueryParams,
) (*CreatedAndReturnedEntity[K, V], error) {
	req, err := c.NewCreateRequest(ctx, rp, query, Method_create, create)
	if err != nil {
		return nil, err
	}

	v, res, err := DoAndUnmarshal(c.RestLiClient, req, c.EntityUnmarshaler)
	if err != nil {
		return nil, err
	}

	k, err := c.unmarshalReturnEntityKey(res)
	if err != nil {
		return nil, err
	}
	return &CreatedAndReturnedEntity[K, V]{
		CreatedEntity: *k,
		Entity:        v,
	}, nil
}

func (c *CollectionClient[K, V, PV]) unmarshalReturnEntityKey(res *http.Response) (result *CreatedEntity[K], err error) {
	if h := res.Header.Get(RestLiHeader_ID); len(h) > 0 {
		var reader restlicodec.Reader
		reader, err = restlicodec.NewRor2Reader(h)
		if err != nil {
			return nil, err
		}

		var k K
		k, err = c.KeyUnmarshaler(reader)
		if _, mfe := err.(*restlicodec.MissingRequiredFieldsError); mfe && !c.RestLiClient.StrictResponseDeserialization {
			err = nil
		}
		if err != nil {
			return nil, err
		}
		return &CreatedEntity[K]{
			Id:     k,
			Status: res.StatusCode,
		}, nil
	} else {
		return nil, &CreateResponseHasNoEntityHeaderError{Response: res}
	}
}

// Find executes a rest.li find request
func (c *CollectionClient[K, V, PV]) Find(
	ctx context.Context,
	rp ResourcePath,
	query QueryParams,
) (*FinderResults[V], error) {
	results, err := FindWithMetadata[K, V, PV, stdtypes.EmptyRecord](c, ctx, rp, query, nil)
	if err != nil {
		return nil, err
	}
	return &results.FinderResults, nil
}

// FindWithMetadata executes a rest.li find request for finders that declare metadata
func FindWithMetadata[K comparable, V, PV restlicodec.Marshaler, M any](
	c *CollectionClient[K, V, PV],
	ctx context.Context,
	rp ResourcePath,
	query QueryParams,
	metadataUnmarshaler restlicodec.GenericUnmarshaler[M],
) (results *FinderResultsWithMetadata[V, M], err error) {
	u, err := c.FormatQueryUrl(rp, query)
	if err != nil {
		return nil, err
	}

	req, err := NewGetRequest(ctx, u, Method_finder)
	if err != nil {
		return nil, err
	}

	results, _, err = DoAndUnmarshal(c.RestLiClient, req, func(reader restlicodec.Reader) (results *FinderResultsWithMetadata[V, M], err error) {
		results = new(FinderResultsWithMetadata[V, M])
		err = reader.ReadRecord(elementsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
			switch field {
			case elementsField:
				results.Results, err = restlicodec.ReadArray(reader, c.EntityUnmarshaler)
				return err
			case metadataField:
				if metadataUnmarshaler != nil {
					results.Metadata, err = metadataUnmarshaler(reader)
					return err
				} else {
					return reader.Skip()
				}
			case pagingField:
				return reader.ReadMap(func(reader restlicodec.Reader, key string) (err error) {
					if key == totalField {
						var t int
						t, err = reader.ReadInt()
						if err != nil {
							return err
						}
						results.Total = &t
					} else {
						err = reader.Skip()
					}
					return nil
				})
			default:
				return reader.Skip()
			}
		})
		return results, err
	})

	return results, err
}

type FinderResults[T any] struct {
	Results []T
	Total   *int
}

type FinderResultsWithMetadata[T, M any] struct {
	FinderResults[T]
	Metadata M
}
