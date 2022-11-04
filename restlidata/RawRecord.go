package restlidata

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/PapaCharlie/go-restli/v2/fnv1a"
	"github.com/PapaCharlie/go-restli/v2/restlicodec"
)

// RawRecord is a container for arbitrary data. Because it gets parsed from raw JSON without any extra type information,
// it is expected that there will be unknown side effects, such as floats turning into integers or vice versa. This is
// designed to fill the gap for Java's DataMap-backed implementation of rest.li that supports untyped messages. To
// attempt deserialization into a know/real rest.li object, use the UnmarshalTo method
// Use at your own risk!
type RawRecord map[string]interface{}

// ComputeHash for a RawRecord always returns the 0-hash
func (r RawRecord) ComputeHash() fnv1a.Hash {
	return fnv1a.ZeroHash()
}

// Equals for a RawRecord always returns false, unless it is being compared with itself
func (r RawRecord) Equals(other RawRecord) bool {
	return reflect.ValueOf(r).Pointer() == reflect.ValueOf(other).Pointer()
}

func (r RawRecord) MarshalRestLi(writer restlicodec.Writer) error {
	return writeInterface(r, writer)
}

func (r RawRecord) NewInstance() RawRecord {
	return make(RawRecord)
}

func (r *RawRecord) UnmarshalRestLi(reader restlicodec.Reader) error {
	i, err := reader.ReadInterface()
	if err != nil {
		return err
	}

	*r = i.(map[string]interface{})
	return nil
}

// UnmarshalTo attempts to deserialize this RawRecord into the given object by serializing it first to JSON then calling
// the given object's unmarshal method
func (r *RawRecord) UnmarshalTo(obj restlicodec.Unmarshaler) error {
	return obj.UnmarshalRestLi(restlicodec.NewInterfaceReader(*r))
}

func writeInterface(value interface{}, writer restlicodec.Writer) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr && value == nil {
		return fmt.Errorf("go-restli: Illegal nil value for RawRecord")
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		writer.WriteFloat64(v.Float())
		return nil

	case reflect.Int32, reflect.Int64, reflect.Int:
		writer.WriteInt64(v.Int())
		return nil

	case reflect.Bool:
		writer.WriteBool(v.Bool())
		return nil

	case reflect.String:
		writer.WriteString(v.String())
		return nil

	case reflect.Array:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("go-restli: Illegal fixed-size array type (%s)", v.Type())
		}
		bytes := make([]byte, v.Len())
		for i := range bytes {
			bytes[i] = v.Index(i).Interface().(byte)
		}
		writer.WriteBytes(bytes)
		return nil
	case reflect.Slice:
		// rest.li doesn't support int8 as a standalone type. In other words, a slice of bytes should be treated as a
		// primitive byte string and not an array of individual bytes.
		if v.Type().Elem().Kind() == reflect.Uint8 {
			writer.WriteBytes(v.Bytes())
			return nil
		}
		return writer.WriteArray(func(itemWriter func() restlicodec.Writer) error {
			for i := 0; i < v.Len(); i++ {
				err := writeInterface(v.Index(i).Interface(), itemWriter())
				if err != nil {
					return err
				}
			}
			return nil
		})

	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("go-restli: Illegal map key type %s (must be `string`)", v.Type().Key())
		}

		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) error {
			for _, k := range keys {
				err := writeInterface(v.MapIndex(k).Interface(), keyWriter(k.String()))
				if err != nil {
					return err
				}
			}
			return nil
		})

	case reflect.Struct:
		marshaler, ok := value.(restlicodec.Marshaler)
		if !ok {
			return fmt.Errorf("go-restli: Received struct type %s that does not implement restlicodec.Marshaller "+
				"(you _must_ pass in a pointer to a generated struct in order to use the marshaler interface)",
				reflect.ValueOf(value).Type())
		}
		return marshaler.MarshalRestLi(writer)
	}

	return fmt.Errorf("go-restli: Unknown type for RawRecord value (%s)", v.Type())
}
