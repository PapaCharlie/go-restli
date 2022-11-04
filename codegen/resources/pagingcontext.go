package resources

import (
	"github.com/PapaCharlie/go-restli/v2/codegen/types"
	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
)

var PagingContext = &types.Record{
	NamedType: types.NamedType{Identifier: utils.PagingContextIdentifier},
	Fields: []types.Field{
		{
			Type:       types.RestliType{Primitive: &types.Int32Primitive},
			Name:       "start",
			IsOptional: true,
		},
		{
			Type:       types.RestliType{Primitive: &types.Int32Primitive},
			Name:       "count",
			IsOptional: true,
		},
	},
}

func init() {
	utils.TypeRegistry.Register(PagingContext, utils.RootPackage)
}
