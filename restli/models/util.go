package models

import (
	"bytes"
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var CommentWrapWidth = 100

func LoadModels(reader io.Reader) ([]*Model, error) {
	snapshot := &struct {
		Models []*Model `json:"models"`
	}{}
	err := json.NewDecoder(reader).Decode(snapshot)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	snapshot.Models = append(snapshot.Models, flattenModels(snapshot.Models)...)
	return snapshot.Models, nil
}

func (m *Model) GenerateModelCode(outputDir, packagePrefix string) (filename string, err error) {
	code, packagePath, typeName := m.generateCode(packagePrefix)
	if code == nil {
		return
	}

	f := jen.NewFilePath(packagePath)
	f.Add(code)
	filename = filepath.Join(outputDir, packagePath, typeName+".go")

	err = write(filename, f)
	return
}

func write(filename string, file *jen.File) error {
	b := bytes.NewBuffer(nil)
	if err := file.Render(b); err != nil {
		return errors.WithStack(err)
	}

	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return errors.WithStack(err)
	}

	if err := ioutil.WriteFile(filename, b.Bytes(), os.ModePerm); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func flattenModels(models []*Model) (innerModels []*Model) {
	for _, m := range models {
		innerModels = append(innerModels, m.InnerModels()...)
	}
	if len(innerModels) > 0 {
		innerModels = append(innerModels, flattenModels(innerModels)...)
	}
	return innerModels
}

func publicFieldName(fieldName string) string {
	return strings.ToUpper(fieldName[:1]) + fieldName[1:]
}

func jsonTag(fieldName string) map[string]string {
	return map[string]string{"json": fieldName}
}

func addWordWrappedComment(statement *jen.Statement, comment string) {
	if comment == "" {
		return
	}

	for len(comment) > CommentWrapWidth {
		index := strings.LastIndexFunc(comment[:CommentWrapWidth], unicode.IsSpace)
		if index > 0 {
			statement.Comment(comment[:index]).Line()
			comment = comment[index+1:]
		} else {
			break
		}
	}

	statement.Comment(comment).Line()
}
