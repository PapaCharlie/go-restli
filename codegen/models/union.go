package models

import (
	"encoding/json"
	"log"
	"strings"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

type UnionModel struct {
	Types []UnionFieldModel
}

type UnionFieldModel struct {
	Model *Model
	Alias string
}

func (u *UnionFieldModel) UnmarshalJSON(data []byte) error {
	m := &Model{}
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	u.Model = m

	type t UnionFieldModel
	if err := json.Unmarshal(data, (*t)(u)); err != nil {
		if !strings.Contains(err.Error(), "json: cannot unmarshal string into Go value of type") {
			return err
		}
	}

	return nil
}

func (u *UnionModel) UnmarshalJSON(data []byte) error {
	var types []UnionFieldModel
	if err := json.Unmarshal(data, &types); err != nil {
		return err
	}
	u.Types = types
	return nil
}

func (u *UnionModel) InnerModels() (models []*Model) {
	for _, t := range u.Types {
		models = append(models, t.Model)
	}
	return
}

func (u *UnionModel) GoType() (def *Statement) {
	if len(u.Types) == 0 {
		log.Panicln("Empty union", u)
	}

	return StructFunc(func(def *Group) {
		for _, t := range u.Types {
			var tag FieldTag
			tag.Json.Name = t.alias()
			tag.Json.Optional = true

			field := def.Empty()
			AddWordWrappedComment(field, t.Model.Doc).Line()
			field.Id(t.name())
			field.Op("*").Add(t.Model.GoType())
			field.Tag(tag.ToMap())
		}
	})
}

func (u *UnionFieldModel) name() string {
	alias := u.alias()
	return ExportedIdentifier(alias[strings.LastIndex(alias, ".")+1:])
}

func (u *UnionFieldModel) alias() string {
	if u.Alias != "" {
		return u.Alias
	}
	if u.Model.Primitive != nil {
		return u.Model.Primitive[0]
	}
	if u.Model.Bytes != nil {
		return "bytes"
	}
	if u.Model.Fixed != nil {
		return FixedModelTypeName
	}
	if u.Model.Array != nil {
		return ArrayModelTypeName
	}
	if u.Model.Map != nil {
		return MapModelTypeName
	}
	return u.Model.Namespace + "." + u.Model.Name
}