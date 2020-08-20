package restlicodec

import (
	"fmt"
	"testing"
)

func TestEncoder(t *testing.T) {
	encoders := []*Encoder{
		NewCompactJsonEncoder(),
		NewPrettyJsonEncoder(),
		NewHeaderEncoder(),
		NewPathEncoder(),
		NewQueryEncoder(),
		NewFinderEncoder(),
	}

	for _, e := range encoders {
		e.WriteObjectStart()
		e.Int32Field("int32", 32)
		e.WriteFieldDelimiter()
		e.Int64Field("int64", 64)
		e.WriteFieldDelimiter()
		e.StringMapField("map", map[string]string{
			"foo": "bar",
		})
		e.WriteFieldDelimiter()
		e.BoolArrayField("arr", []bool{true, false, true})
		e.WriteObjectEnd()
		fmt.Println(e.Finalize())
		fmt.Println(e.Finalize())
	}
}
