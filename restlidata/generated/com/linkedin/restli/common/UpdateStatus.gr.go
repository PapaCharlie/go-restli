/*
Code generated by "github.com/PapaCharlie/go-restli"; DO NOT EDIT.

Source file: /Users/pchesnai/code/personal/go-restli/spec-parser/build/libs/go-restli-spec-parser-2.0.2-SNAPSHOT.jar
*/

package common

import (
	"github.com/PapaCharlie/go-restli/fnv1a"
	"github.com/PapaCharlie/go-restli/restli/patch"
	"github.com/PapaCharlie/go-restli/restlicodec"
)

// A rest.li update status.
type UpdateStatus struct {
	Status int32

	Error *ErrorResponse
}

func (u *UpdateStatus) Equals(other *UpdateStatus) bool {
	if u == other {
		return true
	}
	if u == nil || other == nil {
		return false
	}

	return u.Status == other.Status &&
		u.Error.Equals(other.Error)
}

func (u *UpdateStatus) ComputeHash() fnv1a.Hash {
	if u == nil {
		return fnv1a.ZeroHash()
	}
	hash := fnv1a.NewHash()

	hash.AddInt32(u.Status)

	if u.Error != nil {
		hash.Add(u.Error.ComputeHash())
	}

	return hash
}

func (u *UpdateStatus) MarshalFields(keyWriter func(string) restlicodec.Writer) (err error) {
	if u.Error != nil {
		err = u.Error.MarshalRestLi(keyWriter("error"))
		if err != nil {
			return err
		}
	}
	keyWriter("status").WriteInt32(u.Status)
	return nil
}

func (u *UpdateStatus) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(u.MarshalFields)
}

func (u *UpdateStatus) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = u.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

var UpdateStatusRequiredFields = restlicodec.NewRequiredFields().Add(
	"status",
)

func (u *UpdateStatus) UnmarshalField(reader restlicodec.Reader, field string) (found bool, err error) {
	switch field {
	case "status":
		found = true
		u.Status, err = reader.ReadInt32()
	case "error":
		found = true
		u.Error = new(ErrorResponse)
		err = u.Error.UnmarshalRestLi(reader)
	}
	return found, err
}

func (u *UpdateStatus) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	err = reader.ReadRecord(UpdateStatusRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		found, err := u.UnmarshalField(reader, field)
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

func (u *UpdateStatus) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, u)
}

func (u *UpdateStatus) NewInstance() *UpdateStatus {
	return new(UpdateStatus)
}

/*
================================================================================
PARTIAL UPDATE STRUCTS
================================================================================
*/

type UpdateStatus_PartialUpdate_Delete_Fields struct {
	Error bool
}

func (u *UpdateStatus_PartialUpdate_Delete_Fields) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteArray(func(itemWriter func() restlicodec.Writer) (err error) {
		if u.Error {
			itemWriter().WriteString("error")
		}
		return nil
	})
}

func (u *UpdateStatus_PartialUpdate_Delete_Fields) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = u.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

func (u *UpdateStatus_PartialUpdate_Delete_Fields) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	var field string
	return reader.ReadArray(func(reader restlicodec.Reader) (err error) {
		field, err = reader.ReadString()
		if err != nil {
			return err
		}

		switch field {
		case "error":
			u.Error = true
		}
		return nil
	})
}

func (u *UpdateStatus_PartialUpdate_Delete_Fields) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, u)
}

func (u *UpdateStatus_PartialUpdate_Delete_Fields) NewInstance() *UpdateStatus_PartialUpdate_Delete_Fields {
	return new(UpdateStatus_PartialUpdate_Delete_Fields)
}

type UpdateStatus_PartialUpdate_Set_Fields struct {
	// status
	Status *int32
	// error
	Error *ErrorResponse
}

func (u *UpdateStatus_PartialUpdate_Set_Fields) MarshalFields(keyWriter func(string) restlicodec.Writer) (err error) {
	if u.Error != nil {
		err = u.Error.MarshalRestLi(keyWriter("error"))
		if err != nil {
			return err
		}
	}
	if u.Status != nil {
		keyWriter("status").WriteInt32(*u.Status)
	}
	return nil
}

func (u *UpdateStatus_PartialUpdate_Set_Fields) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(u.MarshalFields)
}

func (u *UpdateStatus_PartialUpdate_Set_Fields) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = u.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

