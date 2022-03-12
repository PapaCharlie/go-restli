package protocol

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/PapaCharlie/go-restli/protocol/batchkeyset"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/PapaCharlie/go-restli/protocol/stdtypes"
)

type queryParamsFunc func() (string, error)

func (q queryParamsFunc) EncodeQueryParams() (string, error) {
	return q()
}

type BatchQueryParamsEncoder[T any] interface {
	EncodeQueryParams(set batchkeyset.BatchKeySet[T]) (string, error)
}

type BatchQueryParamsDecoder[T any] interface {
	DecodeQueryParams(reader restlicodec.QueryParamsReader) ([]T, error)
}

type SliceBatchQueryParams[T any] struct{}

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

type WrappedBatchQueryParamsDecoder[T any, QP BatchQueryParamsDecoder[T]] struct {
	keys []T
	qp   QP
}

func (b *WrappedBatchQueryParamsDecoder[T, QP]) DecodeQueryParams(reader restlicodec.QueryParamsReader) (err error) {
	v := reflect.New(reflect.TypeOf(b.qp).Elem())
	b.qp = v.Interface().(QP)
	b.keys, err = b.qp.DecodeQueryParams(reader)
	return err
}

func batchQueryParams[T any](set batchkeyset.BatchKeySet[T], query BatchQueryParamsEncoder[T]) QueryParamsEncoder {
	if query != nil {
		return queryParamsFunc(func() (string, error) {
			return query.EncodeQueryParams(set)
		})
	} else {
		return set
	}
}

// type CreatedEntities[K any] []*CreatedEntity[K]
//
// func (c CreatedEntities[K]) MarshalRestLi(writer restlicodec.Writer) error {
// 	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
// 		return restlicodec.WriteArray(keyWriter(elementsField), c, (*CreatedEntity[K]).MarshalRestLi)
// 	})
// }
//
// func (c *CreatedEntities[K]) UnmarshalRestLi(reader restlicodec.Reader) error {
// 	return reader.ReadRecord(elementsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
// 		if field == elementsField {
// 			var arr []*CreatedEntity[K]
// 			arr, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[*CreatedEntity[K]])
// 			*c = arr
// 			return err
// 		} else {
// 			return reader.Skip()
// 		}
// 	})
// }

type CreatedEntity[K any] struct {
	Location *string
	Id       K
	Status   int
	Error    *stdtypes.ErrorResponse
}

func (c *CreatedEntity[K]) SetLocation(ctx *RequestContext) error {
	c.Location = new(string)
	id, err := c.marshalIdForLocation()
	if err != nil {
		return err
	}
	*c.Location = strings.TrimSuffix(ctx.Request.RequestURI, "/") + "/" + id
	return nil
}

func (c *CreatedEntity[K]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(batchCreateResponseRequiredFields, c.unmarshalRestLi)
}

func (c *CreatedEntity[K]) unmarshalRestLi(reader restlicodec.Reader, field string) (err error) {
	switch field {
	case locationField:
		c.Location = new(string)
		*c.Location, err = reader.ReadString()
		return err
	case idField:
		var rawKey string
		rawKey, err = reader.ReadString()
		if err != nil {
			return err
		}

		var r restlicodec.Reader
		r, err = restlicodec.NewRor2Reader(rawKey)
		if err != nil {
			return err
		}

		c.Id, err = restlicodec.UnmarshalRestLi[K](r)
		return err
	case statusField:
		c.Status, err = reader.ReadInt()
		return err
	case errorField:
		c.Error, err = restlicodec.UnmarshalRestLi[*stdtypes.ErrorResponse](reader)
		return err
	default:
		return reader.Skip()
	}
}

func (c *CreatedEntity[K]) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		err = c.marshalId(keyWriter)
		if err != nil {
			return err
		}

		c.marshalLocation(keyWriter)
		c.marshalStatus(keyWriter)
		return nil
	})
}

func (c *CreatedEntity[K]) marshalIdForLocation() (id string, err error) {
	w := restlicodec.NewRor2PathWriter()
	err = restlicodec.MarshalRestLi[K](c.Id, w)
	if err != nil {
		return "", err
	}
	return w.Finalize(), nil
}

func (c *CreatedEntity[K]) marshalId(keyWriter func(key string) restlicodec.Writer) (err error) {
	w := restlicodec.NewRor2PathWriter()
	err = restlicodec.MarshalRestLi[K](c.Id, w)
	if err != nil {
		return err
	}
	keyWriter(idField).WriteString(w.Finalize())
	return nil
}

func (c *CreatedEntity[K]) marshalLocation(keyWriter func(key string) restlicodec.Writer) {
	if c.Location != nil {
		keyWriter(locationField).WriteString(*c.Location)
	}
}

