package restli

import (
	"context"
	"net/http"
	"reflect"

	"github.com/PapaCharlie/go-restli/restli/batchkeyset"
	"github.com/PapaCharlie/go-restli/restlicodec"
	"github.com/PapaCharlie/go-restli/restlidata"
)

type queryParamsFunc func() (string, error)

func (q queryParamsFunc) EncodeQueryParams() (string, error) {
	return q()
}

type batchQueryParamsEncoder[T any] interface {
	EncodeQueryParams(set batchkeyset.BatchKeySet[T]) (string, error)
}

type batchQueryParamsDecoder[T any] interface {
	DecodeQueryParams(reader restlicodec.QueryParamsReader) ([]T, error)
}

type SliceBatchQueryParams[T any] struct{}

var entityIdsRequiredResponseFields = restlicodec.RequiredFields{batchkeyset.EntityIDsField}

func (s *SliceBatchQueryParams[T]) DecodeQueryParams(reader restlicodec.QueryParamsReader) (ids []T, err error) {
	err = reader.ReadRecord(entityIdsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == batchkeyset.EntityIDsField {
			ids, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[T])
		} else {
			err = reader.Skip()
		}
		return err
	})
	return ids, err
}

type WrappedBatchQueryParamsDecoder[T any, QP batchQueryParamsDecoder[T]] struct {
	keys []T
	qp   QP
}

func (b *WrappedBatchQueryParamsDecoder[T, QP]) DecodeQueryParams(reader restlicodec.QueryParamsReader) (err error) {
	v := reflect.New(reflect.TypeOf(b.qp).Elem())
	b.qp = v.Interface().(QP)
	b.keys, err = b.qp.DecodeQueryParams(reader)
	return err
}

func batchQueryParams[T any](set batchkeyset.BatchKeySet[T], query batchQueryParamsEncoder[T]) QueryParamsEncoder {
	if query != nil {
		return queryParamsFunc(func() (string, error) {
			return query.EncodeQueryParams(set)
		})
	} else {
		return set
	}
}

// BatchCreate executes a batch_create with the given slice of entities
func BatchCreate[K comparable, V Object[V]](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (createdEntities []*restlidata.CreatedEntity[K], err error) {
	return batchCreate[K, V, *restlidata.CreatedEntity[K]](c, ctx, rp, entities, query, readOnlyFields)
}

func BatchCreateWithReturnEntity[K comparable, V Object[V]](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (createdEntities []*restlidata.CreatedAndReturnedEntity[K, V], err error) {
	return batchCreate[K, V, *restlidata.CreatedAndReturnedEntity[K, V]](c, ctx, rp, entities, query, readOnlyFields)
}

func batchCreate[K comparable, V Object[V], R restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (createdEntities []R, err error) {
	req, err := NewCreateRequest(c, ctx, rp, query, Method_batch_create, restlicodec.MarshalerFunc(func(writer restlicodec.Writer) error {
		return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
			return restlicodec.WriteArray(keyWriter(restlidata.ElementsField), entities, V.MarshalRestLi)
		})
	}), readOnlyFields)
	if err != nil {
		return nil, err
	}

	elements, _, err := DoAndUnmarshal(c, req, restlicodec.UnmarshalRestLi[*restlidata.Elements[R]])
	if err != nil {
		return nil, err
	}
	return elements.Elements, nil
}

func BatchDelete[K comparable](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	keys []K,
	query batchQueryParamsEncoder[K],
) (*restlidata.BatchResponse[K, *restlidata.BatchEntityUpdateResponse], error) {
	keySet := batchkeyset.NewBatchKeySet[K]()
	err := batchkeyset.AddAllKeys(keySet, keys...)
	if err != nil {
		return nil, err
	}

	req, err := NewDeleteRequest(c, ctx, rp, batchQueryParams(keySet, query), Method_batch_delete)
	if err != nil {
		return nil, err
	}

	return doBatchQuery[K, *restlidata.BatchEntityUpdateResponse](c, keySet, req)
}

func BatchGet[K comparable, V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	keys []K,
	query batchQueryParamsEncoder[K],
) (*restlidata.BatchResponse[K, V], error) {
	keySet := batchkeyset.NewBatchKeySet[K]()
	err := batchkeyset.AddAllKeys(keySet, keys...)
	if err != nil {
		return nil, err
	}

	req, err := NewGetRequest(c, ctx, rp, batchQueryParams(keySet, query), Method_batch_get)
	if err != nil {
		return nil, err
	}

	return doBatchQuery[K, V](c, keySet, req)
}

func BatchUpdate[K comparable, V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	entities map[K]V,
	query batchQueryParamsEncoder[K],
	createAndReadOnlyFields restlicodec.PathSpec,
) (*restlidata.BatchResponse[K, *restlidata.BatchEntityUpdateResponse], error) {
	keys := batchkeyset.NewBatchKeySet[K]()
	err := batchkeyset.AddAllMapKeys(keys, entities)
	if err != nil {
		return nil, err
	}

	req, err := NewJsonRequest(
		c,
		ctx,
		rp,
		batchQueryParams(keys, query),
		http.MethodPut,
		Method_batch_update,
		batchEntities[K, V](entities),
		createAndReadOnlyFields,
	)
	if err != nil {
		return nil, err
	}

	return doBatchQuery[K, *restlidata.BatchEntityUpdateResponse](c, keys, req)
}

func BatchPartialUpdate[K comparable, PV restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	entities map[K]PV,
	query batchQueryParamsEncoder[K],
	createAndReadOnlyFields restlicodec.PathSpec,
) (*restlidata.BatchResponse[K, *restlidata.BatchEntityUpdateResponse], error) {
	keys := batchkeyset.NewBatchKeySet[K]()
	err := batchkeyset.AddAllMapKeys(keys, entities)
	if err != nil {
		return nil, err
	}

	req, err := NewJsonRequest(
		c,
		ctx,
		rp,
		batchQueryParams(keys, query),
		http.MethodPost,
		Method_batch_partial_update,
		batchEntities[K, PV](entities),
		createAndReadOnlyFields,
	)
	if err != nil {
		return nil, err
	}

	return doBatchQuery[K, *restlidata.BatchEntityUpdateResponse](c, keys, req)
}

func doBatchQuery[K comparable, V restlicodec.Marshaler](
	c *Client,
	keys batchkeyset.BatchKeySet[K],
	req *http.Request,
) (res *restlidata.BatchResponse[K, V], err error) {
	data, _, err := c.do(req)
	if err != nil {
		return nil, err
	}

	r, err := restlicodec.NewJsonReader(data)
	if err != nil {
		return nil, err
	}
	res = new(restlidata.BatchResponse[K, V])
	return res, res.UnmarshalWithKeyLocator(r, keys)
}

type batchEntities[K comparable, V restlicodec.Marshaler] map[K]V

func (b batchEntities[K, V]) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		return restlidata.MarshalBatchEntities(b, keyWriter(restlidata.EntitiesField))
	})
}

var entitiesRequiredResponseFields = restlicodec.RequiredFields{restlidata.EntitiesField}

func (b batchEntities[K, V]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(entitiesRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == restlidata.EntitiesField {
			return restlidata.UnmarshalBatchEntities(b, reader)
		} else {
			return reader.Skip()
		}
	})
}
