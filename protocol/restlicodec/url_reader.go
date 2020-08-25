package restlicodec

import (
	"fmt"
	"net/url"
	"strconv"
)

type urlReaderState int

const (
	noState = urlReaderState(iota)
	inObject
	inArray
)

func (s *urlReaderState) location() string {
	switch *s {
	case inObject:
		return "object"
	case inArray:
		return "array"
	default:
		return "string"
	}
}

type urlReader struct {
	decoder func(string) (string, error)
	data    string
	pos     int
	state   urlReaderState
}

func NewHeaderReader(header string) Reader {
	return &urlReader{
		decoder: url.QueryUnescape,
		data:    header,
	}
}

func (u *urlReader) readFieldName() string {
	startPos := u.pos
	for ; u.pos < len(u.data); u.pos++ {
		if u.data[u.pos] == ':' {
			break
		}
	}
	s := u.data[startPos:u.pos]
	u.pos++
	return s
}

func (u *urlReader) ReadMap(mapReader func(field string) error) (err error) {
	u.state = inObject
	if len(u.data) <= u.pos || u.data[u.pos] != '(' {
		return fmt.Errorf("invalid encoded %s does not start with '(': %s", u.state.location(), string(u.data))
	}
	u.pos++

loop:
	for {
		fieldName := u.readFieldName()
		switch fieldName {
		case "":
			return fmt.Errorf("invalid encoded string does not end with ')': %s", string(u.data))
		case ")":
			// This should only happen if the empty map/object () was read
			// Consider sanity check that u.data[u.pos-1] == '(' ?
			break loop
		default:
			err = mapReader(fieldName)
			if err != nil {
				return err
			}
			switch u.data[u.pos] {
			case ',':
				u.pos++
				continue loop
			case ')':
				u.pos++
				break loop
			default:
				return fmt.Errorf("invalid encoded string does has incorrect object delimiter at %d: %s", u.pos, string(u.data))
			}
		}
	}

	return nil
}

func (u *urlReader) ReadArray(arrayReader func() error) (err error) {
	u.state = inArray
	const list = "List("
	if len(u.data)-u.pos < len(list) || string(u.data[u.pos:u.pos+len(list)]) != list {
		return fmt.Errorf("invalid encoded %s does not start with "+list+": %s", u.state.location(), string(u.data))
	}
	u.pos += len(list)
loop:
	for {
		err = arrayReader()
		if err != nil {
			return err
		}
		switch u.data[u.pos] {
		case ',':
			u.pos++
			continue loop
		case ')':
			u.pos++
			break loop
		default:
			return fmt.Errorf("invalid encoded string does has incorrect array delimiter at %d: %s", u.pos, string(u.data))
		}
	}

	return nil
}

func (u *urlReader) Skip() error {
	parens := 0
	for ; u.pos < len(u.data); u.pos++ {
		switch u.data[u.pos] {
		case '(':
			parens++
		case ',':
			if parens == 0 {
				return nil
			}
		case ')':
			if parens == 0 {
				return nil
			}
			parens--
		}
	}
	return fmt.Errorf("invalid encoded string does has incorrect object delimiters: %s", string(u.data))
}

func (u *urlReader) unsafeReadPrimitiveFieldValue() (startPos int, err error) {
	startPos = u.pos
	if startPos == 0 {
		u.pos = len(u.data)
	} else {
	loop:
		for ; u.pos < len(u.data); u.pos++ {
			switch u.data[u.pos] {
			case ',', ')':
				break loop
			}
		}
		if u.pos == len(u.data) {
			return 0, fmt.Errorf("invalid encoded string has incorrect field delimiters: %s", string(u.data))
		}
	}
	for i := startPos; i < u.pos; i++ {
		if u.data[i] == '(' {
			return 0, fmt.Errorf("illegal unescaped '(' in primitive field at %d: %s", i, string(u.data))
		}
	}
	return startPos, nil
}

func (u *urlReader) readAndDecodePrimitiveFieldValue() (string, error) {
	startPos, err := u.unsafeReadPrimitiveFieldValue()
	if err != nil {
		return "", err
	}
	return u.decoder(string(u.data[startPos:u.pos]))
}

func (u *urlReader) ReadInt32() (v int32, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	i, err := strconv.ParseInt(decoded, 10, 32)
	return int32(i), err
}

func (u *urlReader) ReadInt64() (v int64, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	return strconv.ParseInt(decoded, 10, 64)
}

func (u *urlReader) ReadFloat32() (v float32, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	f, err := strconv.ParseFloat(decoded, 32)
	return float32(f), err
}

func (u *urlReader) ReadFloat64() (v float64, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	return strconv.ParseFloat(decoded, 64)
}

func (u *urlReader) ReadBool() (v bool, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	return strconv.ParseBool(decoded)
}

func (u *urlReader) ReadString() (v string, err error) {
	bytes, err := u.ReadBytes()
	return string(bytes), err
}

func (u *urlReader) ReadBytes() (v []byte, err error) {
	startPos, err := u.unsafeReadPrimitiveFieldValue()
	if err != nil {
		return nil, err
	}

	s := u.data[startPos:u.pos]
	if len(s) == 0 {
		return nil, fmt.Errorf("invalid empty element at %d: %s", startPos, string(u.data))
	}
	if s == emptyString {
		return nil, nil
	}
	var decoded string
	decoded, err = u.decoder(s)
	return []byte(decoded), nil
}
