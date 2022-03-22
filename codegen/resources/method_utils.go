package resources

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type MethodImplementation interface {
	GetMethod() *Method
	GetPathKeys() []*PathKey
	FuncName() string
	FuncParamNames() []Code
	FuncParamTypes() []Code
	NonErrorFuncReturnParam() Code
	GenerateCode() *utils.CodeFile
	RegisterMethod(server, resource, segments Code) Code
}

type methodImplementation struct {
	Resource *Resource
	*Method
}

func (m *methodImplementation) GetMethod() *Method {
	return m.Method
}

func (m *methodImplementation) GetPathKeys() (keys []*PathKey) {
	for _, segment := range m.Resource.ParentSegments() {
		if segment.PathKey != nil {
			keys = append(keys, segment.PathKey)
		}
	}
	if m.OnEntity {
		keys = append(keys, m.Resource.LastSegment().PathKey)
	}
	return keys
}

func declareRpStruct(m MethodImplementation, def *Group) {
	rp := def.Add(Rp).Op(":=").Op("&")
	if m.GetMethod().OnEntity {
		rp.Id(ResourceEntityPath)
	} else {
		rp.Id(ResourcePath)
	}

	pathParamNames := entityParamNames(m)
	rp.Add(utils.OrderedValues(func(add func(key Code, value Code)) {
		for _, name := range pathParamNames {
			add(name, name)
		}
	}))
}

func splatRpAndParams(m MethodImplementation) []Code {
	params := []Code{Ctx}
	for _, pk := range m.GetPathKeys() {
		params = append(params, Add(Rp).Dot(pk.Name))
	}
	params = append(params, m.FuncParamNames()...)
	return params
}

func methodFuncName(m MethodImplementation, withContext context) string {
	n := m.FuncName()
	if withContext != none {
		n += WithContext
	}
	return n
}

type context int

const (
	none context = iota
	clientContext
	resourceContext
)

func methodParams(m MethodImplementation, ctx context) (params []Code) {
	names, types := methodParamNames(m), methodParamTypes(m)
	switch ctx {
	case clientContext:
		params = append(params, Add(Ctx).Add(Context))
	case resourceContext:
		params = append(params, Add(Ctx).Op("*").Qual(utils.RestLiPackage, "RequestContext"))
	}
	for i, name := range names {
		params = append(params, Add(name).Add(types[i]))
	}
	return params
}

func methodParamNames(m MethodImplementation) []Code {
	return append(entityParamNames(m), m.FuncParamNames()...)
}

func methodParamTypes(m MethodImplementation) []Code {
	return append(entityParamTypes(m), m.FuncParamTypes()...)
}

func methodReturnParams(m MethodImplementation) []Code {
	p := m.NonErrorFuncReturnParam()
	if p == nil {
		return []Code{Err().Error()}
	} else {
		return []Code{p, Err().Error()}
	}
}

func registerParams(m MethodImplementation) []Code {
	params := []Code{RequestContextParam}
	rp := Add(Rp).Op("*")
	if m.GetMethod().OnEntity {
		rp.Id(ResourceEntityPath)
	} else {
		rp.Id(ResourcePath)
	}
	params = append(params, rp)

	names, types := m.FuncParamNames(), m.FuncParamTypes()
	for i, name := range names {
		params = append(params, Add(name).Add(types[i]))
	}

	if len(m.GetMethod().Params) == 0 {
		p := Id("_")
		if rM, ok := m.(*RestMethod); ok && rM.usesBatchQueryParams() {
			p.Op("*").Qual(utils.RestLiPackage, "SliceBatchQueryParams").Index(rM.EntityKeyType())
		} else {
			p.Add(EmptyRecord)
		}
		params = append(params, p)
	}

	return params
}

func entityParamNames(m MethodImplementation) (params []Code) {
	for _, pk := range m.GetPathKeys() {
		params = append(params, Id(pk.Name))
	}
	return params
}

func entityParamTypes(m MethodImplementation) (params []Code) {
	for _, pk := range m.GetPathKeys() {
		params = append(params, pk.GoType())
	}
	return params
}
