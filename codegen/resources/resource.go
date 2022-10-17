package resources

import (
	"log"
	"sort"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func (pk *PathKey) GoType() *Statement {
	return pk.Type.ReferencedType()
}

func (r *Resource) PackagePath() string {
	return utils.FqcpToPackagePath(r.PackageRoot, r.Namespace)
}

func (r *Resource) LocalType(name string) *Statement {
	return Qual(r.PackagePath(), name)
}

func (r *Resource) NewCodeFile(filename string) *utils.CodeFile {
	return &utils.CodeFile{
		PackagePath: r.PackagePath(),
		PackageRoot: r.PackageRoot,
		SourceFile:  r.SourceFile,
		Filename:    filename,
		Code:        Empty(),
	}
}

func (r *Resource) GenerateCode() []*utils.CodeFile {
	// Some resources simply define no methods, so skip generating any code
	if len(r.Methods) == 0 {
		return nil
	}

	resource := r.NewCodeFile("resource")

	for _, m := range r.Methods {
		if !m.GetMethod().OnEntity {
			r.addResourcePathStructs(resource.Code, false)
			break
		}
	}

	for _, m := range r.Methods {
		if m.GetMethod().OnEntity {
			r.addResourcePathStructs(resource.Code, true)
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
		resource.Code.Var().DefsFunc(func(def *Group) {
			if len(r.ReadOnlyFields) > 0 {
				def.Add(ReadOnlyFields).Op("=").Add(newPathSpec(r.ReadOnlyFields))
			}

			var createAndReadOnlyFields []string
			inserted := make(map[string]bool)
			for _, field := range append(append([]string(nil), r.ReadOnlyFields...), r.CreateOnlyFields...) {
				if _, ok := inserted[field]; ok {
					log.Panicf("%q is declared both as read only and create only", field)
				}
				inserted[field] = true
				createAndReadOnlyFields = append(createAndReadOnlyFields, field)
			}
			sort.Strings(createAndReadOnlyFields)
			def.Add(CreateAndReadOnlyFields).Op("=").Add(newPathSpec(createAndReadOnlyFields)).Line().Line()
		}).Line().Line()
	}

	if r.ResourceSchema != nil {
		resource.Code.Type().DefsFunc(func(def *Group) {
			entityType := r.ResourceSchema.ReferencedType()
			def.Id(Elements).Op("=").Qual(utils.RestLiCommonPackage, Elements).Index(entityType)
			if r.LastSegment().PathKey != nil {
				keyType := r.LastSegment().PathKey.Type.ReferencedType()
				def.Id(CreatedEntity).Op("=").Qual(utils.RestLiCommonPackage, CreatedEntity).Index(keyType)
				def.Id(CreatedAndReturnedEntity).Op("=").Qual(utils.RestLiCommonPackage, CreatedAndReturnedEntity).Index(List(keyType, entityType))

				batchResponse := Qual(utils.RestLiCommonPackage, BatchResponse)
				def.Id(BatchEntities).Op("=").Add(batchResponse).Index(List(keyType, entityType))
				def.Id(BatchResponse).Op("=").Add(batchResponse).Index(List(keyType, BatchEntityUpdateResponse))
			}
		}).Line().Line()
	}

	resource.Code.Add(r.generateClientCode(), r.generateResourceCode())
	codeFiles := []*utils.CodeFile{resource}

	for _, m := range r.Methods {
		codeFiles = append(codeFiles, m.GenerateCode())
	}

	codeFiles = append(codeFiles, r.generateTestCode())

	return codeFiles
}

func (r *Resource) generateClientCode() Code {
	def := Empty()

	utils.AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(ClientInterfaceType).InterfaceFunc(func(def *Group) {
		for _, m := range r.Methods {
			if m.GetMethod().MethodType != REST_METHOD {
				utils.AddWordWrappedComment(def.Empty(), m.GetMethod().Doc)
			}
			def.Add(r.clientFuncDeclaration(m, none))
			def.Add(r.clientFuncDeclaration(m, clientContext))
		}
	}).Line().Line()

	c := Code(Id("c"))
	def.Type().Id(ClientType).Struct(Op("*").Add(RestLiClientQual)).Line().Line()
	def.Func().Id("NewClient").Params(Add(c).Op("*").Add(RestLiClientQual)).Id(ClientInterfaceType).
		Block(Return(Op("&").Id(ClientType).Values(c))).Line().Line()

	return def
}

func (r *Resource) generateResourceCode() Code {
	def := Empty()

	utils.AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(ResourceInterfaceType).InterfaceFunc(func(def *Group) {
		for _, m := range r.Methods {
			if m.GetMethod().MethodType != REST_METHOD {
				utils.AddWordWrappedComment(def.Empty(), m.GetMethod().Doc)
			}
			def.Add(r.resourceFuncDeclaration(m))
		}
	}).Line().Line()

	var server, resource Code = Id("server"), Id("resource")
	def.Func().Id("RegisterResource").
		Params(Add(server).Qual(utils.RestLiPackage, "Server"), Add(resource).Id(ResourceInterfaceType)).
		BlockFunc(func(def *Group) {
			segments := Code(Id("segments"))
			def.Add(segments).Op(":=").Index().Qual(utils.RestLiPackage, "ResourcePathSegment").
				ValuesFunc(func(def *Group) {
					for _, rps := range r.ResourcePathSegments {
						def.Line().Qual(utils.RestLiPackage, "NewResourcePathSegment").Call(Lit(rps.ResourceName), Lit(rps.PathKey != nil))
					}
					def.Line()
				})

			for _, m := range r.Methods {
				def.Add(m.RegisterMethod(server, resource, segments))
			}

		})

	return def
}

func (r *Resource) addResourcePathStructs(def *Statement, onEntity bool) {
	var structName string
	if onEntity {
		structName = ResourceEntityPath
	} else {
		structName = ResourcePath
	}

	segments := append([]ResourcePathSegment(nil), r.ParentSegments()...)
	segments = append(segments, ResourcePathSegment{ResourceName: r.LastSegment().ResourceName})
	if onEntity {
		segments[len(segments)-1].PathKey = r.LastSegment().PathKey
	}

	def.Type().Id(structName).StructFunc(func(def *Group) {
		for _, rps := range segments {
			if rps.PathKey != nil {
				def.Id(rps.PathKey.Name).Add(rps.PathKey.GoType())
			}
		}
	}).Line().Line()

	const rp = "rp"
	types.AddNewInstance(def, rp, structName)

	utils.AddFuncOnReceiver(def, rp, structName, "ResourcePath", types.RecordShouldUsePointer).
		Params().
		Params(Id("path").String(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(types.Writer).Op(":=").Qual(utils.RestLiCodecPackage, "NewRor2PathWriter").Call()

			for _, rps := range segments {
				path := "/" + rps.ResourceName
				if rps.PathKey != nil {
					path += "/"
				}
				def.Add(types.Writer).Dot("RawPathSegment").Call(Lit(path))
				if rps.PathKey != nil {
					def.Add(types.Writer.Write(rps.PathKey.Type, types.Writer, Add(Rp).Dot(rps.PathKey.Name), Lit(""), Err()))
				}
			}

			def.Line()

			def.Return(types.Writer.Finalize(), Nil())
		}).Line().Line()

	utils.AddFuncOnReceiver(def, rp, structName, "RootResource", types.RecordShouldUsePointer).
		Params().
		String().
		Block(Return(Lit(r.ResourcePathSegments[0].ResourceName))).Line().Line()

	segmentsSlice := Code(Id("segments"))
	utils.AddFuncOnReceiver(def, rp, structName, "UnmarshalResourcePath", types.RecordShouldUsePointer).
		Params(Add(segmentsSlice).Index().Add(types.ReaderQual)).
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			for i, rps := range segments {
				if rps.PathKey != nil {
					accessor := Add(Rp).Dot(rps.PathKey.Name)
					if rps.PathKey.Type.ShouldReference() {
						def.Add(accessor).Op("=").New(rps.PathKey.Type.GoType())
					}
					def.Add(types.Reader.Read(rps.PathKey.Type, Add(segmentsSlice).Index(Lit(i)), accessor))
					def.Add(utils.IfErrReturn(Err())).Line()
				}
			}
			def.Line()

			def.Return(Nil())
		}).Line().Line()
}