var UpdateStatus_PartialUpdate_Set_FieldsRequiredFields = restlicodec.NewRequiredFields()

func (u *UpdateStatus_PartialUpdate_Set_Fields) UnmarshalField(reader restlicodec.Reader, field string) (found bool, err error) {
	switch field {
	case "status":
		found = true
		u.Status = new(int32)
		*u.Status, err = reader.ReadInt32()
	case "error":
		found = true
		u.Error = new(ErrorResponse)
		err = u.Error.UnmarshalRestLi(reader)
	}
	return found, err
}

func (u *UpdateStatus_PartialUpdate_Set_Fields) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	err = reader.ReadRecord(UpdateStatus_PartialUpdate_Set_FieldsRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		found, err := u.UnmarshalField(reader, field)
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

func (u *UpdateStatus_PartialUpdate_Set_Fields) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, u)
}

func (u *UpdateStatus_PartialUpdate_Set_Fields) NewInstance() *UpdateStatus_PartialUpdate_Set_Fields {
	return new(UpdateStatus_PartialUpdate_Set_Fields)
}

// UpdateStatus_PartialUpdate is used to represent a partial update on UpdateStatus. Toggling the value of a field
// in Delete represents selecting it for deletion in a partial update, while
// setting the value of a field in Update represents setting that field in the
// current struct. Other fields in this struct represent record fields that can
// themselves be partially updated.
type UpdateStatus_PartialUpdate struct {
	Delete_Fields UpdateStatus_PartialUpdate_Delete_Fields
	Set_Fields    UpdateStatus_PartialUpdate_Set_Fields
	Error         *ErrorResponse_PartialUpdate
}

func (u *UpdateStatus_PartialUpdate) MarshalRestLiPatch(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
		checker := patch.PartialUpdateFieldChecker{RecordType: "com.linkedin.restli.common.UpdateStatus"}
		if err = checker.CheckField(writer, "status", false, u.Set_Fields.Status != nil, false); err != nil {
			return err
		}
		if err = checker.CheckField(writer, "error", u.Delete_Fields.Error, u.Set_Fields.Error != nil, u.Error != nil); err != nil {
			return err
		}
		if checker.HasDeletes {
			err = u.Delete_Fields.MarshalRestLi(keyWriter("$delete"))
			if err != nil {
				return err
			}
		}

		if checker.HasSets {
			err = u.Set_Fields.MarshalRestLi(keyWriter("$set"))
			if err != nil {
				return err
			}
		}

		if u.Error != nil {
			err = u.Error.MarshalRestLiPatch(keyWriter("error"))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (u *UpdateStatus_PartialUpdate) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
		return u.MarshalRestLiPatch(keyWriter(patch.PatchField).SetScope())
	})
}

func (u *UpdateStatus_PartialUpdate) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = u.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

func (u *UpdateStatus_PartialUpdate) UnmarshalRestLiPatch(reader restlicodec.Reader) (err error) {
	err = reader.ReadMap(func(reader restlicodec.Reader, key string) (err error) {
		switch key {
		case "$delete":
			err = u.Delete_Fields.UnmarshalRestLi(reader)
		case "$set":
			err = u.Set_Fields.UnmarshalRestLi(reader)
		case "error":
			u.Error = new(ErrorResponse_PartialUpdate)
			err = u.Error.UnmarshalRestLiPatch(reader)
		default:
			err = reader.Skip()
		}
		return err
	})
	if err != nil {
		return err
	}
	checker := patch.PartialUpdateFieldChecker{RecordType: "com.linkedin.restli.common.UpdateStatus"}
	if err = checker.CheckField(reader, "status", false, u.Set_Fields.Status != nil, false); err != nil {
		return err
	}
	if err = checker.CheckField(reader, "error", u.Delete_Fields.Error, u.Set_Fields.Error != nil, u.Error != nil); err != nil {
		return err
	}
	return nil
}

func (u *UpdateStatus_PartialUpdate) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	return reader.ReadRecord(patch.RequiredPatchRecordFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == patch.PatchField {
			return u.UnmarshalRestLiPatch(reader)
		} else {
			return reader.Skip()
		}
	})
}

func (u *UpdateStatus_PartialUpdate) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, u)
}

func (u *UpdateStatus_PartialUpdate) NewInstance() *UpdateStatus_PartialUpdate {
	return new(UpdateStatus_PartialUpdate)
}
