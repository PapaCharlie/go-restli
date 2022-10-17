/*
Code generated by "github.com/PapaCharlie/go-restli"; DO NOT EDIT.

Source file: /Users/pchesnai/code/personal/go-restli/spec-parser/build/libs/go-restli-spec-parser-2.0.0-SNAPSHOT.jar
*/

package common

import (
	"github.com/PapaCharlie/go-restli/fnv1a"
	"github.com/PapaCharlie/go-restli/restli/patch"
	"github.com/PapaCharlie/go-restli/restlicodec"
)

// A "marker" data schema for data that is itself a data schema (a "PDSC for PDSCs"). Because PDSC is not expressive enough to describe it's own format, this is only a marker, and has no fields. Despite having no fields, it is required that data marked with this schema be non-empty. Specifically, is required that data marked as using this schema fully conform to the PDSC format (https://github.com/linkedin/rest.li/wiki/DATA-Data-Schema-and-Templates#schema-definition).
type PegasusSchema struct{}

func (p *PegasusSchema) Equals(other *PegasusSchema) bool {
	if p == other {
		return true
	}
	if p == nil || other == nil {
		return false
	}

	return true
}

func (p *PegasusSchema) ComputeHash() fnv1a.Hash {
	if p == nil {
		return fnv1a.ZeroHash()
	}
	hash := fnv1a.NewHash()

	return hash
}

func (p *PegasusSchema) MarshalFields(keyWriter func(string) restlicodec.Writer) (err error) {
	return nil
}

func (p *PegasusSchema) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(p.MarshalFields)
}

func (p *PegasusSchema) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = p.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

var PegasusSchemaRequiredFields = restlicodec.NewRequiredFields()

func (p *PegasusSchema) UnmarshalField(reader restlicodec.Reader, field string) (found bool, err error) {
	return false, nil
}

func (p *PegasusSchema) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	err = reader.ReadRecord(PegasusSchemaRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		found, err := p.UnmarshalField(reader, field)
		if err != nil {
			return err
		}
		if !found {
			err = reader.Skip()
		}
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *PegasusSchema) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, p)
}

func (p *PegasusSchema) NewInstance() *PegasusSchema {
	return new(PegasusSchema)
}

/*
================================================================================
PARTIAL UPDATE STRUCTS
================================================================================
*/

type PegasusSchema_PartialUpdate_Set_Fields struct{}

func (p *PegasusSchema_PartialUpdate_Set_Fields) MarshalFields(keyWriter func(string) restlicodec.Writer) (err error) {
	return nil
}

func (p *PegasusSchema_PartialUpdate_Set_Fields) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(p.MarshalFields)
}

func (p *PegasusSchema_PartialUpdate_Set_Fields) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = p.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

var PegasusSchema_PartialUpdate_Set_FieldsRequiredFields = restlicodec.NewRequiredFields()

func (p *PegasusSchema_PartialUpdate_Set_Fields) UnmarshalField(reader restlicodec.Reader, field string) (found bool, err error) {
	return false, nil
}

func (p *PegasusSchema_PartialUpdate_Set_Fields) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	err = reader.ReadRecord(PegasusSchema_PartialUpdate_Set_FieldsRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		found, err := p.UnmarshalField(reader, field)
		if err != nil {
			return err
		}
		if !found {
			err = reader.Skip()
		}
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *PegasusSchema_PartialUpdate_Set_Fields) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, p)
}

func (p *PegasusSchema_PartialUpdate_Set_Fields) NewInstance() *PegasusSchema_PartialUpdate_Set_Fields {
	return new(PegasusSchema_PartialUpdate_Set_Fields)
}

// PegasusSchema_PartialUpdate is used to represent a partial update on PegasusSchema. Toggling the value of a field
// in Delete represents selecting it for deletion in a partial update, while
// setting the value of a field in Update represents setting that field in the
// current struct. Other fields in this struct represent record fields that can
// themselves be partially updated.
type PegasusSchema_PartialUpdate struct{}

func (p *PegasusSchema_PartialUpdate) MarshalRestLiPatch(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
		return nil
	})
}

func (p *PegasusSchema_PartialUpdate) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
		return p.MarshalRestLiPatch(keyWriter(patch.PatchField).SetScope())
	})
}

func (p *PegasusSchema_PartialUpdate) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = p.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

func (p *PegasusSchema_PartialUpdate) UnmarshalRestLiPatch(reader restlicodec.Reader) (err error) {
	return reader.ReadMap(func(reader restlicodec.Reader, key string) (err error) {
		return reader.Skip()
	})
}

func (p *PegasusSchema_PartialUpdate) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	return reader.ReadRecord(patch.RequiredPatchRecordFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == patch.PatchField {
			return p.UnmarshalRestLiPatch(reader)
		} else {
			return reader.Skip()
		}
	})
}

func (p *PegasusSchema_PartialUpdate) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, p)
}

func (p *PegasusSchema_PartialUpdate) NewInstance() *PegasusSchema_PartialUpdate {
	return new(PegasusSchema_PartialUpdate)
}