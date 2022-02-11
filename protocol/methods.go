package protocol

import (
	"context"
	"net/http"
	"net/url"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

// DoGetRequest executes a rest.li Get request against the given url and parses the results in the given
// restlicodec.Unmarshaler
func DoGetRequest[V any](c *RestLiClient, ctx context.Context, url *url.URL, unmarshaler restlicodec.GenericUnmarshaler[V]) (v V, err error) {
	req, err := GetRequest(ctx, url, Method_get)
	if err != nil {
		return v, err
	}

	v, _, err = DoAndUnmarshal(c, req, unmarshaler)
	return v, err
}

// DoCreateRequest executes a rest.li Create request against the given url and places the given restlicodec.Marshaler in
// the request's body. The X-RestLi-Id header field will be parsed into id (though a
// CreateResponseHasNoEntityHeaderError will be returned if the header is not set) and if returnEntity is non-nil, it
// will be used to unmarhsal the body of the response.
func DoCreateRequest[K any](
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	create restlicodec.Marshaler,
	readOnlyFields restlicodec.PathSpec,
	unmarshaler restlicodec.GenericUnmarshaler[K],
) (k K, err error) {
	req, err := CreateRequest(ctx, url, Method_create, create, readOnlyFields)
	if err != nil {
		return k, err
	}

	res, err := c.DoAndIgnore(req)
	if err != nil {
		return k, err
	}

	return unmarshalReturnEntityKey(c, res, unmarshaler)
}

// DoCreateRequestWithReturnEntity executes a rest.li Create request against the given url and places the given
// restlicodec.Marshaler in the request's body. The X-RestLi-Id header field will be parsed into id (though a
// CreateResponseHasNoEntityHeaderError will be returned if the header is not set) and entityUnmarshaler will be used to
// unmarhsal the body of the response.
func DoCreateRequestWithReturnEntity[K, V any](
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	create restlicodec.Marshaler,
	readOnlyFields restlicodec.PathSpec,
	keyUnmarshaler restlicodec.GenericUnmarshaler[K],
	entityUnmarshaler restlicodec.GenericUnmarshaler[V],
) (k K, v V, err error) {
	req, err := CreateRequest(ctx, url, Method_create, create, readOnlyFields)
	if err != nil {
		return k, v, err
	}

	v, res, err := DoAndUnmarshal(c, req, entityUnmarshaler)
	if err != nil {
		return k, v, err
	}

	k, err = unmarshalReturnEntityKey(c, res, keyUnmarshaler)
	return k, v, err
}

func unmarshalReturnEntityKey[K any](
	c *RestLiClient,
	res *http.Response,
	unmarshaler restlicodec.GenericUnmarshaler[K],
) (k K, err error) {
	if h := res.Header.Get(RestLiHeader_ID); len(h) > 0 {
		var reader restlicodec.Reader
		reader, err = restlicodec.NewRor2Reader(h)
		if err != nil {
			return k, err
		}

		k, err = unmarshaler(reader)
		if _, mfe := err.(*restlicodec.MissingRequiredFieldsError); mfe && !c.StrictResponseDeserialization {
			err = nil
		}
		return k, err
	} else {
		return k, &CreateResponseHasNoEntityHeaderError{Response: res}
	}
}

// DoUpdateRequest executes a rest.li Update request and places the given restlicodec.Marshaler in the request's body.
func DoUpdateRequest(c *RestLiClient, ctx context.Context, url *url.URL, update restlicodec.Marshaler) error {
	req, err := JsonRequest(ctx, url, http.MethodPut, Method_update, update, nil)
	if err != nil {
		return err
	}

	_, err = c.DoAndIgnore(req)
	return err
}

// DoPartialUpdateRequest executes a rest.li Partial Update request and places the given patch objects wrapped in a
// PartialUpdate in the request's body.
func DoPartialUpdateRequest(
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	patch restlicodec.Marshaler,
	createAndReadOnlyFields restlicodec.PathSpec,
) error {
	req, err := JsonRequest(ctx, url, http.MethodPost, Method_partial_update,
		restlicodec.MarshalerFunc(func(writer restlicodec.Writer) error {
			return writer.WriteMap(func(fieldNameWriter func(fieldName string) restlicodec.Writer) error {
				return patch.MarshalRestLi(fieldNameWriter("patch").SetScope())
			})
		}), createAndReadOnlyFields)
	if err != nil {
		return err
	}

	_, err = c.DoAndIgnore(req)
	return err
}

// DoDeleteRequest executes a rest.li Delete request
func DoDeleteRequest(c *RestLiClient, ctx context.Context, url *url.URL) error {
	req, err := DeleteRequest(ctx, url, Method_delete)
	if err != nil {
		return err
	}

	_, err = c.DoAndIgnore(req)
	return err
}

// DoFinderRequest executes a rest.li Finder request and uses the given restlicodec.Unmarshaler to unmarshal the
// response's body.
func DoFinderRequest[T any](
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	unmarshaler restlicodec.GenericUnmarshaler[T],
) (results []T, total *int, err error) {
	req, err := GetRequest(ctx, url, Method_finder)
	if err != nil {
		return nil, nil, err
	}

	results, _, err = DoAndUnmarshal(c, req, func(reader restlicodec.Reader) (results []T, err error) {
		err = reader.ReadRecord(elementsRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
			switch field {
			case elementsField:
				results, err = restlicodec.ReadArray(reader, unmarshaler)
				return err
			case "paging":
				return reader.ReadMap(func(reader restlicodec.Reader, key string) (err error) {
					if key == "total" {
						var t int
						t, err = reader.ReadInt()
						if err != nil {
							return err
						}
						total = &t
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

	return results, total, err
}

// DoActionRequest executes a rest.li Action request and places the given restlicodec.Marshaler in the request's body
// and discards the response body. Actions with no params are expected to use the EmptyRecord.
func DoActionRequest(
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	params restlicodec.Marshaler,
) (err error) {
	req, err := newActionRequest(ctx, url, params)
	if err != nil {
		return err
	}

	_, err = c.DoAndIgnore(req)
	return err
}

// DoActionRequestWithResults executes a rest.li Action request and places the given restlicodec.Marshaler in the
// request's body, and returns the results after deserialization. Actions with no params are expected to use the
// EmptyRecord.
func DoActionRequestWithResults[T any](
	c *RestLiClient,
	ctx context.Context,
	url *url.URL,
	params restlicodec.Marshaler,
	results restlicodec.GenericUnmarshaler[T],
) (t T, err error) {
	req, err := newActionRequest(ctx, url, params)
	if err != nil {
		return t, err
	}

	t, _, err = DoAndUnmarshal(c, req, func(reader restlicodec.Reader) (t T, err error) {
		err = reader.ReadRecord(actionRequiredResponseFields, func(reader restlicodec.Reader, field string) (err error) {
			switch field {
			case valueField:
				t, err = results(reader)
				return err
			default:
				return reader.Skip()
			}
		})
		return t, err
	})
	return t, err
}

func newActionRequest(ctx context.Context, url *url.URL, params restlicodec.Marshaler) (*http.Request, error) {
	return JsonRequest(ctx, url, http.MethodPost, Method_action, params, nil)
}
