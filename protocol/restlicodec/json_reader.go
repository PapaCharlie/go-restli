package restlicodec

import (
	"github.com/mailru/easyjson/jlexer"
)

type jsonReader struct {
	lexer jlexer.Lexer
}

func NewJsonReader(data []byte) Reader {
	return &jsonReader{lexer: jlexer.Lexer{Data: data}}
}

func (j *jsonReader) ReadMap(mapReader func(field string) error) (err error) {
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
		fieldName := j.lexer.UnsafeFieldName(false)
		j.lexer.WantColon()
		if j.lexer.IsNull() {
			j.lexer.Skip()
			j.lexer.WantComma()
			continue
		}

		err = mapReader(fieldName)
		if err != nil {
			return err
		}

		j.lexer.WantComma()
	}

	j.lexer.Delim('}')
	if isTopLevel {
		j.lexer.Consumed()
	}

	return j.lexer.Error()
}

func (j *jsonReader) ReadArray(arrayReader func() error) (err error) {
	if j.lexer.IsNull() {
		j.lexer.Skip()
	} else {
		j.lexer.Delim('[')
		for !j.lexer.IsDelim(']') {
			err = arrayReader()
			if err != nil {
				return err
			}
			j.lexer.WantComma()
		}
		j.lexer.Delim(']')
	}
	return j.lexer.Error()
}

func (j *jsonReader) Skip() error {
	j.lexer.SkipRecursive()
	return j.lexer.Error()
}

func (j *jsonReader) ReadInt32() (int32, error) {
	return j.lexer.Int32(), j.lexer.Error()
}

func (j *jsonReader) ReadInt64() (int64, error) {
	return j.lexer.Int64(), j.lexer.Error()
}

func (j *jsonReader) ReadFloat32() (float32, error) {
	return j.lexer.Float32(), j.lexer.Error()
}

func (j *jsonReader) ReadFloat64() (float64, error) {
	return j.lexer.Float64(), j.lexer.Error()
}

func (j *jsonReader) ReadBool() (bool, error) {
	return j.lexer.Bool(), j.lexer.Error()
}

func (j *jsonReader) ReadString() (string, error) {
	return j.lexer.String(), j.lexer.Error()
}

func (j *jsonReader) ReadBytes() ([]byte, error) {
	return []byte(j.lexer.String()), j.lexer.Error()
}
