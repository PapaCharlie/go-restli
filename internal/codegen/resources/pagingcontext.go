package resources

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/types"
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
)

var PagingContext = &types.Record{
	NamedType: types.NamedType{Identifier: utils.PagingContextIdentifier},
	Fields: []types.Field{
		{
			Type:         types.RestliType{Primitive: &types.Int32Primitive},
			Name:         "start",
			IsOptional:   true,
			IncludedFrom: &utils.PagingContextIdentifier,
		},
		{
			Type:         types.RestliType{Primitive: &types.Int32Primitive},
			Name:         "count",
			IsOptional:   true,
			IncludedFrom: &utils.PagingContextIdentifier,
		},
	},
}

func init() {
	utils.TypeRegistry.Register(PagingContext)
}

func addPagingContextFields(record *types.Record) {
	record.IncludedRecords = append([]utils.Identifier{utils.PagingContextIdentifier}, record.IncludedRecords...)
	record.Fields = append(record.Fields, PagingContext.Fields...)
}
