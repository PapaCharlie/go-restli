package restlicodec

import (
	"github.com/mailru/easyjson/jlexer"
)

type jsonReader struct {
	missingFieldsTracker
	lexer jlexer.Lexer
}

func NewJsonReader(data []byte) Reader {
	return &jsonReader{lexer: jlexer.Lexer{Data: data}}
}

func (j *jsonReader) ReadMap(mapReader MapReader) (err error) {
	isTopLevel := j.lexer.IsStart()
	if j.lexer.IsNull() {
		if isTopLevel {
			j.lexer.Consumed()
		}
		j.lexer.Skip()
		return j.checkError()
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

		j.enterMapScope(fieldName)
		err = mapReader(j, fieldName)
		if err != nil {
			return err
		}
		j.exitScope()

		j.lexer.WantComma()
	}

	j.lexer.Delim('}')
	if isTopLevel {
		j.lexer.Consumed()
	}

	return j.checkError()
}

func (j *jsonReader) ReadArray(arrayReader ArrayReader) (err error) {
	if j.lexer.IsNull() {
		j.lexer.Skip()
	} else {
		j.lexer.Delim('[')
		index := 0
		for !j.lexer.IsDelim(']') {
			j.enterArrayScope(index)
			err = arrayReader(j)
			if err != nil {
				return err
			}
			j.exitScope()
			index++
			j.lexer.WantComma()
		}
		j.lexer.Delim(']')
	}

	return j.checkError()
}

func (j *jsonReader) ReadInterface() (interface{}, error) {
	return j.lexer.Interface(), j.checkError()
}

func (j *jsonReader) Skip() error {
	j.lexer.SkipRecursive()
	return j.checkError()
}

func (j *jsonReader) ReadRawBytes() ([]byte, error) {
	return j.lexer.Raw(), j.checkError()
}

func (j *jsonReader) ReadInt() (int, error) {
	return j.lexer.Int(), j.checkError()
}

func (j *jsonReader) ReadInt32() (int32, error) {
	return j.lexer.Int32(), j.checkError()
}

func (j *jsonReader) ReadInt64() (int64, error) {
	return j.lexer.Int64(), j.checkError()
}

func (j *jsonReader) ReadFloat32() (float32, error) {
	f, err := j.ReadFloat64()
	return float32(f), err
}

func (j *jsonReader) ReadFloat64() (float64, error) {
	n, err := j.lexer.JsonNumber(), j.checkError()
	if err != nil {
		return 0, err
	}
	return n.Float64()
}

func (j *jsonReader) ReadBool() (bool, error) {
	return j.lexer.Bool(), j.checkError()
}

func (j *jsonReader) ReadString() (string, error) {
	return j.lexer.String(), j.checkError()
}

func (j *jsonReader) ReadBytes() ([]byte, error) {
	return []byte(j.lexer.String()), j.checkError()
}

func (j *jsonReader) checkError() error {
	return j.wrapDeserializationError(j.lexer.Error())
}

func (j *jsonReader) AtInputStart() bool {
	return j.lexer.IsStart()
}
