package protocol

import (
	"context"
	"net/http"
	"net/url"

	"github.com/PapaCharlie/go-restli/protocol/batchkeyset"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/PapaCharlie/go-restli/protocol/stdstructs"
)

type CreatedEntity[K any] struct {
	Id     K
	Status int
}

type CreatedAndReturnedEntity[K, V any] struct {
	CreatedEntity[K]
	Entity V
}

// DoBatchCreateRequest executes a batch_create with the given slice of entities
func DoBatchCreateRequest[K any, V restlicodec.Marshaler](
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	entities []V,
	readOnlyFields restlicodec.PathSpec,
	keyUnmarshaler restlicodec.GenericUnmarshaler[K],
) (createdEntities []*CreatedEntity[K], err error) {
	res, err := DoBatchCreateRequestWithReturnEntity[K, V](
		c,
		ctx,
		url,
		entities,
		readOnlyFields,
		keyUnmarshaler,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		createdEntities = append(createdEntities, &e.CreatedEntity)
	}

	return createdEntities, nil
}

func DoBatchCreateRequestWithReturnEntity[K any, V restlicodec.Marshaler](
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	entities []V,
	readOnlyFields restlicodec.PathSpec,
	keyUnmarshaler restlicodec.GenericUnmarshaler[K],
	entityUnmarshaler restlicodec.GenericUnmarshaler[V],
) (createdEntities []*CreatedAndReturnedEntity[K, V], err error) {
	req, err := BatchCreateRequest(ctx, url, entities, readOnlyFields)
	if err != nil {
		return nil, err
	}

	createdEntities, _, err = DoAndUnmarshal(c, req, func(reader restlicodec.Reader) (entities []*CreatedAndReturnedEntity[K, V], err error) {
		err = reader.ReadRecord(elementsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
			var requiredFields restlicodec.RequiredFields
			if entityUnmarshaler != nil {
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

							e.Id, err = keyUnmarshaler(r)
						case statusField:
							e.Status, err = reader.ReadInt()
						case entityField:
							if entityUnmarshaler != nil {
								e.Entity, err = entityUnmarshaler(reader)
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

func DoBatchDeleteRequest[K comparable](
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	keys batchkeyset.BatchKeySet[K],
) (map[K]*BatchEntityUpdateResponse, error) {
	req, err := DeleteRequest(ctx, url, Method_batch_delete)
	if err != nil {
		return nil, err
	}

	return doBatchQuery(c, keys, UnmarshalBatchEntityUpdateResponse, req)
}

func DoBatchGetRequest[K comparable, V any](
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	keys batchkeyset.BatchKeySet[K],
	unmarshaler restlicodec.GenericUnmarshaler[V],
) (map[K]V, error) {
	req, err := GetRequest(ctx, url, Method_batch_get)
	if err != nil {
		return nil, err
	}

	return doBatchQuery(c, keys, unmarshaler, req)
}

func DoBatchUpdateRequest[K comparable, V restlicodec.Marshaler](
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	keys batchkeyset.BatchKeySet[K],
	entities map[K]V,
	createAndReadOnlyFields restlicodec.PathSpec,
) (map[K]*BatchEntityUpdateResponse, error) {
	req, err := JsonRequest(ctx, url, http.MethodPut, Method_batch_update, &batchEntities[K, V]{
		keys:      keys,
		entities:  entities,
		isPartial: false,
	}, createAndReadOnlyFields)
	if err != nil {
		return nil, err
	}

	return doBatchQuery(c, keys, UnmarshalBatchEntityUpdateResponse, req)
}

func DoBatchPartialUpdateRequest[K comparable, V restlicodec.Marshaler](
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	keys batchkeyset.BatchKeySet[K],
	entities map[K]V,
	createAndReadOnlyFields restlicodec.PathSpec,
) (map[K]*BatchEntityUpdateResponse, error) {
	req, err := JsonRequest(ctx, url, http.MethodPost, Method_batch_partial_update, &batchEntities[K, V]{
		keys:      keys,
		entities:  entities,
		isPartial: true,
	}, createAndReadOnlyFields)
	if err != nil {
		return nil, err
	}

	return doBatchQuery(c, keys, UnmarshalBatchEntityUpdateResponse, req)
}

func doBatchQuery[K comparable, V any](
	c *RestLiClient,
	keys batchkeyset.BatchKeySet[K],
	unmarshaler restlicodec.GenericUnmarshaler[V],
	req *http.Request,
) (entities map[K]V, err error) {
	data, _, err := c.do(req)
	if err != nil {
		return nil, err
	}

	reader := restlicodec.NewJsonReader(data)
	entities = make(map[K]V)
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
					errors[originalKey], err = stdstructs.UnmarshalErrorResponse(valueReader)
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
