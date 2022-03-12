package protocol

import (
	"context"
	"log"
	"net/http"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/PapaCharlie/go-restli/protocol/stdtypes"
)

type ActionQueryParam string

func (a ActionQueryParam) EncodeQueryParams() (string, error) {
	return "action=" + string(a), nil
}

// DoActionRequest executes a rest.li Action request and places the given restlicodec.Marshaler in the request's body
// and discards the response body. Actions with no params are expected to use the EmptyRecord.
func DoActionRequest(
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	params restlicodec.Marshaler,
) (err error) {
	req, err := newActionRequest(c, ctx, rp, query, params)
	if err != nil {
		return err
	}

	_, err = DoAndIgnore(c, req)
	return err
}

// DoActionRequestWithResults executes a rest.li Action request and places the given restlicodec.Marshaler in the
// request's body, and returns the results after deserialization. Actions with no params are expected to use the
// EmptyRecord.
func DoActionRequestWithResults[T any](
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	params restlicodec.Marshaler,
	unmarshaler restlicodec.GenericUnmarshaler[T],
) (t T, err error) {
	req, err := newActionRequest(c, ctx, rp, query, params)
	if err != nil {
		return t, err
	}

	t, _, err = DoAndUnmarshal(c, req, func(reader restlicodec.Reader) (t T, err error) {
		err = reader.ReadRecord(actionRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
			switch field {
			case valueField:
				t, err = unmarshaler(reader)
				return err
			default:
				return reader.Skip()
			}
		})
		return t, err
	})
	return t, err
}

func newActionRequest(
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	params restlicodec.Marshaler,
) (*http.Request, error) {
	return NewJsonRequest(c, ctx, rp, query, http.MethodPost, Method_action, params, nil)
}

func RegisterAction[RP ResourcePathUnmarshaler, P any](
	s Server,
	segments []ResourcePathSegment,
	name string,
	action func(*RequestContext, RP, P) error,
) {
	registerAction[RP, P, struct{}](s, segments, name, nil,
		func(ctx *RequestContext, rp RP, params P) (struct{}, error) {
			return struct{}{}, action(ctx, rp, params)
		})
}

func RegisterActionWithResults[RP ResourcePathUnmarshaler, P, R any](
	s Server,
	segments []ResourcePathSegment,
	name string,
	resultsMarshaler restlicodec.GenericMarshaler[R],
	action func(*RequestContext, RP, P) (R, error),
) {
	registerAction(s, segments, name, resultsMarshaler,
		func(ctx *RequestContext, rp RP, params P) (R, error) {
			return action(ctx, rp, params)
		})
}

func registerAction[RP ResourcePathUnmarshaler, P, R any](
	s Server,
	segments []ResourcePathSegment,
	name string,
	resultsMarshaler restlicodec.GenericMarshaler[R],
	h func(ctx *RequestContext, rp RP, params P) (R, error),
) {
	p := s.subNode(segments)
	if _, ok := p.actions[name]; ok {
		log.Panicf("go-restli: Cannot register action %q twice for %v", name, segments)
	}
	p.actions[name] = func(
		ctx *RequestContext,
		segments []restlicodec.Reader,
		body []byte,
	) (responseBody restlicodec.Marshaler, err error) {
		rp, err := UnmarshalResourcePath[RP](segments)
		if err != nil {
			return newErrorResponsef(nil, http.StatusBadRequest, "Invalid path for action %q: %s", name)
		}

		// Special case for action params that are EmptyRecord: skip reading the body altogether, as it's valid for an
		// action with no parameters to supply an empty POST body
		var params P
		if !stdtypes.IsEmptyRecord(params) {
			var r restlicodec.Reader
			r, err = restlicodec.NewJsonReader(body)
			if err == nil {
				params, err = restlicodec.UnmarshalRestLi[P](r)
			}
			if err != nil {
				return newErrorResponsef(err, http.StatusBadRequest, "Invalid arguments for action %q: %s", name)
			}
		}

		results, err := h(ctx, rp, params)
		if err != nil {
			return newErrorResponsef(err, http.StatusBadRequest, "Action %q failed: %s", name)
		}
		if resultsMarshaler != nil {
			return restlicodec.MarshalerFunc(func(writer restlicodec.Writer) error {
				return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) (err error) {
					return resultsMarshaler(results, keyWriter(valueField))
				})
			}), nil
		} else {
			return nil, nil
		}
	}
}
