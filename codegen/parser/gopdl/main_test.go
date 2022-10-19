package main

import (
	"reflect"
	"testing"

	"github.com/PapaCharlie/go-restli/codegen/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

func TestShit(t *testing.T) {
	in := loadPDLs("/Users/pchesnai/code/personal/go-restli/internal/tests/testdata/extra-test-suite/schemas/extras/SinglePrimitiveField.pdl")[0]
	lexer := parser.NewPdlLexer(in)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewPdlParser(stream)
	 new(parser.ErrorListener)
	p.AddErrorListener()
	p.BuildParseTrees = true
	tree := p.Document()
	if
	for _, child := range tree.GetChildren() {
		t.Log(reflect.ValueOf(child).Type().String())
	}
	t.Log(tree.GetChildren()[2].(*parser.TypeDeclarationContext).GetText())
}
