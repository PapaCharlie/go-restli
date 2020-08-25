package protocol

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

func (c *RestLiClient) DoGetRequest(ctx context.Context, url *url.URL, result restlicodec.Unmarshaler) (*http.Response, error) {
	req, err := GetRequest(ctx, url, Method_get)
	if err != nil {
		return nil, err
	}

	return c.DoAndDecode(req, result)
}

var CreateResponseHasNoEntityHeaderError = errors.New("response from CREATE request did not specify a " + RestLiHeader_ID + " header")

func (c *RestLiClient) DoCreateRequest(ctx context.Context, url *url.URL, create restlicodec.Marshaler, id restlicodec.Unmarshaler, returnEntity restlicodec.Unmarshaler) (res *http.Response, err error) {
	req, err := JsonRequest(ctx, url, http.MethodPost, Method_create, create)
	if err != nil {
		return nil, err
	}

	if returnEntity != nil {
		res, err = c.DoAndDecode(req, returnEntity)
	} else {
		res, err = c.DoAndIgnore(req)
	}
	if err != nil {
		return res, err
	}

	if res.StatusCode/100 != 2 {
		return res, fmt.Errorf("invalid response code from %s: %d", url, res.StatusCode)
	}

	if h := res.Header.Get(RestLiHeader_ID); len(h) > 0 {
		err = id.UnmarshalRestLi(restlicodec.NewHeaderReader(h))
		if err != nil {
			return res, err
		}
	} else {
		return res, CreateResponseHasNoEntityHeaderError
	}

	return res, nil
}

func (c *RestLiClient) DoUpdateRequest(ctx context.Context, url *url.URL, create restlicodec.Marshaler) (*http.Response, error) {
	req, err := JsonRequest(ctx, url, http.MethodPut, Method_update, create)
	if err != nil {
		return nil, err
	}

	res, err := c.DoAndIgnore(req)
	if err != nil {
		return res, err
	}

	if res.StatusCode/100 != 2 {
		return res, fmt.Errorf("invalid response code from %s: %d", url, res.StatusCode)
	}

	return res, nil
}

func (c *RestLiClient) DoPartialUpdateRequest(ctx context.Context, url *url.URL, patch restlicodec.Marshaler) (*http.Response, error) {
	req, err := JsonRequest(ctx, url, http.MethodPost, Method_partial_update, &PartialUpdate{Patch: patch})
	if err != nil {
		return nil, err
	}

	res, err := c.DoAndIgnore(req)
	if err != nil {
		return res, err
	}

	if res.StatusCode/100 != 2 {
		return res, fmt.Errorf("invalid response code from %s: %d", url, res.StatusCode)
	}

	return res, nil
}

func (c *RestLiClient) DoDeleteRequest(ctx context.Context, url *url.URL) (*http.Request, error) {
	req, err := DeleteRequest(ctx, url, Method_delete)
	if err != nil {
		return nil, err
	}

	res, err := c.DoAndIgnore(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("invalid response code from %s: %d", url, res.StatusCode)
	}

	return req, nil
}

func (c *RestLiClient) DoFinderRequest(ctx context.Context, url *url.URL, results restlicodec.Unmarshaler) (*http.Response, error) {
	req, err := GetRequest(ctx, url, Method_finder)
	if err != nil {
		return nil, err
	}

	return c.DoAndDecode(req, results)
}

func (c *RestLiClient) DoActionRequest(ctx context.Context, url *url.URL, params restlicodec.Marshaler, results restlicodec.Unmarshaler) (res *http.Response, err error) {
	req, err := JsonRequest(ctx, url, http.MethodPost, Method_action, params)
	if err != nil {
		return nil, err
	}

	if results != nil {
		res, err = c.DoAndDecode(req, results)
	} else {
		res, err = c.DoAndIgnore(req)
	}
	return res, err
}
