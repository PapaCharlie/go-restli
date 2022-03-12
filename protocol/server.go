package protocol

import (
	"net/http"
	"path"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

func RegisterCreate[K any, RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder, V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	readOnlyFields restlicodec.PathSpec,
	create func(*RequestContext, RP, V, QP) (*CreatedEntity[K], error),
) {
	registerMethodWithBody(s, segments, Method_create, readOnlyFields, 0,
		restlicodec.UnmarshalRestLi[V],
		func(ctx *RequestContext, rp RP, v V, qp QP) (responseBody restlicodec.Marshaler, err error) {
			ctx.ResponseStatus = http.StatusCreated
			createdEntity, err := create(ctx, rp, v, qp)
			if err != nil {
				return nil, err
			}

			err = writeIdHeaders(ctx, createdEntity.Id)
			if err != nil {
				return newErrorResponsef(err, http.StatusInternalServerError,
					"%q failed, could not serialize ID header: %s", Method_create)
			}

			if createdEntity.Status != 0 {
				ctx.ResponseStatus = createdEntity.Status
			}
			return nil, nil
		})
}

func writeIdHeaders[K any](ctx *RequestContext, id K) (err error) {
	w := restlicodec.NewRor2HeaderWriter()
	err = restlicodec.MarshalRestLi(id, w)
	if err != nil {
		return err
	}
	s := w.Finalize()
	ctx.ResponseHeaders.Add(RestLiHeader_ID, s)
	ctx.ResponseHeaders.Add("Location", path.Join(ctx.Request.RequestURI, s))
	return nil
}

func RegisterCreateWithReturnEntity[K any, RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder, V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	readOnlyFields restlicodec.PathSpec,
	create func(*RequestContext, RP, V, QP) (*CreatedAndReturnedEntity[K, V], error),
) {
	registerMethodWithBody(s, segments, Method_create, readOnlyFields, 0,
		restlicodec.UnmarshalRestLi[V],
		func(ctx *RequestContext, rp RP, v V, qp QP) (responseBody restlicodec.Marshaler, err error) {
			ctx.ResponseStatus = http.StatusCreated
			createdEntity, err := create(ctx, rp, v, qp)
			if err != nil {
				return nil, err
			}

			err = writeIdHeaders(ctx, createdEntity.Id)
			if err != nil {
				return newErrorResponsef(err, http.StatusInternalServerError,
					"%q failed, could not serialize ID header: %s", Method_create)
			}

			if createdEntity.Status != 0 {
				ctx.ResponseStatus = createdEntity.Status
			}
			return createdEntity.Entity, nil
		})
}

func RegisterBatchCreate[K any, RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder, V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	readOnlyFields restlicodec.PathSpec,
	batchCreate func(*RequestContext, RP, []V, QP) ([]*CreatedEntity[K], error),
) {
	registerMethodWithBody(s, segments, Method_batch_create, readOnlyFields, 2,
		restlicodec.UnmarshalRestLi[*Elements[V]],
		func(ctx *RequestContext, rp RP, v *Elements[V], qp QP) (responseBody restlicodec.Marshaler, err error) {
			entities, err := batchCreate(ctx, rp, v.Elements, qp)
			if err != nil {
				return nil, err
			} else {
				return &Elements[*CreatedEntity[K]]{Elements: entities}, nil
			}
		},
	)
}

func RegisterBatchCreateWithReturnEntity[K any, RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder, V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	readOnlyFields restlicodec.PathSpec,
	batchCreate func(*RequestContext, RP, []V, QP) ([]*CreatedAndReturnedEntity[K, V], error),
) {
	registerMethodWithBody(s, segments, Method_batch_create, readOnlyFields, 2,
		restlicodec.UnmarshalRestLi[*Elements[V]],
		func(ctx *RequestContext, rp RP, v *Elements[V], qp QP) (responseBody restlicodec.Marshaler, err error) {
			entities, err := batchCreate(ctx, rp, v.Elements, qp)
			if err != nil {
				return nil, err
			} else {
				return &Elements[*CreatedAndReturnedEntity[K, V]]{Elements: entities}, nil
			}
		},
	)
}

func RegisterGet[RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder, V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	get func(*RequestContext, RP, QP) (V, error),
) {
	registerMethodWithNoBody(s, segments, Method_get,
		func(ctx *RequestContext, rp RP, qp QP) (responseBody restlicodec.Marshaler, err error) {
			return get(ctx, rp, qp)
		})
}

func RegisterBatchGet[K comparable, RP ResourcePathUnmarshaler, QP BatchQueryParamsDecoder[K], V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	batchGet func(*RequestContext, RP, []K, QP) (*BatchResponse[K, V], error),
) {
	registerMethodWithNoBody(s, segments, Method_batch_get,
		func(ctx *RequestContext, rp RP, qp *WrappedBatchQueryParamsDecoder[K, QP]) (responseBody restlicodec.Marshaler, err error) {
			return batchGet(ctx, rp, qp.keys, qp.qp)
		})
}

func RegisterGetAll[RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder, V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	getAll func(*RequestContext, RP, QP) (*Elements[V], error),
) {
	registerMethodWithNoBody(s, segments, Method_get_all,
		func(ctx *RequestContext, rp RP, qp QP) (responseBody restlicodec.Marshaler, err error) {
			return getAll(ctx, rp, qp)
		})
}

func RegisterDelete[RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder](
	s Server,
	segments []ResourcePathSegment,
	deleteF func(*RequestContext, RP, QP) error,
) {
	registerMethodWithNoBody(s, segments, Method_delete,
		func(ctx *RequestContext, rp RP, qp QP) (responseBody restlicodec.Marshaler, err error) {
			ctx.ResponseStatus = http.StatusNoContent
			return nil, deleteF(ctx, rp, qp)
		})
}

func RegisterBatchDelete[K comparable, RP ResourcePathUnmarshaler, QP BatchQueryParamsDecoder[K]](
	s Server,
	segments []ResourcePathSegment,
	batchDelete func(*RequestContext, RP, []K, QP) (*BatchResponse[K, *BatchEntityUpdateResponse], error),
) {
	registerMethodWithNoBody(s, segments, Method_batch_delete,
		func(ctx *RequestContext, rp RP, qp *WrappedBatchQueryParamsDecoder[K, QP]) (responseBody restlicodec.Marshaler, err error) {
			return batchDelete(ctx, rp, qp.keys, qp.qp)
		})
}

func RegisterUpdate[RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder, V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	readAndCreateOnlyFields restlicodec.PathSpec,
	update func(*RequestContext, RP, V, QP) error,
) {
	registerMethodWithBody(s, segments, Method_update, readAndCreateOnlyFields, 0,
		restlicodec.UnmarshalRestLi[V],
		func(ctx *RequestContext, rp RP, v V, qp QP) (responseBody restlicodec.Marshaler, err error) {
			ctx.ResponseStatus = http.StatusNoContent
			return nil, update(ctx, rp, v, qp)
		})
}

func RegisterBatchUpdate[K comparable, RP ResourcePathUnmarshaler, QP BatchQueryParamsDecoder[K], V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	readAndCreateOnlyFields restlicodec.PathSpec,
	batchUpdate func(*RequestContext, RP, map[K]V, QP) (*BatchResponse[K, *BatchEntityUpdateResponse], error),
) {
	registerMethodWithBody(s, segments, Method_batch_update, readAndCreateOnlyFields, 2,
		func(reader restlicodec.Reader) (entities map[K]V, err error) {
			entities = make(map[K]V)
			return entities, batchEntities[K, V](entities).UnmarshalRestLi(reader)
		},
		func(ctx *RequestContext, rp RP, v map[K]V, qp *WrappedBatchQueryParamsDecoder[K, QP]) (responseBody restlicodec.Marshaler, err error) {
			return batchUpdate(ctx, rp, v, qp.qp)
		})
}

func RegisterPartialUpdate[RP ResourcePathUnmarshaler, QP restlicodec.QueryParamsDecoder, V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	excludedFields restlicodec.PathSpec,
	partialUpdate func(*RequestContext, RP, V, QP) error,
) {
	registerMethodWithBody(s, segments, Method_partial_update, excludedFields, 0,
		restlicodec.UnmarshalRestLi[V],
		func(ctx *RequestContext, rp RP, v V, qp QP) (responseBody restlicodec.Marshaler, err error) {
			ctx.ResponseStatus = http.StatusNoContent
			return nil, partialUpdate(ctx, rp, v, qp)
		})
}

func RegisterBatchPartialUpdate[K comparable, RP ResourcePathUnmarshaler, QP BatchQueryParamsDecoder[K], V restlicodec.Marshaler](
	s Server,
	segments []ResourcePathSegment,
	readAndCreateOnlyFields restlicodec.PathSpec,
	batchPartialUpdate func(*RequestContext, RP, map[K]V, QP) (*BatchResponse[K, *BatchEntityUpdateResponse], error),
) {
	registerMethodWithBody(s, segments, Method_batch_partial_update, readAndCreateOnlyFields, 3,
		func(reader restlicodec.Reader) (entities map[K]V, err error) {
			entities = make(map[K]V)
			return entities, batchEntities[K, V](entities).UnmarshalRestLi(reader)
		},
		func(ctx *RequestContext, rp RP, v map[K]V, qp *WrappedBatchQueryParamsDecoder[K, QP]) (responseBody restlicodec.Marshaler, err error) {
			return batchPartialUpdate(ctx, rp, v, qp.qp)
		})
}
