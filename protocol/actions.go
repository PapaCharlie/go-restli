package protocol

import (
	"context"
	"net/http"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

// DoActionRequest executes a rest.li Action request and places the given restlicodec.Marshaler in the request's body
// and discards the response body. Actions with no params are expected to use the EmptyRecord.
func DoActionRequest(
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParams,
	params restlicodec.Marshaler,
) (err error) {
	req, err := newActionRequest(c, ctx, rp, query, params)
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
	rp ResourcePath,
	query QueryParams,
	params restlicodec.Marshaler,
	results restlicodec.GenericUnmarshaler[T],
) (t T, err error) {
	req, err := newActionRequest(c, ctx, rp, query, params)
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

func newActionRequest(
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParams,
	params restlicodec.Marshaler,
) (*http.Request, error) {
	u, err := c.FormatQueryUrl(rp, query)
	if err != nil {
		return nil, err
	}
	return NewJsonRequest(ctx, u, http.MethodPost, Method_action, params, nil)
}

type ActionQueryParam string

func (a ActionQueryParam) EncodeQueryParams() (string, error) {
	return "action=" + string(a), nil
}
