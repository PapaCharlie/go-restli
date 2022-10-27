package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func init() {
	utils.TypeRegistry.Register(&customType{
		identifier: utils.RawRecordIdentifier,
		sourceFile: "https://github.com/PapaCharlie/go-restli/blob/master/restlidata/RawRecord.go",
	}, utils.RootPackage)
	utils.TypeRegistry.Register(&customType{
		identifier: utils.EmptyRecordIdentifier,
		sourceFile: "https://github.com/PapaCharlie/go-restli/blob/master/restlidata/EmptyRecord.go",
	}, utils.RestLiDataPackage+"/generated")
}

type customType struct {
	identifier utils.Identifier
	sourceFile string
}

func (c *customType) GetIdentifier() utils.Identifier {
	return c.identifier
}

func (c *customType) GetSourceFile() string {
	return c.sourceFile
}

func (c *customType) InnerTypes() utils.IdentifierSet {
	return nil
}

func (c *customType) ShouldReference() utils.ShouldUsePointer {
	return utils.No
}

func (c *customType) GenerateCode() *Statement {
	return Empty()
}
