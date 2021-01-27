package types

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type ComplexKey struct {
	NamedType
	Key    utils.Identifier
	Params utils.Identifier
}

func (ck *ComplexKey) InnerTypes() utils.IdentifierSet {
	return utils.NewIdentifierSet(ck.Key, ck.Params)
}

func (ck *ComplexKey) GenerateCode() *Statement {
	record := &Record{
		NamedType: ck.NamedType,
		Fields: []Field{
			{
				Name:               utils.ComplexKeyParams,
				IsOptional:         true,
				Type:               RestliType{Reference: &ck.Params},
				isComplexKeyParams: true,
			},
		},
		IncludedRecords: []utils.Identifier{ck.Key},
	}
	for _, f := range utils.TypeRegistry.Resolve(ck.Key).(*Record).Fields {
		f.IncludedFrom = &ck.Key
		record.Fields = append(record.Fields, f)
	}

	return Empty().
		Add(record.GenerateStruct()).Line().Line().
		Add(record.GenerateEquals()).Line().Line().
		Add(record.GenerateComputeHash()).Line().Line().
		Add(record.GenerateMarshalRestLi()).Line().Line().
		Add(record.GenerateUnmarshalRestLi()).Line().Line()
}

func (ck *ComplexKey) KeyAccessor() Code {
	return Dot(ck.Key.Name)
}
