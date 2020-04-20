package codegen

import (
	. "github.com/dave/jennifer/jen"
)

func (t *RestliType) RestLiURLDecodeModel(input *Statement, output *Statement) (def *Statement, hasError bool) {
	return t.RestLiDecodeModel(RestLiUrlEncoder, input, output)
}

func (t *RestliType) RestLiReducedDecodeModel(input *Statement, output *Statement) (def *Statement, hasError bool) {
	return t.RestLiDecodeModel(RestLiReducedEncoder, input, output)
}

func (t *RestliType) RestLiDecodeModel(encoder string, input *Statement, output *Statement) (*Statement, bool) {
	decoderRef := Qual(ProtocolPackage, encoder)

	if t.Reference != nil {
		return Add(output).Dot(RestLiDecode).Call(decoderRef, input), true
	}

	if t.Primitive != nil {
		return Add(decoderRef).Dot("Decode"+ExportedIdentifier(t.Primitive.Type)).Call(input, output), true
	}

	Logger.Panicf("%+v cannot be url decoded", t)
	return nil, false
}
