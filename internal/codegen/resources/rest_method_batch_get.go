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
	isComplexKey := ck != nil ||
		(r.EntityPathKey.Type.UnderlyingPrimitive() == nil && r.EntityPathKey.Type.NativeTyperef == nil)
	keyAccessor := func(accessor Code) *Statement {
		if ck != nil {
			return Add(accessor).Add(ck.KeyAccessor())
		} else {
			return Add(accessor)
		}
	}

	originalKeys := Id("originalKeys")
	if isComplexKey {
		def.Add(originalKeys).Op(":=").Make(Map(utils.Hash).Index().Add(r.EntityPathKey.Type.ReferencedType()))
		def.For().List(Id("_"), key).Op(":=").Range().Add(Keys).BlockFunc(func(def *Group) {
			keyHash := Id("keyHash")
			def.Add(keyHash).Op(":=").Add(keyAccessor(key)).Dot(utils.ComputeHash).Call()
			index := Add(originalKeys).Index(keyHash)
			def.Add(index).Op("=").Append(index, key)
		}).Line()
	}

	def.Add(Entities).Op("=").Make(r.batchGetReturnType())
	rawKey := Id("rawKey")
	resultsReader := Func().Params(Add(types.Reader).Add(types.ReaderQual), Add(rawKey).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		v := Code(Id("v"))
		if isComplexKey {
			def.Add(v).Op(":=").New(r.EntityPathKey.Type.GoType())
		} else {
			def.Var().Add(v).Add(r.EntityPathKey.Type.GoType())
		}
		keyReader := Id("keyReader")
		def.List(keyReader, Err()).Op(":=").Add(utils.NewRor2Reader).Call(rawKey)
		def.Add(utils.IfErrReturn(Err()))
		def.Add(types.Reader.Read(r.EntityPathKey.Type, keyReader, v))
		def.Add(utils.IfErrReturn(Err())).Line()

		if isComplexKey {
			originalKey := Code(Id("originalKey"))
			def.Var().Add(originalKey).Add(r.EntityPathKey.Type.ReferencedType())
			def.For().List(Id("_"), key).Op(":=").Range().Add(originalKeys).Index(keyAccessor(v).Dot(utils.ComputeHash).Call()).BlockFunc(func(def *Group) {
				right := keyAccessor(key)
				if ck != nil {
					right = Op("&").Add(right)
				}

				def.If(keyAccessor(v).Dot(utils.Equals).Call(right)).Block(
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
	def.Return(Entities, Err())
}
