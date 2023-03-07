/*
Code generated by "github.com/PapaCharlie/go-restli/v2"; DO NOT EDIT.

Source file: /Users/pchesnai/code/personal/go-restli/v2/go-restli-spec-parser.jar
*/

package common

import (
	"github.com/PapaCharlie/go-restli/v2/fnv1a"
	"github.com/PapaCharlie/go-restli/v2/restli/equals"
	"github.com/PapaCharlie/go-restli/v2/restli/patch"
	"github.com/PapaCharlie/go-restli/v2/restlicodec"
)

// Specifies the extension schema annotation configuration for defining the entity relationship among entities.
type ExtensionSchemaAnnotation struct {
	// Specifies only for one to many relationship. It can use either GET_ALL("get_all") or a FINDER("finder:<method name>").
	Using *string
	// Specifies parameters if any of them are specified. It is also used for FINDER.
	Params *map[string]string
	// Specifies versionSuffix in multi-version scenario. If is is not provided, will pick first version by default.
	VersionSuffix *string
}

func (e *ExtensionSchemaAnnotation) Equals(other *ExtensionSchemaAnnotation) bool {
	if e == other {
		return true
	}
	if e == nil || other == nil {
		return false
	}

	return equals.ComparablePointer(e.Using, other.Using) &&
		equals.ComparableMapPointer(e.Params, other.Params) &&
		equals.ComparablePointer(e.VersionSuffix, other.VersionSuffix)
}

func (e *ExtensionSchemaAnnotation) ComputeHash() fnv1a.Hash {
	if e == nil {
		return fnv1a.ZeroHash()
	}
	hash := fnv1a.NewHash()

	if e.Using != nil {
		hash.AddString(*e.Using)
	}

	if e.Params != nil {
		fnv1a.AddMap(hash, *e.Params, fnv1a.Hash.AddString)
	}

	if e.VersionSuffix != nil {
		hash.AddString(*e.VersionSuffix)
	}

	return hash
}

func (e *ExtensionSchemaAnnotation) MarshalFields(keyWriter func(string) restlicodec.Writer) (err error) {
	if e.Params != nil {
		err = restlicodec.WriteMap(keyWriter("params"), *e.Params, restlicodec.WriteString)
		if err != nil {
			return err
		}
	}
	if e.Using != nil {
		keyWriter("using").WriteString(*e.Using)
	}
	if e.VersionSuffix != nil {
		keyWriter("versionSuffix").WriteString(*e.VersionSuffix)
	}
	return nil
}

func (e *ExtensionSchemaAnnotation) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(e.MarshalFields)
}

func (e *ExtensionSchemaAnnotation) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = e.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

var ExtensionSchemaAnnotationRequiredFields = restlicodec.NewRequiredFields()

func (e *ExtensionSchemaAnnotation) UnmarshalField(reader restlicodec.Reader, field string) (found bool, err error) {
	switch field {
	case "using":
		found = true
		e.Using = new(string)
		*e.Using, err = reader.ReadString()
	case "params":
		found = true
		e.Params = new(map[string]string)
		*e.Params, err = restlicodec.ReadMap(reader, restlicodec.Reader.ReadString)
	case "versionSuffix":
		found = true
		e.VersionSuffix = new(string)
		*e.VersionSuffix, err = reader.ReadString()
	}
	return found, err
}

func (e *ExtensionSchemaAnnotation) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	err = reader.ReadRecord(ExtensionSchemaAnnotationRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		found, err := e.UnmarshalField(reader, field)
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

func (e *ExtensionSchemaAnnotation) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, e)
}

func (e *ExtensionSchemaAnnotation) NewInstance() *ExtensionSchemaAnnotation {
	return new(ExtensionSchemaAnnotation)
}

/*
================================================================================
PARTIAL UPDATE STRUCTS
================================================================================
*/

type ExtensionSchemaAnnotation_PartialUpdate_Delete_Fields struct {
	Using         bool
	Params        bool
	VersionSuffix bool
}

func (e *ExtensionSchemaAnnotation_PartialUpdate_Delete_Fields) MarshalDeleteFields(write func(string)) {
	if e.Params {
		write("params")
	}
	if e.Using {
		write("using")
	}
	if e.VersionSuffix {
		write("versionSuffix")
	}
}

func (e *ExtensionSchemaAnnotation_PartialUpdate_Delete_Fields) UnmarshalDeleteField(field string) (err error) {
	switch field {
	case "using":
		e.Using = true
		return nil
	case "params":
		e.Params = true
		return nil
	case "versionSuffix":
		e.VersionSuffix = true
		return nil
	default:
		return patch.NoSuchFieldErr
	}
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) MarshalDeleteFields(itemWriter func() restlicodec.Writer) (err error) {
	write := func(name string) {
		itemWriter().WriteString(name)
	}
	e.Delete_Fields.MarshalDeleteFields(write)
	return nil
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) UnmarshalDeleteField(field string) (err error) {
	return e.Delete_Fields.UnmarshalDeleteField(field)
}

