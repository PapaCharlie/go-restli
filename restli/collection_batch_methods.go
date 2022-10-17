package restli

import (
	"context"
	"net/http"

	"github.com/PapaCharlie/go-restli/restli/batchkeyset"
	"github.com/PapaCharlie/go-restli/restlicodec"
	"github.com/PapaCharlie/go-restli/restlidata/generated/com/linkedin/restli/common"
)

type queryParamsFunc func() (string, error)

func (q queryParamsFunc) EncodeQueryParams() (string, error) {
	return q()
}

type batchQueryParamsEncoder[T any] interface {
	EncodeQueryParams(set batchkeyset.BatchKeySet[T]) (string, error)
}

type batchQueryParamsDecoder[T, QP any] interface {
	NewInstance() QP
	DecodeQueryParams(reader restlicodec.QueryParamsReader) ([]T, error)
}

type SliceBatchQueryParams[T any] struct{}

var entityIdsRequiredResponseFields = restlicodec.NewRequiredFields().Add(batchkeyset.EntityIDsField)

func (s *SliceBatchQueryParams[T]) NewInstance() *SliceBatchQueryParams[T] {
	return new(SliceBatchQueryParams[T])
}

func (s *SliceBatchQueryParams[T]) DecodeQueryParams(reader restlicodec.QueryParamsReader) (ids []T, err error) {
	err = reader.ReadRecord(entityIdsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == batchkeyset.EntityIDsField {
			ids, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[T])
		} else {
			err = restlicodec.NoSuchFieldErr
		}
		return err
	})
	return ids, err
}

type wrappedBatchQueryParamsDecoder[T any, QP batchQueryParamsDecoder[T, QP]] struct {
	keys []T
	qp   QP
}

func (b *wrappedBatchQueryParamsDecoder[T, QP]) NewInstance() *wrappedBatchQueryParamsDecoder[T, QP] {
	return new(wrappedBatchQueryParamsDecoder[T, QP])
}

func (b *wrappedBatchQueryParamsDecoder[T, QP]) DecodeQueryParams(reader restlicodec.QueryParamsReader) (err error) {
	b.qp = b.qp.NewInstance()
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
func BatchCreate[K comparable, V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (createdEntities []*common.CreatedEntity[K], err error) {
	return batchCreate[V, *common.CreatedEntity[K]](c, ctx, rp, entities, query, readOnlyFields)
}

func BatchCreateWithReturnEntity[K comparable, V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (createdEntities []*common.CreatedAndReturnedEntity[K, V], err error) {
	return batchCreate[V, *common.CreatedAndReturnedEntity[K, V]](c, ctx, rp, entities, query, readOnlyFields)
}

func batchCreate[V restlicodec.Marshaler, R restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (createdEntities []R, err error) {
	req, err := NewCreateRequest(c, ctx, rp, query, Method_batch_create, restlicodec.MarshalerFunc(func(writer restlicodec.Writer) error {
		return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
			return restlicodec.WriteArray(keyWriter(common.ElementsField), entities, V.MarshalRestLi)
		})
	}), readOnlyFields)
	if err != nil {
		return nil, err
	}

	elements, _, err := DoAndUnmarshal(c, req, restlicodec.UnmarshalRestLi[*common.Elements[R]])
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
) (*common.BatchResponse[K, *common.BatchEntityUpdateResponse], error) {
	keySet := batchkeyset.NewBatchKeySet[K]()
	err := batchkeyset.AddAllKeys(keySet, keys...)
	if err != nil {
		return nil, err
	}

	req, err := NewDeleteRequest(c, ctx, rp, batchQueryParams(keySet, query), Method_batch_delete)
	if err != nil {
		return nil, err
	}

	return doBatchQuery[K, *common.BatchEntityUpdateResponse](c, keySet, req)
}

func BatchGet[K comparable, V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	keys []K,
	query batchQueryParamsEncoder[K],
) (*common.BatchResponse[K, V], error) {
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
) (*common.BatchResponse[K, *common.BatchEntityUpdateResponse], error) {
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

	return doBatchQuery[K, *common.BatchEntityUpdateResponse](c, keys, req)
}

func BatchPartialUpdate[K comparable, PV restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	entities map[K]PV,
	query batchQueryParamsEncoder[K],
	createAndReadOnlyFields restlicodec.PathSpec,
) (*common.BatchResponse[K, *common.BatchEntityUpdateResponse], error) {
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

	return doBatchQuery[K, *common.BatchEntityUpdateResponse](c, keys, req)
}

func doBatchQuery[K comparable, V restlicodec.Marshaler](
	c *Client,
	keys batchkeyset.BatchKeySet[K],
	req *http.Request,
) (res *common.BatchResponse[K, V], err error) {
	data, _, err := c.do(req)
	if err != nil {
		return nil, err
	}

	r, err := restlicodec.NewJsonReader(data)
	if err != nil {
		return nil, err
	}
	res = new(common.BatchResponse[K, V])
	return res, res.UnmarshalWithKeyLocator(r, keys)
}

type batchEntities[K comparable, V restlicodec.Marshaler] map[K]V

func (b batchEntities[K, V]) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		return common.MarshalBatchEntities(b, keyWriter(common.EntitiesField))
	})
}

var entitiesRequiredResponseFields = restlicodec.NewRequiredFields().Add(common.EntitiesField)

func (b batchEntities[K, V]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(entitiesRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == common.EntitiesField {
			return common.UnmarshalBatchEntities(b, reader)
		} else {
			return restlicodec.NoSuchFieldErr
		}
	})
}
