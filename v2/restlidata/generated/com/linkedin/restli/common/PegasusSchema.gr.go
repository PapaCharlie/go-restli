/*
Code generated by "github.com/PapaCharlie/go-restli/v2"; DO NOT EDIT.

Source file: /Users/pchesnai/code/personal/go-restli/v2/go-restli-spec-parser.jar
*/

package common

import (
	"github.com/PapaCharlie/go-restli/v2/fnv1a"
	"github.com/PapaCharlie/go-restli/v2/restli/patch"
	"github.com/PapaCharlie/go-restli/v2/restlicodec"
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

type PegasusSchema_PartialUpdate_Delete_Fields struct{}

func (p *PegasusSchema_PartialUpdate_Delete_Fields) MarshalDeleteFields(write func(string)) {}

func (p *PegasusSchema_PartialUpdate_Delete_Fields) UnmarshalDeleteField(field string) (err error) {
	switch field {
	default:
		return patch.NoSuchFieldErr
	}
}

func (p *PegasusSchema_PartialUpdate) MarshalDeleteFields(itemWriter func() restlicodec.Writer) (err error) {
	write := func(name string) {
		itemWriter().WriteString(name)
	}
	p.Delete_Fields.MarshalDeleteFields(write)
	return nil
}

func (p *PegasusSchema_PartialUpdate) UnmarshalDeleteField(field string) (err error) {
	return p.Delete_Fields.UnmarshalDeleteField(field)
}

type PegasusSchema_PartialUpdate_Set_Fields struct{}

func (p *PegasusSchema_PartialUpdate_Set_Fields) MarshalFields(keyWriter func(string) restlicodec.Writer) (err error) {
	return nil
}

func (p *PegasusSchema_PartialUpdate_Set_Fields) UnmarshalField(reader restlicodec.Reader, field string) (found bool, err error) {
	return false, nil
}

func (p *PegasusSchema_PartialUpdate) MarshalSetFields(keyWriter func(string) restlicodec.Writer) (err error) {
	err = p.Set_Fields.MarshalFields(keyWriter)
	return err
}

func (p *PegasusSchema_PartialUpdate) UnmarshalSetField(reader restlicodec.Reader, field string) (found bool, err error) {
	return p.Set_Fields.UnmarshalField(reader, field)
}

// PegasusSchema_PartialUpdate is used to represent a partial update on PegasusSchema. Toggling the value of a field
// in Delete_Field represents selecting it for deletion in a partial update, while
// setting the value of a field in Set_Fields represents setting that field in the
// current struct. Other fields in this struct represent record fields that can
// themselves be partially updated.
type PegasusSchema_PartialUpdate struct {
	Delete_Fields PegasusSchema_PartialUpdate_Delete_Fields
	Set_Fields    PegasusSchema_PartialUpdate_Set_Fields
}

func (p *PegasusSchema_PartialUpdate) CheckFields(fieldChecker *patch.PartialUpdateFieldChecker, keyChecker restlicodec.KeyChecker) (err error) {
	return nil
}

func (p *PegasusSchema_PartialUpdate) MarshalRestLiPatch(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
		fieldChecker := &patch.PartialUpdateFieldChecker{
			RecordType: "com.linkedin.restli.common.PegasusSchema",
		}
		err = p.CheckFields(fieldChecker, writer)
		if err != nil {
			return err
		}
		if fieldChecker.HasDeletes {
			err = keyWriter("$delete").WriteArray(func(itemWriter func() restlicodec.Writer) (err error) {
				p.MarshalDeleteFields(itemWriter)
				return nil
			})
			if err != nil {
				return err
			}
		}

		if fieldChecker.HasSets {
			err = keyWriter("$set").WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
				return p.MarshalSetFields(keyWriter)
			})
			if err != nil {
				return err
			}
		}

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
	err = reader.ReadMap(func(reader restlicodec.Reader, key string) (err error) {
		switch key {
		case "$delete":
			err = reader.ReadArray(func(reader restlicodec.Reader) (err error) {
				var field string
				field, err = reader.ReadString()
				if err != nil {
					return err
				}

				err = p.UnmarshalDeleteField(field)
				if err == patch.NoSuchFieldErr {
					err = nil
				}
				return err
			})
		case "$set":
			err = reader.ReadMap(func(reader restlicodec.Reader, key string) (err error) {
				found, err := p.UnmarshalSetField(reader, key)
				if !found {
					err = reader.Skip()
				}
				return err
			})
		default:
			err = reader.Skip()
		}
		return err
	})
	if err != nil {
		return err
	}
	fieldChecker := &patch.PartialUpdateFieldChecker{
		RecordType: "com.linkedin.restli.common.PegasusSchema",
	}
	err = p.CheckFields(fieldChecker, reader)
	if err != nil {
		return err
	}
	return nil
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
