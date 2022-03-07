package protocol

import (
	"context"
	"net/http"

	"github.com/PapaCharlie/go-restli/protocol/batchkeyset"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/PapaCharlie/go-restli/protocol/stdtypes"
)

type queryParamsFunc func() (string, error)

func (q queryParamsFunc) EncodeQueryParams() (string, error) {
	return q()
}

type BatchQueryParams[T any] interface {
	EncodeQueryParams(set batchkeyset.BatchKeySet[T]) (string, error)
}

func batchQueryParams[T any](set batchkeyset.BatchKeySet[T], query BatchQueryParams[T]) QueryParams {
	if query != nil {
		return queryParamsFunc(func() (string, error) {
			return query.EncodeQueryParams(set)
		})
	} else {
		return set
	}
}

type CreatedEntity[K any] struct {
	Id     K
	Status int
}

type CreatedAndReturnedEntity[K, V any] struct {
	CreatedEntity[K]
	Entity V
}

func (c *CollectionClient[K, V, PV]) BatchCreateRequest(
	ctx context.Context,
	rp ResourcePath,
	query QueryParams,
	method RestLiMethod,
	create restlicodec.Marshaler,
) (*http.Request, error) {
	return c.NewJsonRequest(ctx, rp, query, http.MethodPost, method, create, c.ReadOnlyFields)
}

// BatchCreate executes a batch_create with the given slice of entities
func (c *CollectionClient[K, V, PV]) BatchCreate(
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParams,
) (createdEntities []*CreatedEntity[K], err error) {
	res, err := c.batchCreate(ctx, rp, entities, query, false)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		createdEntities = append(createdEntities, &e.CreatedEntity)
	}

	return createdEntities, nil
}

func (c *CollectionClient[K, V, PV]) BatchCreateWithReturnEntity(
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParams,
) (createdEntities []*CreatedAndReturnedEntity[K, V], err error) {
	return c.batchCreate(ctx, rp, entities, query, true)
}

func (c *CollectionClient[K, V, PV]) batchCreate(
	ctx context.Context,
	rp ResourcePath,
	entities []V,
	query QueryParams,
	returnEntity bool,
) (createdEntities []*CreatedAndReturnedEntity[K, V], err error) {
	u, err := c.FormatQueryUrl(rp, query)
	if err != nil {
		return nil, err
	}

	req, err := BatchCreateRequest(ctx, u, entities, c.ReadOnlyFields)
	if err != nil {
		return nil, err
	}

	createdEntities, _, err = DoAndUnmarshal(c.RestLiClient, req, func(reader restlicodec.Reader) (entities []*CreatedAndReturnedEntity[K, V], err error) {
		err = reader.ReadRecord(elementsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
			var requiredFields restlicodec.RequiredFields
			if returnEntity {
				requiredFields = batchCreateWithReturnResponseRequiredFields
			} else {
				requiredFields = batchCreateResponseRequiredFields
			}
			switch field {
			case elementsField:
				return reader.ReadArray(func(reader restlicodec.Reader) (err error) {
					e := &CreatedAndReturnedEntity[K, V]{}
					err = reader.ReadRecord(requiredFields, func(reader restlicodec.Reader, field string) (err error) {
						switch field {
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

							e.Id, err = c.KeyUnmarshaler(r)
						case statusField:
							e.Status, err = reader.ReadInt()
						case entityField:
							if returnEntity {
								e.Entity, err = c.EntityUnmarshaler(reader)
							} else {
								err = reader.Skip()
							}
						default:
							err = reader.Skip()
						}
						return err
					})
					if err != nil {
						return err
					}
					entities = append(entities, e)
					return nil
				})
			default:
				return reader.Skip()
			}
		})
		return entities, err
	})
	return createdEntities, err
}

func (c *CollectionClient[K, V, PV]) BatchDelete(
	ctx context.Context,
	rp ResourcePath,
	keys []K,
	query BatchQueryParams[K],
) (map[K]*BatchEntityUpdateResponse, error) {
	keySet := c.BatchKeySetProvider()
	err := batchkeyset.AddAllKeys(keySet, keys...)
	if err != nil {
		return nil, err
	}

	req, err := c.NewDeleteRequest(ctx, rp, batchQueryParams(keySet, query), Method_batch_delete)
	if err != nil {
		return nil, err
	}

	return doBatchQuery(c.RestLiClient, keySet, UnmarshalBatchEntityUpdateResponse, req)
}

