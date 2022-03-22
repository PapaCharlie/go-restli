package restlicodec

import "io"

var NoopWriter noopWriter

type noopWriter struct{}

func (n noopWriter) WriteInt(int)         {}
func (n noopWriter) WriteInt32(int32)     {}
func (n noopWriter) WriteInt64(int64)     {}
func (n noopWriter) WriteFloat32(float32) {}
func (n noopWriter) WriteFloat64(float64) {}
func (n noopWriter) WriteBool(bool)       {}
func (n noopWriter) WriteString(string)   {}
func (n noopWriter) WriteBytes([]byte)    {}

func (n noopWriter) WriteRawBytes([]byte) {
}

func (n noopWriter) WriteMap(MapWriter) error {
	return nil
}

func (n noopWriter) WriteArray(ArrayWriter) error {
	return nil
}

func (n noopWriter) IsKeyExcluded(string) bool {
	return false
}

func (n noopWriter) SetScope(...string) Writer {
	return n
}

var EmptyReader emptyReader

type emptyReader struct{}

func (e emptyReader) String() string { return "" }

func (e emptyReader) ReadInt() (int, error)         { return 0, io.EOF }
func (e emptyReader) ReadInt32() (int32, error)     { return 0, io.EOF }
func (e emptyReader) ReadInt64() (int64, error)     { return 0, io.EOF }
func (e emptyReader) ReadFloat32() (float32, error) { return 0, io.EOF }
func (e emptyReader) ReadFloat64() (float64, error) { return 0, io.EOF }
func (e emptyReader) ReadBool() (bool, error)       { return false, io.EOF }
func (e emptyReader) ReadString() (string, error)   { return "", io.EOF }
func (e emptyReader) ReadBytes() ([]byte, error)    { return nil, io.EOF }

func (e emptyReader) IsKeyExcluded(string) bool { return false }

func (e emptyReader) ReadMap(MapReader) error                    { return io.EOF }
func (e emptyReader) ReadRecord(RequiredFields, MapReader) error { return io.EOF }
func (e emptyReader) ReadArray(ArrayReader) error                { return io.EOF }
func (e emptyReader) ReadInterface() (interface{}, error)        { return nil, io.EOF }
func (e emptyReader) ReadRawBytes() ([]byte, error)              { return nil, io.EOF }
func (e emptyReader) Skip() error                                { return io.EOF }
