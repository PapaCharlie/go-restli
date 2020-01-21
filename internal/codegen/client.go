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
	c := &CodeFile{
		SourceFile:  r.SourceFile,
		PackagePath: r.PackagePath(),
		Filename:    "client",
		Code:        Empty(),
	}

	var generatedRestMethods []Code

	AddWordWrappedComment(c.Code, r.Doc).Line()
	c.Code.Type().Id(ClientInterfaceType).InterfaceFunc(func(def *Group) {
		for _, m := range r.Methods {
			if m.MethodType != REST_METHOD {
				AddWordWrappedComment(def.Empty(), m.Doc)
			} else {
				if code := r.GenerateRestMethodCode(m); code != nil {
					generatedRestMethods = append(generatedRestMethods, code.Line().Line())
				} else {
					// this method is not currently supported, don't add it to the interface
					continue
				}
			}
			def.Add(r.clientFunc(m))
		}
	}).Line().Line()
	c.Code.Type().Id(ClientType).Struct(Op("*").Qual(ProtocolPackage, RestLiClient)).Line().Line()
	c.Code.Func().Id("NewClient").Params(Id("c").Op("*").Qual(ProtocolPackage, RestLiClient)).Id("Client").
		Block(Return(Op("&").Id(ClientType).Values(Id("c")))).
		Line().Line()

	for _, m := range r.Methods {
		if !m.OnEntity {
			r.addResourcePathFunc(c.Code, ResourcePath, m)
			break
		}
	}

	for _, m := range r.Methods {
		if m.OnEntity {
			r.addResourcePathFunc(c.Code, ResourceEntityPath, m)
			break
		}
	}

	c.Code.Add(generatedRestMethods...)

	codeFiles := []*CodeFile{c}

	for _, m := range r.Methods {
		switch m.MethodType {
		case REST_METHOD:
			// This is generated during the interface definition
		case ACTION:
			codeFiles = append(codeFiles, r.GenerateActionCode(m))
		case FINDER:
			codeFiles = append(codeFiles, r.GenerateFinderCode(m))
		}
	}

	return codeFiles
}

func (r *Resource) addResourcePathFunc(def *Statement, funcName string, m *Method) {
	def.Func().Id(funcName).
		ParamsFunc(func(def *Group) { m.addEntityTypes(def) }).
		Params(String(), Error()).BlockFunc(func(def *Group) {

		def.Var().Id(PathVar).String()
		path := m.Path
		for _, pk := range m.PathKeys {
			encodedVariableName := pk.Name + "Str"
			assignment, hasError := pk.Type.RestLiURLEncodeModel(Id(pk.Name))
			if hasError {
				def.List(Id(encodedVariableName), Err()).Op(":=").Add(assignment)
				IfErrReturn(def, Lit(""), Err())
			} else {
				def.Id(encodedVariableName).Op(":=").Add(assignment)
			}

			pattern := fmt.Sprintf("{%s}", pk.Name)
			idx := strings.Index(path, pattern)
			if idx < 0 {
				Logger.Panicf("%s does not appear in %s", pattern, path)
			}
			def.Id(PathVar).Op("+=").Lit(path[:idx]).Op("+").Id(encodedVariableName)
			path = path[idx+len(pattern):]
		}
		def.Line()

		if path != "" {
			def.Id(PathVar).Op("+=").Lit(path)
		}

		def.Return(Id(PathVar), Nil())
	}).Line().Line()
}
