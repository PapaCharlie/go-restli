package codegen

import (
	. "github.com/dave/jennifer/jen"
)

func (t *RestliType) RestLiURLDecodeModel(input *Statement, output *Statement) *Statement {
	return t.RestLiDecodeModel(RestLiUrlEncoder, input, output)
}

func (t *RestliType) RestLiReducedDecodeModel(input *Statement, output *Statement) *Statement {
	return t.RestLiDecodeModel(RestLiReducedEncoder, input, output)
}

func (t *RestliType) RestLiDecodeModel(encoder string, input *Statement, output *Statement) *Statement {
	decoderRef := Qual(ProtocolPackage, encoder)

	if t.Reference != nil {
		return Add(output).Dot(RestLiDecode).Call(decoderRef, input)
	}

	if t.Primitive != nil {
		return Add(decoderRef).Dot("Decode"+ExportedIdentifier(t.Primitive.Type)).Call(input, output)
	}

	Logger.Panicf("%+v cannot be url decoded", t)
	return nil
}
