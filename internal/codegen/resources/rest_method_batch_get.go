package resources

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/types"
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func (r *RestMethod) generateBatchGet(def *Group) {
	returns := []Code{Nil(), Err()}
	key := Code(Id("key"))

	formatQueryUrl(r, def, func(itemWriter Code, def *Group) {
		def.For(List(Id("_"), key).Op(":=").Range().Add(Keys)).BlockFunc(func(def *Group) {
			def.Add(types.Writer.Write(r.EntityPathKey.Type, Add(itemWriter).Call(), key, Err()))
		})
		def.Return(Nil())
	}, returns...)

	ck := r.EntityPathKey.Type.ComplexKey()
	originalKeys := Id("originalKeys")
	if ck != nil {
		def.Add(originalKeys).Op(":=").Make(Map(types.Hash).Index().Add(r.EntityPathKey.Type.ReferencedType()))
		def.For().List(Id("_"), key).Op(":=").Range().Add(Keys).BlockFunc(func(def *Group) {
			keyHash := Id("keyHash")
			def.Add(keyHash).Op(":=").Add(key).Add(ck.KeyAccessor()).Dot(types.ComputeHash).Call()
			index := Add(originalKeys).Index(keyHash)
			def.Add(index).Op("=").Append(index, key)
		}).Line()
	}

	def.Add(Entities).Op("=").Make(r.batchGetReturnType())
	rawKey := Id("rawKey")
	resultsReader := Func().Params(Add(types.Reader).Add(types.ReaderQual), Add(rawKey).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		v := Code(Id("v"))
		if ck != nil {
			def.Add(v).Op(":=").New(r.EntityPathKey.Type.GoType())
		} else {
			def.Var().Add(v).Add(r.EntityPathKey.Type.GoType())
		}
		keyReader := Id("keyReader")
		def.List(keyReader, Err()).Op(":=").Add(types.NewRor2Reader).Call(rawKey)
		def.Add(utils.IfErrReturn(Err()))
		def.Add(types.Reader.Read(r.EntityPathKey.Type, keyReader, v))
		def.Add(utils.IfErrReturn(Err())).Line()

		if ck != nil {
			originalKey := Code(Id("originalKey"))
			def.Var().Add(originalKey).Add(r.EntityPathKey.Type.ReferencedType())
			def.For().List(Id("_"), key).Op(":=").Range().Add(originalKeys).Index(Add(v).Add(ck.KeyAccessor()).Dot(types.ComputeHash).Call()).BlockFunc(func(def *Group) {
				def.If(Add(v).Add(ck.KeyAccessor()).Dot(types.Equals).Call(Op("&").Add(key).Add(ck.KeyAccessor()))).Block(
					Add(originalKey).Op("=").Add(key),
					Break(),
				)
			})
			def.If(Add(originalKey).Op("==").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("unknown key returned by batch get: %q"), rawKey)),
			)
			def.Line()
			v = originalKey
		}

		accessor := Add(Entities).Index(v)
		if r.Return.ShouldReference() {
			def.Add(accessor).Op("=").New(r.Return.GoType())
		}
		def.Add(types.Reader.Read(*r.Return, types.Reader, accessor))
		def.Add(Return(Err()))
	})
	def.Err().Op("=").Id(ClientReceiver).Dot("DoBatchGetRequest").Call(Ctx, Url, resultsReader)

	def.Add(utils.IfErrReturn(returns...)).Line()

	def.Return()
}
