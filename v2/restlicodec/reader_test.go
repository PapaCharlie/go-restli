package restlicodec

import (
	"encoding/json"
	"testing"
)

var (
	Err           error
	unmarshalData = []byte(`{"status":500}`)
	Obj           = &Object{Status: 10}
)

func BenchmarkMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Err = Obj.MarshalRestLi(NoopWriter)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reader, _ := NewJsonReader(unmarshalData)
		Err = new(Object).UnmarshalRestLi(reader)
	}
}

func BenchmarkReflectMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Err = MarshalRestLi[*Object](Obj, NoopWriter)
	}
}

func BenchmarkReflectUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reader, _ := NewJsonReader(unmarshalData)
		_, Err = UnmarshalRestLi[*Object](reader)
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, Err = json.Marshal(Obj)
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Err = json.Unmarshal(unmarshalData, new(Object))
	}
}

type Object struct {
	Status int32
}

func (e *Object) NewInstance() *Object {
	return new(Object)
}

func (e *Object) MarshalRestLi(writer Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) Writer) (err error) {
		keyWriter("status").WriteInt32(e.Status)
		return nil
	})
}

func (e *Object) UnmarshalRestLi(reader Reader) (err error) {
	err = reader.ReadRecord(nil, func(reader Reader, field string) (err error) {
		switch field {
		case "status":
			e.Status, err = reader.ReadInt32()
		default:
			err = NoSuchFieldErr
		}
		return err
	})
	if err != nil {
		return err
	}

	return err
}
