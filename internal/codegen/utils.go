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

	FindBy = "FindBy"

	ReqVar  = "req"
	ResVar  = "res"
	UrlVar  = "url"
	PathVar = "path"

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

func canonicalizeAccessor(accessor *Statement) string {
	label := ExportedIdentifier(accessor.GoString())
	for i, c := range label {
		if !(unicode.IsDigit(c) || unicode.IsLetter(c)) {
			label = label[:i] + "_" + label[i+1:]
		}
	}
	return label
}

func (r *Resource) clientFunc(m *Method) *Statement {
	var name string
	var params func(*Group)
	var returnParams func(*Group)

	switch m.MethodType {
	case REST_METHOD:
		name = m.restMethodFuncName()
		params = func(def *Group) { m.restMethodFuncParams(def, r.ResourceSchema) }
		returnParams = m.restMethodFuncReturnParams
	case ACTION:
		name = m.actionFuncName()
		params = m.actionFuncParams
		returnParams = m.actionFuncReturnParams
	case FINDER:
		name = m.finderFuncName()
		params = m.finderFuncParams
		returnParams = m.finderFuncReturnParams
	}

	return Id(name).ParamsFunc(params).ParamsFunc(returnParams)
}

func (r *Resource) addClientFunc(def *Statement, m *Method) *Statement {
	return def.Func().Params(Id(ClientReceiver).Op("*").Id(ClientType)).Add(r.clientFunc(m))
}

func callDoAndDecode(def *Group) {
	def.List(Id("_"), Err()).Op("=").Id(ClientReceiver).Dot(DoAndDecode).Call(Id(ReqVar), Op("&").Id(DoAndDecodeResult))
	IfErrReturn(def, Nil(), Err()).Line()
}
