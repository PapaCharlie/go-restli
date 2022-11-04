package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PapaCharlie/go-restli/v2/codegen/resources"
	"github.com/PapaCharlie/go-restli/v2/codegen/types"
	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
)

func ReadManifest(data []byte) (*GoRestliManifest, error) {
	manifest := new(GoRestliManifest)
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&manifest)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

func RegisterManifests(manifests []*GoRestliManifest) (err error) {
	for _, m := range manifests {
		for _, dt := range m.InputDataTypes {
			err = utils.TypeRegistry.Register(dt.GetComplexType(), m.PackageRoot)
			if err != nil {
				return err
			}
			if dt.Typeref != nil && dt.Typeref.IsCustom {
				utils.TypeRegistry.SetCustomTyperef(dt.Typeref.Identifier)
			}
		}
	}

	// Only use the final (input) manifest to resolve dependency data types, that way types that are not defined as
	// input types by this manifest or upstream manifests always come from the same place. If the dependency types of
	// previous packages is used, the import location may spuriously change depending on the order in which upstream
	// manifests are discovered
	inputManifest := manifests[len(manifests)-1]
	for _, dt := range inputManifest.DependencyDataTypes {
		_ = utils.TypeRegistry.Register(dt.GetComplexType(), inputManifest.PackageRoot)
	}
	return utils.TypeRegistry.Finalize()
}

func LocateCustomTyperefs(manifest *GoRestliManifest, outputDir string) error {
	for _, dt := range manifest.InputDataTypes {
		if dt.Typeref == nil {
			continue
		}
		typeref := dt.Typeref
		expectedLocation := filepath.Join(
			outputDir,
			strings.TrimPrefix(typeref.PackagePath(), manifest.PackageRoot),
			typeref.TypeName()+".go",
		)
		_, err := os.Stat(expectedLocation)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		// TODO: Improve validation by reading the file and checking that it correctly declares the right functions etc
		utils.TypeRegistry.SetCustomTyperef(typeref.Identifier)
		typeref.IsCustom = true
	}
	return nil
}

type DataType struct {
	Enum            *types.Enum            `json:"enum,omitempty"`
	Fixed           *types.Fixed           `json:"fixed,omitempty"`
	Record          *types.Record          `json:"record,omitempty"`
	ComplexKey      *types.ComplexKey      `json:"complexKey,omitempty"`
	StandaloneUnion *types.StandaloneUnion `json:"standaloneUnion,omitempty"`
	Typeref         *types.Typeref         `json:"typeref,omitempty"`
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

func (dt *DataType) UnmarshalJSON(data []byte) error {
	type t DataType
	err := json.Unmarshal(data, (*t)(dt))
	if err != nil {
		return err
	}

	if dt.GetComplexType() == nil {
		return fmt.Errorf("go-restli: Must declare at least one underlying type")
	}

	return nil
}

type GoRestliManifest struct {
	PackageRoot         string                `json:"packageRoot"`
	InputDataTypes      []DataType            `json:"inputDataTypes"`
	DependencyDataTypes []DataType            `json:"dependencyDataTypes"`
	Resources           []*resources.Resource `json:"resources"`
}

func (m *GoRestliManifest) UnmarshalJSON(data []byte) error {
	type t GoRestliManifest
	err := json.Unmarshal(data, (*t)(m))
	if err != nil {
		return err
	}

	for _, r := range m.Resources {
		r.PackageRoot = m.PackageRoot
	}

	return nil
}

func (m *GoRestliManifest) GenerateResourceCode() (codeFiles []*utils.CodeFile) {
	for _, r := range m.Resources {
		codeFiles = append(codeFiles, r.GenerateCode()...)
	}
	return codeFiles
}
