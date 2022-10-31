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
	Primitive *PrimitiveType    `json:"primitive,omitempty"`
	Reference *utils.Identifier `json:"reference,omitempty"`
	Array     *RestliType       `json:"array,omitempty"`
	Map       *RestliType       `json:"map,omitempty"`
	RawRecord bool              `json:"rawRecord,omitempty"`
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
		t.Reference = &utils.RawRecordIdentifier
		return nil
	default:
		return errors.Errorf("go-restli: RestliType declares no underlying type! (%s)", string(data))
	}
}
