package models

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"strings"
)

type Union struct {
	Types []UnionFieldType
}

type UnionFieldType struct {
	*Model
	Alias string `json:"alias"`
}

func (u *UnionFieldType) UnmarshalJSON(data []byte) error {
	u.Model = &Model{}
	type t UnionFieldType
	if err := json.Unmarshal(data, (*t)(u)); err != nil {
		return err
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

func (u *Union) generateCode(destinationPackage string) (def *jen.Statement) {
	return
}

func (u *Union) GoType(destinationPackage string) *jen.Statement {
	if len(u.Types) == 0 {
		panic(u)
	}
	var fields []jen.Code
	for _, t := range u.Types {
		def := jen.Empty()
		addWordWrappedComment(def, t.Doc)
		def.Id(t.name())
		def.Add(t.GoType(destinationPackage))
		def.Tag(jsonTag(t.alias()))
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
	if u.Primitive != nil {
		return string(*u.Primitive)
	}
	if u.Fixed != nil {
		return FixedType
	}
	if u.Array != nil {
		return ArrayType
	}
	if u.Map != nil {
		return MapType
	}
	return u.Namespace + "." + u.Name
}
