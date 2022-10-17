package utils

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

const (
	ReadOnlyPermissions = os.FileMode(0444)
	GeneratedFileSuffix = ".gr.go"
	ManifestFile        = "go-restli-manifest.gr.json"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("[go-restli] ")
}

var (
	// HeaderTemplate is a template for the header of generated files. According to golang.org/s/generatedcode,
	// generated files must have a specific header syntax.
	HeaderTemplate = template.Must(template.New("header").Parse(
		`Code generated by "github.com/PapaCharlie/go-restli"; DO NOT EDIT.

Source file: {{.SourceFile}}`))
)

type CodeFile struct {
	SourceFile  string
	PackagePath string
	PackageRoot string
	Filename    string
	Code        *Statement
}

//go:generate go run ../../internal/importnames

var importsRegex = regexp.MustCompile(`[^a-z0-9]`)

func PackageName(pkg string) string {
	pkg = strings.ToLower(path.Base(pkg))
	pkg = importsRegex.ReplaceAllString(pkg, "")
	return pkg
}

func addImportName(pkg string) {
	if _, ok := jenImportNames[pkg]; !ok {
		jenImportNames[pkg] = PackageName(pkg)
	}
}

func (f *CodeFile) Write(outputDir string, generateWithPackageRoot bool) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.Errorf("go-restli: Could not generate model: %+v", e)
		}
	}()
	file := NewFilePathName(f.PackagePath, PackageName(f.PackagePath))
	file.ImportNames(jenImportNames)

	header := bytes.NewBuffer(nil)
	err = HeaderTemplate.Execute(header, f)
	if err != nil {
		return err
	}

	file.HeaderComment(header.String())
	file.Add(f.Code)

	filename := filepath.Join(f.PackagePath, f.Filename+GeneratedFileSuffix)
	if !generateWithPackageRoot {
		filename = strings.TrimPrefix(filename, f.PackageRoot)
	}
	filename = filepath.Join(outputDir, filename)

	err = WriteJenFile(filename, file)
	if err != nil {
		return errors.Wrapf(err, "go-restli: Failed to write code file to %q", filename)
	}

	return nil
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
			if targetDir != "." {
				return os.Remove(targetDir)
			} else {
				return nil
			}
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
			if targetDir != "." {
				return os.Remove(targetDir)
			} else {
				return nil
			}
		}
		return nil
	}

	if _, err = os.Stat(targetDir); os.IsNotExist(err) {
		return nil
	} else {
		err = os.Remove(filepath.Join(targetDir, ManifestFile))
		if err == nil || os.IsNotExist(err) {
			return cleanTargetDir(targetDir)
		} else {
			return err
		}
	}
}

func AddWordWrappedComment(code *Statement, comment string) *Statement {
	if comment != "" {
		// By splitting the comment by string, we prevent jen's default behavior of formatting multiline comments
		// using /* */, which ends up not being well formatted according to gofmt.
		for i, line := range strings.Split(strings.TrimSpace(comment), "\n") {
			if i != 0 {
				code.Line()
			}
			code.Comment(line)
		}
		// TODO: Find a way to pretty-wrap comments to 120 chars
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
	Yes = ShouldUsePointer(0)
	No  = ShouldUsePointer(1)
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
