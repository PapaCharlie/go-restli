package models

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	. "go-restli/codegen"
	"strings"
)

type Union struct {
	Types []UnionFieldType
}

type UnionFieldType struct {
	Model  *Model
	Alias string //`json:"alias"`
	Thing bool
}

func (u *UnionFieldType) UnmarshalJSON(data []byte) error {
	m := &Model{}
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	u.Model = m

	type t UnionFieldType
	if err := json.Unmarshal(data, (*t)(u)); err != nil {
		if !strings.Contains(err.Error(), "json: cannot unmarshal string into Go value of type") {
			return err
		}
	}

	return nil
}

func (u *Union) UnmarshalJSON(data []byte) error {
	var types []UnionFieldType
	if err := json.Unmarshal(data, &types); err != nil {
		return err
	}
	u.Types = types
	return nil
}

func (u *Union) InnerModels() (models []*Model) {
	for _, t := range u.Types {
		models = append(models, t.Model)
	}
	return
}

func (u *Union) generateCode(packagePrefix string) (def *jen.Statement) {
	return
}

func (u *Union) GoType(packagePrefix string) *jen.Statement {
	if len(u.Types) == 0 {
		panic(u)
	}
	var fields []jen.Code
	for _, t := range u.Types {
		def := jen.Empty()
		AddWordWrappedComment(def, t.Model.Doc).Line()
		def.Id(t.name())
		def.Op("*").Add(t.Model.GoType(packagePrefix))
		def.Tag(JsonTag(t.alias()))
		fields = append(fields, def)
	}
	return jen.Struct(fields...)
}

func (u *UnionFieldType) name() string {
	alias := u.alias()
	alias = strings.ToUpper(alias[:1]) + alias[1:]
	return alias[strings.LastIndex(alias, ".")+1:]
}

func (u *UnionFieldType) alias() string {
	if u.Alias != "" {
		return u.Alias
	}
	if u.Model.Primitive != nil {
		return string(*u.Model.Primitive)
	}
	if u.Model.Fixed != nil {
		return FixedType
	}
	if u.Model.Array != nil {
		return ArrayType
	}
	if u.Model.Map != nil {
		return MapType
	}
	return u.Model.Namespace + "." + u.Model.Name
}
