package main

import (
	"log"

	"github.com/PapaCharlie/go-restli/v2/codegen/resources"
	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func main() {
	f := &utils.CodeFile{
		SourceFile:  "https://github.com/PapaCharlie/go-restli/v2/blob/master/codegen/resources/pagingcontext.go",
		PackagePath: resources.PagingContext.PackagePath(),
		PackageRoot: resources.PagingContext.PackageRoot(),
		Filename:    resources.PagingContext.TypeName(),
		Code: Empty().
			Add(resources.PagingContext.GenerateStruct()).Line().Line().
			Add(resources.PagingContext.GenerateEquals()).Line().Line().
			Add(resources.PagingContext.GenerateComputeHash()).Line().Line().
			Add(resources.PagingContext.GenerateQueryParamMarshaler(nil, nil)).Line().Line().
			Add(resources.PagingContext.GenerateQueryParamUnmarshaler(nil)).Line().Line(),
	}

	err := f.Write("..")
	if err != nil {
		log.Panic(err)
	}
}
