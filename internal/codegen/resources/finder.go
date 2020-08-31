package resources

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/internal/codegen/types"
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
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

func (f *Finder) NonErrorFuncReturnParams() []Code {
	return []Code{Id("results").Index().Add(f.Return.ReferencedType())}
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
	c.Code.Add(params.GenerateStruct()).Line().Line()
	c.Code.Add(params.GenerateQueryParamMarshaler(&f.Name)).Line().Line()

	results := &types.Record{
		NamedType: types.NamedType{
			Identifier: utils.Identifier{
				Name:      f.resultsStructType(),
				Namespace: f.Resource.Namespace,
			},
			Doc: fmt.Sprintf("This struct deserializes the response from the %s finder", f.Name),
		},
		Fields: []types.Field{{
			Type: types.RestliType{Array: f.Return},
			Name: "elements",
		}},
	}
	c.Code.Add(results.GenerateStruct()).Line().Line()
	c.Code.Add(results.GenerateUnmarshalRestLi()).Line().Line()

	f.Resource.addClientFuncDeclarations(c.Code, ClientType, f, func(def *Group) {
		formatQueryUrl(f, def, Nil(), Err())

		accessor := Id("elements")
		def.Var().Add(accessor).Id(f.resultsStructType())

		def.Err().Op("=").Id(ClientReceiver).Dot("DoFinderRequest").Call(Ctx, Url, Op("&").Add(accessor))
		def.Add(utils.IfErrReturn(Nil(), Err())).Line()

		def.Return(Add(accessor).Dot("Elements"), Nil())
	})

	return c
}
