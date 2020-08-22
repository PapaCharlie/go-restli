package protocol

import (
	"github.com/PapaCharlie/go-restli/protocol/restliencoding"
)

type PartialUpdate struct {
	Patch restliencoding.Encodable
}

func (p *PartialUpdate) RestLiEncode(encoder *restliencoding.Encoder) error {
	encoder.WriteObjectStart()
	encoder.WriteFieldNameAndDelimiter("patch")
	err := encoder.Encodable(p.Patch)
	if err != nil {
		return err
	}
	encoder.WriteObjectEnd()
	return nil
}

func (p *PartialUpdate) String() string {
	encoder := restliencoding.NewCompactJsonEncoder()
	_ = p.RestLiEncode(encoder)
	return encoder.Finalize()
}

type EmptyRecord struct{}

func (e *EmptyRecord) RestLiEncode(encoder *restliencoding.Encoder) error {
	encoder.WriteObjectStart()
	encoder.WriteObjectEnd()
	return nil
}
