/*
DO NOT EDIT

Code automatically generated by github.com/PapaCharlie/go-restli
Source file: https://github.com/PapaCharlie/go-restli/blob/master/codegen/resources/restlidata.go
*/

package restlidata

import (
	fnv1a "github.com/PapaCharlie/go-restli/fnv1a"
	equals "github.com/PapaCharlie/go-restli/restli/equals"
	restlicodec "github.com/PapaCharlie/go-restli/restlicodec"
)

type CollectionMedata struct {
	// The start index of this collection
	Start int32
	// The number of elements in this collection segment
	Count int32
	// The total number of elements in the entire collection (not just this segment)
	Total *int32

	Links []*Link
}

func (c *CollectionMedata) Equals(other *CollectionMedata) bool {
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

func (c *CollectionMedata) ComputeHash() fnv1a.Hash {
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

func (c *CollectionMedata) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
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
	})
}

func (c *CollectionMedata) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = c.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

// Sanity check NewCollectionMedataWithDefaultValues has no illegal default values
var _ = NewCollectionMedataWithDefaultValues()

func NewCollectionMedataWithDefaultValues() (c *CollectionMedata) {
	c = new(CollectionMedata)
	c.populateLocalDefaultValues()
	return
}

func (c *CollectionMedata) populateLocalDefaultValues() {
	if c.Total == nil {
		val := int32(0)
		c.Total = &val
	}

}

var _CollectionMedataRequiredFields = restlicodec.RequiredFields{
	"start",
	"count",
	"links",
}

func (c *CollectionMedata) NewInstance() *CollectionMedata {
	return new(CollectionMedata)
}

func (c *CollectionMedata) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	err = reader.ReadRecord(_CollectionMedataRequiredFields, func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case "start":
			c.Start, err = reader.ReadInt32()
		case "count":
			c.Count, err = reader.ReadInt32()
		case "total":
			c.Total = new(int32)
			*c.Total, err = reader.ReadInt32()
		case "links":
			c.Links, err = restlicodec.ReadArray(reader, restlicodec.UnmarshalRestLi[*Link])
		default:
			err = reader.Skip()
		}
		return err
	})
	if err != nil {
		return err
	}

	c.populateLocalDefaultValues()
	return err
}

func (c *CollectionMedata) UnmarshalJSON(data []byte) error {
	return restlicodec.UnmarshalJSON(data, c)
}
