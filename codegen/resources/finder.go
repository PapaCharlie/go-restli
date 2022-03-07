package resources

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

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

func (f *Finder) NonErrorFuncReturnParam() Code {
	results := Id("results").Op("*")
	if f.Metadata != nil {
		results.Qual(utils.ProtocolPackage, "FinderResultsWithMetadata").Index(List(f.Return.ReferencedType(), f.Metadata.ReferencedType()))
	} else {
		results.Qual(utils.ProtocolPackage, "FinderResults").Index(f.Return.ReferencedType())
	}
	return results
}

func (f *Finder) paramsStructType() string {
	return FindBy + utils.ExportedIdentifier(f.Name) + "Params"
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
		declareRpStruct(f, def)

		if f.Metadata != nil {
			def.Return(Qual(utils.ProtocolPackage, "FindWithMetadata").Call(
				Op("&").Id(ClientReceiver).Dot(CollectionClient),
				Ctx,
				Rp,
				QueryParams,
				types.ReaderUtils.UnmarshalerFunc(*f.Metadata),
			))
		} else {
			def.Return(Id(ClientReceiver).Dot("Find").Call(Ctx, Rp, QueryParams))
		}
	})

	return c
}
