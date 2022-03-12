package main

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen/resources"
	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func main() {
	for i := range resources.PagingContext.Fields {
		resources.PagingContext.Fields[i].IncludedFrom = nil
	}

	files := []*utils.CodeFile{{
		SourceFile:  "https://github.com/PapaCharlie/go-restli/blob/master/internal/codegen/resources/pagingcontext.go",
		PackagePath: resources.PagingContext.Namespace,
		Filename:    resources.PagingContext.Name,
		Code: Empty().
			Add(resources.PagingContext.GenerateStruct()).Line().Line().
			Add(resources.PagingContext.GenerateEquals()).Line().Line().
			Add(resources.PagingContext.GenerateComputeHash()).Line().Line().
			Add(resources.PagingContext.GenerateQueryParamMarshaler(nil, nil)).Line().Line(),
	}}

	for _, r := range []*types.Record{resources.ErrorResponse, resources.CollectionMetadata, resources.Link} {
		r.GenerateCode()
		files = append(files, &utils.CodeFile{
			SourceFile:  "https://github.com/PapaCharlie/go-restli/blob/master/internal/codegen/resources/stdtypes.go",
			PackagePath: r.Namespace,
			Filename:    r.Name,
			Code: Empty().
				Add(r.GenerateStruct()).Line().Line().
				Add(r.GenerateEquals()).Line().Line().
				Add(r.GenerateComputeHash()).Line().Line().
				Add(r.GenerateMarshalRestLi()).Line().Line().
				Add(r.GeneratePopulateDefaultValues()).Line().Line().
				Add(r.GenerateUnmarshalRestLi()).Line().Line(),
		})
	}

	for _, f := range files {
		err := f.Write("protocol/stdtypes", false)
		if err != nil {
			log.Panic(err)
		}
	}
}
