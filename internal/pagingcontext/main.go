package main

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen/resources"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/dave/jennifer/jen"
)

func main() {
	for i := range resources.PagingContext.Fields {
		resources.PagingContext.Fields[i].IncludedFrom = nil
	}

	code := jen.Empty().
		Add(resources.PagingContext.GenerateStruct()).Line().Line().
		Add(resources.PagingContext.GenerateEquals()).Line().Line().
		Add(resources.PagingContext.GenerateComputeHash()).Line().Line().
		Add(resources.PagingContext.GenerateQueryParamMarshaler(nil, false)).Line().Line()

	pagingContext := &utils.CodeFile{
		SourceFile:  "https://github.com/PapaCharlie/go-restli/blob/master/internal/codegen/resources/pagingcontext.go",
		PackagePath: "protocol",
		Filename:    resources.PagingContext.Name,
		Code:        code,
	}
	err := pagingContext.Write(".")
	if err != nil {
		log.Panic(err)
	}
}
