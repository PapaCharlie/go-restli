package resources

import (
	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type MethodImplementation interface {
	GetMethod() *Method
	GetResource() *Resource
	IsSupported() bool
	FuncName() string
	FuncParamNames() []Code
	FuncParamTypes() []Code
	NonErrorFuncReturnParams() []Code
	GenerateCode() *utils.CodeFile
}

type methodImplementation struct {
	Resource *Resource
	*Method
}

func (m *methodImplementation) GetMethod() *Method {
	return m.Method
}

func (m *methodImplementation) GetResource() *Resource {
	return m.Resource
}

func (m *methodImplementation) addEntityId(def *Group, accessor Code, returns []Code) {
	def.Add(types.Writer.Write(m.EntityPathKey.Type, Add(utils.EntityIDsEncoder).Dot("AddEntityID").Call(), accessor, returns...))
}

func formatQueryUrl(m MethodImplementation, def *Group, populateEntityIDs func(def *Group), returns ...Code) {
	if m.GetMethod().OnEntity {
		def.List(Path, Err()).Op(":=").Id(ResourceEntityPath).Call(m.GetMethod().entityParamNames()...)
	} else {
		def.List(Path, Err()).Op(":=").Id(ResourcePath).Call(m.GetMethod().entityParamNames()...)
	}

	def.Add(utils.IfErrReturn(returns...)).Line()

	encodeQueryParams := Code(Add(QueryParams).Dot(utils.EncodeQueryParams))
	callEncodeQueryParams := func(encoder Code) {
		rawQuery := Id("rawQuery")
		def.Var().Add(rawQuery).String()
		def.List(rawQuery, Err()).Op("=").Add(encoder)
		def.Add(utils.IfErrReturn(returns...))
		def.Add(Path).Op("+=").Lit("?").Op("+").Add(rawQuery)
		def.Line()
	}

	switch m.(type) {
	case *Action:
		def.Add(Path).Op("+=").Lit("?action=" + m.GetMethod().Name)
	case *Finder:
		callEncodeQueryParams(Add(encodeQueryParams).Call())
	case *RestMethod:
		r := m.(*RestMethod)
		hasParams := len(m.GetMethod().Params) > 0
		if r.isBatch() {
			def.Add(utils.EntityIDsEncoder).Op(":=").New(utils.BatchEntityIDsEncoder)
			populateEntityIDs(def)
			def.Line()

			if hasParams {
				callEncodeQueryParams(Add(encodeQueryParams).Call(utils.EntityIDsEncoder))
			} else {
				callEncodeQueryParams(Add(utils.EntityIDsEncoder).Dot("GenerateRawQuery").Call())
			}
		} else {
			if hasParams {
				callEncodeQueryParams(Add(encodeQueryParams).Call())
			}
		}
	}

	def.List(Url, Err()).
		Op(":=").
		Id(ClientReceiver).Dot("FormatQueryUrl").
		Call(Lit(m.GetResource().RootResourceName), Path)
	def.Add(utils.IfErrReturn(returns...)).Line()
}

func methodFuncName(m MethodImplementation, withContext bool) string {
	n := m.FuncName()
	if withContext {
		n += WithContext
	}
	return n
}

func addParams(def *Group, names, types []Code) {
	for i, name := range names {
		def.Add(name).Add(types[i])
	}
}

func methodParamNames(m MethodImplementation) []Code {
	return append(m.GetMethod().entityParamNames(), m.FuncParamNames()...)
}

func methodParamTypes(m MethodImplementation) []Code {
	return append(m.GetMethod().entityParamTypes(), m.FuncParamTypes()...)
}

func methodReturnParams(m MethodImplementation) []Code {
	return append(m.NonErrorFuncReturnParams(), Err().Error())
}

func (m *Method) entityParamNames() (params []Code) {
	for _, pk := range m.PathKeys {
		params = append(params, Id(pk.Name))
	}
	return params
}

func (m *Method) entityParamTypes() (params []Code) {
	for _, pk := range m.PathKeys {
		params = append(params, pk.Type.ReferencedType())
	}
	return params
}

func (r *RestMethod) batchMethodBoilerplate(def *Group, valueSetter func(def *Group, keyAccessor, valueReader Code)) Code {
	ck := r.EntityPathKey.Type.ComplexKey()
	isComplexKey := ck != nil || r.EntityPathKey.Type.UnderlyingPrimitive() == nil
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

		var entityKeyIterator Code
		if r.usesBatchMapInput() {
			entityKeyIterator = List(Key).Op(":=").Range().Add(Entities)
		} else {
			entityKeyIterator = List(Id("_"), Key).Op(":=").Range().Add(Keys)
		}

		def.For().Add(entityKeyIterator).BlockFunc(func(def *Group) {
			keyHash := Id("keyHash")
			def.Add(keyHash).Op(":=").Add(keyAccessor(Key)).Dot(utils.ComputeHash).Call()
			index := Add(originalKeys).Index(keyHash)
			def.Add(index).Op("=").Append(index, Key)
		}).Line()
	}

	keyReader, valueReader := Id("keyReader"), Id("valueReader")
	return Func().Params(Add(keyReader).Add(types.ReaderQual), Add(valueReader).Add(types.ReaderQual)).Params(Err().Error()).BlockFunc(func(def *Group) {
		v := Code(Id("v"))
		if isComplexKey {
			def.Add(v).Op(":=").New(r.EntityPathKey.Type.GoType())
		} else {
			def.Var().Add(v).Add(r.EntityPathKey.Type.GoType())
		}
		def.Add(types.Reader.Read(r.EntityPathKey.Type, keyReader, v))
		def.Add(utils.IfErrReturn(Err())).Line()

		if isComplexKey {
			originalKey := Code(Id("originalKey"))
			def.Var().Add(originalKey).Add(r.EntityPathKey.Type.ReferencedType())
			def.For().List(Id("_"), Key).Op(":=").Range().Add(originalKeys).Index(keyAccessor(v).Dot(utils.ComputeHash).Call()).BlockFunc(func(def *Group) {
				right := keyAccessor(Key)
				if ck != nil {
					right = Op("&").Add(right)
				}

				def.If(keyAccessor(v).Dot(utils.Equals).Call(right)).Block(
					Add(originalKey).Op("=").Add(Key),
					Break(),
				)
			})
			def.If(Add(originalKey).Op("==").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("unknown key returned by batch get: %q"), keyReader)),
			)
			def.Line()
			v = originalKey
		}

		valueSetter(def, v, valueReader)
	})
}
