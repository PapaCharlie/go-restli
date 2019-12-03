package internal

import (
	"encoding/json"
	"log"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

const TyperefModelTypeName = "typeref"

type TyperefModel struct {
	Identifier
	Doc string
	Ref *Model
}

func (r *TyperefModel) CopyWithAlias(alias string) ComplexType {
	rCopy := *r
	rCopy.Name = alias
	return &rCopy
}

func (r *TyperefModel) UnmarshalJSON(data []byte) error {
	t := &struct {
		typeField
		docField
		Identifier
		Ref *Model `json:"ref"`
	}{}
	t.Namespace = currentNamespace // default to the current namespace if none is specified
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	if t.Type != TyperefModelTypeName {
		return &WrongTypeError{Expected: TyperefModelTypeName, Actual: t.Type}
	}
	r.Identifier = t.Identifier
	r.Identifier.Name = ExportedIdentifier(r.Identifier.Name)
	r.Doc = t.Doc
	r.Ref = t.Ref
	return nil
}

func (r *TyperefModel) innerModels() []*Model {
	return []*Model{r.Ref}
}

func (r *TyperefModel) GenerateCode() (def *Statement) {
	def = Empty()

	if ref := r.Ref.ComplexType; ref != nil {
		// TODO
		log.Printf("Warning: type references to non-primitive types are not yet supported (%s)", r.Identifier)
		return def
	}

	AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(r.Name).Add(r.Ref.GoType()).Line().Line()

	if r.Ref.IsBytesOrPrimitive() {
		var accessor *Statement
		var encoder func(*Statement) *Statement
		var decoder func(*Statement) *Statement

		if bytes, ok := r.Ref.BuiltinType.(*BytesModel); ok {
			accessor = Bytes().Call(Op("*").Id(r.receiver()))
			encoder = bytes.encode
			decoder = bytes.decode
		}
		if primitive, ok := r.Ref.BuiltinType.(*PrimitiveModel); ok {
			accessor = Id(primitive[1]).Call(Op("*").Id(r.receiver()))
			encoder = primitive.encode
			decoder = primitive.decode
		}

		AddRestLiEncode(def, r.receiver(), r.Name, func(def *Group) {
			def.Return(encoder(accessor), Nil())
		}).Line().Line()
		AddRestLiDecode(def, r.receiver(), r.Name, func(def *Group) {
			def.Return(decoder(Id(r.receiver())))
		}).Line().Line()

		return def
	}

	if union, ok := r.Ref.BuiltinType.(*UnionModel); ok {
		AddRestLiEncode(def, r.receiver(), r.Name, func(def *Group) {
			def.Err().Op("=").Id(r.receiver()).Dot(ValidateUnionFields).Call()
			def.If(Err().Op("!=").Nil()).Block(Return()).Line()
			def.Var().Id("buf").Qual("strings", "Builder")
			union.restLiWriteToBuf(def, Id(r.receiver()))
			def.Id("data").Op("=").Id("buf").Dot("String").Call()
			def.Return()
		}).Line().Line()

		AddFuncOnReceiver(def, r.receiver(), r.Name, ValidateUnionFields).
			Params().
			Params(Err().Error()).
			BlockFunc(func(def *Group) {
				union.validateUnionFields(def, Id(r.receiver()))
				def.Line().Return()
			})

		return def
	}

	log.Panicf("Illegal typeref type %s defined in %s", r.Ref, ModelRegistry.GetSourceFileFilename(r.Identifier))
	return nil
}
