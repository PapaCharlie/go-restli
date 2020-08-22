package codegen

import (
	"fmt"
	"strings"

	. "github.com/dave/jennifer/jen"
)

type Resource struct {
	Namespace        string
	Doc              string
	SourceFile       string
	RootResourceName string
	ResourceSchema   *RestliType
	Methods          []*Method
}

func (r *Resource) PackagePath() string {
	return FqcpToPackagePath(r.Namespace)
}

func (r *Resource) GenerateCode() []*CodeFile {
	client := r.NewCodeFile("client")

	AddWordWrappedComment(client.Code, r.Doc).Line()
	client.Code.Type().Id(ClientInterfaceType).InterfaceFunc(func(def *Group) {
		for _, m := range r.Methods {
			AddWordWrappedComment(def.Empty(), m.Doc)
			if m.MethodType == REST_METHOD && !isMethodSupported(m.RestLiMethod()) {
				Logger.Printf("Warning: %s method is not currently implemented", m.Name)
				continue
			}
			def.Add(r.clientFuncDeclaration(m, false))
			def.Add(r.clientFuncDeclaration(m, true))
		}
	}).Line().Line()
	client.Code.Type().Id(ClientType).Struct(Op("*").Qual(ProtocolPackage, RestLiClient)).Line().Line()
	client.Code.Func().Id("NewClient").Params(Id("c").Op("*").Qual(ProtocolPackage, RestLiClient)).Id("Client").
		Block(Return(Op("&").Id(ClientType).Values(Id("c")))).
		Line().Line()

	for _, m := range r.Methods {
		if !m.OnEntity {
			r.addResourcePathFunc(client.Code, ResourcePath, m)
			break
		}
	}

	for _, m := range r.Methods {
		if m.OnEntity {
			r.addResourcePathFunc(client.Code, ResourceEntityPath, m)
			break
		}
	}

	codeFiles := []*CodeFile{client}

	for _, m := range r.Methods {
		switch m.MethodType {
		case REST_METHOD:
			if isMethodSupported(m.RestLiMethod()) {
				codeFiles = append(codeFiles, r.GenerateRestMethodCode(m))
			}
		case ACTION:
			codeFiles = append(codeFiles, r.GenerateActionCode(m))
		case FINDER:
			codeFiles = append(codeFiles, r.GenerateFinderCode(m))
		}
	}

	codeFiles = append(codeFiles, r.generateTestCode())

	return codeFiles
}

func (r *Resource) addResourcePathFunc(def *Statement, funcName string, m *Method) {
	def.Func().Id(funcName).
		ParamsFunc(func(def *Group) { m.addEntityTypes(def) }).
		Params(Id("path").String(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(Encoder).Op(":=").Qual(RestLiEncodingPackage, "NewPathEncoder").Call()

			path := m.Path
			for _, pk := range m.PathKeys {
				pattern := fmt.Sprintf("{%s}", pk.Name)
				idx := strings.Index(path, pattern)
				if idx < 0 {
					Logger.Panicf("%s does not appear in %s", pattern, path)
				}
				def.Add(Encoder).Dot("RawPathSegment").Call(Lit(path[:idx]))
				path = path[idx+len(pattern):]

				accessor := Id(pk.Name)
				if pk.Type.Reference != nil && !pk.Type.ShouldReference() {
					accessor = Op("&").Add(accessor)
				}
				Encoder.Write(def, pk.Type, accessor, Lit(""))
			}
			def.Line()

			if path != "" {
				def.Add(Encoder).Dot("RawPathSegment").Call(Lit(path))
			}

			def.Return(Encoder.Finalize(), Nil())
		}).Line().Line()
}

func (r *Resource) generateTestCode() *CodeFile {
	const (
		mock       = "Mock"
		structName = mock + ClientInterfaceType
	)

	var structFields []Code

	funcs := Empty()

	for _, m := range r.Methods {
		if m.MethodType == REST_METHOD && !isMethodSupported(m.RestLiMethod()) {
			continue
		}
		structFields = append(structFields,
			Id(mock+r.methodFuncName(m, false)).Func().ParamsFunc(func(def *Group) {
				def.Id("ctx").Qual("context", "Context")
				r.methodFuncParams(m, def)
			}).ParamsFunc(r.methodReturnParams(m)),
		)
		r.addClientFuncDeclarations(funcs, structName, m, func(def *Group) {
			def.Return(Id(ClientReceiver).Dot(mock + r.methodFuncName(m, false)).CallFunc(func(def *Group) {
				def.Id(ContextVar)
				for _, p := range append(m.entityParams(), m.methodCallParams()...) {
					def.Add(p)
				}
			}))
		}).Line().Line()
	}

	clientTest := r.NewCodeFile("client")
	clientTest.PackagePath += "_test"

	clientTest.Code.Type().Id(structName).Struct(structFields...).Line().Line()
	clientTest.Code.Add(funcs)

	return clientTest
}
