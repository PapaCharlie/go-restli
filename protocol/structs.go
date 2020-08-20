package protocol

import "github.com/PapaCharlie/go-restli/protocol/restlicodec"

type PartialUpdate struct {
	Patch restlicodec.Marshaler
}

func (p *PartialUpdate) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(fieldNameWriter func(fieldName string) restlicodec.Writer) error {
		return p.Patch.MarshalRestLi(fieldNameWriter("patch"))
	})
}

var EmptyRecord emptyRecord

type emptyRecord struct{}

func (e emptyRecord) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadMap(func(string) error { return reader.Skip() })
}

func (e emptyRecord) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(func(string) restlicodec.Writer) error { return nil })
}
