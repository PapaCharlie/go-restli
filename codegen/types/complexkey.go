package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
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

func (ck *ComplexKey) ShouldReference() utils.ShouldUsePointer {
	return RecordShouldUsePointer
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
	}
	for _, f := range utils.TypeRegistry.Resolve(ck.Key).(*Record).Fields {
		f.IncludedFrom = &ck.Key
		record.Fields = append(record.Fields, f)
	}

	def := Empty().
		Add(record.GenerateStruct()).Line().Line().
		Add(record.GenerateEquals()).Line().Line()

	receiver := record.Receiver()
	other := Code(Id("other"))
	utils.AddFuncOnReceiver(def, record.Receiver(), record.Name, "ComplexKeyEquals", RecordShouldUsePointer).
		Params(Add(other).Op("*").Add(record.Qual())).
		Bool().
		BlockFunc(func(def *Group) {
			def.Return(Id(receiver).Add(ck.KeyAccessor())).Dot(utils.Equals).Call(Op("&").Add(other).Add(ck.KeyAccessor()))
		}).Line().Line()

	utils.AddFuncOnReceiver(def, record.Receiver(), record.Name, "ComputeComplexKeyHash", RecordShouldUsePointer).
		Params().
		Add(utils.Hash).
		BlockFunc(func(def *Group) {
			def.Return(Id(receiver).Add(ck.KeyAccessor())).Dot(utils.ComputeHash).Call()
		}).Line().Line()

	def.Add(record.GenerateComputeHash()).Line().Line().
		Add(record.GenerateMarshalRestLi()).Line().Line().
		Add(record.GenerateUnmarshalRestLi()).Line().Line()

	return def
}

func (ck *ComplexKey) KeyAccessor() Code {
	return Dot(ck.Key.Name)
}
