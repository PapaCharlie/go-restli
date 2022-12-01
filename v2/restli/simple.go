package restli

import (
	"context"
	"net/http"

	"github.com/PapaCharlie/go-restli/v2/restlicodec"
)

func Get[V any](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
) (v V, err error) {
	req, err := NewGetRequest(c, ctx, rp, query, Method_get)
	if err != nil {
		return v, err
	}

	v, _, err = DoAndUnmarshal(c, req, restlicodec.UnmarshalRestLi[V])
	return v, err
}

// Update executes a rest.li update request with the given update object
func Update[V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	update V,
	query QueryParamsEncoder,
	createAndReadOnlyFields restlicodec.PathSpec,
) error {
	req, err := NewJsonRequest(c, ctx, rp, query, http.MethodPut, Method_update, update, createAndReadOnlyFields)
	if err != nil {
		return err
	}

	_, err = DoAndIgnore(c, req)
	return err
}

// PartialUpdate executes a rest.li partial update request with the given patch object
func PartialUpdate[PV restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	patch PV,
	query QueryParamsEncoder,
	createAndReadOnlyFields restlicodec.PathSpec,
) error {
	req, err := NewJsonRequest(c, ctx, rp, query, http.MethodPost, Method_partial_update, patch, createAndReadOnlyFields)
	if err != nil {
		return err
	}

	_, err = DoAndIgnore(c, req)
	return err
}

// PartialUpdateWithReturnEntity is like PartialUpdate, except it parses the returned entity from the response.
func PartialUpdateWithReturnEntity[PV, V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	patch PV,
	query QueryParamsEncoder,
	createAndReadOnlyFields restlicodec.PathSpec,
) (v V, err error) {
	req, err := NewJsonRequest(c, ctx, rp, query, http.MethodPost, Method_partial_update, patch, createAndReadOnlyFields)
	if err != nil {
		return v, err
	}

	v, _, err = DoAndUnmarshal(c, req, restlicodec.UnmarshalRestLi[V])
	return v, err
}

// Delete executes a rest.li delete request
func Delete(
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
) error {
	req, err := NewDeleteRequest(c, ctx, rp, query, Method_delete)
	if err != nil {
		return err
	}

	_, err = DoAndIgnore(c, req)
	return err
}
