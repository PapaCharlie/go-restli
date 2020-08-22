package codegen

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
)

func (m *Method) actionFuncName() string {
	return ExportedIdentifier(m.Name + "Action")
}

func (r *Resource) actionFuncParams(m *Method, def *Group) {
	m.addEntityTypes(def)
	if len(m.Params) > 0 {
		def.Id("params").Op("*").Qual(r.PackagePath(), m.actionStructType())
	}
}

func (m *Method) actionMethodCallParams() (params []Code) {
	if len(m.Params) > 0 {
		params = append(params, Id("params"))
	}
	return params
}

func (m *Method) actionStructType() string {
	return m.actionFuncName() + "Params"
}

func (m *Method) actionFuncReturnParams(def *Group) {
	if m.Return != nil {
		def.Add(m.Return.ReferencedType())
	}
	def.Error()
}

func (r *Resource) GenerateActionCode(a *Method) *CodeFile {
	actionName := a.Name + "Action"
	c := r.NewCodeFile(actionName)

	actionNameConst := ExportedIdentifier(actionName)
	c.Code.Const().Id(actionNameConst).Op("=").Lit(a.Name).Line()

	hasParams := len(a.Params) > 0
	if hasParams {
		record := &Record{
			NamedType: NamedType{
				Identifier: Identifier{
					Name:      a.actionStructType(),
					Namespace: r.Namespace,
				},
				Doc: fmt.Sprintf("This struct provides the parameters to the %s action", a.Name),
			},
			Fields: a.Params,
		}
		c.Code.Add(record.generateStruct()).Line()
		c.Code.Add(record.generateRestliEncoder()).Line()
	}

	r.addClientFuncDeclarations(c.Code, ClientType, a, func(def *Group) {
		var pathFunc string
		if a.OnEntity {
			pathFunc = ResourceEntityPath
		} else {
			pathFunc = ResourcePath
		}

		returns := a.Return != nil
		var errReturnParams []Code
		if returns {
			errReturnParams = []Code{a.Return.ZeroValueReference(), Err()}
		} else {
			errReturnParams = []Code{Err()}
		}

		def.List(Id(PathVar), Err()).Op(":=").Id(pathFunc).Call(a.entityParams()...)
		IfErrReturn(def, errReturnParams...).Line()

		r.callFormatQueryUrl(def)
		IfErrReturn(def, errReturnParams...).Line()
		def.Id(UrlVar).Dot("RawQuery").Op("=").Lit("action=" + a.Name)

		req := def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver)
		var params *Statement
		if hasParams {
			params = Id("params")
		} else {
			params = Op("&").Qual(ProtocolPackage, "EmptyRecord").Block()
		}
		req.Dot("ActionRequest").Call(Id(ContextVar), Id(UrlVar), params)
		IfErrReturn(def, errReturnParams...).Line()

		if returns {
			result := Id("actionResult")
			def.Var().Add(result).Struct(Id("Value").Add(a.Return.GoType()))
			callDoAndDecode(def, Op("&").Add(result), a.Return.ZeroValueReference())
			returnValue := Add(result).Dot("Value")
			if a.Return.ShouldReference() {
				returnValue = Op("&").Add(returnValue)
			}
			def.Return(returnValue, Nil())
		} else {
			def.List(Id("_"), Err()).Op("=").Id(ClientReceiver).Dot("DoAndIgnore").Call(Id(ReqVar))
			def.Return(Err())
		}
	})

	return c
}
