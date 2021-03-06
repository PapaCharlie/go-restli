package types

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func init() {
	utils.TypeRegistry.Register(new(RawRecord))
}

type RawRecord struct{}

func (r *RawRecord) GetIdentifier() utils.Identifier {
	return utils.RawRecordContextIdentifier
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
