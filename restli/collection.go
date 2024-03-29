package restli

import (
	"context"
	"net/http"

	"github.com/PapaCharlie/go-restli/restlicodec"
	"github.com/PapaCharlie/go-restli/restlidata"
)

// Create executes a rest.li create request with the given object. The X-RestLi-Id header field will be parsed into id
// (though a CreateResponseHasNoEntityHeaderError will be returned if the header is not set). The response body will
// always be ignored.
func Create[K any, V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	create V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (*restlidata.CreatedEntity[K], error) {
	req, err := NewCreateRequest(c, ctx, rp, query, Method_create, create, readOnlyFields)
	if err != nil {
		return nil, err
	}

	res, err := DoAndIgnore(c, req)
	if err != nil {
		return nil, err
	}

	return unmarshalReturnEntityKey[K](c, res)
}

// CreateWithReturnEntity is like CollectionClient.Create, except it parses the returned entity from the response.
func CreateWithReturnEntity[K any, V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	create V,
	query QueryParamsEncoder,
	readOnlyFields restlicodec.PathSpec,
) (*restlidata.CreatedAndReturnedEntity[K, V], error) {
	req, err := NewCreateRequest(c, ctx, rp, query, Method_create, create, readOnlyFields)
	if err != nil {
		return nil, err
	}

	v, res, err := DoAndUnmarshal(c, req, restlicodec.UnmarshalRestLi[V])
	if err != nil {
		return nil, err
	}

	k, err := unmarshalReturnEntityKey[K](c, res)
	if err != nil {
		return nil, err
	}
	return &restlidata.CreatedAndReturnedEntity[K, V]{
		CreatedEntity: *k,
		Entity:        v,
	}, nil
}

func unmarshalReturnEntityKey[K any](c *Client, res *http.Response) (result *restlidata.CreatedEntity[K], err error) {
	if h := res.Header.Get(IDHeader); len(h) > 0 {
		var reader restlicodec.Reader
		reader, err = restlicodec.NewRor2Reader(h)
		if err != nil {
			return nil, err
		}

		var k K
		k, err = restlicodec.UnmarshalRestLi[K](reader)
		if _, mfe := err.(*restlicodec.MissingRequiredFieldsError); mfe && !c.StrictResponseDeserialization {
			err = nil
		}
		if err != nil {
			return nil, err
		}
		return &restlidata.CreatedEntity[K]{
			Id:     k,
			Status: res.StatusCode,
		}, nil
	} else {
		return nil, &CreateResponseHasNoEntityHeaderError{Response: res}
	}
}

func GetAll[V restlicodec.Marshaler](
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
) (results *restlidata.Elements[V], err error) {
	req, err := NewGetRequest(c, ctx, rp, query, Method_get_all)
	if err != nil {
		return nil, err
	}

	results, _, err = DoAndUnmarshal(c, req, restlicodec.UnmarshalRestLi[*restlidata.Elements[V]])
	return results, err
}
