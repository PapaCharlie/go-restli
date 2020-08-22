package restlidecoding

import (
	"errors"

	"github.com/mailru/easyjson/jlexer"
)

type jsonDecoder struct {
	lexer jlexer.Lexer
}

func NewJsonDecoder(data []byte) Decoder {
	return &jsonDecoder{lexer: jlexer.Lexer{Data: data}}
}

// FieldIgnoredError is used as a return value from ReadObject to indicate that the given field was unknown and
// subsequently ignored
var FieldIgnoredError = errors.New("a non-fatal error when an unknown field was encountered and subsequently ignored")

func (j *jsonDecoder) ReadObject(objectReader func(field string) error) (err error) {
	isTopLevel := j.lexer.IsStart()
	if j.lexer.IsNull() {
		if isTopLevel {
			j.lexer.Consumed()
		}
		j.lexer.Skip()
		return j.lexer.Error()
	}

	j.lexer.Delim('{')
	for !j.lexer.IsDelim('}') {
		key := j.lexer.UnsafeFieldName(false)
		j.lexer.WantColon()
		if j.lexer.IsNull() {
			j.lexer.Skip()
			j.lexer.WantComma()
			continue
		}

		err = objectReader(key)
		if err == FieldIgnoredError {
			j.lexer.SkipRecursive()
		} else if err != nil {
			return err
		}

		j.lexer.WantComma()
	}
	j.lexer.Delim('}')
	if isTopLevel {
		j.lexer.Consumed()
	}

	return nil
}

func (j *jsonDecoder) ReadArray(arrayReader func(index int) error) (err error) {
	if j.lexer.IsNull() {
		j.lexer.Skip()
	} else {
		j.lexer.Delim('[')
		i := 0
		for !j.lexer.IsDelim(']') {
			err = arrayReader(i)
			if err != nil {
				return err
			}
			i++
			j.lexer.WantComma()
		}
		j.lexer.Delim(']')
	}
	return nil
}

func (j *jsonDecoder) ReadMap(mapReader func(key string) error) error {
	return j.ReadObject(mapReader)
}

func (j *jsonDecoder) Int32() (int32, error) {
	return j.lexer.Int32(), j.lexer.Error()
}

func (j *jsonDecoder) Int64() (int64, error) {
	return j.lexer.Int64(), j.lexer.Error()
}

func (j *jsonDecoder) Float32() (float32, error) {
	return j.lexer.Float32(), j.lexer.Error()
}

func (j *jsonDecoder) Float64() (float64, error) {
	return j.lexer.Float64(), j.lexer.Error()
}

func (j *jsonDecoder) Bool() (bool, error) {
	return j.lexer.Bool(), j.lexer.Error()
}

func (j *jsonDecoder) String() (string, error) {
	return j.lexer.String(), j.lexer.Error()
}

func (j *jsonDecoder) Bytes() ([]byte, error) {
	return []byte(j.lexer.String()), j.lexer.Error()
}

func (j *jsonDecoder) Decodable(decodable Decodable) error {
	return decodable.RestliDecode(j.SubDecoder())
}

func (j *jsonDecoder) SubDecoder() Decoder {
	return j
}
