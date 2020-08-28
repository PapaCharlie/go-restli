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

func (m *Method) actionResultsStructType() string {
	return m.actionFuncName() + "Results"
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
		c.Code.Add(record.generateMarshalRestLi()).Line()
	}

	if a.Return != nil {
		results := &Record{
			NamedType: NamedType{
				Identifier: Identifier{
					Name:      a.actionResultsStructType(),
					Namespace: r.Namespace,
				},
				Doc: fmt.Sprintf("This struct deserializes the response from the %s action", a.Name),
			},
			Fields: []Field{{
				Type: *a.Return,
				Name: "value",
			}},
		}
		c.Code.Add(results.generateStruct()).Line().Line()
		c.Code.Add(results.generateUnmarshalRestLi()).Line().Line()
	}

	r.addClientFuncDeclarations(c.Code, ClientType, a, func(def *Group) {
		returns := a.Return != nil
		var errReturnParams []Code
		if returns {
			errReturnParams = []Code{a.Return.ZeroValueReference(), Err()}
		} else {
			errReturnParams = []Code{Err()}
		}

		formatQueryUrl(r, a, def, errReturnParams...)
		def.Add(UrlVar).Dot("RawQuery").Op("=").Lit("action=" + a.Name)

		var params *Statement
		if hasParams {
			params = Id("params")
		} else {
			params = Qual(ProtocolPackage, "EmptyRecord")
		}

		callParams := []Code{
			ContextVar,
			UrlVar,
			params,
		}

		result := Id("actionResult")
		var resultsAccessor Code
		if returns {
			def.Var().Add(result).Id(a.actionResultsStructType())
			callParams = append(callParams, result)
			resultsAccessor = Op("&").Add(result)
		} else {
			resultsAccessor = Nil()
		}

		def.Err().Op("=").Id(ClientReceiver).Dot("DoActionRequest").Call(ContextVar, UrlVar, params, resultsAccessor)

		if returns {
			def.Add(IfErrReturn(errReturnParams...)).Line()
			returnValue := Add(result).Dot("Value")
			if a.Return.ShouldReference() {
				returnValue = Op("&").Add(returnValue)
			}
			def.Return(returnValue, Nil())
		} else {
			def.Return(Err())
		}
	})

	return c
}
