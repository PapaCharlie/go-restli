package resources

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func (r *Resource) PackagePath() string {
	return utils.FqcpToPackagePath(r.Namespace)
}

func (r *Resource) NewCodeFile(filename string) *utils.CodeFile {
	return &utils.CodeFile{
		PackagePath: r.PackagePath(),
		SourceFile:  r.SourceFile,
		Filename:    filename,
		Code:        Empty(),
	}
}

func (r *Resource) GenerateCode() []*utils.CodeFile {
	client := r.NewCodeFile("client")

	utils.AddWordWrappedComment(client.Code, r.Doc).Line()
	client.Code.Type().Id(ClientInterfaceType).InterfaceFunc(func(def *Group) {
		for _, m := range r.Methods {
			if !m.IsSupported() {
				utils.Logger.Printf("Warning: %s method is not currently implemented", m.GetMethod().Name)
				continue
			}
			if m.GetMethod().MethodType != REST_METHOD {
				utils.AddWordWrappedComment(def.Empty(), m.GetMethod().Doc)
			}
			def.Add(r.clientFuncDeclaration(m, false))
			def.Add(r.clientFuncDeclaration(m, true))
		}
	}).Line().Line()
	client.Code.Type().Id(ClientType).Struct(Op("*").Add(RestLiClient)).Line().Line()
	client.Code.Func().Id("NewClient").Params(Id("c").Op("*").Add(RestLiClient)).Id("Client").
		Block(Return(Op("&").Id(ClientType).Values(Id("c")))).
		Line().Line()

	for _, m := range r.Methods {
		if !m.GetMethod().OnEntity {
			r.addResourcePathFunc(client.Code, ResourcePath, m.GetMethod())
			break
		}
	}

	for _, m := range r.Methods {
		if m.GetMethod().OnEntity {
			r.addResourcePathFunc(client.Code, ResourceEntityPath, m.GetMethod())
			break
		}
	}

	newPathSpec := func(directives []string) Code {
		return Qual(utils.RestLiCodecPackage, "NewPathSpec").CallFunc(func(def *Group) {
			for _, d := range directives {
				def.Line().Add(Lit(d))
			}
			def.Line()
		})
	}

	if len(r.ReadOnlyFields) > 0 || len(r.CreateOnlyFields) > 0 {
		client.Code.Var().DefsFunc(func(def *Group) {
			if len(r.ReadOnlyFields) > 0 {
				def.Add(ReadOnlyFields).Op("=").Add(newPathSpec(r.ReadOnlyFields))
			}
			if len(r.CreateOnlyFields) > 0 {
				def.Add(CreateOnlyFields).Op("=").Add(newPathSpec(r.CreateOnlyFields))
			}

			var createAndReadOnlyFields []string
			inserted := make(map[string]bool)
			for _, d := range append(append([]string(nil), r.ReadOnlyFields...), r.CreateOnlyFields...) {
				if _, ok := inserted[d]; ok {
					continue
				}
				inserted[d] = true
				createAndReadOnlyFields = append(createAndReadOnlyFields, d)
			}
			sort.Strings(createAndReadOnlyFields)
			def.Add(CreateAndReadOnlyFields).Op("=").Add(newPathSpec(createAndReadOnlyFields))
		})
	}

	codeFiles := []*utils.CodeFile{client}

	for _, m := range r.Methods {
		if m.IsSupported() {
			codeFiles = append(codeFiles, m.GenerateCode())
		}
	}

	codeFiles = append(codeFiles, r.generateTestCode())

	return codeFiles
}

func (r *Resource) addResourcePathFunc(def *Statement, funcName string, m *Method) {
	def.Func().Id(funcName).
		ParamsFunc(func(def *Group) { addParams(def, m.entityParamNames(), m.entityParamTypes()) }).
		Params(Id("path").String(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(types.Writer).Op(":=").Qual(utils.RestLiCodecPackage, "NewRor2PathWriter").Call()

			path := m.Path
			for _, pk := range m.PathKeys {
				pattern := fmt.Sprintf("{%s}", pk.Name)
				idx := strings.Index(path, pattern)
				if idx < 0 {
					utils.Logger.Panicf("%s does not appear in %s", pattern, path)
				}
				def.Add(types.Writer).Dot("RawPathSegment").Call(Lit(path[:idx]))
				path = path[idx+len(pattern):]

				def.Add(types.Writer.Write(pk.Type, types.Writer, Id(pk.Name), Lit(""), Err()))
			}
			def.Line()

			if path != "" {
				def.Add(types.Writer).Dot("RawPathSegment").Call(Lit(path))
			}

			def.Return(types.Writer.Finalize(), Nil())
		}).Line().Line()
}

func (r *Resource) generateTestCode() *utils.CodeFile {
	const (
		mock       = "Mock"
		structName = mock + ClientInterfaceType
	)

	var structFields []Code

	funcs := Empty()

	for _, m := range r.Methods {
		if !m.IsSupported() {
			continue
		}
		structFields = append(structFields,
			Id(mock+methodFuncName(m, false)).Func().ParamsFunc(func(def *Group) {
				def.Add(Ctx).Add(Context)
				addParams(def, methodParamNames(m), methodParamTypes(m))
			}).Params(methodReturnParams(m)...),
		)
		r.addClientFuncDeclarations(funcs, structName, m, func(def *Group) {
			def.Return(Id(ClientReceiver).Dot(mock + methodFuncName(m, false)).CallFunc(func(def *Group) {
				def.Add(Ctx)
				for _, p := range methodParamNames(m) {
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

func (r *Resource) clientFuncDeclaration(m MethodImplementation, withContext bool) *Statement {
	params := func(def *Group) {
		if withContext {
			def.Add(Ctx).Add(Context)
		}
		addParams(def, methodParamNames(m), methodParamTypes(m))
	}

	return Id(methodFuncName(m, withContext)).ParamsFunc(params).Params(methodReturnParams(m)...)
}

func (r *Resource) addClientFuncDeclarations(def *Statement, clientType string, m MethodImplementation, block func(*Group)) *Statement {
	clientFuncDeclarationStart := Func().Params(Id(ClientReceiver).Op("*").Id(clientType))

	utils.AddWordWrappedComment(def, m.GetMethod().Doc).Line().
		Add(clientFuncDeclarationStart).
		Add(r.clientFuncDeclaration(m, false)).
		Block(Return(Id(ClientReceiver).Dot(methodFuncName(m, true)).CallFunc(func(def *Group) {
			def.Qual("context", "Background").Call()
			for _, p := range methodParamNames(m) {
				def.Add(p)
			}
		}))).
		Line().Line()

	utils.AddWordWrappedComment(def, m.GetMethod().Doc).Line().
		Add(clientFuncDeclarationStart).
		Add(r.clientFuncDeclaration(m, true)).
		BlockFunc(block)

	return def
}

func (r *Resource) readOnlyFields() Code {
	if len(r.ReadOnlyFields) > 0 {
		return ReadOnlyFields
	} else {
		return NoExcludedFields
	}
}

func (r *Resource) createOnlyFields() Code {
	if len(r.CreateOnlyFields) > 0 {
		return CreateOnlyFields
	} else {
		return NoExcludedFields
	}
}
func (r *Resource) createAndReadOnlyFields() Code {
	if len(r.ReadOnlyFields) > 0 || len(r.CreateOnlyFields) > 0 {
		return CreateAndReadOnlyFields
	} else {
		return NoExcludedFields
	}
}
