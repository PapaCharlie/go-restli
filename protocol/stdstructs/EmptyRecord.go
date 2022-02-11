package stdstructs

import "github.com/PapaCharlie/go-restli/protocol/restlicodec"

var EmptyRecord emptyRecord

type emptyRecord struct{}

func (e emptyRecord) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadMap(func(restlicodec.Reader, string) error { return reader.Skip() })
}

func (e emptyRecord) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(func(string) restlicodec.Writer) error { return nil })
}
