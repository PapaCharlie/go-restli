package main

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen/resources"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func main() {
	for i := range resources.PagingContext.Fields {
		resources.PagingContext.Fields[i].IncludedFrom = nil
	}

	files := []*utils.CodeFile{
		{
			SourceFile:  "https://github.com/PapaCharlie/go-restli/blob/master/internal/codegen/resources/pagingcontext.go",
			PackagePath: resources.PagingContext.Namespace,
			Filename:    resources.PagingContext.Name,
			Code: Empty().
				Add(resources.PagingContext.GenerateStruct()).Line().Line().
				Add(resources.PagingContext.GenerateEquals()).Line().Line().
				Add(resources.PagingContext.GenerateComputeHash()).Line().Line().
				Add(resources.PagingContext.GenerateQueryParamMarshaler(nil, nil)).Line().Line().
				Add(GenerateNewPagingContext()).Line().Line(),
		},
		{
			SourceFile:  "https://github.com/PapaCharlie/go-restli/blob/master/internal/codegen/resources/errorresponse.go",
			PackagePath: resources.ErrorResponse.Namespace,
			Filename:    resources.ErrorResponse.Name,
			Code: Empty().
				Add(resources.ErrorResponse.GenerateStruct()).Line().Line().
				Add(resources.ErrorResponse.GenerateEquals()).Line().Line().
				Add(resources.ErrorResponse.GenerateComputeHash()).Line().Line().
				Add(resources.ErrorResponse.GenerateMarshalRestLi()).Line().Line().
				Add(resources.ErrorResponse.GenerateUnmarshalRestLi()).Line().Line().
				Add(resources.ErrorResponse.GenerateUnmarshalerFunc()).Line().Line(),
		},
	}

	for _, f := range files {
		err := f.Write("protocol/stdstructs", false)
		if err != nil {
			log.Panic(err)
		}
	}
}

func GenerateNewPagingContext() Code {
	start, count := Id("start"), Id("count")
	return Func().Id("NewPagingContext").
		Params(Add(start), Add(count).Int32()).
		Add(utils.PagingContextIdentifier.Qual()).
		Block(
			Return(Add(utils.PagingContextIdentifier.Qual()).Values(Dict{
				Id("Start"): Op("&").Add(start),
				Id("Count"): Op("&").Add(count),
			})),
		)
}
