package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const EnumShouldUsePointer = utils.No

type Enum struct {
	NamedType
	Symbols     []string
	SymbolToDoc map[string]string
}

func (e *Enum) InnerTypes() utils.IdentifierSet {
	return nil
}

func (e *Enum) ShouldReference() utils.ShouldUsePointer {
	return EnumShouldUsePointer
}

func (e *Enum) GenerateCode() (def *Statement) {
	def = Empty()
	utils.AddWordWrappedComment(def, e.Doc).Line()
	def.Type().Id(e.TypeName()).Int32().Line()

	unknownEnum := Code(Id("_" + e.SymbolIdentifier("unknown")))

	def.Const().DefsFunc(func(def *Group) {
		def.Add(unknownEnum).Op("=").Id(e.TypeName()).Call(Iota())
		for _, symbol := range e.Symbols {
			def.Add(utils.AddWordWrappedComment(Empty(), e.SymbolToDoc[symbol]))
			def.Id(e.SymbolIdentifier(symbol))
		}
	}).Line()

	values := "_" + e.TypeName() + "_values"
	def.Var().Id(values).Op("=").Map(String()).Id(e.TypeName()).Values(DictFunc(func(dict Dict) {
		for _, s := range e.Symbols {
			dict[Lit(s)] = Id(e.SymbolIdentifier(s))
		}
	})).Line()

	strings := "_" + e.TypeName() + "_strings"
	def.Var().Id(strings).Op("=").Map(Id(e.TypeName())).String().Values(DictFunc(func(dict Dict) {
		for _, s := range e.Symbols {
			dict[Id(e.SymbolIdentifier(s))] = Lit(s)
		}
	})).Line().Line()

	val, ok := Code(Id("val")), Code(Id("ok"))
	receiver := e.Receiver()
	getter := "Get" + e.TypeName() + "FromString"
	getEnumString := func(def *Group) {
		def.List(val, ok).Op(":=").Id(strings).Index(Id(receiver))
	}

	const isValid = "IsValid"
	utils.AddFuncOnReceiver(def, receiver, e.TypeName(), isValid, EnumShouldUsePointer).Params().Bool().BlockFunc(func(def *Group) {
		def.Return(Lit(0).Op("<").Id(receiver).Op("&&").Id(receiver).Op("<=").Lit(len(e.Symbols)))
	}).Line().Line()

	AddEquals(def, receiver, e.TypeName(), EnumShouldUsePointer, func(other Code, def *Group) {
		def.If(Id(receiver).Dot(isValid).Call().Op("&&").Add(other).Dot(isValid).Call()).
			Block(Return(Id(receiver).Op("==").Add(other))).
			Else().
			Block(Return(False()))
	})

	AddCustomComputeHash(def, receiver, e.TypeName(), EnumShouldUsePointer, func(def *Group) {
		def.If(Id(receiver).Dot(isValid).Call()).
			Block(Return(Qual(utils.HashPackage, "HashInt32").Call(Int32().Call(Id(receiver))))).
			Else().
			Block(Return(utils.ZeroHash))
	})

	def.Func().Id("All" + e.TypeName() + "Values").Params().Index().Id(e.TypeName()).BlockFunc(func(def *Group) {
		def.Return(Index().Id(e.TypeName()).ValuesFunc(func(def *Group) {
			for _, s := range e.Symbols {
				def.Id(e.SymbolIdentifier(s))
			}
		}))
	}).Line().Line()

	def.Func().Id(getter).Params(Add(val).String()).Params(Id(receiver).Id(e.TypeName()), Err().Error()).
		BlockFunc(func(def *Group) {
			def.List(Id(receiver), ok).Op(":=").Id(values).Index(val)
			def.If(Op("!").Add(ok)).BlockFunc(func(def *Group) {
				def.Err().Op("=").Op("&").Add(utils.UnknownEnumValue).Add(utils.OrderedValues(
					func(add func(key Code, value Code)) {
						add(Id("Enum"), Lit(e.Identifier.FullName()))
						add(Id("Value"), val)
					}))
			})
			def.Return(Id(receiver), Err())
		}).Line().Line()

	utils.AddStringer(def, receiver, e.TypeName(), EnumShouldUsePointer, func(def *Group) {
		getEnumString(def)
		def.If(Op("!").Add(ok)).Block(
			Return(Lit("$UNKNOWN$")),
		)
		def.Return(val)
	}).Line().Line()

	utils.AddPointer(def, receiver, e.TypeName())

	AddMarshalRestLi(def, receiver, e.TypeName(), EnumShouldUsePointer, func(def *Group) {
		getEnumString(def)
		def.If(Op("!").Add(ok)).Block(
			Return(Op("&").Add(utils.IllegalEnumConstant).Add(utils.OrderedValues(func(add func(key Code, value Code)) {
				add(Id("Enum"), Lit(e.Identifier.FullName()))
				add(Id("Constant"), Int().Call(Id(receiver)))
			}))))
		def.List(Writer).Dot("WriteString").Call(val)
		def.Return(Nil())
	})

	AddUnmarshalRestli(def, receiver, e.TypeName(), EnumShouldUsePointer, func(def *Group) {
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
	return utils.ExportedIdentifier(e.TypeName() + "_" + symbol)
}

func (e *Enum) isValidSymbol(v string) bool {
	for _, symbol := range e.Symbols {
		if symbol == v {
			return true
		}
	}
	return false
}
