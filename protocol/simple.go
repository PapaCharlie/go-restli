package protocol

import (
	"context"
	"net/http"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

type SimpleClient[V RestLiObject[V], PV restlicodec.Marshaler] struct {
	*RestLiClient
	EntityUnmarshaler restlicodec.GenericUnmarshaler[V]

	CreateAndReadOnlyFields restlicodec.PathSpec
}

func (c *SimpleClient[V, PV]) NewJsonRequest(
	ctx context.Context,
	rp ResourcePath,
	query QueryParams,
	httpMethod string,
	restLiMethod RestLiMethod,
	contents restlicodec.Marshaler,
	excludedFields restlicodec.PathSpec,
) (*http.Request, error) {
	u, err := c.FormatQueryUrl(rp, query)
	if err != nil {
		return nil, err
	}

	return NewJsonRequest(ctx, u, httpMethod, restLiMethod, contents, excludedFields)
}

func (c *SimpleClient[V, PV]) NewGetRequest(
	ctx context.Context,
	rp ResourcePath,
	query QueryParams,
	method RestLiMethod,
) (*http.Request, error) {
	u, err := c.FormatQueryUrl(rp, query)
	if err != nil {
		return nil, err
	}

	return NewGetRequest(ctx, u, method)
}

func (c *SimpleClient[V, PV]) Get(ctx context.Context, rp ResourcePath, query QueryParams) (v V, err error) {
	req, err := c.NewGetRequest(ctx, rp, query, Method_get)
	if err != nil {
		return v, err
	}

	v, _, err = DoAndUnmarshal(c.RestLiClient, req, c.EntityUnmarshaler)
	return v, err
}

// Update executes a rest.li update request with the given update object
func (c *SimpleClient[V, PV]) Update(
	ctx context.Context,
	rp ResourcePath,
	update V,
	query QueryParams,
) error {
	req, err := c.NewJsonRequest(ctx, rp, query, http.MethodPut, Method_update, update, nil)
	if err != nil {
		return err
	}

	_, err = c.RestLiClient.DoAndIgnore(req)
	return err
}

// PartialUpdate executes a rest.li partial update request with the given patch object
func (c *SimpleClient[V, PV]) PartialUpdate(
	ctx context.Context,
	rp ResourcePath,
	patch PV,
	query QueryParams,
) error {
	req, err := c.NewJsonRequest(ctx, rp, query, http.MethodPost, Method_partial_update,
		restlicodec.MarshalerFunc(func(writer restlicodec.Writer) error {
			return writer.WriteMap(func(fieldNameWriter func(fieldName string) restlicodec.Writer) error {
				return patch.MarshalRestLi(fieldNameWriter("patch").SetScope())
			})
		}), c.CreateAndReadOnlyFields)
	if err != nil {
		return err
	}

	_, err = c.RestLiClient.DoAndIgnore(req)
	return err
}

func (c *SimpleClient[V, PV]) NewDeleteRequest(
	ctx context.Context,
	rp ResourcePath,
	query QueryParams,
	method RestLiMethod,
) (*http.Request, error) {
	u, err := c.FormatQueryUrl(rp, query)
	if err != nil {
		return nil, err
	}

	return NewDeleteRequest(ctx, u, method)
}

// Delete executes a rest.li delete request
func (c *SimpleClient[V, PV]) Delete(ctx context.Context, rp ResourcePath, query QueryParams) error {
	req, err := c.NewDeleteRequest(ctx, rp, query, Method_delete)
	if err != nil {
		return err
	}

	_, err = c.RestLiClient.DoAndIgnore(req)
	return err
}
