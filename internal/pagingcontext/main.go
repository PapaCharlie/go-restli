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

	files := []*utils.CodeFile{{
		SourceFile:  "https://github.com/PapaCharlie/go-restli/blob/master/codegen/resources/pagingcontext.go",
		PackagePath: resources.PagingContext.Namespace,
		Filename:    resources.PagingContext.Name,
		Code: Empty().
			Add(resources.PagingContext.GenerateStruct()).Line().Line().
			Add(resources.PagingContext.GenerateEquals()).Line().Line().
			Add(resources.PagingContext.GenerateComputeHash()).Line().Line().
			Add(resources.PagingContext.GenerateQueryParamMarshaler(nil, nil)).Line().Line(),
	}}

	for _, f := range files {
		err := f.Write("..", false)
		if err != nil {
			log.Panic(err)
		}
	}
}
