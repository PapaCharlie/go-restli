package codegen

import (
	. "github.com/dave/jennifer/jen"
)

func (t *RestliType) RestLiURLEncodeModel(accessor *Statement) (def *Statement, hasError bool) {
	return t.RestLiEncodeModel(RestLiUrlEncoder, accessor)
}

func (t *RestliType) RestLiReducedEncodeModel(accessor *Statement) (def *Statement, hasError bool) {
	return t.RestLiEncodeModel(RestLiReducedEncoder, accessor)
}

func (t *RestliType) RestLiEncodeModel(encoder string, accessor *Statement) (*Statement, bool) {
	encoderRef := Qual(ProtocolPackage, encoder)

	if t.Primitive != nil {
		return encoderRef.Dot("Encode" + ExportedIdentifier(t.Primitive.Type)).Call(accessor), false
	}

	if t.Reference != nil {
		return Add(accessor).Dot(MarshalRestLi).Call(encoderRef), true
	}

	Logger.Panicf("%+v cannot be url encoded", t)
	return nil, false
}