func (c *CreatedEntity[K]) marshalStatus(keyWriter func(key string) restlicodec.Writer) {
	s := c.Status
	if s == 0 {
		s = http.StatusCreated
	}
	keyWriter(statusField).WriteInt(s)
}

// type CreatedAndReturnedEntities[K any, V restlicodec.Marshaler] []*CreatedAndReturnedEntity[K, V]
//
// func (c CreatedAndReturnedEntities[K, V]) MarshalRestLi(writer restlicodec.Writer) error {
// 	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
// 		return restlicodec.WriteArray(keyWriter(elementsField), c, (*CreatedAndReturnedEntity[K, V]).MarshalRestLi)
// 	})
// }
//
// func (c *CreatedAndReturnedEntities[K, V]) UnmarshalRestLi(reader restlicodec.Reader) error {
// 	return reader.ReadRecord(elementsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
// 		if field == elementsField {
// 			var arr []*CreatedAndReturnedEntity[K, V]
// 			arr, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[*CreatedAndReturnedEntity[K, V]])
// 			*c = arr
// 			return err
// 		} else {
// 			return reader.Skip()
// 		}
// 	})
// }

type CreatedAndReturnedEntity[K any, V restlicodec.Marshaler] struct {
	CreatedEntity[K]
	Entity V
}

func (c *CreatedAndReturnedEntity[K, V]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(batchCreateWithReturnResponseRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == entityField {
			c.Entity, err = restlicodec.UnmarshalRestLi[V](reader)
			return err
		} else {
			return c.unmarshalRestLi(reader, field)
		}
	})
}

func (c *CreatedAndReturnedEntity[K, V]) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		err = c.Entity.MarshalRestLi(keyWriter(entityField))
		if err != nil {
			return err
		}

		err = c.marshalId(keyWriter)
		if err != nil {
			return err
		}

		c.marshalLocation(keyWriter)
		c.marshalStatus(keyWriter)
		return nil
	})
}

// BatchCreate executes a batch_create with the given slice of entities
func BatchCreate[K comparable, V RestLiObject[V]](
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (createdEntities []*CreatedEntity[K], err error) {
	return batchCreate[K, V, *CreatedEntity[K]](c, ctx, rp, entities, query, readOnlyFields)
}

func BatchCreateWithReturnEntity[K comparable, V RestLiObject[V]](
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (createdEntities []*CreatedAndReturnedEntity[K, V], err error) {
	return batchCreate[K, V, *CreatedAndReturnedEntity[K, V]](c, ctx, rp, entities, query, readOnlyFields)
}

func batchCreate[K comparable, V RestLiObject[V], R any](
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (createdEntities []R, err error) {
	req, err := NewCreateRequest(c, ctx, rp, query, Method_batch_create, restlicodec.MarshalerFunc(func(writer restlicodec.Writer) error {
		return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
			return restlicodec.WriteArray(keyWriter(elementsField), entities, V.MarshalRestLi)
		})
	}), readOnlyFields)
	if err != nil {
		return nil, err
	}

	createdEntities, _, err = DoAndUnmarshal(c, req, func(reader restlicodec.Reader) (entities []R, err error) {
		err = reader.ReadRecord(elementsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
			if field == elementsField {
				entities, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[R])
				return err
			} else {
				return reader.Skip()
			}
		})
		return entities, err
	})
	return createdEntities, err
}

func BatchDelete[K comparable](
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	keys []K,
	query BatchQueryParamsEncoder[K],
) (*BatchResponse[K, *BatchEntityUpdateResponse], error) {
	keySet := batchkeyset.NewBatchKeySet[K]()
	err := batchkeyset.AddAllKeys(keySet, keys...)
	if err != nil {
		return nil, err
	}

	req, err := NewDeleteRequest(c, ctx, rp, batchQueryParams(keySet, query), Method_batch_delete)
	if err != nil {
		return nil, err
	}

	return doBatchQuery[K, *BatchEntityUpdateResponse](c, keySet, req)
}

func BatchGet[K comparable, V restlicodec.Marshaler](
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	keys []K,
	query BatchQueryParamsEncoder[K],
) (*BatchResponse[K, V], error) {
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
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	entities map[K]V,
	query BatchQueryParamsEncoder[K],
	createAndReadOnlyFields restlicodec.PathSpec,
) (*BatchResponse[K, *BatchEntityUpdateResponse], error) {
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

	return doBatchQuery[K, *BatchEntityUpdateResponse](c, keys, req)
}

func BatchPartialUpdate[K comparable, PV restlicodec.Marshaler](
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	entities map[K]PV,
	query BatchQueryParamsEncoder[K],
	createAndReadOnlyFields restlicodec.PathSpec,
) (*BatchResponse[K, *BatchEntityUpdateResponse], error) {
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

	return doBatchQuery[K, *BatchEntityUpdateResponse](c, keys, req)
}

