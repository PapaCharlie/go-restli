package restlicodec

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var readRawBytesUnsupported = errors.New("go-restli: ReadRawBytes not supported when reading from `any`")

func NewInterfaceReader(v any) Reader {
	return NewInterfaceReaderWithExcludedFields(v, nil, 0)
}

func NewInterfaceReaderWithExcludedFields(v any, excludedFields PathSpec, leadingScopeToIgnore int) Reader {
	return &anyReader{
		missingFieldsTracker: newMissingFieldsTracker(excludedFields, leadingScopeToIgnore),
		atStart:              true,
		value:                v,
	}
}

type anyReader struct {
	missingFieldsTracker
	atStart bool
	value   any
}

func (a *anyReader) val() (v reflect.Value, err error) {
	v = reflect.ValueOf(a.value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v, fmt.Errorf("go-restli: Illegal nil pointer")
		}
		v = v.Elem()
	}

	return v, nil
}

func (a *anyReader) atInputStart() bool {
	return a.atStart
}

func (a *anyReader) String() string {
	return fmt.Sprint(a.value)
}

func (a *anyReader) ReadInt() (int, error) {
	return readInt[int](a)
}

func (a *anyReader) ReadInt32() (int32, error) {
	return readInt[int32](a)
}

func (a *anyReader) ReadInt64() (i int64, err error) {
	return readInt[int64](a)
}

func readInt[T int | int32 | int64](a *anyReader) (i T, err error) {
	v, err := a.val()
	if err != nil {
		return i, err
	}

	s, isString := readString(v)
	switch {
	case isString:
		var i64 int64
		i64, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("go-restli: Invalid integer string %q: %w", v.String(), err)
		}
		return T(i64), nil
	case v.CanInt():
		return T(v.Int()), nil
	case v.CanFloat():
		return T(v.Float()), nil
	default:
		return 0, a.cannotPrimitive(reflect.TypeOf(i))
	}
}

func (a *anyReader) ReadFloat32() (float32, error) {
	return readFloat[float32](a)
}

func (a *anyReader) ReadFloat64() (f float64, err error) {
	return readFloat[float64](a)
}

func readFloat[T float32 | float64](a *anyReader) (f T, err error) {
	v, err := a.val()
	if err != nil {
		return f, err
	}

	s, isString := readString(v)
	switch {
	case isString:
		var f64 float64
		f64, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("go-restli: Invalid floating point string %q: %w", v.String(), err)
		}
		return T(f64), nil
	case v.CanInt():
		return T(v.Int()), nil
	case v.CanFloat():
		return T(v.Float()), nil
	default:
		return 0, a.cannotPrimitive(reflect.TypeOf(f))
	}
}

func (a *anyReader) ReadBool() (b bool, err error) {
	v, err := a.val()
	if err != nil {
		return b, err
	}

	s, isString := readString(v)
	switch {
	case isString:
		b, err = strconv.ParseBool(s)
		if err != nil {
			return false, fmt.Errorf("go-restli: Invalid boolean string %q: %w", v.String(), err)
		}
		return b, nil
	case v.Kind() == reflect.Bool:
		return v.Bool(), nil
	default:
		return false, a.cannotPrimitive(reflect.TypeOf(b))
	}
}

func (a *anyReader) ReadString() (string, error) {
	v, err := a.val()
	if err != nil {
		return "", err
	}

	s, isString := readString(v)
	if isString {
		return s, nil
	} else {
		return "", a.cannotPrimitive(reflect.TypeOf(""))
	}
}

func readString(v reflect.Value) (string, bool) {
	switch {
	case v.Kind() == reflect.String:
		return v.String(), true
	case v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8:
		return string(v.Interface().([]byte)), true
	default:
		return "", false
	}
}

func (a *anyReader) ReadBytes() ([]byte, error) {
	return readBytes(a.ReadString())
}

func (a *anyReader) ReadMap(mapReader MapReader) (err error) {
	v, err := a.val()
	if err != nil {
		return err
	}

	if v.Kind() != reflect.Map || v.Type().Key().Kind() != reflect.String {
		return &InvalidTypeError{
			DesiredType: reflect.TypeOf(map[string]any(nil)),
			ActualType:  v.Type(),
			Location:    a.scopeString(),
		}
	}

	oldValue, oldAtStart := a.value, a.atStart
	for r := v.MapRange(); r.Next(); {
		a.value = r.Value().Interface()
		key := r.Key().String()
		err = a.enterMapScope(key)
		if err != nil {
			return err
		}
		err = mapReader(a, r.Key().String())
		if err != nil {
			return err
		}
	}
	a.value, a.atStart = oldValue, oldAtStart

	return nil
}

func (a *anyReader) ReadRecord(requiredFields RequiredFields, recordReader MapReader) error {
	return readRecord(a, requiredFields, recordReader)
}

func (a *anyReader) ReadArray(arrayReader ArrayReader) (err error) {
	v, err := a.val()
	if err != nil {
		return err
	}

	if v.Kind() != reflect.Slice {
		return &InvalidTypeError{
			DesiredType: reflect.TypeOf([]any(nil)),
			ActualType:  v.Type(),
			Location:    a.scopeString(),
		}
	}

	oldValue, oldAtStart := a.value, a.atStart
	for i := 0; i < v.Len(); i++ {
		a.value = v.Index(i).Interface()
		a.enterArrayScope(i)
		err = arrayReader(a)
		if err != nil {
			return err
		}
	}
	a.value, a.atStart = oldValue, oldAtStart

	return nil
}

func (a *anyReader) ReadInterface() (interface{}, error) {
	return a.value, nil
}

func (a *anyReader) ReadRawBytes() ([]byte, error) {
	return nil, readRawBytesUnsupported
}

func (a *anyReader) Skip() error {
	return nil
}

type InvalidTypeError struct {
	DesiredType, ActualType reflect.Type
	Location                string
}

func (c *InvalidTypeError) Error() string {
	return fmt.Sprintf("go-restli: Cannot read %q from %q value at %q", c.DesiredType, c.ActualType, c.Location)
}

func (a *anyReader) cannotPrimitive(desiredType reflect.Type) error {
	return &InvalidTypeError{
		DesiredType: desiredType,
		ActualType:  reflect.TypeOf(a.value),
		Location:    a.scopeString(),
	}
}