func (c *CollectionClient[K, V, PV]) BatchGet(
	ctx context.Context,
	rp ResourcePath,
	keys []K,
	query BatchQueryParams[K],
) (map[K]V, error) {
	keySet := c.BatchKeySetProvider()
	err := batchkeyset.AddAllKeys(keySet, keys...)
	if err != nil {
		return nil, err
	}

	req, err := c.NewGetRequest(ctx, rp, batchQueryParams(keySet, query), Method_batch_get)
	if err != nil {
		return nil, err
	}

	return doBatchQuery(c.RestLiClient, keySet, c.EntityUnmarshaler, req)
}

func (c *CollectionClient[K, V, PV]) BatchUpdate(
	ctx context.Context,
	rp ResourcePath,
	entities map[K]V,
	query BatchQueryParams[K],
) (map[K]*BatchEntityUpdateResponse, error) {
	keys := c.BatchKeySetProvider()
	err := batchkeyset.AddAllMapKeys(keys, entities)
	if err != nil {
		return nil, err
	}

	req, err := c.NewJsonRequest(
		ctx,
		rp,
		batchQueryParams(keys, query),
		http.MethodPut,
		Method_batch_update,
		&batchEntities[K, V]{
			keys:      keys,
			entities:  entities,
			isPartial: false,
		},
		c.CreateAndReadOnlyFields,
	)
	if err != nil {
		return nil, err
	}

	return doBatchQuery(c.RestLiClient, keys, UnmarshalBatchEntityUpdateResponse, req)
}

func (c *CollectionClient[K, V, PV]) BatchPartialUpdate(
	ctx context.Context,
	rp ResourcePath,
	entities map[K]PV,
	query BatchQueryParams[K],
) (map[K]*BatchEntityUpdateResponse, error) {
	keys := c.BatchKeySetProvider()
	err := batchkeyset.AddAllMapKeys(keys, entities)
	if err != nil {
		return nil, err
	}

	req, err := c.NewJsonRequest(
		ctx,
		rp,
		batchQueryParams(keys, query),
		http.MethodPost,
		Method_batch_partial_update,
		&batchEntities[K, PV]{
			keys:      keys,
			entities:  entities,
			isPartial: true,
		},
		c.CreateAndReadOnlyFields,
	)
	if err != nil {
		return nil, err
	}

	return doBatchQuery(c.RestLiClient, keys, UnmarshalBatchEntityUpdateResponse, req)
}

func doBatchQuery[K comparable, T any](
	c *RestLiClient,
	keys batchkeyset.BatchKeySet[K],
	unmarshaler restlicodec.GenericUnmarshaler[T],
	req *http.Request,
) (entities map[K]T, err error) {
	data, _, err := c.do(req)
	if err != nil {
		return nil, err
	}

	reader := restlicodec.NewJsonReader(data)
	entities = make(map[K]T)
	errors := make(BatchRequestResponseError[K])

	err = reader.ReadRecord(batchRequestResponseRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		isResults, isErrors := field == resultsField, field == errorsField
		if isResults || isErrors {
			return reader.ReadMap(func(valueReader restlicodec.Reader, rawKey string) (err error) {
				keyReader, err := restlicodec.NewRor2Reader(rawKey)
				if err != nil {
					return err
				}

				originalKey, err := keys.LocateOriginalKey(keyReader)
				if err != nil {
					return err
				}

				if isResults {
					entities[originalKey], err = unmarshaler(valueReader)
				} else {
					errors[originalKey], err = stdtypes.UnmarshalRestLiErrorResponse(valueReader)
				}
				return err
			})
		} else {
			return reader.Skip()
		}
	})
	if err != nil {
		return nil, err
	}
	if len(errors) > 0 {
		return entities, errors
	}
	return entities, nil
}

type batchEntities[K comparable, V restlicodec.Marshaler] struct {
	keys      batchkeyset.BatchKeySet[K]
	entities  map[K]V
	isPartial bool
}

func (b *batchEntities[K, V]) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		return keyWriter("entities").WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
			for k, v := range b.entities {
				w := restlicodec.NewRor2HeaderWriter()
				err = b.keys.MarshalKey(w, k)
				if err != nil {
					return err
				}

				writer := keyWriter(w.Finalize())
				if b.isPartial {
					err = writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
						return v.MarshalRestLi(keyWriter("patch"))
					})
				} else {
					err = v.MarshalRestLi(writer)
				}
				if err != nil {
					return err
				}
			}
			return nil
		})
	})
}

type BatchEntityUpdateResponse struct {
	Status int
}

func UnmarshalBatchEntityUpdateResponse(reader restlicodec.Reader) (r *BatchEntityUpdateResponse, err error) {
	r = new(BatchEntityUpdateResponse)
	err = r.UnmarshalRestLi(reader)
	return r, err
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
