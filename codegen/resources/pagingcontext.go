package resources

import (
	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
)

var PagingContextIdentifier = utils.Identifier{
	IsNativeIdentifier: true,
	Name:               "PagingContext",
	Namespace:          utils.ProtocolPackage,
}

var PagingContext = &types.Record{
	NamedType: types.NamedType{Identifier: PagingContextIdentifier},
	Fields: []types.Field{
		{
			Type:         types.RestliType{Primitive: &types.Int32Primitive},
			Name:         "start",
			IsOptional:   true,
			IncludedFrom: &PagingContextIdentifier,
		},
		{
			Type:         types.RestliType{Primitive: &types.Int32Primitive},
			Name:         "count",
			IsOptional:   true,
			IncludedFrom: &PagingContextIdentifier,
		},
	},
}

func init() {
	utils.TypeRegistry.Register(PagingContext)
}

func addPagingContextFields(record *types.Record) {
	record.IncludedRecords = append([]utils.Identifier{PagingContextIdentifier}, record.IncludedRecords...)
	record.Fields = append(record.Fields, PagingContext.Fields...)
}
