package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"unicode"

	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

const (
	RestLiClient = "RestLiClient"

	FormatQueryUrl     = "FormatQueryUrl"
	ResourcePath       = "ResourcePath"
	ResourceEntityPath = "ResourceEntityPath"
	DoAndIgnore        = "DoAndIgnore"
	DoAndDecode        = "DoAndDecode"
	DoAndDecodeResult  = "doAndDecodeResult"

	WithContext = "WithContext"
	FindBy      = "FindBy"

	ReqVar     = "req"
	ResVar     = "res"
	UrlVar     = "url"
	PathVar    = "path"
	ContextVar = "ctx"

	ClientReceiver      = "c"
	ClientType          = "client"
	ClientInterfaceType = "Client"
)

var Logger = log.New(os.Stderr, "[go-restli] ", log.LstdFlags|log.Lshortfile)

func GenerateCode(specBytes []byte, outputDir string) error {
	var schemas GoRestliSpec

	// Use a Decode regardless since it'll handle leading/trailing whitespace and other niceties
	err := json.NewDecoder(bytes.NewBuffer(specBytes)).Decode(&schemas)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Could not deserialize GoRestliSpec")
	}

	codeFiles := append(TypeRegistry.GenerateTypeCode(), schemas.GenerateClientCode()...)

	for _, code := range codeFiles {
		file, err := code.Write(outputDir)
		if err != nil {
			return errors.Wrapf(err, "go-restli: Could not generate code for %+v:\n%s", code, code.Code.GoString())
		} else {
			fmt.Println(file)
		}
	}

	return GenerateAllImportsFile(outputDir, codeFiles)
}

func GenerateAllImportsFile(outputDir string, codeFiles []*CodeFile) error {
	imports := make(map[string]bool)
	for _, code := range codeFiles {
		if code == nil {
			continue
		}
		imports[code.PackagePath] = true
	}
	f := NewFile("main")
	for p := range imports {
		f.Anon(p)
	}
	f.Func().Id("TestAllImports").Params(Op("*").Qual("testing", "T")).Block()

	err := WriteJenFile(filepath.Join(outputDir, PackagePrefix, "all_imports_test.go"), f)
	if err != nil {
		return errors.Wrapf(err, "Could not write all imports file: %+v", err)
	}
	return nil
}

func writeStringToBuf(def *Group, s *Statement) *Statement {
	return def.Id("buf").Dot("WriteString").Call(s)
}

func writeArrayToBuf(def *Group, accessor *Statement, items *RestliType, returnOnError ...Code) {
	writeStringToBuf(def, Lit("List("))

	def.For(List(Id("idx"), Id("val")).Op(":=").Range().Add(accessor)).BlockFunc(func(def *Group) {
		def.If(Id("idx").Op("!=").Lit(0)).Block(Id("buf").Dot("WriteByte").Call(LitRune(','))).Line()
		items.WriteToBuf(def, Id("val"), returnOnError...)
	})

	def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
}

func writeMapToBuf(def *Group, accessor *Statement, values *RestliType, returnOnError ...Code) {
	def.Id("buf").Dot("WriteByte").Call(LitRune('('))

	def.Id("idx").Op(":=").Lit(0)
	def.For(List(Id("key"), Id("val")).Op(":=").Range().Add(accessor)).BlockFunc(func(def *Group) {
		def.If(Id("idx").Op("!=").Lit(0)).Block(Id("buf").Dot("WriteByte").Call(LitRune(','))).Line()
		def.Id("idx").Op("++")
		writeStringToBuf(def, Id(Codec).Dot("EncodeString").Call(Id("key")))
		def.Id("buf").Dot("WriteByte").Call(LitRune(':'))
		values.WriteToBuf(def, Id("val"), returnOnError...)
	})

	def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
}

func canonicalizeAccessor(accessor *Statement) string {
	label := ExportedIdentifier(accessor.GoString())
	for i, c := range label {
		if !(unicode.IsDigit(c) || unicode.IsLetter(c)) {
			label = label[:i] + "_" + label[i+1:]
		}
	}
	return label
}

func (r *Resource) methodFuncName(m *Method, withContext bool) string {
	var name string

	switch m.MethodType {
	case REST_METHOD:
		name = m.restMethodFuncName()
	case ACTION:
		name = m.actionFuncName()
	case FINDER:
		name = m.finderFuncName()
	}

	if withContext {
		name += WithContext
	}

	return name
}

func (r *Resource) methodFuncParams(m *Method, def *Group) {
	switch m.MethodType {
	case REST_METHOD:
		r.restMethodFuncParams(m, def)
	case ACTION:
		r.actionFuncParams(m, def)
	case FINDER:
		r.finderFuncParams(m, def)
	}
}

func (r *Resource) methodReturnParams(m *Method) func(*Group) {
	var returnParams func(*Group)

	switch m.MethodType {
	case REST_METHOD:
		returnParams = m.restMethodFuncReturnParams
	case ACTION:
		returnParams = m.actionFuncReturnParams
	case FINDER:
		returnParams = m.finderFuncReturnParams
	}

	return returnParams
}

func (m *Method) methodCallParams() []Code {
	var methodCallParams []Code

	switch m.MethodType {
	case REST_METHOD:
		methodCallParams = m.restMethodCallParams()
	case ACTION:
		methodCallParams = m.actionMethodCallParams()
	case FINDER:
		methodCallParams = m.finderMethodCallParams()
	}

	return methodCallParams
}

func (r *Resource) clientFuncDeclaration(m *Method, withContext bool) *Statement {
	params := func(def *Group) {
		if withContext {
			def.Id("ctx").Qual("context", "Context")
		}
		r.methodFuncParams(m, def)
	}

	return Id(r.methodFuncName(m, withContext)).ParamsFunc(params).ParamsFunc(r.methodReturnParams(m))
}

func (r *Resource) addClientFuncDeclarations(def *Statement, clientType string, m *Method, block func(*Group)) *Statement {
	clientFuncDeclarationStart := Func().Params(Id(ClientReceiver).Op("*").Id(clientType))

	AddWordWrappedComment(def, m.Doc).Line().
		Add(clientFuncDeclarationStart).
		Add(r.clientFuncDeclaration(m, false)).
		Block(Return(Id(ClientReceiver).Dot(r.methodFuncName(m, true)).CallFunc(func(def *Group) {
			def.Qual("context", "Background").Call()
			for _, p := range append(m.entityParams(), m.methodCallParams()...) {
				def.Add(p)
			}
		}))).
		Line().Line()

	AddWordWrappedComment(def, m.Doc).Line().
		Add(clientFuncDeclarationStart).
		Add(r.clientFuncDeclaration(m, true)).
		BlockFunc(block)

	return def
}

func callDoAndDecode(def *Group, accessor *Statement, zeroValueInCaseOfError *Statement) {
	def.List(Id("_"), Err()).Op("=").Id(ClientReceiver).Dot(DoAndDecode).Call(Id(ReqVar), accessor)
	IfErrReturn(def, zeroValueInCaseOfError, Err()).Line()
}
