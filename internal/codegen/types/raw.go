package types

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func init() {
	utils.TypeRegistry.Register(new(RawRecord))
}

var RawRecordIdentifier = utils.Identifier{
	IsNativeIdentifier: true,
	Name:               "RawRecord",
	Namespace:          utils.ProtocolPackage,
}

type RawRecord struct{}

func (r *RawRecord) GetIdentifier() utils.Identifier {
	return RawRecordIdentifier
}

func (r *RawRecord) GetSourceFile() string {
	return "https://github.com/PapaCharlie/go-restli/blob/master/protocol/RawRecord.go"
}

func (r *RawRecord) InnerTypes() utils.IdentifierSet {
	return nil
}

func (r *RawRecord) GenerateCode() (def *Statement) {
	return Empty()
}
