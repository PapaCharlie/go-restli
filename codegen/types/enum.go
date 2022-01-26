package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type Enum struct {
	NamedType
	Symbols     []string
	SymbolToDoc map[string]string
}

func (e *Enum) InnerTypes() utils.IdentifierSet {
	return nil
}

func (e *Enum) GenerateCode() (def *Statement) {
	def = Empty()
	utils.AddWordWrappedComment(def, e.Doc).Line()
	def.Type().Id(e.Name).Int().Line()

	unknownEnum := Id("_" + e.SymbolIdentifier("unknown"))

	def.Const().DefsFunc(func(def *Group) {
		def.Add(unknownEnum).Op("=").Id(e.Name).Call(Iota())
		for _, symbol := range e.Symbols {
			def.Add(utils.AddWordWrappedComment(Empty(), e.SymbolToDoc[symbol]))
			def.Id(e.SymbolIdentifier(symbol))
		}
	}).Line()

	values := "_" + e.Name + "_values"
	def.Var().Id(values).Op("=").Map(String()).Id(e.Name).Values(DictFunc(func(dict Dict) {
		for _, s := range e.Symbols {
			dict[Lit(s)] = Id(e.SymbolIdentifier(s))
		}
	})).Line()

	strings := "_" + e.Name + "_strings"
	def.Var().Id(strings).Op("=").Map(Id(e.Name)).String().Values(DictFunc(func(dict Dict) {
		for _, s := range e.Symbols {
			dict[Id(e.SymbolIdentifier(s))] = Lit(s)
		}
	})).Line().Line()

	val, ok := Code(Id("val")), Code(Id("ok"))
	receiver := utils.ReceiverName(e.Name)
	getter := "Get" + e.Name + "FromString"
	accessor := Op("*").Id(receiver)
	getEnumString := func(def *Group) {
		def.List(val, ok).Op(":=").Id(strings).Index(accessor)
	}

	AddEquals(def, receiver, e.Name, func(other Code, def *Group) {
		def.Return(Op("*").Id(receiver).Op("== *").Add(other))
	})

	AddComputeHash(def, receiver, e.Name, func(h Code, def *Group) {
		def.Add(h).Op("=").Add(utils.Hash).Call(Op("*").Id(receiver))
		def.Return(h)
	})

	def.Func().Id("All" + e.Name + "Values").Params().Index().Id(e.Name).BlockFunc(func(def *Group) {
		def.Return(Index().Id(e.Name).ValuesFunc(func(def *Group) {
			for _, s := range e.Symbols {
				def.Id(e.SymbolIdentifier(s))
			}
		}))
	}).Line().Line()

	def.Func().Id(getter).Params(Add(val).String()).Params(Id(receiver).Id(e.Name), Err().Error()).
		BlockFunc(func(def *Group) {
			def.List(Id(receiver), ok).Op(":=").Id(values).Index(val)
			def.If(Op("!").Add(ok)).BlockFunc(func(def *Group) {
				def.Err().Op("=").Op("&").Add(utils.UnknownEnumValue).Values(Dict{
					Id("Enum"):  Lit(e.Identifier.String()),
					Id("Value"): val,
				})
			})
			def.Return(Id(receiver), Err())
		}).Line().Line()

	utils.AddStringer(def, receiver, e.Name, func(def *Group) {
		getEnumString(def)
		def.If(Op("!").Add(ok)).Block(
			Return(Lit("$UNKNOWN$")),
		)
		def.Return(val)
	}).Line().Line()

	utils.AddPointer(def, receiver, e.Name)

	AddMarshalRestLi(def, receiver, e.Name, func(def *Group) {
		getEnumString(def)
		def.If(Op("!").Add(ok)).Block(
			Return(Op("&").Add(utils.IllegalEnumConstant).Values(Dict{
				Id("Enum"):     Lit(e.Identifier.String()),
				Id("Constant"): Int().Call(accessor),
			})),
		)
		def.List(Writer).Dot("WriteString").Call(val)
		def.Return(Nil())
	})
	AddUnmarshalRestli(def, receiver, e.Name, func(def *Group) {
		value := Id("value")
		def.Var().Add(value).String()
		def.Add(Reader.Read(RestliType{Primitive: &StringPrimitive}, Reader, value))
		def.Add(utils.IfErrReturn(Err()))
		def.Line()

		def.Op("*").Id(receiver).Op("=").Id(values).Index(Id("value"))
		def.Return(Nil())
	})

	return def
}

func (e *Enum) SymbolIdentifier(symbol string) string {
	return utils.ExportedIdentifier(e.Name + "_" + symbol)
}

func (e *Enum) zeroValueLit() *Statement {
	return e.Qual().Call(Lit(0))
}

func (e *Enum) isValidSymbol(v string) bool {
	for _, symbol := range e.Symbols {
		if symbol == v {
			return true
		}
	}
	return false
}
