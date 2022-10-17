package common

import (
	"net/http"

	"github.com/PapaCharlie/go-restli/restlicodec"
)

const (
	ElementsField = "elements"
	ValueField    = "value"
	StatusField   = "status"
	StatusesField = "statuses"
	ResultsField  = "results"
	ErrorField    = "error"
	ErrorsField   = "errors"
	IdField       = "id"
	LocationField = "location"
	PagingField   = "paging"
	MetadataField = "metadata"
	EntityField   = "entity"
	EntitiesField = "entities"
)

var (
	elementsRequiredResponseFields              = restlicodec.NewRequiredFields().Add(ElementsField)
	batchEntityUpdateResponseRequiredFields     = restlicodec.NewRequiredFields().Add(StatusField)
	batchResponseRequiredFields                 = restlicodec.NewRequiredFields().Add(ResultsField)
	batchCreateResponseRequiredFields           = restlicodec.NewRequiredFields().Add(StatusField)
	batchCreateWithReturnResponseRequiredFields = restlicodec.NewRequiredFields().Add(StatusField, EntityField)
)

type BatchEntityUpdateResponse struct {
	Status int
}

func (b *BatchEntityUpdateResponse) NewInstance() *BatchEntityUpdateResponse {
	return new(BatchEntityUpdateResponse)
}

func (b *BatchEntityUpdateResponse) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		s := b.Status
		if s == 0 {
			s = http.StatusNoContent
		}
		keyWriter(StatusField).WriteInt(s)
		return nil
	})
}

func (b *BatchEntityUpdateResponse) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(batchEntityUpdateResponseRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case StatusField:
			b.Status, err = reader.ReadInt()
		default:
			err = restlicodec.NoSuchFieldErr
		}
		return err
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

type BatchResponse[K comparable, V restlicodec.Marshaler] struct {
	Statuses map[K]int
	Results  map[K]V
	Errors   map[K]*ErrorResponse
}

func (b *BatchResponse[K, V]) AddStatus(key K, value int) {
	if b.Statuses == nil {
		b.Statuses = make(map[K]int)
	}
	b.Statuses[key] = value
}

func (b *BatchResponse[K, V]) AddResult(key K, value V) {
	if b.Results == nil {
		b.Results = make(map[K]V)
	}
	b.Results[key] = value
}

func (b *BatchResponse[K, V]) AddError(key K, value *ErrorResponse) {
	if b.Errors == nil {
		b.Errors = make(map[K]*ErrorResponse)
	}
	b.Errors[key] = value
}

func (b *BatchResponse[K, V]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return b.UnmarshalWithKeyLocator(reader, nil)
}

type KeyLocator[T any] interface {
	LocateOriginalKeyFromReader(keyReader restlicodec.Reader) (originalKey T, err error)
}

func (b *BatchResponse[K, V]) UnmarshalWithKeyLocator(reader restlicodec.Reader, keys KeyLocator[K]) error {
	return reader.ReadRecord(batchResponseRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case ResultsField:
			b.Results = make(map[K]V)
		case StatusesField:
			b.Statuses = make(map[K]int)
		case ErrorsField:
			b.Errors = make(map[K]*ErrorResponse)
		default:
			return restlicodec.NoSuchFieldErr
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
			case ResultsField:
				b.Results[originalKey], err = restlicodec.UnmarshalRestLi[V](valueReader)
			case StatusesField:
				b.Statuses[originalKey], err = valueReader.ReadInt()
			case ErrorsField:
				b.Errors[originalKey], err = restlicodec.UnmarshalRestLi[*ErrorResponse](valueReader)
			}
			return err
		})
	})
}

func (b *BatchResponse[K, V]) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		err = MarshalBatchEntities(b.Errors, keyWriter(ErrorsField))
		if err != nil {
			return err
		}

		err = MarshalBatchEntities(b.Results, keyWriter(ResultsField))
		if err != nil {
			return err
		}

		err = MarshalBatchEntities(b.Statuses, keyWriter(StatusesField))
		if err != nil {
			return err
		}

		return nil
	})
}

type CreatedEntity[K any] struct {
	Location *string
	Id       K
	Status   int
}

func (c *CreatedEntity[K]) NewInstance() *CreatedEntity[K] {
	return new(CreatedEntity[K])
}

