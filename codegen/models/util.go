package models

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	. "go-restli/codegen"
	"io"
	"log"
	"strings"
)

func LoadModels(reader io.Reader) ([]*Model, error) {
	snapshot := &struct {
		Models map[string]*Model `json:"models"`
	}{}
	err := json.NewDecoder(reader).Decode(snapshot)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var models []*Model
	for _, m := range snapshot.Models {
		models = append(models, m)
	}

	models = append(models, flattenModels(models)...)
	replaceReferences(models)
	return models, nil
}

func flattenModels(models []*Model) (innerModels []*Model) {
	for _, m := range models {
		m.register()
		for _, im := range m.InnerModels() {
			im.register()
			innerModels = append(innerModels, im)
		}
	}
	if len(innerModels) > 0 {
		innerModels = append(innerModels, flattenModels(innerModels)...)
	}
	return innerModels
}

func replaceReferences(models []*Model) {
	for _, m := range models {
		if m.Reference != nil {
			*m = *GetRegisteredModel(m.Ns, m.Name)
		}
	}
}

func escapeNamespace(namespace string) string {
	return strings.Replace(namespace, "internal", "_internal", -1)
}

var loadedModels = make(map[string]*Model)

func (m *Model) register() {
	if m.Primitive != nil || m.Reference != nil {
		return
	}

	loadedModels[m.PackagePath()+"."+m.Name] = m
}

func GetRegisteredModel(ns Ns, name string) *Model {
	return loadedModels[ns.PackagePath()+"."+name]
}

func SetDefaultValue(def *jen.Group, receiver, name, rawJson string, model *Model) {
	def.If(jen.Id(receiver).Dot(name).Op("==").Nil()).BlockFunc(func(def *jen.Group) {
		// Special case for primitives, instead of parsing them from JSON every time, we can leave them as literals
		if model.Primitive != nil {
			def.Id("v").Op(":=").Lit(model.Primitive.GetLit(rawJson))
			def.Id(receiver).Dot(name).Op("=").Op("&").Id("v")
			return
		}

		// Empty arrays and maps can be initialized directly, regardless of type
		if (model.Array != nil && rawJson == "[]") || (model.Map != nil && rawJson == "{}") {
			def.Id(receiver).Dot(name).Op("=").Make(model.GoType(), jen.Lit(0))
			return
		}

		// Enum values can also be added as literals
		if model.Enum != nil {
			var v string
			err := json.Unmarshal([]byte(rawJson), &v)
			if err != nil {
				log.Panicln("illegal enum", err)
			}
			def.Id("v").Op(":=").Qual(model.PackagePath(), model.Enum.SymbolIdentifier(v))
			def.Id(receiver).Dot(name).Op("= &").Id("v")
			return
		}

		if !model.IsMapOrArray() {
			def.Id(receiver).Dot(name).Op("=").New(model.GoType())
		}

		field := jen.Empty()
		if model.IsMapOrArray() {
			field.Op("&")
		}
		field.Id(receiver).Dot(name)

		def.Err().Op(":=").Qual(EncodingJson, Unmarshal).Call(jen.Index().Byte().Call(jen.Lit(rawJson)), field)
		def.If(jen.Err().Op("!=").Nil()).Block(jen.Qual("log", "Panicln").Call(jen.Lit("Illegal default value"), jen.Err()))
	})
}
