package protocol

import "github.com/PapaCharlie/go-restli/protocol/restlicodec"

const (
	PatchField = "patch"
)

var RequiredPatchRecordFields = restlicodec.RequiredFields{PatchField}

type PartialUpdateFieldChecker struct {
	RecordType string
	HasDeletes bool
	HasSets    bool
}

func (c *PartialUpdateFieldChecker) CheckField(
	checker restlicodec.KeyChecker,
	fieldName string,
	isDeleteSet bool,
	isSetSet bool,
	isPartialUpdateSet bool,
) error {
	if !(isDeleteSet || isSetSet || isPartialUpdateSet) {
		return nil
	}

	if checker.IsKeyExcluded(fieldName) {
		return &IllegalPartialUpdateError{
			Message:    "Cannot delete/update/partial update read-only or create-ony",
			Field:      fieldName,
			RecordType: c.RecordType,
		}
	}

	if (isDeleteSet && isSetSet) || (isDeleteSet && isPartialUpdateSet) || (isSetSet && isPartialUpdateSet) {
		return &IllegalPartialUpdateError{
			Message:    "Only one of set/update/partial update can be specified for",
			Field:      fieldName,
			RecordType: c.RecordType,
		}
	}

	if isDeleteSet {
		c.HasDeletes = true
	}

	if isSetSet {
		c.HasSets = true
	}

	return nil
}
