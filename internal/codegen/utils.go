package codegen

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

	_ = os.MkdirAll(outputDir, os.ModePerm)
	parsedSpecs := filepath.Join(outputDir, "parsed-specs.json")
	_ = os.Remove(parsedSpecs)
	err := ioutil.WriteFile(parsedSpecs, specBytes, ReadOnlyPermissions)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Failed to write parsed specs to %q", parsedSpecs)
	}

	// Use a Reader regardless since it'll handle leading/trailing whitespace and other niceties
	err = json.NewDecoder(bytes.NewBuffer(specBytes)).Decode(&schemas)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Could not deserialize GoRestliSpec")
	}

	tmpOutputDir, err := ioutil.TempDir("", "go-restli_*")
	if err != nil {
		return errors.Wrapf(err, "go-restli: Failed to create temporary directory")
	}
	defer os.RemoveAll(tmpOutputDir)

	codeFiles := append(TypeRegistry.GenerateTypeCode(), schemas.GenerateClientCode()...)

	for _, code := range codeFiles {
		err = code.Write(tmpOutputDir)
		if err != nil {
			return errors.Wrapf(err, "go-restli: Could not generate code for %+v:\n%s", code, code.Code.GoString())
		}
	}

	err = GenerateAllImportsTest(tmpOutputDir, codeFiles)
	if err != nil {
		return err
	}

	children, err := ioutil.ReadDir(tmpOutputDir)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Could not list %q", tmpOutputDir)
	}

	for _, c := range children {
		source := filepath.Join(tmpOutputDir, c.Name())
		destination := filepath.Join(outputDir, c.Name())

		err = os.RemoveAll(destination)
		if err != nil {
			return errors.Wrapf(err, "go-restli: Failed to delete %q", destination)
		}

		err = os.Rename(source, destination)
		if err != nil {
			return errors.Wrapf(err, "go-restli: Failed to move %q to %q", source, destination)
		}
	}

	return nil
}

func GenerateAllImportsTest(outputDir string, codeFiles []*CodeFile) error {
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

	out := filepath.Join(outputDir, "all_imports_test.go")
	_ = os.Remove(out)
	err := WriteJenFile(out, f)
	if err != nil {
		return errors.Wrapf(err, "Could not write all imports file: %+v", err)
	}
	return nil
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
