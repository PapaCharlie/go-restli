package protocol

import (
	"encoding/json"
)

type Bytes []byte

func (b *Bytes) MarshalJSON() (data []byte, err error) {
	return json.Marshal(string(*b))
}

func (b *Bytes) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*b = Bytes(s)
	return nil
}
