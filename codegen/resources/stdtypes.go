package resources

import (
	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/PapaCharlie/go-restli/protocol"
)

// ErrorResponse is manually parsed from https://github.com/linkedin/rest.li/blob/master/restli-common/src/main/pegasus/com/linkedin/restli/common/ErrorResponse.pdl
var ErrorResponse = &types.Record{
	NamedType: types.NamedType{Identifier: utils.Identifier{
		Namespace: utils.StdTypesPackage,
		Name:      "ErrorResponse",
	}},
	Fields: []types.Field{
		{
			Type:       types.RestliType{Primitive: &types.Int32Primitive},
			Name:       "status",
			Doc:        "The HTTP status code.",
			IsOptional: true,
		},
		{
			Type:       types.RestliType{Primitive: &types.StringPrimitive},
			Name:       "message",
			Doc:        "A human-readable explanation of the error.",
			IsOptional: true,
		},
		{
			Type:       types.RestliType{Primitive: &types.StringPrimitive},
			Name:       "exceptionClass",
			Doc:        "The FQCN of the exception thrown by the server.",
			IsOptional: true,
		},
		{
			Type:       types.RestliType{Primitive: &types.StringPrimitive},
			Name:       "stackTrace",
			Doc:        "The full stack trace of the exception thrown by the server.",
			IsOptional: true,
		},
	},
}

// CollectionMetadata is manually parsed from https://github.com/linkedin/rest.li/blob/master/restli-common/src/main/pegasus/com/linkedin/restli/common/CollectionMetadata.pdl
var CollectionMetadata = &types.Record{
	NamedType: types.NamedType{Identifier: utils.Identifier{
		Namespace: utils.StdTypesPackage,
		Name:      "CollectionMedata",
	}},
	Fields: []types.Field{
		{
			Type: types.RestliType{Primitive: &types.Int32Primitive},
			Name: "start",
			Doc:  "The start index of this collection",
		},
		{
			Type: types.RestliType{Primitive: &types.Int32Primitive},
			Name: "count",
			Doc:  "The number of elements in this collection segment",
		},
		{
			Type:         types.RestliType{Primitive: &types.Int32Primitive},
			Name:         "total",
			Doc:          "The total number of elements in the entire collection (not just this segment)",
			DefaultValue: protocol.StringPointer("0"),
		},
		{
			Type: types.RestliType{Array: &types.RestliType{Reference: &Link.Identifier}},
			Name: "links",
		},
	},
}

// Link is manually parsed from https://github.com/linkedin/rest.li/blob/master/restli-common/src/main/pegasus/com/linkedin/restli/common/Link.pdl
var Link = &types.Record{
	NamedType: types.NamedType{Identifier: utils.Identifier{
		Namespace: utils.StdTypesPackage,
		Name:      "Link",
	}},
	Fields: []types.Field{
		{
			Type: types.RestliType{Primitive: &types.StringPrimitive},
			Name: "rel",
			Doc:  "The link relation e.g. 'self' or 'next'",
		},
		{
			Type: types.RestliType{Primitive: &types.StringPrimitive},
			Name: "href",
			Doc:  "The link URI",
		},
		{
			Type: types.RestliType{Primitive: &types.StringPrimitive},
			Name: "type",
			Doc:  "The type (media type) of the resource",
		},
	},
}

func init() {
	utils.TypeRegistry.Register(ErrorResponse)
	utils.TypeRegistry.Register(CollectionMetadata)
	utils.TypeRegistry.Register(Link)
}
