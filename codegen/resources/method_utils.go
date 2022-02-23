package resources

import (
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
	NonErrorFuncReturnParam() Code
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

func declareRpStruct(m MethodImplementation, def *Group) {
	rp := def.Add(Rp).Op(":=").Op("&")
	if m.GetMethod().OnEntity {
		rp.Id(ResourceEntityPath)
	} else {
		rp.Id(ResourcePath)
	}

	pathParamNames := m.GetMethod().entityParamNames()
	rp.Add(utils.OrderedValues(func(add func(key Code, value Code)) {
		for _, name := range pathParamNames {
			add(name, name)
		}
	}))
}

func methodFuncName(m MethodImplementation, withContext bool) string {
	n := m.FuncName()
	if withContext {
		n += WithContext
	}
	return n
}

func addEntityParams(def *Group, m MethodImplementation) {
	names, types := methodParamNames(m), methodParamTypes(m)
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
	p := m.NonErrorFuncReturnParam()
	if p == nil {
		return []Code{Err().Error()}
	} else {
		return []Code{p, Err().Error()}
	}
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
