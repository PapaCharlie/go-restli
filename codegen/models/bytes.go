package models

import (
	"encoding/json"

	"github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

const BytesModelTypeName = "bytes"

type BytesModel struct {
}

func (b *BytesModel) UnmarshalJSON(data []byte) error {
	var bytes string
	if err := json.Unmarshal(data, &bytes); err != nil {
		return errors.WithStack(err)
	}

	if bytes != BytesModelTypeName {
		return errors.New("Not a \"bytes\" type")
	}
	return nil
}

func (b *BytesModel) GoType() *Statement {
	return Qual(codegen.GetRestLiProtocolPackage(), "Bytes")
}

func (b *BytesModel) encode(accessor *Statement) *Statement {
	return Id(codegen.Codec).Dot("EncodeBytes").Call(accessor)
}

func (b *BytesModel) decode(accessor *Statement) *Statement {
	return Id(codegen.Codec).Dot("DecodeBytes").Call(Id("data"), Call(Op("*").Add(codegen.Bytes())).Call(accessor))
}
