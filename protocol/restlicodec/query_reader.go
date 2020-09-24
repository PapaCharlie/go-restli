package restlicodec

import (
	"fmt"
	"net/url"
)

type RestLiQueryParamsReader interface {
	ReadParams(paramsReader MapReader) error
}

type queryParamsReader struct {
	pos  int
	data string
}

func NewRestLiQueryParamsReader(rawQuery string) RestLiQueryParamsReader {
	return &queryParamsReader{data: rawQuery}
}

func (r *queryParamsReader) ReadParams(paramsReader MapReader) (err error) {
	for {
		fieldName, err := r.readFieldName()
		if err != nil {
			return err
		}
		if fieldName == "" {
			break
		}
		data := r.readField()
		err = validateRor2Input(data)
		if err != nil {
			return err
		}
		reader := &ror2Reader{
			decoder: url.QueryUnescape,
			data:    []byte(data),
		}
		err = paramsReader(reader, fieldName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *queryParamsReader) readFieldName() (string, error) {
	startPos := r.pos
	for ; r.pos < len(r.data); r.pos++ {
		if r.data[r.pos] == '=' {
			name := r.data[startPos:r.pos]
			r.pos++
			return name, nil
		}
	}
	if startPos == r.pos {
		return "", nil
	} else {
		return "", fmt.Errorf("illegal query has incorrect field delimiter: %q", r.data)
	}
}

func (r *queryParamsReader) readField() string {
	startPos := r.pos

	for ; r.pos < len(r.data); r.pos++ {
		if r.data[r.pos] == '&' {
			break
		}
	}

	field := r.data[startPos:r.pos]
	r.pos++
	return field
}
