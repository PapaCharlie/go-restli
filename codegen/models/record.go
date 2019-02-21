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

func (r *Record) GenerateCode() (def *jen.Statement) {
	def = jen.Empty()

	AddWordWrappedComment(def, r.Doc).Line()

	def.Type().Id(r.Name).StructFunc(func(def *jen.Group) {
		for _, i := range r.Include {
			if rec := i.Record; rec != nil {
				def.Add(i.GoType())
				continue
			}
			log.Panic("Illegal included type:", i)
		}

		for _, f := range r.Fields {
			field := def.Empty()
			AddWordWrappedComment(field, f.Doc).Line()
			field.Id(ExportedIdentifier(f.Name))
			if f.Optional || f.Default != nil {
				field.Add(f.Type.PointerType()).Tag(JsonTag(f.Name, true))
			} else {
				field.Add(f.Type.GoType()).Tag(JsonTag(f.Name, false))
			}
		}
	}).Line().Line()

	receiver := PrivateIdentifier(r.Name[:1])

	def.Func().
		Id("New" + r.Name).Params().
		Params(jen.Id(receiver).Op("*").Id(r.Name))
	def.BlockFunc(func(def *jen.Group) {
		def.Id(receiver).Op("=").New(jen.Id(r.Name))
		for _, f := range r.Fields {
			if f.Type.Record != nil && !f.Optional && f.Default == nil {
				def.Id(receiver).Dot(ExportedIdentifier(f.Name)).Op("=").Op("*").Qual(f.Type.PackagePath(), "New"+f.Type.Record.Name).Call()
			}
		}
		def.Id(receiver).Dot(PopulateDefaultValues).Call()
		def.Return()
	}).Line().Line()

	def.Func().
		Params(jen.Id(receiver).Op("*").Id(r.Name)).
		Id(PopulateDefaultValues).Params().
		Params()
	def.BlockFunc(func(def *jen.Group) {
		for _, f := range r.Fields {
			name := ExportedIdentifier(f.Name)
			if f.Default != nil {
				SetDefaultValue(def, receiver, name, string(f.Default), f.Type)
				def.Line()
			}
		}
	}).Line().Line()

	AddMarshalJSON(def, receiver, r.Name, func(def *jen.Group) {
		def.Id(receiver).Dot(PopulateDefaultValues).Call()
		def.Type().Id("_t").Id(r.Name)
		def.Return(jen.Qual(EncodingJson, Marshal).Call(jen.Call(jen.Op("*").Id("_t")).Call(jen.Id(receiver))))
	}).Line().Line()

	AddUnmarshalJSON(def, receiver, r.Name, func(def *jen.Group) {
		def.Type().Id("_t").Id(r.Name)
		def.Err().Op("=").Qual(EncodingJson, Unmarshal).Call(jen.Id("data"), jen.Call(jen.Op("*").Id("_t")).Call(jen.Id(receiver)))
		IfErrReturn(def).Line()
		def.Id(receiver).Dot(PopulateDefaultValues).Call()
		def.Return()
	}).Line().Line()

	return
}
