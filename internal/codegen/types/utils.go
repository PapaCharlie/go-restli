package types

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddUnmarshalRestli(def *Statement, receiver, typeName string, f func(def *Group)) *Statement {
	utils.AddFuncOnReceiver(def, receiver, typeName, UnmarshalRestLi).
		Params(Add(Reader).Add(ReaderQual)).
		Params(Err().Error()).
		BlockFunc(f).
		Line().Line()

	data := Id("data")
	utils.AddFuncOnReceiver(def, receiver, typeName, UnmarshalJSON).
		Params(Add(data).Index().Byte()).
		Params(Error()).
		BlockFunc(func(def *Group) {
			def.Add(Reader).Op(":=").Qual(RestLiCodecPackage, "NewJsonReader").Call(data)
			def.Return(Id(receiver).Dot(UnmarshalRestLi).Call(Reader))
		})

	return def
}