func doBatchQuery[K comparable, V restlicodec.Marshaler](
	c *RestLiClient,
	keys batchkeyset.BatchKeySet[K],
	req *http.Request,
) (res *BatchResponse[K, V], err error) {
	data, _, err := c.do(req)
	if err != nil {
		return nil, err
	}

	r, err := restlicodec.NewJsonReader(data)
	if err != nil {
		return nil, err
	}
	res = new(BatchResponse[K, V])
	return res, res.unmarshalRestLi(r, keys)
}

type batchEntities[K comparable, V restlicodec.Marshaler] map[K]V

func (b batchEntities[K, V]) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		return MarshalBatchEntities(b, keyWriter(entitiesField))
	})
}

func (b batchEntities[K, V]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(entitiesRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == entitiesField {
			return UnmarshalBatchEntities(b, reader)
		} else {
			return reader.Skip()
		}
	})
}

func MarshalBatchEntities[K comparable, V any](entities map[K]V, writer restlicodec.Writer) (err error) {
	return restlicodec.WriteGenericMap(writer, entities,
		func(k K) (string, error) {
			w := restlicodec.NewRor2HeaderWriter()
			err = restlicodec.MarshalRestLi(k, w)
			if err != nil {
				return "", err
			}
			return w.Finalize(), nil
		},
		func(v V, writer restlicodec.Writer) error {
			return restlicodec.MarshalRestLi(v, writer.SetScope())
		},
	)
}

func UnmarshalBatchEntities[K comparable, V restlicodec.Marshaler](entities map[K]V, reader restlicodec.Reader) (err error) {
	return reader.ReadMap(func(reader restlicodec.Reader, field string) (err error) {
		r, err := restlicodec.NewRor2Reader(field)
		if err != nil {
			return err
		}
		k, err := restlicodec.UnmarshalRestLi[K](r)
		if err != nil {
			return err
		}

		var v V
		v, err = restlicodec.UnmarshalRestLi[V](reader)
		if err != nil {
			return err
		}
		entities[k] = v
		return nil
	})
}

func NewBatchResponse[K comparable, V restlicodec.Marshaler]() *BatchResponse[K, V] {
	return &BatchResponse[K, V]{
		Statuses: map[K]int{},
		Results:  map[K]V{},
		Errors:   map[K]*stdtypes.ErrorResponse{},
	}
}

type BatchResponse[K comparable, V restlicodec.Marshaler] struct {
	Statuses map[K]int
	Results  map[K]V
	Errors   map[K]*stdtypes.ErrorResponse
}

func (b *BatchResponse[K, V]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return b.unmarshalRestLi(reader, nil)
}

func (b *BatchResponse[K, V]) unmarshalRestLi(reader restlicodec.Reader, keys batchkeyset.BatchKeySet[K]) error {
	return reader.ReadRecord(batchResponseRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case resultsField:
			b.Results = make(map[K]V)
		case statusesField:
			b.Statuses = make(map[K]int)
		case errorsField:
			b.Errors = make(map[K]*stdtypes.ErrorResponse)
		default:
			return reader.Skip()
		}

		return reader.ReadMap(func(valueReader restlicodec.Reader, rawKey string) (err error) {
			keyReader, err := restlicodec.NewRor2Reader(rawKey)
			if err != nil {
				return err
			}

			var originalKey K
			if keys != nil {
				originalKey, err = keys.LocateOriginalKeyFromReader(keyReader)
			} else {
				originalKey, err = restlicodec.UnmarshalRestLi[K](keyReader)
			}
			if err != nil {
				return err
			}

			switch field {
			case resultsField:
				b.Results[originalKey], err = restlicodec.UnmarshalRestLi[V](valueReader)
			case statusesField:
				b.Statuses[originalKey], err = valueReader.ReadInt()
			case errorsField:
				b.Errors[originalKey], err = restlicodec.UnmarshalRestLi[*stdtypes.ErrorResponse](valueReader)
			}
			return err
		})
	})
}

func (b *BatchResponse[K, V]) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		err = MarshalBatchEntities(b.Errors, keyWriter(errorsField))
		if err != nil {
			return err
		}

		err = MarshalBatchEntities(b.Results, keyWriter(resultsField))
		if err != nil {
			return err
		}

		err = MarshalBatchEntities(b.Statuses, keyWriter(statusesField))
		if err != nil {
			return err
		}

		return nil
	})
}

type BatchEntityUpdateResponse struct {
	Status int
}

func (b *BatchEntityUpdateResponse) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		s := b.Status
		if s == 0 {
			s = http.StatusNoContent
		}
		keyWriter(statusField).WriteInt(s)
		return nil
	})
}

func (b *BatchEntityUpdateResponse) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(batchEntityUpdateResponseRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case statusField:
			b.Status, err = reader.ReadInt()
		default:
			err = reader.Skip()
		}
		return err
	})
}
