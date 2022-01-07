package protocol

import (
	"context"
	"net/http"
	"net/url"
	"sort"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

// DoGetRequest executes a rest.li Get request against the given url and parses the results in the given
// restlicodec.Unmarshaler
func (c *RestLiClient) DoGetRequest(ctx context.Context, url *url.URL, result restlicodec.Unmarshaler) (err error) {
	req, err := GetRequest(ctx, url, Method_get)
	if err != nil {
		return err
	}

	_, err = c.DoAndDecode(req, result)
	return err
}

// CreateResponseHasNoEntityHeaderError is used specifically when a Create request succeeds but the resource
// implementation does not set the X-RestLi-Id header. This error is recoverable and can be ignored if the response id
// is not required
type CreateResponseHasNoEntityHeaderError struct {
	Request  *http.Request
	Response *http.Response
}

func (c CreateResponseHasNoEntityHeaderError) Error() string {
	return "response from CREATE request did not specify a " + RestLiHeader_ID + " header"
}

// DoCreateRequest executes a rest.li Create request against the given url and places the given restlicodec.Marshaler in
// the request's body. The X-RestLi-Id header field will be parsed into id (though a
// CreateResponseHasNoEntityHeaderError will be returned if the header is not set) and if returnEntity is non-nil, it
// will be used to unmarhsal the body of the response.
func (c *RestLiClient) DoCreateRequest(
	ctx context.Context,
	url *url.URL,
	create restlicodec.Marshaler,
	readOnlyFields restlicodec.PathSpec,
	id restlicodec.Unmarshaler,
	returnEntity restlicodec.Unmarshaler,
) (err error) {
	req, err := JsonRequest(ctx, url, http.MethodPost, Method_create, create, readOnlyFields)
	if err != nil {
		return err
	}

	var res *http.Response
	if returnEntity != nil {
		res, err = c.DoAndDecode(req, returnEntity)
	} else {
		res, err = c.DoAndIgnore(req)
	}
	if err != nil {
		return err
	}

	if h := res.Header.Get(RestLiHeader_ID); len(h) > 0 {
		var reader restlicodec.Reader
		reader, err = restlicodec.NewRor2Reader(h)
		if err != nil {
			return err
		}
		err = id.UnmarshalRestLi(reader)
		if _, mfe := err.(*restlicodec.MissingRequiredFieldsError); mfe && !c.StrictResponseDeserialization {
			err = nil
		}
		if err != nil {
			return err
		}
	} else {
		return &CreateResponseHasNoEntityHeaderError{
			Request:  req,
			Response: res,
		}
	}

	return nil
}

// DoUpdateRequest executes a rest.li Update request and places the given restlicodec.Marshaler in the request's body.
func (c *RestLiClient) DoUpdateRequest(ctx context.Context, url *url.URL, update restlicodec.Marshaler) error {
	req, err := JsonRequest(ctx, url, http.MethodPut, Method_update, update, nil)
	if err != nil {
		return err
	}

	_, err = c.DoAndIgnore(req)
	return err
}

// DoPartialUpdateRequest executes a rest.li Partial Update request and places the given patch objects wrapped in a
// PartialUpdate in the request's body.
func (c *RestLiClient) DoPartialUpdateRequest(
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
func (c *RestLiClient) DoDeleteRequest(ctx context.Context, url *url.URL) error {
	req, err := DeleteRequest(ctx, url, Method_delete)
	if err != nil {
		return err
	}

	_, err = c.DoAndIgnore(req)
	return err
}

// DoFinderRequest executes a rest.li Finder request and uses the given restlicodec.Unmarshaler to unmarshal the
// response's body.
func (c *RestLiClient) DoFinderRequest(ctx context.Context, url *url.URL, elements restlicodec.ArrayReader) (total *int, err error) {
	req, err := GetRequest(ctx, url, Method_finder)
	if err != nil {
		return nil, err
	}

	results := restlicodec.UnmarshalerFunc(func(reader restlicodec.Reader) error {
		const elementsField = "elements"
		hasElements := false

		err = reader.ReadMap(func(reader restlicodec.Reader, key string) error {
			switch key {
			case elementsField:
				hasElements = true
				return reader.ReadArray(elements)
			case "paging":
				return reader.ReadMap(func(reader restlicodec.Reader, key string) (err error) {
					if key == "total" {
						var t int64
						t, err = reader.ReadInt64()
						if err != nil {
							return err
						}
						tInt := int(t)
						total = &tInt
					} else {
						err = reader.Skip()
					}
					return nil
				})
			default:
				return reader.Skip()
			}
		})

		if err != nil {
			return err
		}

		if !hasElements {
			reader.RecordMissingRequiredFields(map[string]struct{}{elementsField: {}})
		}

		return reader.CheckMissingFields()
	})

	_, err = c.DoAndDecode(req, results)
	return total, err
}

// DoActionRequest executes a rest.li Action request and places the given restlicodec.Marshaler in the request's body.
// Actions with no params are expected to use the EmptyRecord instead. If the given restlicodec.Unmarshaler for the
// results is non-nil, it will be used to unmarshal the request's body, otherwise the body will be discarded.
func (c *RestLiClient) DoActionRequest(ctx context.Context, url *url.URL, params restlicodec.Marshaler, results restlicodec.Unmarshaler) error {
	req, err := JsonRequest(ctx, url, http.MethodPost, Method_action, params, nil)
	if err != nil {
		return err
	}

	if results != nil {
		_, err = c.DoAndDecode(req, results)
	} else {
		_, err = c.DoAndIgnore(req)
	}
	return err
}

type BatchEntityIDsEncoder []restlicodec.WriteCloser

func (e *BatchEntityIDsEncoder) AddEntityID() restlicodec.Writer {
	writer := restlicodec.NewRestLiQueryParamsWriter()
	*e = append(*e, writer)
	return writer
}

func (e *BatchEntityIDsEncoder) Encode(paramNameWriter func(string) restlicodec.Writer) error {
	encodedKeys := make([]string, len(*e))
	for i, w := range *e {
		encodedKeys[i] = w.Finalize()
	}
	sort.Strings(encodedKeys)

	return paramNameWriter("ids").WriteArray(func(itemWriter func() restlicodec.Writer) error {
		for _, k := range encodedKeys {
			itemWriter().WriteRawBytes([]byte(k))
		}
		return nil
	})
}

func (e *BatchEntityIDsEncoder) GenerateRawQuery() (string, error) {
	writer := restlicodec.NewRestLiQueryParamsWriter()
	err := writer.WriteParams(func(paramNameWriter func(key string) restlicodec.Writer) error {
		return e.Encode(paramNameWriter)
	})
	if err != nil {
		return "", err
	}
	return writer.Finalize(), nil
}

func (c *RestLiClient) DoBatchGetRequest(ctx context.Context, url *url.URL, reader BatchResultsReader) (err error) {
	req, err := GetRequest(ctx, url, Method_batch_get)
	if err != nil {
		return err
	}

	res := &batchRequestResponse{
		Results: reader,
		Errors: &BatchMethodError{
			Request: req,
		},
	}
	res.Errors.Response, err = c.DoAndDecode(req, res)
	return err
}

func (c *RestLiClient) DoBatchPartialUpdateRequest(
	ctx context.Context,
	url *url.URL,
	entities BatchEntities,
	reader BatchResultsReader,
	createAndReadOnlyFields restlicodec.PathSpec,
) (err error) {
	req, err := JsonRequest(ctx, url, http.MethodPost, Method_batch_partial_update, entities, createAndReadOnlyFields)
	if err != nil {
		return err
	}

	res := &batchRequestResponse{
		Results: reader,
		Errors: &BatchMethodError{
			Request: req,
		},
	}
	res.Errors.Response, err = c.DoAndDecode(req, res)
	return err
}
