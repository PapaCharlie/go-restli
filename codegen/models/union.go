package models

import (
	"encoding/json"
	"log"
	"strings"
	"unicode"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type UnionModel struct {
	Types      []UnionFieldModel
	IsOptional bool
}

type UnionFieldModel struct {
	Type  *Model
	Alias string
}

func (u *UnionFieldModel) UnmarshalJSON(data []byte) error {
	m := &Model{}
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	u.Type = m

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
	nullIndex := -1
	for i, m := range types {
		if primitive, ok := m.Type.BuiltinType.(*PrimitiveModel); ok && *primitive == NullPrimitive {
			nullIndex = i
			u.IsOptional = true
		}
	}
	if nullIndex >= 0 {
		types = append(types[:nullIndex], types[nullIndex+1:]...)
	}
	if len(types) == 0 {
		return errors.Errorf("Empty union: %s", string(data))
	}
	u.Types = types

	return nil
}

func (u *UnionModel) innerModels() (models []*Model) {
	for _, t := range u.Types {
		models = append(models, t.Type)
	}
	return models
}

func (u *UnionModel) GoType() (def *Statement) {
	return StructFunc(func(def *Group) {
		for _, t := range u.Types {
			var tag FieldTag
			tag.Json.Name = t.alias()
			tag.Json.Optional = true

			field := def.Empty()
			field.Id(t.name())
			field.Op("*").Add(t.Type.GoType())
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
	if u.Type.BuiltinType != nil {
		switch u.Type.BuiltinType.(type) {
		case *PrimitiveModel:
			return u.Type.BuiltinType.(*PrimitiveModel)[0]
		case *BytesModel:
			return BytesModelTypeName
		case *MapModel:
			return MapModelTypeName
		case *ArrayModel:
			return ArrayModelTypeName
		default:
			log.Panicln("Unknown builtin type", u.Type.BuiltinType)
		}
	}
	if _, isFixed := u.Type.ComplexType.(*FixedModel); isFixed {
		return FixedModelTypeName
	}
	id := u.Type.ComplexType.GetIdentifier()
	return (&id).GetQualifiedClasspath()
}

func (u *UnionModel) restLiWriteToBuf(def *Group, accessor *Statement) {
	label := "end" + ExportedIdentifier(accessor.GoString())
	for i, c := range label {
		if !(unicode.IsDigit(c) || unicode.IsLetter(c)) {
			label = label[:i] + "_" + label[i+1:]
		}
	}

	for _, t := range u.Types {
		def.If(Add(accessor).Dot(t.name()).Op("!=").Nil()).BlockFunc(func(def *Group) {
			writeStringToBuf(def, Lit("("+t.alias()+":"))
			fieldAccessor := Add(accessor).Dot(t.name())
			if t.Type.IsMapOrArray() || t.Type.IsBytesOrPrimitive() {
				fieldAccessor = Op("*").Add(fieldAccessor)
			}
			t.Type.restLiWriteToBuf(def, fieldAccessor)
			def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
			def.Goto().Id(label)
		}).Line()
	}

	def.Id(label).Op(":")
}
