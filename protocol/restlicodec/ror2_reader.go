package restlicodec

import (
	"fmt"
	"net/url"
	"strconv"
)

type ror2ReaderState int

const (
	noState = ror2ReaderState(iota)
	inMap
	inArray
)

const list = "List("

func (s *ror2ReaderState) location() string {
	switch *s {
	case inMap:
		return "map"
	case inArray:
		return "array"
	default:
		return "string"
	}
}

type ror2Reader struct {
	missingFieldsTracker
	decoder func(string) (string, error)
	data    []byte
	pos     int
	state   ror2ReaderState
}

func validateRor2Input(data string) error {
	parens := 0
	for i, c := range data {
		switch c {
		case '(':
			parens++
		case ')':
			parens--
			if parens < 0 {
				return &DeserializationError{
					Err: fmt.Errorf("illegal ROR2 string has unbalanced delimiters at %d: %q", i, data),
				}
			}
		}
	}
	return nil
}

// NewRor2Reader returns a new Reader that reads objects serialized using the rest.li protocol 2.0 object and array
// representation (ROR2), whose spec is defined here:
// https://linkedin.github.io/rest.li/spec/protocol#restli-protocol-20-object-and-listarray-representation
// Because the "reduced" URL encoding used for rest.li headers is a subset of the standard URL encoding, this Reader can
// be used for both the "full" URL encoding and the "reduced" URL encoding (though query parameters should be read with
// the reader returned by NewRor2QueryReader as there exist encoding differences, namely " " being encoding as `%20` and
// `+` respectively)
// An error will be returned if an upfront validation of the given string reveals it is not a valid ROR2 string. Note
// that if this function does not return an error, it does _not_ mean subsequent calls to the Read* functions will not
// return an error
func NewRor2Reader(data string) (Reader, error) {
	err := validateRor2Input(data)
	if err != nil {
		return nil, err
	}
	return &ror2Reader{
		decoder: url.PathUnescape,
		data:    []byte(data),
	}, nil
}

// readFieldName advances the current position until a ':' is encountered or the end of the data is reached, and returns
// the string found between the starting position and the end position.
func (u *ror2Reader) readFieldName() (string, error) {
	if u.data[u.pos] == ')' {
		return ")", nil
	}
	startPos := u.pos
loop:
	for ; u.pos < len(u.data); u.pos++ {
		switch u.data[u.pos] {
		case ':':
			if u.pos == startPos {
				return "", u.errorf("invalid ROR2 %s string has empty field name at %d: %q", u.state.location(), u.pos, string(u.data))
			}
			break loop
		case ',', ')':
			return "", u.errorf("invalid ROR2 %s string has invalid field name at %d: %q", u.state.location(), u.pos, string(u.data))
		}
	}
	s := string(u.data[startPos:u.pos])
	u.pos++
	return s, nil
}

func (u *ror2Reader) ReadMap(mapReader MapReader) (err error) {
	u.state = inMap
	if !u.atMap() {
		return u.errorf("invalid ROR2 %s string does not start with '(': %q", u.state.location(), string(u.data))
	}
	u.pos++
	if len(u.data) < u.pos {
		return u.errorf("invalid ROR2 %s is unclosed: %q", u.state.location(), string(u.data))
	}

loop:
	for {
		var fieldName string
		fieldName, err = u.readFieldName()
		if err != nil {
			return err
		}
		switch fieldName {
		case ")":
			u.pos++
			break loop
		default:
			u.enterMapScope(fieldName)
			err = mapReader(u, fieldName)
			if err != nil {
				return err
			}
			u.exitScope()
			switch u.data[u.pos] {
			case ',':
				u.pos++
				continue loop
			case ')':
				u.pos++
				break loop
			default:
				return u.errorf("invalid ROR2 string does has incorrect object delimiter at %d: %q", u.pos, string(u.data))
			}
		}
	}

	return nil
}

func (u *ror2Reader) atMap() bool {
	return len(u.data) > u.pos && u.data[u.pos] == '('
}

func (u *ror2Reader) ReadArray(arrayReader ArrayReader) (err error) {
	u.state = inArray
	if !u.atArray() {
		return u.errorf("invalid ROR2 %s string does not start with "+list+": %q", u.state.location(), string(u.data))
	}
	u.pos += len(list)
	if len(u.data) < u.pos {
		return u.errorf("invalid ROR2 %s string is not closed %q", u.state.location(), string(u.data))
	}
	if u.data[u.pos] == ')' {
		u.pos++
		return nil
	}

	index := 0
loop:
	for {
		u.enterArrayScope(index)
		err = arrayReader(u)
		if err != nil {
			return err
		}
		u.exitScope()
		index++
		switch u.data[u.pos] {
		case ',':
			u.pos++
			continue loop
		case ')':
			u.pos++
			break loop
		default:
			return u.errorf("invalid ROR2 string has incorrect array delimiter at %d: %q", u.pos, string(u.data))
		}
	}

	return nil
}

func (u *ror2Reader) atArray() bool {
	return len(u.data)-u.pos > len(list) && string(u.data[u.pos:u.pos+len(list)]) == list
}

