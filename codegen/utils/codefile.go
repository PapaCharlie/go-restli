package utils

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

const (
	ReadOnlyPermissions = os.FileMode(0444)
	GeneratedFileSuffix = ".gr.go"
	ParsedSpecsFile     = "parsed-specs.gr.json"
)

var (
	Logger = log.New(os.Stderr, "[go-restli] ", log.LstdFlags|log.Lshortfile)

	PackagePrefix string

	HeaderTemplate = template.Must(template.New("header").Parse(`DO NOT EDIT

Code automatically generated by github.com/PapaCharlie/go-restli
Source file: {{.SourceFile}}`))
)

type CodeFile struct {
	SourceFile  string
	PackagePath string
	Filename    string
	Code        *Statement
}

func (f *CodeFile) Write(outputDir string, writeInPackageDirs bool) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.Errorf("go-restli: Could not generate model: %+v", e)
		}
	}()
	file := NewFilePathName(f.PackagePath, strings.ToLower(filepath.Base(f.PackagePath)))

	header := bytes.NewBuffer(nil)
	err = HeaderTemplate.Execute(header, f)
	if err != nil {
		return err
	}

	file.HeaderComment(header.String())
	file.Add(f.Code)

	var filename string
	if writeInPackageDirs {
		relpath, err := filepath.Rel(PackagePrefix, f.PackagePath)
		if err != nil {
			return err
		}

		filename = filepath.Join(outputDir, relpath, f.Filename+GeneratedFileSuffix)
	} else {
		filename = filepath.Join(outputDir, f.Filename+GeneratedFileSuffix)
	}

	err = WriteJenFile(filename, file)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Failed to write code file to %q", filename)
	}

	return nil
}

func (f *CodeFile) Identifier() string {
	return f.PackagePath + "." + f.Filename
}

func WriteJenFile(filename string, file *File) error {
	b := bytes.NewBuffer(nil)
	if err := file.Render(b); err != nil {
		return errors.WithStack(err)
	}

	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return errors.WithStack(err)
	}

	_ = os.Remove(filename)

	if err := ioutil.WriteFile(filename, b.Bytes(), ReadOnlyPermissions); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func CleanTargetDir(targetDir string) (err error) {
	var cleanTargetDir func(targetDir string) error
	cleanTargetDir = func(targetDir string) (err error) {
		children, err := os.ReadDir(targetDir)
		if err != nil {
			return err
		}

		if len(children) == 0 {
			return os.Remove(targetDir)
		}

		for _, c := range children {
			sub := filepath.Join(targetDir, c.Name())
			if c.IsDir() {
				err = CleanTargetDir(sub)
			} else {
				if strings.HasSuffix(c.Name(), GeneratedFileSuffix) {
					err = os.Remove(sub)
				}
			}
			if err != nil {
				return err
			}
		}

		children, err = os.ReadDir(targetDir)
		if err != nil {
			return err
		}

		if len(children) == 0 {
			return os.Remove(targetDir)
		}
		return nil
	}

	if _, err = os.Stat(targetDir); os.IsNotExist(err) {
		return nil
	} else {
		err = os.Remove(filepath.Join(targetDir, ParsedSpecsFile))
		if err == nil || os.IsNotExist(err) {
			return cleanTargetDir(targetDir)
		} else {
			return err
		}
	}
}

func AddWordWrappedComment(code *Statement, comment string) *Statement {
	if comment != "" {
		// TODO: Find a way to pretty-wrap comments to 120 chars
		code.Comment(comment)
		return code
	} else {
		return code
	}
}

func ExportedIdentifier(identifier string) string {
	buf := new(strings.Builder)
	for i, c := range identifier {
		switch {
		case unicode.IsLetter(c):
			if i == 0 {
				buf.WriteRune(unicode.ToUpper(c))
			} else {
				buf.WriteRune(c)
			}
		case unicode.IsNumber(c):
			if i == 0 {
				buf.WriteString("Exported_")
			}
			buf.WriteRune(c)
		case c == '_':
			if i == 0 {
				buf.WriteString("Exported")
			}
			buf.WriteRune(c)
		// Because $ is a valid identifier character in Java, it technically does not cause compile error and is
		// therefore occasionally used. To support this usecase, explicitly handle $ by replacing it with DOLLAR. All
		// other non-alphanumeric (plus _) characters are considered illegal
		case c == '$':
			if i != 0 {
				buf.WriteRune('_')
			}
			buf.WriteString("DOLLAR_")
		default:
			log.Panicf("Illegal identifier character %q in %q", c, identifier)
		}
	}
	return buf.String()
}

func ReceiverName(typeName string) string {
	return strings.ToLower(typeName[:1])
}

type ShouldUsePointer int

const (
	No  = ShouldUsePointer(0)
	Yes = ShouldUsePointer(1)
)

func (p ShouldUsePointer) ShouldUsePointer() bool {
	return p == Yes
}

func AddFuncOnReceiver(def *Statement, receiver, typeName, funcName string, pointer ShouldUsePointer) *Statement {
	r := Id(receiver)
	if pointer.ShouldUsePointer() {
		r.Op("*")
	}
	r.Id(typeName)
	return def.Func().Params(r).Id(funcName)
}

func AddStringer(def *Statement, receiver, typeName string, pointer ShouldUsePointer, f func(def *Group)) *Statement {
	return AddFuncOnReceiver(def, receiver, typeName, "String", pointer).
		Params().
		String().
		BlockFunc(f)
}

func AddPointer(def *Statement, receiver, typeName string) *Statement {
	def.Comment("Pointer returns a pointer to the given receiver, useful for inlining setting optional fields.").Line()
	def.Func().
		Params(Id(receiver).Id(typeName)).
		Id("Pointer").Params().
		Op("*").Id(typeName).
		BlockFunc(func(def *Group) {
			def.Return(Op("&").Id(receiver))
		}).Line().Line()
	return def
}

func IfErrReturn(results ...Code) *Statement {
	return If(Err().Op("!=").Nil()).Block(Return(results...))
}

func JsonFieldTag(name string, optional bool) map[string]string {
	tags := map[string]string{"json": name}
	if optional {
		tags["json"] += ",omitempty"
	}
	return tags
}

func OrderedValues(f func(add func(key, value Code))) *Statement {
	return ValuesFunc(func(def *Group) {
		f(func(key, value Code) {
			def.Line().Add(key).Op(":").Add(value)
		})
		def.Line()
	})
}
