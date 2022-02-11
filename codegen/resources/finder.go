package resources

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

var total Code = Id("total")
var results Code = Id("results")

type Finder struct{ methodImplementation }

func (f *Finder) IsSupported() bool {
	return true
}

func (f *Finder) FuncName() string {
	return FindBy + utils.ExportedIdentifier(f.Name)
}

func (f *Finder) FuncParamNames() []Code {
	return []Code{QueryParams}
}

func (f *Finder) FuncParamTypes() []Code {
	return []Code{Op("*").Qual(f.Resource.PackagePath(), f.paramsStructType())}
}

func (f *Finder) NonErrorFuncReturnParams() []Code {
	params := []Code{Add(results).Index().Add(f.Return.ReferencedType())}
	if f.PagingSupported {
		params = append(params, Add(total).Op("*").Int())
	}
	return params
}

func (f *Finder) paramsStructType() string {
	return FindBy + utils.ExportedIdentifier(f.Name) + "Params"
}

func (f *Finder) resultsStructType() string {
	return FindBy + utils.ExportedIdentifier(f.Name) + "Results"
}

func (f *Finder) GenerateCode() *utils.CodeFile {
	c := f.Resource.NewCodeFile("findBy" + utils.ExportedIdentifier(f.Name))

	c.Code.Const().Id(utils.ExportedIdentifier(FindBy + utils.ExportedIdentifier(f.Name))).Op("=").Lit(f.Name).Line()

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
	if f.PagingSupported {
		addPagingContextFields(params)
	}
	c.Code.Add(params.GenerateStruct()).Line().Line()
	c.Code.Add(params.GenerateQueryParamMarshaler(&f.Name, nil)).Line().Line()

	f.Resource.addClientFuncDeclarations(c.Code, ClientType, f, func(def *Group) {
		returnsOnErr := []Code{Nil()}
		if f.PagingSupported {
			returnsOnErr = append(returnsOnErr, Nil())
		}
		returnsOnErr = append(returnsOnErr, Err())
		formatQueryUrl(f, def, returnsOnErr...)

		call := Qual(utils.ProtocolPackage, "DoFinderRequest").Call(
			RestLiClientReceiver,
			Ctx,
			Url,
			types.Reader.UnmarshalerFunc(*f.Return),
		)
		if f.PagingSupported {
			def.Return(call)
		} else {
			def.List(results, Id("_"), Err()).Op("=").Add(call)
			def.Return(results, Err())
		}
	})

	return c
}