func (u *ror2Reader) ReadInterface() (interface{}, error) {
	switch {
	case u.atMap():
		m := make(map[string]interface{})
		err := u.ReadMap(func(reader Reader, k string) (err error) {
			var v interface{}
			v, err = reader.ReadInterface()
			if err != nil {
				return err
			}

			m[k] = v
			return nil
		})
		if err != nil {
			return nil, err
		}
		return m, nil
	case u.atArray():
		var a []interface{}
		err := u.ReadArray(func(reader Reader) (err error) {
			var v interface{}
			v, err = reader.ReadInterface()
			if err != nil {
				return err
			}
			a = append(a, v)
			return nil
		})
		if err != nil {
			return nil, err
		}
		return a, nil
	default:
		return u.ReadString()
	}
}

func (u *ror2Reader) Skip() error {
	if u.pos == 0 {
		u.pos = len(u.data)
		return nil
	}

	inMapOrArray := u.atArray() || u.atMap()
	parens := 0
	for ; u.pos < len(u.data); u.pos++ {
		switch u.data[u.pos] {
		case '(':
			if !inMapOrArray {
				return u.errorf("unescaped '(' at %d: %q", u.pos, string(u.data))
			} else {
				parens++
			}
		case ',':
			if !inMapOrArray {
				return nil
			} else {
				if parens == 0 {
					return nil
				}
			}
		case ')':
			if !inMapOrArray {
				return nil
			} else {
				if parens == 0 {
					return nil
				}
				parens--
			}
		}
	}
	return u.errorf("invalid ROR2 string has incorrect object delimiters: %q", string(u.data))
}

func (u *ror2Reader) ReadRawBytes() ([]byte, error) {
	startPos := u.pos
	err := u.Skip()
	if err != nil {
		return nil, u.wrapDeserializationError(err)
	} else {
		return u.data[startPos:u.pos], nil
	}
}

// unsafeReadPrimitiveFieldValue moves the current position forward until an end-of-field delimiter is reached, i.e.
// either ',' or ')'. Only primitive values can be read by this method, therefore any unescaped delimiter characters are
// considered illegal and an error is returned. If the current position is 0, then it is assumed that the data is one
// single top-level primitive and the position will be moved to the end of the data. Otherwise if the end of the data is
// reached without encountering an end-of-field delimiter, an error will be returned.
func (u *ror2Reader) unsafeReadPrimitiveFieldValue() (startPos int, err error) {
	startPos = u.pos
	if startPos == 0 {
		u.pos = len(u.data)
	} else {
		for ; u.pos < len(u.data); u.pos++ {
			if c := u.data[u.pos]; c == ',' || c == ')' {
				break
			}
		}
		if u.pos == len(u.data) {
			return 0, u.errorf("invalid ROR2 string has incorrect field delimiters: %q", string(u.data))
		}
	}
	for i := startPos; i < u.pos; i++ {
		// Check that the field didn't contain any illegal unescaped characters
		if c := u.data[i]; c == '(' || c == ',' || c == ')' {
			return 0, u.errorf("illegal unescaped '%s' in primitive field at %d: %q", string(c), i, string(u.data))
		}
	}
	return startPos, nil
}

// readAndDecodePrimitiveFieldValue calls unsafeReadPrimitiveFieldValue and decodes the field's value using the given
// decoder.
func (u *ror2Reader) readAndDecodePrimitiveFieldValue() (string, error) {
	startPos, err := u.unsafeReadPrimitiveFieldValue()
	if err != nil {
		return "", err
	}
	s, err := u.decoder(string(u.data[startPos:u.pos]))
	return s, u.wrapDeserializationError(err)
}

func (u *ror2Reader) ReadInt() (v int, err error) {
	var v64 int64
	v64, err = u.ReadInt64()
	if err != nil {
		return v, err
	}
	return int(v64), nil
}

func (u *ror2Reader) ReadInt32() (v int32, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	i, err := strconv.ParseInt(decoded, 10, 32)
	return int32(i), u.wrapDeserializationError(err)
}

func (u *ror2Reader) ReadInt64() (v int64, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	v, err = strconv.ParseInt(decoded, 10, 64)
	return v, u.wrapDeserializationError(err)
}

func (u *ror2Reader) ReadFloat32() (v float32, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	f, err := strconv.ParseFloat(decoded, 32)
	return float32(f), u.wrapDeserializationError(err)
}

func (u *ror2Reader) ReadFloat64() (v float64, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	v, err = strconv.ParseFloat(decoded, 64)
	return v, u.wrapDeserializationError(err)
}

func (u *ror2Reader) ReadBool() (v bool, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	v, err = strconv.ParseBool(decoded)
	return v, u.wrapDeserializationError(err)
}

func (u *ror2Reader) ReadString() (v string, err error) {
	bytes, err := u.ReadBytes()
	return string(bytes), err
}

func (u *ror2Reader) ReadBytes() (v []byte, err error) {
	startPos, err := u.unsafeReadPrimitiveFieldValue()
	if err != nil {
		return nil, err
	}

	s := string(u.data[startPos:u.pos])
	if len(s) == 0 {
		return nil, u.errorf("invalid empty element at %d: %q", startPos, string(u.data))
	}
	if s == emptyString {
		return nil, nil
	}
	var decoded string
	decoded, err = u.decoder(s)
	return []byte(decoded), u.wrapDeserializationError(err)
}

func (u *ror2Reader) errorf(format string, args ...interface{}) error {
	return &DeserializationError{
		Scope: u.scopeString(),
		Err:   fmt.Errorf(format, args...),
	}
}

func (u *ror2Reader) AtInputStart() bool {
	return u.pos == 0
}
