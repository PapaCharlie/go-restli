/*
Code generated by "github.com/PapaCharlie/go-restli"; DO NOT EDIT.

Source file: /Users/pchesnai/code/personal/go-restli/spec-parser/build/libs/go-restli-spec-parser-2.0.2.jar
*/

package common

import (
	"github.com/PapaCharlie/go-restli/fnv1a"
	"github.com/PapaCharlie/go-restli/restli/equals"
	"github.com/PapaCharlie/go-restli/restli/patch"
	"github.com/PapaCharlie/go-restli/restlicodec"
)

// Metadata and pagination links for this collection
type CollectionMetadata struct {
	// The start index of this collection
	Start int32
	// The number of elements in this collection segment
	Count int32
	// The total number of elements in the entire collection (not just this segment)
	Total *int32

	Links []*Link
}

// Sanity check NewCollectionMetadataWithDefaultValues has no illegal default values
var _ = NewCollectionMetadataWithDefaultValues()

func NewCollectionMetadataWithDefaultValues() (c *CollectionMetadata) {
	c = new(CollectionMetadata)
	c.populateLocalDefaultValues()
	return
}

func (c *CollectionMetadata) populateLocalDefaultValues() {
	if c.Total == nil {
		val := int32(0)
		c.Total = &val
	}

}

func (c *CollectionMetadata) Equals(other *CollectionMetadata) bool {
	if c == other {
		return true
	}
	if c == nil || other == nil {
		return false
	}

	return c.Start == other.Start &&
		c.Count == other.Count &&
		equals.ComparablePointer(c.Total, other.Total) &&
		equals.ObjectArray(c.Links, other.Links)
}

func (c *CollectionMetadata) ComputeHash() fnv1a.Hash {
	if c == nil {
		return fnv1a.ZeroHash()
	}
	hash := fnv1a.NewHash()

	hash.AddInt32(c.Start)

	hash.AddInt32(c.Count)

	if c.Total != nil {
		hash.AddInt32(*c.Total)
	}

	fnv1a.AddHashableArray(hash, c.Links)

	return hash
}

func (c *CollectionMetadata) MarshalFields(keyWriter func(string) restlicodec.Writer) (err error) {
	keyWriter("count").WriteInt32(c.Count)
	err = restlicodec.WriteArray(keyWriter("links"), c.Links, (*Link).MarshalRestLi)
	if err != nil {
		return err
	}
	keyWriter("start").WriteInt32(c.Start)
	if c.Total != nil {
		keyWriter("total").WriteInt32(*c.Total)
	}
	return nil
}

func (c *CollectionMetadata) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(c.MarshalFields)
}

func (c *CollectionMetadata) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = c.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

var CollectionMetadataRequiredFields = restlicodec.NewRequiredFields().Add(
	"start",
	"count",
	"links",
)

func (c *CollectionMetadata) UnmarshalField(reader restlicodec.Reader, field string) (found bool, err error) {
	switch field {
	case "start":
		found = true
		c.Start, err = reader.ReadInt32()
	case "count":
		found = true
		c.Count, err = reader.ReadInt32()
	case "total":
		found = true
		c.Total = new(int32)
		*c.Total, err = reader.ReadInt32()
	case "links":
		found = true
		c.Links, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[*Link])
	}
	return found, err
}

func (c *CollectionMetadata) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	err = reader.ReadRecord(CollectionMetadataRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		found, err := c.UnmarshalField(reader, field)
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

	c.populateLocalDefaultValues()
	return nil
}

func (c *CollectionMetadata) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, c)
}

func (c *CollectionMetadata) NewInstance() *CollectionMetadata {
	return new(CollectionMetadata)
}

/*
================================================================================
PARTIAL UPDATE STRUCTS
================================================================================
*/

type CollectionMetadata_PartialUpdate_Delete_Fields struct {
	Total bool
}

func (c *CollectionMetadata_PartialUpdate_Delete_Fields) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteArray(func(itemWriter func() restlicodec.Writer) (err error) {
		if c.Total {
			itemWriter().WriteString("total")
		}
		return nil
	})
}

func (c *CollectionMetadata_PartialUpdate_Delete_Fields) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = c.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

func (c *CollectionMetadata_PartialUpdate_Delete_Fields) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	var field string
	return reader.ReadArray(func(reader restlicodec.Reader) (err error) {
		field, err = reader.ReadString()
		if err != nil {
			return err
		}

		switch field {
		case "total":
			c.Total = true
		}
		return nil
	})
}

func (c *CollectionMetadata_PartialUpdate_Delete_Fields) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, c)
}

func (c *CollectionMetadata_PartialUpdate_Delete_Fields) NewInstance() *CollectionMetadata_PartialUpdate_Delete_Fields {
	return new(CollectionMetadata_PartialUpdate_Delete_Fields)
}

type CollectionMetadata_PartialUpdate_Set_Fields struct {
	// start
	Start *int32
	// count
	Count *int32
	// total
	Total *int32
	// links
	Links *[]*Link
}

