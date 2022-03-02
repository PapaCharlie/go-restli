package cmd

import (
	"encoding/json"

	"github.com/PapaCharlie/go-restli/codegen/resources"
	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/pkg/errors"
)

type GoRestliSpec struct {
	DataTypes []struct {
		Enum            *types.Enum            `json:"enum"`
		Fixed           *types.Fixed           `json:"fixed"`
		Record          *types.Record          `json:"record"`
		ComplexKey      *types.ComplexKey      `json:"complexKey"`
		StandaloneUnion *types.StandaloneUnion `json:"standaloneUnion"`
		Typeref         *types.Typeref         `json:"typeref"`
	} `json:"dataTypes"`
	Resources []resources.Resource
}

func (s *GoRestliSpec) UnmarshalJSON(data []byte) error {
	type t GoRestliSpec
	err := json.Unmarshal(data, (*t)(s))
	if err != nil {
		return err
	}

	for _, t := range s.DataTypes {
		var complexType utils.ComplexType
		switch {
		case t.Enum != nil:
			complexType = t.Enum
		case t.Fixed != nil:
			complexType = t.Fixed
		case t.Record != nil:
			complexType = t.Record
		case t.ComplexKey != nil:
			complexType = t.ComplexKey
		case t.StandaloneUnion != nil:
			complexType = t.StandaloneUnion
		case t.Typeref != nil:
			complexType = t.Typeref
			t.Typeref.CheckNativeTyperef()
		default:
			return errors.New("go-restli: Must declare at least one underlying type")
		}
		utils.TypeRegistry.Register(complexType)
	}

	utils.TypeRegistry.FlagCyclicDependencies()
	return nil
}

func (s *GoRestliSpec) GenerateClientCode() (codeFiles []*utils.CodeFile) {
	for _, r := range s.Resources {
		codeFiles = append(codeFiles, r.GenerateCode()...)
	}
	return codeFiles
}
