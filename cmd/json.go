package cmd

import (
	"bytes"
	"encoding/json"

	"github.com/PapaCharlie/go-restli/codegen/resources"
	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/pkg/errors"
)

func ReadManifest(data []byte) (*GoRestliManifest, error) {
	manifest := new(GoRestliManifest)
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&manifest)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

type DataType struct {
	Enum            *types.Enum            `json:"enum"`
	Fixed           *types.Fixed           `json:"fixed"`
	Record          *types.Record          `json:"record"`
	ComplexKey      *types.ComplexKey      `json:"complexKey"`
	StandaloneUnion *types.StandaloneUnion `json:"standaloneUnion"`
	Typeref         *types.Typeref         `json:"typeref"`
}

func (dt *DataType) GetComplexType() utils.ComplexType {
	switch {
	case dt.Enum != nil:
		return dt.Enum
	case dt.Fixed != nil:
		return dt.Fixed
	case dt.Record != nil:
		return dt.Record
	case dt.ComplexKey != nil:
		return dt.ComplexKey
	case dt.StandaloneUnion != nil:
		return dt.StandaloneUnion
	case dt.Typeref != nil:
		return dt.Typeref
	default:
		return nil
	}
}

type GoRestliManifest struct {
	PackageRoot string     `json:"packageRoot"`
	DataTypes   []DataType `json:"dataTypes"`
	Resources   []*resources.Resource
}

func (m *GoRestliManifest) UnmarshalJSON(data []byte) error {
	type t GoRestliManifest
	err := json.Unmarshal(data, (*t)(m))
	if err != nil {
		return err
	}

	var complexTypes []utils.ComplexType
	for _, dt := range m.DataTypes {
		if ct := dt.GetComplexType(); ct == nil {
			return errors.New("go-restli: Must declare at least one underlying type")
		} else {
			complexTypes = append(complexTypes, ct)
		}
	}
	utils.TypeRegistry.RegisterTypes(complexTypes, m.PackageRoot)

	for _, r := range m.Resources {
		r.PackageRoot = m.PackageRoot
	}

	return nil
}

func (m *GoRestliManifest) GenerateClientCode() (codeFiles []*utils.CodeFile) {
	for _, r := range m.Resources {
		codeFiles = append(codeFiles, r.GenerateCode()...)
	}
	return codeFiles
}
