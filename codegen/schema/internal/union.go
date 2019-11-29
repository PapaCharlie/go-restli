package internal

import (
	"encoding/json"
	"fmt"
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

func (u *UnionModel) UnmarshalJSON(data []byte) error {
	var types []UnionFieldModel
	if err := json.Unmarshal(data, &types); err != nil {
		return errors.WithStack(err)
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

func (u *UnionModel) restLiWriteToBuf(def *Group, accessor *Statement) {
	label := "end" + canonicalizeAccessor(accessor)

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

func (u *UnionModel) validateUnionFields(def *Group, accessor *Statement) {
	isSet := "is" + canonicalizeAccessor(accessor) + "Set"
	def.Id(isSet).Op(":=").False().Line()
	errorMessage := fmt.Sprintf("must specify exactly one member of %s", accessor.GoString())

	for i, t := range u.Types {
		def.If(Add(accessor).Dot(t.name()).Op("!=").Nil()).
			BlockFunc(func(def *Group) {
				if i == 0 {
					def.Id(isSet).Op("=").True()
				} else {
					def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
						def.Id(isSet).Op("=").True()
					}).Else().BlockFunc(func(def *Group) {
						def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(errorMessage))
						def.Return()
					})
				}
			}).Line()
	}
	def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
		def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(errorMessage))
		def.Return()
	})
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

func (u *UnionFieldModel) UnmarshalJSON(data []byte) error {
	m := &struct {
		Type  *Model
		Alias string
	}{}
	if err := json.Unmarshal(data, m); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal string into Go value of type struct") {
			sm := &Model{}
			err = json.Unmarshal(data, sm)
			if err != nil {
				return errors.Wrapf(err, "cannot deserialize %s as a union field", string(data))
			}
			m.Type = sm
		} else {
			return errors.Wrapf(err, "Could not deserialize %s as a union field", string(data))
		}
	}
	u.Alias = m.Alias
	u.Type = m.Type

	type t UnionFieldModel
	if err := json.Unmarshal(data, (*t)(u)); err != nil {
		if !strings.Contains(err.Error(), "json: cannot unmarshal string into Go value of type") {
			return errors.WithStack(err)
		}
	}

	return nil
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
	return u.Type.ComplexType.GetIdentifier().GetQualifiedClasspath()
}
