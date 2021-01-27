package resources

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/internal/codegen/types"
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type Action struct{ methodImplementation }

func (a *Action) IsSupported() bool {
	return true
}

func (a *Action) FuncName() string {
	return utils.ExportedIdentifier(a.Name + "Action")
}

func (a *Action) FuncParamNames() []Code {
	if len(a.Params) > 0 {
		return []Code{ActionParams}
	} else {
		return nil
	}
}

func (a *Action) FuncParamTypes() []Code {
	if len(a.Params) > 0 {
		return []Code{Op("*").Qual(a.Resource.PackagePath(), a.paramsStructType())}
	} else {
		return nil
	}
}

func (a *Action) NonErrorFuncReturnParams() []Code {
	if a.Return != nil {
		return []Code{Id("actionResult").Add(a.Return.ReferencedType())}
	} else {
		return nil
	}
}

func (a *Action) paramsStructType() string {
	return a.FuncName() + "Params"
}

func (a *Action) resultsStructType() string {
	return a.FuncName() + "Result"
}

func (m *Method) actionFuncReturnParams(def *Group) {
	if m.Return != nil {
		def.Add(m.Return.ReferencedType())
	}
	def.Error()
}

func (a *Action) GenerateCode() *utils.CodeFile {
	actionName := a.Name + "Action"
	c := a.Resource.NewCodeFile(actionName)

	actionNameConst := utils.ExportedIdentifier(actionName)
	c.Code.Const().Id(actionNameConst).Op("=").Lit(a.Name).Line()

	hasParams := len(a.Params) > 0
	if hasParams {
		record := &types.Record{
			NamedType: types.NamedType{
				Identifier: utils.Identifier{
					Name:      a.paramsStructType(),
					Namespace: a.Resource.Namespace,
				},
				Doc: fmt.Sprintf("This struct provides the parameters to the %s action", a.Name),
			},
			Fields: a.Params,
		}
		c.Code.Add(record.GenerateStruct()).Line()
		c.Code.Add(record.GenerateMarshalRestLi()).Line()
	}

	if a.Return != nil {
		results := &types.Record{
			NamedType: types.NamedType{
				Identifier: utils.Identifier{
					Name:      a.resultsStructType(),
					Namespace: a.Resource.Namespace,
				},
				Doc: fmt.Sprintf("This struct deserializes the response from the %s action", a.Name),
			},
			Fields: []types.Field{{
				Type: *a.Return,
				Name: "value",
			}},
		}
		c.Code.Add(results.GenerateStruct()).Line().Line()
		c.Code.Add(results.GenerateUnmarshalRestLi()).Line().Line()
	}

	a.Resource.addClientFuncDeclarations(c.Code, ClientType, a, func(def *Group) {
		returns := a.Return != nil
		var errReturnParams []Code
		if returns {
			errReturnParams = []Code{a.Return.ZeroValueReference(), Err()}
		} else {
			errReturnParams = []Code{Err()}
		}

		formatQueryUrl(a, def, nil, errReturnParams...)

		var params Code
		if hasParams {
			params = ActionParams
		} else {
			params = Qual(utils.ProtocolPackage, "EmptyRecord")
		}

		callParams := []Code{
			Ctx,
			Url,
			params,
		}

		result := Id("actionResultUnmarshaler")
		var resultsAccessor Code
		if returns {
			def.Var().Add(result).Id(a.resultsStructType())
			callParams = append(callParams, result)
			resultsAccessor = Op("&").Add(result)
		} else {
			resultsAccessor = Nil()
		}

		def.Err().Op("=").Id(ClientReceiver).Dot("DoActionRequest").Call(Ctx, Url, params, resultsAccessor)

		if returns {
			returnValue := Add(result).Dot("Value")
			if a.Return.ShouldReference() {
				returnValue = Op("&").Add(returnValue)
			}
			def.Return(returnValue, Err())
		} else {
			def.Return(Err())
		}
	})

	return c
}
