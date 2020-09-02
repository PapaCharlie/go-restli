package restlicodec

import (
	"fmt"
	"net/url"
	"strconv"
)

type ror2ReaderState int

const (
	noState = ror2ReaderState(iota)
	inObject
	inArray
)

func (s *ror2ReaderState) location() string {
	switch *s {
	case inObject:
		return "object"
	case inArray:
		return "array"
	default:
		return "string"
	}
}

type ror2Reader struct {
	decoder func(string) (string, error)
	data    []byte
	pos     int
	state   ror2ReaderState
}

// NewRor2Reader returns a new Reader that reads objects serialized using the rest.li protocol 2.0 object and array
// representation (ROR2), whose spec is defined here:
// https://linkedin.github.io/rest.li/spec/protocol#restli-protocol-20-object-and-listarray-representation
// Because the "reduced" URL encoding used for rest.li headers is a subset of the standard URL encoding, this Reader can
// be used for both the "full" URL encoding and the "reduced" URL encoding.
// An error will be returned if an upfront validation of the given string reveals it is not a valid ROR2 string. Note
// that if this function does not return an error, it does _not_ mean subsequent calls to the Read* functions will not
// return an error
func NewRor2Reader(data string) (Reader, error) {
	parens := 0
	for i, c := range data {
		switch c {
		case '(':
			parens++
		case ')':
			parens--
			if parens < 0 {
				return nil, fmt.Errorf("illegal ROR2 string has unbalanced delimiters at %d: %s", i, data)
			}
		}
	}
	return &ror2Reader{
		decoder: url.QueryUnescape,
		data:    []byte(data),
	}, nil
}

// readFieldName advances the current position until a ':' is encountered or the end of the data is reached, and returns
// the string found between the starting position and the end position
func (u *ror2Reader) readFieldName() string {
	startPos := u.pos
	for ; u.pos < len(u.data); u.pos++ {
		if u.data[u.pos] == ':' {
			break
		}
	}
	s := string(u.data[startPos:u.pos])
	u.pos++
	return s
}

func (u *ror2Reader) ReadMap(mapReader MapReader) (err error) {
	u.state = inObject
	if len(u.data) <= u.pos || u.data[u.pos] != '(' {
		return fmt.Errorf("invalid ROR2 %s string does not start with '(': %s", u.state.location(), string(u.data))
	}
	u.pos++

loop:
	for {
		fieldName := u.readFieldName()
		switch fieldName {
		case "":
			return fmt.Errorf("invalid ROR2 %s string does not end with ')': %s", u.state.location(), string(u.data))
		case ")":
			// This should only happen if the empty map/object () was read
			// Consider sanity check that u.data[u.pos-1] == '(' ?
			break loop
		default:
			err = mapReader(u, fieldName)
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
				return fmt.Errorf("invalid ROR2 string does has incorrect object delimiter at %d: %s", u.pos, string(u.data))
			}
		}
	}

	return nil
}

func (u *ror2Reader) ReadArray(arrayReader ArrayReader) (err error) {
	u.state = inArray
	const list = "List("
	if len(u.data)-u.pos < len(list) || string(u.data[u.pos:u.pos+len(list)]) != list {
		return fmt.Errorf("invalid ROR2 %s string does not start with "+list+": %s", u.state.location(), string(u.data))
	}
	u.pos += len(list)
loop:
	for {
		err = arrayReader(u)
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
			return fmt.Errorf("invalid ROR2 string has incorrect array delimiter at %d: %s", u.pos, string(u.data))
		}
	}

	return nil
}

func (u *ror2Reader) Skip() error {
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
	return fmt.Errorf("invalid ROR2 string has incorrect object delimiters: %s", string(u.data))
}

func (u *ror2Reader) Raw() ([]byte, error) {
	startPos := u.pos
	err := u.Skip()
	if err != nil {
		return nil, err
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
			return 0, fmt.Errorf("invalid ROR2 string has incorrect field delimiters: %s", string(u.data))
		}
	}
	for i := startPos; i < u.pos; i++ {
		// Check that the field didn't contain any illegal unescaped characters
		if c := u.data[i]; c == '(' || c == ',' || c == ')' {
			return 0, fmt.Errorf("illegal unescaped '%s' in primitive field at %d: %s", string(c), i, string(u.data))
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
	return u.decoder(string(u.data[startPos:u.pos]))
}

func (u *ror2Reader) ReadInt32() (v int32, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	i, err := strconv.ParseInt(decoded, 10, 32)
	return int32(i), err
}

func (u *ror2Reader) ReadInt64() (v int64, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	return strconv.ParseInt(decoded, 10, 64)
}

func (u *ror2Reader) ReadFloat32() (v float32, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	f, err := strconv.ParseFloat(decoded, 32)
	return float32(f), err
}

func (u *ror2Reader) ReadFloat64() (v float64, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	return strconv.ParseFloat(decoded, 64)
}

func (u *ror2Reader) ReadBool() (v bool, err error) {
	var decoded string
	decoded, err = u.readAndDecodePrimitiveFieldValue()
	if err != nil {
		return v, err
	}
	return strconv.ParseBool(decoded)
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
		return nil, fmt.Errorf("invalid empty element at %d: %s", startPos, string(u.data))
	}
	if s == emptyString {
		return nil, nil
	}
	var decoded string
	decoded, err = u.decoder(s)
	return []byte(decoded), nil
}