func (c *CreatedEntity[K]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(batchCreateResponseRequiredFields, c.unmarshalRestLi)
}

func (c *CreatedEntity[K]) unmarshalRestLi(reader restlicodec.Reader, field string) (err error) {
	switch field {
	case LocationField:
		c.Location = new(string)
		*c.Location, err = reader.ReadString()
		return err
	case IdField:
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
	case StatusField:
		c.Status, err = reader.ReadInt()
		return err
	default:
		return restlicodec.NoSuchFieldErr
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

func (c *CreatedEntity[K]) marshalId(keyWriter func(key string) restlicodec.Writer) (err error) {
	w := restlicodec.NewRor2PathWriter()
	err = restlicodec.MarshalRestLi[K](c.Id, w)
	if err != nil {
		return err
	}
	keyWriter(IdField).WriteString(w.Finalize())
	return nil
}

func (c *CreatedEntity[K]) marshalLocation(keyWriter func(key string) restlicodec.Writer) {
	if c.Location != nil {
		keyWriter(LocationField).WriteString(*c.Location)
	}
}

func (c *CreatedEntity[K]) marshalStatus(keyWriter func(key string) restlicodec.Writer) {
	s := c.Status
	if s == 0 {
		s = http.StatusCreated
	}
	keyWriter(StatusField).WriteInt(s)
}

type CreatedAndReturnedEntity[K any, V restlicodec.Marshaler] struct {
	CreatedEntity[K]
	Entity V
}

func (c *CreatedAndReturnedEntity[K, V]) NewInstance() *CreatedAndReturnedEntity[K, V] {
	return new(CreatedAndReturnedEntity[K, V])
}

func (c *CreatedAndReturnedEntity[K, V]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(batchCreateWithReturnResponseRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == EntityField {
			c.Entity, err = restlicodec.UnmarshalRestLi[V](reader)
			return err
		} else {
			return c.unmarshalRestLi(reader, field)
		}
	})
}

func (c *CreatedAndReturnedEntity[K, V]) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		err = c.Entity.MarshalRestLi(keyWriter(EntityField))
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

type Elements[V restlicodec.Marshaler] struct {
	Elements []V
	Paging   *CollectionMetadata
}

func (f *Elements[V]) NewInstance() *Elements[V] {
	return new(Elements[V])
}

func (f *Elements[V]) marshalElements(keyWriter func(key string) restlicodec.Writer) (err error) {
	return restlicodec.WriteArray(keyWriter(ElementsField), f.Elements, V.MarshalRestLi)
}

func (f *Elements[V]) marshalPaging(keyWriter func(key string) restlicodec.Writer) error {
	if f.Paging != nil {
		return f.Paging.MarshalRestLi(keyWriter(PagingField))
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
	case ElementsField:
		f.Elements, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[V])
		return err
	case PagingField:
		f.Paging = new(CollectionMetadata)
		return f.Paging.UnmarshalRestLi(reader)
	default:
		return restlicodec.NoSuchFieldErr
	}
}

func (f *Elements[V]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(elementsRequiredResponseFields, f.unmarshalRestLi)
}

type ElementsWithMetadata[V, M restlicodec.Marshaler] struct {
	Elements []V
	Paging   *CollectionMetadata
	Metadata M
}

func (f *ElementsWithMetadata[V, M]) NewInstance() *ElementsWithMetadata[V, M] {
	return new(ElementsWithMetadata[V, M])
}

func (f *ElementsWithMetadata[V, M]) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
		err = restlicodec.WriteArray(keyWriter(ElementsField), f.Elements, V.MarshalRestLi)
		if err != nil {
			return err
		}

		err = f.Metadata.MarshalRestLi(keyWriter(MetadataField))
		if err != nil {
			return err
		}

		if f.Paging != nil {
			err = f.Paging.MarshalRestLi(keyWriter(PagingField))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (f *ElementsWithMetadata[V, M]) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadRecord(elementsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case MetadataField:
			f.Metadata, err = restlicodec.UnmarshalRestLi[M](reader)
			return err
		case ElementsField:
			f.Elements, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[V])
			return err
		case PagingField:
			f.Paging = new(CollectionMetadata)
			return f.Paging.UnmarshalRestLi(reader)
		default:
			return restlicodec.NoSuchFieldErr
		}
	})
}