func (r *Resource) generateTestCode() *utils.CodeFile {
	const (
		mock           = "Mock"
		clientStruct   = mock + ClientInterfaceType
		resourceStruct = mock + ResourceInterfaceType
	)

	var clientStructFields []Code
	clientFuncs := Empty()

	for _, m := range r.Methods {
		clientStructFields = append(clientStructFields,
			Id(mock+m.FuncName()).Func().Params(methodParams(m, clientContext)...).Params(methodReturnParams(m)...),
		)
		r.addClientFuncDeclarations(clientFuncs, clientStruct, m, func(def *Group) {
			def.Return(Id(ClientReceiver).Dot(mock + methodFuncName(m, none)).CallFunc(func(def *Group) {
				def.Add(Ctx)
				for _, p := range methodParamNames(m) {
					def.Add(p)
				}
			}))
		}).Line().Line()
	}

	var resourceStructFields []Code
	resourceFuncs := Empty()

	receiver := Id("r")
	for _, m := range r.Methods {
		name := mock + m.FuncName()
		resourceStructFields = append(resourceStructFields,
			Id(name).Func().Params(methodParams(m, resourceContext)...).Params(methodReturnParams(m)...),
		)
		resourceFuncs.Func().Params(Add(receiver).Op("*").Id(resourceStruct)).
			Add(r.resourceFuncDeclaration(m)).
			BlockFunc(func(def *Group) {
				def.Return(Add(receiver).Dot(name).Call(append([]Code{Ctx}, methodParamNames(m)...)...))
			}).Line().Line()
	}

	clientTest := r.NewCodeFile("resource")
	clientTest.PackagePath += "_test"

	clientTest.Code.
		Type().Id(clientStruct).Struct(clientStructFields...).Line().Line().
		Add(clientFuncs).Line().Line().
		Type().Id(resourceStruct).Struct(resourceStructFields...).Line().Line().
		Add(resourceFuncs)

	return clientTest
}

