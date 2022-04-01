package resources

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type Finder struct{ methodImplementation }

func (f *Finder) FuncName() string {
	return FindBy + utils.ExportedIdentifier(f.Name)
}

func (f *Finder) FuncParamNames() []Code {
	if len(f.Params) > 0 {
		return []Code{QueryParams}
	} else {
		return nil
	}
}

func (f *Finder) FuncParamTypes() []Code {
	if len(f.Params) > 0 {
		return []Code{Op("*").Qual(f.Resource.PackagePath(), f.paramsStructType())}
	} else {
		return nil
	}
}

func (f *Finder) NonErrorFuncReturnParam() Code {
	results := Add(Results).Op("*")
	if f.Metadata != nil {
		results.Add(f.Resource.LocalType(f.returnTypeAliasName()))
	} else {
		results.Add(f.Resource.LocalType(Elements))
	}
	return results
}

func (f *Finder) returnTypeAliasName() string {
	return f.FuncName() + "Elements"
}

func (f *Finder) paramsStructType() string {
	return FindBy + utils.ExportedIdentifier(f.Name) + "Params"
}

func (f *Finder) GenerateCode() *utils.CodeFile {
	c := f.Resource.NewCodeFile("findBy" + utils.ExportedIdentifier(f.Name))

	if len(f.Params) > 0 {
		params := &types.Record{
			NamedType: types.NamedType{
				Identifier: utils.Identifier{
					Name:      f.paramsStructType(),
					Namespace: f.Resource.Namespace,
				},
				Doc: fmt.Sprintf("%s provides the parameters to the %s finder", f.paramsStructType(), f.Name),
			},
			Fields: f.Params,
		}
		c.Code.Add(params.GenerateStruct()).Line().Line().
			Add(params.GenerateQueryParamMarshaler(&f.Name, nil)).Line().Line().
			Add(params.GenerateQueryParamUnmarshaler(nil)).Line().Line().
			Add(params.GeneratePopulateDefaultValues()).Line().Line()
	} else {
		c.Code.Const().Id(f.paramsStructType()).Op("=").Qual(utils.RestLiPackage, "QueryParamsString").
			Call(Lit("q=" + f.Name)).Line()
	}

	if f.Metadata != nil {
		c.Code.Type().Id(f.returnTypeAliasName()).Op("=").
			Add(ElementsWithMetadata).Index(List(f.Return.ReferencedType(), f.Metadata.ReferencedType())).Line().Line()
	}

	f.Resource.addClientFuncDeclarations(c.Code, ClientType, f, func(def *Group) {
		declareRpStruct(f, def)

		name := "Find"
		genericParams := []Code{f.Return.ReferencedType()}
		if f.Metadata != nil {
			name += "WithMetadata"
			genericParams = append(genericParams, f.Metadata.ReferencedType())
		}

		var qp Code
		if len(f.Params) > 0 {
			qp = QueryParams
		} else {
			qp = Id(f.paramsStructType())
		}

		def.Return(Qual(utils.RestLiPackage, name).Index(List(genericParams...)).
			Call(RestLiClientReceiver, Ctx, Rp, qp))
	})

	return c
}

func (f *Finder) RegisterMethod(server, resource, segments Code) Code {
	name := "RegisterFinder"
	if f.Metadata != nil {
		name += "WithMetadata"
	}

	return Qual(utils.RestLiPackage, name).Call(
		Add(server), Add(segments), Lit(f.Name),
		Line().Func().
			Params(registerParams(f)...).
			Params(methodReturnParams(f)...).
			BlockFunc(func(def *Group) {
				def.Return(resource).Dot(f.FuncName()).Call(splatRpAndParams(f)...)
			}),
	)
}
