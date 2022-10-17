package resources

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type Action struct{ methodImplementation }

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

func (a *Action) NonErrorFuncReturnParam() Code {
	if a.Return != nil {
		return Id("actionResult").Add(a.Return.ReferencedType())
	} else {
		return nil
	}
}

func (a *Action) paramsStructType() string {
	return a.FuncName() + "Params"
}

func (a *Action) GenerateCode() *utils.CodeFile {
	actionName := a.Name + "Action"
	c := a.Resource.NewCodeFile(actionName)

	actionNameConst := utils.ExportedIdentifier(actionName)
	c.Code.Const().Id(actionNameConst).Op("=").Qual(utils.RestLiPackage, "QueryParamsString").
		Call(Lit("action=" + a.Name)).Line()

	hasParams := len(a.Params) > 0
	if hasParams {
		record := &types.Record{
			NamedType: types.NamedType{
				Identifier: utils.Identifier{
					Name:      a.paramsStructType(),
					Namespace: a.Resource.Namespace,
				},
				Doc: fmt.Sprintf("%s provides the parameters to the %s action", a.paramsStructType(), a.Name),
			},
			Fields: a.Params,
		}
		c.Code.
			Add(record.GenerateStruct()).Line().
			Add(record.GenerateMarshalRestLi()).Line().
			Add(record.GenerateUnmarshalRestLi()).Line()
	}

	a.Resource.addClientFuncDeclarations(c.Code, ClientType, a, func(def *Group) {
		returns := a.Return != nil
		declareRpStruct(a, def)

		var params Code
		if hasParams {
			params = ActionParams
		} else {
			params = Add(EmptyRecord).Values()
		}

		f := "DoActionRequest"
		callParams := []Code{RestLiClientReceiver, Ctx, Rp, Id(actionNameConst), params}
		if returns {
			f += "WithResults"
			callParams = append(callParams, types.Reader.UnmarshalerFunc(*a.Return))
		}

		def.Return(Qual(utils.RestLiPackage, f).Call(callParams...))
	})

	return c
}

func (a *Action) RegisterMethod(server, resource, segments Code) Code {
	name := "RegisterAction"
	if a.Return != nil {
		name += "WithResults"
	}

	return Qual(utils.RestLiPackage, name).CallFunc(func(def *Group) {
		def.Add(server)
		def.Add(segments)
		def.Lit(a.Name)

		if a.Return != nil {
			def.Line().Add(types.Writer.MarshalerFunc(*a.Return))
		}

		def.Line().Func().Params(registerParams(a)...).Params(methodReturnParams(a)...).BlockFunc(func(def *Group) {
			def.Return(resource).Dot(a.FuncName()).Call(splatRpAndParams(a)...)
		})

	})
}
