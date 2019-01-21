package models

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	. "go-restli/codegen"
	"log"
)

const RecordType = "record"

type Record struct {
	NameAndDoc
	Include []*Model
	Fields  []Field
}

type Field struct {
	NameAndDoc
	Type     *Model          `json:"type"`
	Optional bool            `json:"optional"`
	Default  json.RawMessage `json:"default"`
}

func (r *Record) InnerModels() (models []*Model) {
	models = append(models, r.Include...)
	for _, f := range r.Fields {
		models = append(models, f.Type)
	}
	return
}

func (r *Record) generateCode(packagePrefix string) (def *jen.Statement) {
	def = jen.Empty()

	AddWordWrappedComment(def, r.Doc).Line()

	var fields []jen.Code
	for _, i := range r.Include {
		if ref := i.Reference; ref != nil {
			fields = append(fields, i.GoType(packagePrefix))
			continue
		}
		if rec := i.Record; rec != nil {
			fields = append(fields, i.GoType(packagePrefix))
			continue
		}
		log.Panic("Illegal included type:", i)
	}

	for _, f := range r.Fields {
		field := jen.Empty()
		AddWordWrappedComment(field, f.Doc).Line()
		field.Id(ExportedIdentifier(f.Name))
		field.Add(f.Type.GoType(packagePrefix))
		field.Tag(JsonTag(f.Name))
		fields = append(fields, field)
	}

	def.Type().Id(r.Name).Struct(fields...)

	return
}
