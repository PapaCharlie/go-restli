package protocol

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type RestLiCodec struct {
	encoder func(string) string
	decoder func(string) (string, error)
}

var RestLiUrlEncoder = RestLiCodec{
	encoder: url.PathEscape,
	decoder: url.PathUnescape,
}

var RestLiReducedEncoder = RestLiCodec{
	encoder: strings.NewReplacer(
		",", url.PathEscape(","),
		"(", url.PathEscape("("),
		")", url.PathEscape(")"),
		"'", url.PathEscape("'"),
		":", url.PathEscape(":")).Replace,
	decoder: url.PathUnescape,
}

type RestLiEncodableType int

const (
	Record = RestLiEncodableType(iota)
	String
	Number
)

type RestLiEncodable interface {
	RestLiEncode(codec RestLiCodec) (data string, err error)
	RestLiDecode(codec RestLiCodec, data string) (err error)
}

func (r *RestLiCodec) encodeJson(v interface{}) (data string) {
	if n, ok := v.(json.Number); ok {
		return n.String()
	}
	if b, ok := v.(bool); ok {
		return fmt.Sprintf("%t", b)
	}
	if s, ok := v.(string); ok {
		return r.EncodeString(s)
	}

	var buf strings.Builder

	if a, ok := v.([]interface{}); ok {
		buf.WriteString("List(")
		for i, val := range a {
			if i != 0 {
				buf.WriteString(",")
			}
			buf.WriteString(r.encodeJson(val))
		}
		buf.WriteString(")")
	}

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString("(")
		first := true
		for key, val := range m {
			if first {
				first = false
			} else {
				buf.WriteByte(',')
			}
			buf.WriteString(r.encodeJson(key))
			buf.WriteByte(':')
			buf.WriteString(r.encodeJson(val))
		}
		buf.WriteString(")")
	}

	return buf.String()
}

//// This is currently implemented by serializing into JSON, deserializing into a map[string]interface{} and re-encoding
//// it in the rest.li format. This is convenient since we know the JSON serialization is correct, we can largely just
//// reuse it. It is possible to statically generate this, but would be difficult and possibly unwieldy
//func (r *RestLiCodec) Encode(v RestLiEncodable) (data string, err error) {
//	var jsonData []byte
//	jsonData, err = v.MarshalJSON()
//	if err != nil {
//		return
//	}
//
//	switch v.RestLiEncodable() {
//	case Number:
//		data = string(jsonData)
//		return
//	case String:
//		var s string
//		err = json.Unmarshal(jsonData, &s)
//		if err != nil {
//			return
//		}
//		data = r.EncodeString(s)
//		return
//	case Record:
//		var m map[string]interface{}
//		decoder := json.NewDecoder(bytes.NewBuffer(jsonData))
//		decoder.UseNumber()
//		err = decoder.Decode(&m)
//		if err != nil {
//			return
//		}
//		data = r.encodeJson(m)
//		return
//	default:
//		log.Panicln("unknown RestLiEncodableType", v.RestLiEncodable())
//	}
//
//	return
//}
//
//// This is currently implemented by serializing into JSON, deserializing into a map[string]interface{} and re-encoding
//// it in the rest.li format. This is convenient since we know the JSON serialization is correct, we can largely just
//// reuse it. It is possible to statically generate this, but would be difficult and possibly unwieldy
//func (r *RestLiCodec) Decode(data string, v RestLiEncodable) (err error) {
//	switch v.RestLiEncodable() {
//	case Number:
//		data = string(jsonData)
//		return
//	case String:
//		var s string
//		err = json.Unmarshal(jsonData, &s)
//		if err != nil {
//			return
//		}
//		data = r.EncodeString(s)
//		return
//	case Record:
//		var m map[string]interface{}
//		decoder := json.NewDecoder(bytes.NewBuffer(jsonData))
//		decoder.UseNumber()
//		err = decoder.Decode(&m)
//		if err != nil {
//			return
//		}
//		data = r.encodeJson(m)
//		return
//	default:
//		log.Panicln("unknown RestLiEncodableType", v.RestLiEncodable())
//	}
//
//	return
//}
