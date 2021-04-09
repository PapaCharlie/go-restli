package types

import (
	"encoding/json"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/pkg/errors"
)

type NamedType struct {
	utils.Identifier
	SourceFile string `json:"sourceFile"`
	Doc        string `json:"doc"`
}

type RestliType struct {
	Primitive     *PrimitiveType    `json:"primitive"`
	Reference     *utils.Identifier `json:"reference"`
	Array         *RestliType       `json:"array"`
	Map           *RestliType       `json:"map"`
	RawRecord     bool              `json:"rawRecord"`
	NativeTyperef *NativeTyperef    `json:"nativeTyperef"`
}

func (t *RestliType) UnmarshalJSON(data []byte) error {
	type _t RestliType
	err := json.Unmarshal(data, (*_t)(t))
	if err != nil {
		return err
	}

	switch {
	case t.Primitive != nil:
		return nil
	case t.Reference != nil:
		return nil
	case t.Array != nil:
		return nil
	case t.Map != nil:
		return nil
	case t.RawRecord:
		t.Reference = &RawRecordIdentifier
		return nil
	case t.NativeTyperef != nil:
		// Because the NativeTyperef type is reused in the CLI parameters, it's best to simply validate here that the
		// original type name and underlying primitive is defined here, rather than in a custom UnmarshalJSON on
		// NativeTyperef itself
		if t.NativeTyperef.OriginalTypeName == nil {
			return errors.Errorf("go-restli: NativeTyperef does not declare the original type name! (%s)", string(data))
		}
		if t.NativeTyperef.Primitive == nil {
			return errors.Errorf("go-restli: NativeTyperef does not declare underlying primitive type! (%s)", string(data))
		}
		return nil
	default:
		return errors.Errorf("go-restli: RestliType declares no underlying type! (%s)", string(data))
	}
}
