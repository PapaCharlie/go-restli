package types

import (
	"encoding/json"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type ObjectCodeGenerator struct {
	identifier utils.Identifier
	pointer    utils.ShouldUsePointer
	receiver   string
}

func NewObjectCodeGenerator(identifier utils.Identifier, pointer utils.ShouldUsePointer) *ObjectCodeGenerator {
	return &ObjectCodeGenerator{
		identifier: identifier,
		pointer:    pointer,
		receiver:   utils.ReceiverName(identifier.Name),
	}
}

func NewObjectCodeGeneratorWithCustomReceiver(
	identifier utils.Identifier,
	pointer utils.ShouldUsePointer,
	receiver string,
) *ObjectCodeGenerator {
	return &ObjectCodeGenerator{
		identifier: identifier,
		pointer:    pointer,
		receiver:   receiver,
	}
}

func (o *ObjectCodeGenerator) ReceiverName() string {
	return o.receiver
}

func (o *ObjectCodeGenerator) Receiver() *Statement {
	return Id(o.receiver)
}

type GeneratedTypeManifest struct {
	utils.Identifier
	Package          string  `json:"package"`
	NativeIdentifier *string `json:"nativeIdentifier,omitempty"`
}

func (m *GeneratedTypeManifest) ExternalIdentifier() utils.Identifier {
	id := utils.Identifier{Namespace: m.Package}
	if m.NativeIdentifier != nil {
		id.Name = *m.NativeIdentifier
	} else {
		id.Name = m.Name
	}
	return id
}

func (m *GeneratedTypeManifest) AddTypeManifest(def *Statement) Code {
	manifest, _ := json.MarshalIndent(m, "", "  ")
	return def.Comment("go-restli:generated " + string(manifest)).Line().Line()
}

func AddTypeManifest(def *Statement, internal utils.Identifier, pkg string) Code {
	manifest := &GeneratedTypeManifest{
		Identifier: internal,
		Package:    pkg,
	}
	return manifest.AddTypeManifest(def)
}

func (o *ObjectCodeGenerator) DeclareType(def *Statement, doc string, t Code) Code {
	return DeclareType(def, o.identifier.Name, doc, t)
}

func DeclareType(def *Statement, name, doc string, t Code) Code {
	utils.AddWordWrappedComment(def, doc).Line()
	return def.Type().Id(name).Add(t).Line()
}

func (o *ObjectCodeGenerator) Equals(def *Statement, equals func(receiver, other Code, def *Group)) Code {
	return AddEquals(def, o.receiver, o.identifier.Name, o.pointer, equals)
}

func AddEquals(def *Statement, receiver, typeName string, pointer utils.ShouldUsePointer, f func(receiver, other Code, def *Group)) *Statement {
	other := Id("other")
	rightHandType := Id(typeName)
	if pointer.ShouldUsePointer() {
		rightHandType = Op("*").Add(rightHandType)
	}
	r := Id(receiver)

	return utils.AddFuncOnReceiver(def, receiver, typeName, utils.Equals, pointer).
		Params(Add(other).Add(rightHandType)).Bool().
		BlockFunc(func(def *Group) {
			if pointer.ShouldUsePointer() {
				def.If(Add(r).Op("==").Add(other)).Block(Return(True()))
				def.If(Add(r).Op("==").Nil().Op("||").Add(other).Op("==").Nil()).Block(Return(False())).Line()
			}
			f(r, other, def)
		}).Line().Line()
}

func (o *ObjectCodeGenerator) ComputeHash(def *Statement, hash func(receiver, h Code, def *Group)) Code {
	return AddComputeHash(def, o.receiver, o.identifier.Name, o.pointer, hash)
}

func AddComputeHash(def *Statement, receiver, typeName string, pointer utils.ShouldUsePointer, f func(receiver, h Code, def *Group)) *Statement {
	return AddCustomComputeHash(def, receiver, typeName, pointer, func(receiver Code, def *Group) {
		if pointer.ShouldUsePointer() {
			def.Add(If(Add(receiver).Op("==").Nil()).Block(Return(utils.ZeroHash)))
		}
		h := Id("hash")
		def.Add(h).Op(":=").Add(utils.NewHash).Line()
		f(receiver, h, def)
		def.Return(h)
	})
}

func (o *ObjectCodeGenerator) CustomComputeHash(def *Statement, hash func(receiver Code, def *Group)) Code {
	return AddCustomComputeHash(def, o.receiver, o.identifier.Name, o.pointer, hash)
}

func AddCustomComputeHash(def *Statement, receiver, typeName string, pointer utils.ShouldUsePointer, f func(receiver Code, def *Group)) *Statement {
	return utils.AddFuncOnReceiver(def, receiver, typeName, utils.ComputeHash, pointer).
		Params().Params(utils.Hash).
		BlockFunc(func(def *Group) { f(Id(receiver), def) }).Line().Line()
}

func (o *ObjectCodeGenerator) MarshalRestLi(def *Statement, marshal func(receiver, writer Code, def *Group)) Code {
	return AddMarshalRestLi(def, o.receiver, o.identifier.Name, o.pointer, marshal)
}

func AddMarshalRestLi(def *Statement, receiver, typeName string, pointer utils.ShouldUsePointer, f func(receiver, writer Code, def *Group)) *Statement {
	r := Id(receiver)
	utils.AddFuncOnReceiver(def, receiver, typeName, utils.MarshalRestLi, pointer).
		Params(WriterParam).
		Params(Err().Error()).
		BlockFunc(func(def *Group) { f(r, Writer, def) }).
		Line().Line()

	utils.AddFuncOnReceiver(def, receiver, typeName, "MarshalJSON", pointer).
		Params().
		Params(Id("data").Index().Byte(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(Writer).Op(":=").Qual(utils.RestLiCodecPackage, "NewCompactJsonWriter").Call()
			def.Err().Op("=").Add(r).Dot(utils.MarshalRestLi).Call(Writer)
			def.Add(utils.IfErrReturn(Nil(), Err()))
			def.Return(Index().Byte().Call(Add(WriterUtils.Finalize(Writer))), Nil())
		}).Line().Line()

	return def
}

func (o *ObjectCodeGenerator) UnmarshalRestLi(def *Statement, unmarshal func(receiver, reader Code, def *Group)) Code {
	AddUnmarshalerFunc(def, o.receiver, o.identifier, o.pointer)
	return AddUnmarshalRestLi(def, o.receiver, o.identifier.Name, unmarshal)
}

func AddUnmarshalRestLi(def *Statement, receiver, typeName string, f func(receiver, reader Code, def *Group)) *Statement {
	r := Id(receiver)
	utils.AddFuncOnReceiver(def, receiver, typeName, utils.UnmarshalRestLi, utils.Yes).
		Params(ReaderParam).
		Params(Err().Error()).
		BlockFunc(func(def *Group) { f(r, Reader, def) }).
		Line().Line()

	data := Id("data")
	utils.AddFuncOnReceiver(def, receiver, typeName, "UnmarshalJSON", utils.Yes).
		Params(Add(data).Index().Byte()).
		Params(Error()).
		BlockFunc(func(def *Group) {
			def.Add(Reader).Op(":=").Add(utils.NewJsonReader).Call(data)
			def.Return(Add(r).Dot(utils.UnmarshalRestLi).Call(Reader))
		}).Line().Line()

	return def
}

func AddUnmarshalerFunc(def *Statement, receiver string, id utils.Identifier, pointer utils.ShouldUsePointer) *Statement {
	r := Code(Id(receiver))
	param := Add(r)
	if pointer.ShouldUsePointer() {
		param.Op("*")
	}

	return def.Func().Id(id.UnmarshalerFuncName()).
		Params(ReaderParam).
		Params(Add(param).Add(id.Qual()), Err().Error()).
		BlockFunc(func(def *Group) {
			if pointer.ShouldUsePointer() {
				def.Add(Add(r).Op("=").New(id.Qual()))
			}
			def.Add(ReaderUtils.Read(RestliType{Reference: &id}, Reader, r))
			def.Return(r, Err())
		}).Line().Line()
}
