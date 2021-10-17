package resources

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

var total Code = Id("total")

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
	params := []Code{Id("results").Index().Add(f.Return.ReferencedType())}
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
			Doc: fmt.Sprintf("This struct provides the parameters to the %s finder", f.Name),
		},
		Fields: f.Params,
	}
	if f.PagingSupported {
		addPagingContextFields(params)
	}
	c.Code.Add(params.GenerateStruct()).Line().Line()
	c.Code.Add(params.GenerateQueryParamMarshaler(&f.Name, false)).Line().Line()

	f.Resource.addClientFuncDeclarations(c.Code, ClientType, f, func(def *Group) {
		returnsOnErr := []Code{Nil()}
		if f.PagingSupported {
			returnsOnErr = append(returnsOnErr, Nil())
		}
		returnsOnErr = append(returnsOnErr, Err())
		formatQueryUrl(f, def, nil, returnsOnErr...)

		elementType := types.RestliType{Array: f.Return}
		elements := Id("elements")
		def.Var().Add(elements).Add(elementType.GoType())

		var finderReturns []Code
		if f.PagingSupported {
			finderReturns = append(finderReturns, total)
		} else {
			finderReturns = append(finderReturns, Id("_"))
		}
		finderReturns = append(finderReturns, Err())

		def.List(finderReturns...).Op("=").Id(ClientReceiver).Dot("DoFinderRequest").Call(Ctx, Url,
			Func().Params(Add(types.Reader).Add(types.ReaderQual)).Params(Err().Error()).BlockFunc(func(def *Group) {
				types.Reader.ReadArrayFunc(elementType, types.Reader, elements, def)
			}),
		)

		returns := []Code{elements}
		if f.PagingSupported {
			returns = append(returns, total)
		}
		returns = append(returns, Err())
		def.Return(returns...)
	})

	return c
}
