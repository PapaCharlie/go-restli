package protocol

import (
	"fmt"
	"strconv"
)

func (r *RestLiCodec) EncodeInt32(v int32) string {
	return fmt.Sprintf("%d", v)
}

func (r *RestLiCodec) EncodeInt64(v int64) string {
	return fmt.Sprintf("%d", v)
}

func (r *RestLiCodec) EncodeFloat32(v float32) string {
	return fmt.Sprintf("%g", v)
}

func (r *RestLiCodec) EncodeFloat64(v float64) string {
	return fmt.Sprintf("%g", v)
}

func (r *RestLiCodec) EncodeBool(v bool) string {
	return fmt.Sprintf("%t", v)
}

func (r *RestLiCodec) EncodeString(v string) string {
	return r.encoder(v)
}

func (r *RestLiCodec) EncodeBytes(v Bytes) string {
	return r.EncodeString(string(v))
}

func (r *RestLiCodec) DecodeInt32(data string, v *int32) error {
	i, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		return err
	}
	*v = int32(i)
	return nil
}

func (r *RestLiCodec) DecodeInt64(data string, v *int64) error {
	i, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return err
	}
	*v = int64(i)
	return nil
}

func (r *RestLiCodec) DecodeFloat32(data string, v *float32) error {
	f, err := strconv.ParseFloat(data, 32)
	if err != nil {
		return err
	}
	*v = float32(f)
	return nil
}

func (r *RestLiCodec) DecodeFloat64(data string, v *float64) error {
	f, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return err
	}
	*v = float64(f)
	return nil
}

func (r *RestLiCodec) DecodeBool(data string, v *bool) error {
	b, err := strconv.ParseBool(data)
	if err != nil {
		return err
	}
	*v = bool(b)
	return nil
}

func (r *RestLiCodec) DecodeString(data string, v *string) error {
	s, err := r.decoder(data)
	if err != nil {
		return err
	}
	*v = s
	return nil
}

func (r *RestLiCodec) DecodeBytes(data string, v *Bytes) error {
	var s string
	err := r.DecodeString(data, &s)
	if err != nil {
		return err
	}
	*v = Bytes(s)
	return nil
}