type ExtensionSchemaAnnotation_PartialUpdate_Set_Fields struct {
	// using
	Using *string
	// params
	Params *map[string]string
	// versionSuffix
	VersionSuffix *string
}

func (e *ExtensionSchemaAnnotation_PartialUpdate_Set_Fields) MarshalFields(keyWriter func(string) restlicodec.Writer) (err error) {
	if e.Params != nil {
		err = restlicodec.WriteMap(keyWriter("params"), *e.Params, restlicodec.WriteString)
		if err != nil {
			return err
		}
	}
	if e.Using != nil {
		keyWriter("using").WriteString(*e.Using)
	}
	if e.VersionSuffix != nil {
		keyWriter("versionSuffix").WriteString(*e.VersionSuffix)
	}
	return nil
}

func (e *ExtensionSchemaAnnotation_PartialUpdate_Set_Fields) UnmarshalField(reader restlicodec.Reader, field string) (found bool, err error) {
	switch field {
	case "using":
		found = true
		e.Using = new(string)
		*e.Using, err = reader.ReadString()
	case "params":
		found = true
		e.Params = new(map[string]string)
		*e.Params, err = restlicodec.ReadMap(reader, restlicodec.Reader.ReadString)
	case "versionSuffix":
		found = true
		e.VersionSuffix = new(string)
		*e.VersionSuffix, err = reader.ReadString()
	}
	return found, err
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) MarshalSetFields(keyWriter func(string) restlicodec.Writer) (err error) {
	err = e.Set_Fields.MarshalFields(keyWriter)
	return err
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) UnmarshalSetField(reader restlicodec.Reader, field string) (found bool, err error) {
	return e.Set_Fields.UnmarshalField(reader, field)
}

// ExtensionSchemaAnnotation_PartialUpdate is used to represent a partial update on ExtensionSchemaAnnotation. Toggling the value of a field
// in Delete_Field represents selecting it for deletion in a partial update, while
// setting the value of a field in Set_Fields represents setting that field in the
// current struct. Other fields in this struct represent record fields that can
// themselves be partially updated.
type ExtensionSchemaAnnotation_PartialUpdate struct {
	Delete_Fields ExtensionSchemaAnnotation_PartialUpdate_Delete_Fields
	Set_Fields    ExtensionSchemaAnnotation_PartialUpdate_Set_Fields
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) CheckFields(fieldChecker *patch.PartialUpdateFieldChecker, keyChecker restlicodec.KeyChecker) (err error) {
	if err = fieldChecker.CheckField(keyChecker, "using", e.Delete_Fields.Using, e.Set_Fields.Using != nil, false); err != nil {
		return err
	}
	if err = fieldChecker.CheckField(keyChecker, "params", e.Delete_Fields.Params, e.Set_Fields.Params != nil, false); err != nil {
		return err
	}
	if err = fieldChecker.CheckField(keyChecker, "versionSuffix", e.Delete_Fields.VersionSuffix, e.Set_Fields.VersionSuffix != nil, false); err != nil {
		return err
	}
	return nil
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) MarshalRestLiPatch(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
		fieldChecker := &patch.PartialUpdateFieldChecker{
			RecordType: "com.linkedin.restli.common.ExtensionSchemaAnnotation",
		}
		err = e.CheckFields(fieldChecker, writer)
		if err != nil {
			return err
		}
		if fieldChecker.HasDeletes {
			err = keyWriter("$delete").WriteArray(func(itemWriter func() restlicodec.Writer) (err error) {
				e.MarshalDeleteFields(itemWriter)
				return nil
			})
			if err != nil {
				return err
			}
		}

		if fieldChecker.HasSets {
			err = keyWriter("$set").WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
				return e.MarshalSetFields(keyWriter)
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
		return e.MarshalRestLiPatch(keyWriter(patch.PatchField).SetScope())
	})
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = e.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) UnmarshalRestLiPatch(reader restlicodec.Reader) (err error) {
	err = reader.ReadMap(func(reader restlicodec.Reader, key string) (err error) {
		switch key {
		case "$delete":
			err = reader.ReadArray(func(reader restlicodec.Reader) (err error) {
				var field string
				field, err = reader.ReadString()
				if err != nil {
					return err
				}

				err = e.UnmarshalDeleteField(field)
				if err == patch.NoSuchFieldErr {
					err = nil
				}
				return err
			})
		case "$set":
			err = reader.ReadMap(func(reader restlicodec.Reader, key string) (err error) {
				found, err := e.UnmarshalSetField(reader, key)
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
		RecordType: "com.linkedin.restli.common.ExtensionSchemaAnnotation",
	}
	err = e.CheckFields(fieldChecker, reader)
	if err != nil {
		return err
	}
	return nil
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	return reader.ReadRecord(patch.RequiredPatchRecordFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == patch.PatchField {
			return e.UnmarshalRestLiPatch(reader)
		} else {
			return reader.Skip()
		}
	})
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, e)
}

func (e *ExtensionSchemaAnnotation_PartialUpdate) NewInstance() *ExtensionSchemaAnnotation_PartialUpdate {
	return new(ExtensionSchemaAnnotation_PartialUpdate)
}