func (r *Resource) clientFuncDeclaration(m MethodImplementation, withContext context) *Statement {
	return Id(methodFuncName(m, withContext)).Params(methodParams(m, withContext)...).Params(methodReturnParams(m)...)
}

func (r *Resource) resourceFuncDeclaration(m MethodImplementation) *Statement {
	return Id(m.FuncName()).Params(methodParams(m, resourceContext)...).Params(methodReturnParams(m)...)
}

func (r *Resource) addClientFuncDeclarations(def *Statement, clientType string, m MethodImplementation, block func(*Group)) *Statement {
	clientFuncDeclarationStart := Func().Params(Id(ClientReceiver).Op("*").Id(clientType))

	utils.AddWordWrappedComment(def, m.GetMethod().Doc).Line().
		Add(clientFuncDeclarationStart).
		Add(r.clientFuncDeclaration(m, none)).
		Block(Return(Id(ClientReceiver).Dot(methodFuncName(m, clientContext)).CallFunc(func(def *Group) {
			def.Qual("context", "Background").Call()
			for _, p := range methodParamNames(m) {
				def.Add(p)
			}
		}))).
		Line().Line()

	utils.AddWordWrappedComment(def, m.GetMethod().Doc).Line().
		Add(clientFuncDeclarationStart).
		Add(r.clientFuncDeclaration(m, clientContext)).
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

func (r *Resource) createAndReadOnlyFields() Code {
	if len(r.ReadOnlyFields) > 0 || len(r.CreateOnlyFields) > 0 {
		return CreateAndReadOnlyFields
	} else {
		return NoExcludedFields
	}
}

func (r *Resource) ParentSegments() []ResourcePathSegment {
	return r.ResourcePathSegments[:len(r.ResourcePathSegments)-1]
}

func (r *Resource) LastSegment() ResourcePathSegment {
	return r.ResourcePathSegments[len(r.ResourcePathSegments)-1]
}