func (c *CollectionMetadata_PartialUpdate_Set_Fields) MarshalFields(keyWriter func(string) restlicodec.Writer) (err error) {
	if c.Count != nil {
		keyWriter("count").WriteInt32(*c.Count)
	}
	if c.Links != nil {
		err = restlicodec.WriteArray(keyWriter("links"), *c.Links, (*Link).MarshalRestLi)
		if err != nil {
			return err
		}
	}
	if c.Start != nil {
		keyWriter("start").WriteInt32(*c.Start)
	}
	if c.Total != nil {
		keyWriter("total").WriteInt32(*c.Total)
	}
	return nil
}

func (c *CollectionMetadata_PartialUpdate_Set_Fields) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(c.MarshalFields)
}

func (c *CollectionMetadata_PartialUpdate_Set_Fields) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = c.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

var CollectionMetadata_PartialUpdate_Set_FieldsRequiredFields = restlicodec.NewRequiredFields()

func (c *CollectionMetadata_PartialUpdate_Set_Fields) UnmarshalField(reader restlicodec.Reader, field string) (found bool, err error) {
	switch field {
	case "start":
		found = true
		c.Start = new(int32)
		*c.Start, err = reader.ReadInt32()
	case "count":
		found = true
		c.Count = new(int32)
		*c.Count, err = reader.ReadInt32()
	case "total":
		found = true
		c.Total = new(int32)
		*c.Total, err = reader.ReadInt32()
	case "links":
		found = true
		c.Links = new([]*Link)
		*c.Links, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[*Link])
	}
	return found, err
}

func (c *CollectionMetadata_PartialUpdate_Set_Fields) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	err = reader.ReadRecord(CollectionMetadata_PartialUpdate_Set_FieldsRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		found, err := c.UnmarshalField(reader, field)
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

func (c *CollectionMetadata_PartialUpdate_Set_Fields) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, c)
}

func (c *CollectionMetadata_PartialUpdate_Set_Fields) NewInstance() *CollectionMetadata_PartialUpdate_Set_Fields {
	return new(CollectionMetadata_PartialUpdate_Set_Fields)
}

// CollectionMetadata_PartialUpdate is used to represent a partial update on CollectionMetadata. Toggling the value of a field
// in Delete represents selecting it for deletion in a partial update, while
// setting the value of a field in Update represents setting that field in the
// current struct. Other fields in this struct represent record fields that can
// themselves be partially updated.
type CollectionMetadata_PartialUpdate struct {
	Delete_Fields CollectionMetadata_PartialUpdate_Delete_Fields
	Set_Fields    CollectionMetadata_PartialUpdate_Set_Fields
}

func (c *CollectionMetadata_PartialUpdate) MarshalRestLiPatch(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
		checker := patch.PartialUpdateFieldChecker{RecordType: "com.linkedin.restli.common.CollectionMetadata"}
		if err = checker.CheckField(writer, "start", false, c.Set_Fields.Start != nil, false); err != nil {
			return err
		}
		if err = checker.CheckField(writer, "count", false, c.Set_Fields.Count != nil, false); err != nil {
			return err
		}
		if err = checker.CheckField(writer, "total", c.Delete_Fields.Total, c.Set_Fields.Total != nil, false); err != nil {
			return err
		}
		if err = checker.CheckField(writer, "links", false, c.Set_Fields.Links != nil, false); err != nil {
			return err
		}
		if checker.HasDeletes {
			err = c.Delete_Fields.MarshalRestLi(keyWriter("$delete"))
			if err != nil {
				return err
			}
		}

		if checker.HasSets {
			err = c.Set_Fields.MarshalRestLi(keyWriter("$set"))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (c *CollectionMetadata_PartialUpdate) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
		return c.MarshalRestLiPatch(keyWriter(patch.PatchField).SetScope())
	})
}

func (c *CollectionMetadata_PartialUpdate) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = c.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

func (c *CollectionMetadata_PartialUpdate) UnmarshalRestLiPatch(reader restlicodec.Reader) (err error) {
	err = reader.ReadMap(func(reader restlicodec.Reader, key string) (err error) {
		switch key {
		case "$delete":
			err = c.Delete_Fields.UnmarshalRestLi(reader)
		case "$set":
			err = c.Set_Fields.UnmarshalRestLi(reader)
		default:
			err = reader.Skip()
		}
		return err
	})
	if err != nil {
		return err
	}
	checker := patch.PartialUpdateFieldChecker{RecordType: "com.linkedin.restli.common.CollectionMetadata"}
	if err = checker.CheckField(reader, "start", false, c.Set_Fields.Start != nil, false); err != nil {
		return err
	}
	if err = checker.CheckField(reader, "count", false, c.Set_Fields.Count != nil, false); err != nil {
		return err
	}
	if err = checker.CheckField(reader, "total", c.Delete_Fields.Total, c.Set_Fields.Total != nil, false); err != nil {
		return err
	}
	if err = checker.CheckField(reader, "links", false, c.Set_Fields.Links != nil, false); err != nil {
		return err
	}
	return nil
}

func (c *CollectionMetadata_PartialUpdate) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	return reader.ReadRecord(patch.RequiredPatchRecordFields, func(reader restlicodec.Reader, field string) (err error) {
		if field == patch.PatchField {
			return c.UnmarshalRestLiPatch(reader)
		} else {
			return reader.Skip()
		}
	})
}

func (c *CollectionMetadata_PartialUpdate) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, c)
}

func (c *CollectionMetadata_PartialUpdate) NewInstance() *CollectionMetadata_PartialUpdate {
	return new(CollectionMetadata_PartialUpdate)
}
