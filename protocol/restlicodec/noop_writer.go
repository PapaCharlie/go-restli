package restlicodec

type NoopWriter struct{}

var noopWriter NoopWriter

func (n NoopWriter) WriteInt32(int32)     {}
func (n NoopWriter) WriteInt64(int64)     {}
func (n NoopWriter) WriteFloat32(float32) {}
func (n NoopWriter) WriteFloat64(float64) {}
func (n NoopWriter) WriteBool(bool)       {}
func (n NoopWriter) WriteString(string)   {}
func (n NoopWriter) WriteBytes([]byte)    {}

func (n NoopWriter) WriteRawBytes([]byte) {
}

func (n NoopWriter) WriteMap(MapWriter) error {
	return nil
}

func (n NoopWriter) WriteArray(ArrayWriter) error {
	return nil
}

func (n NoopWriter) IsKeyExcluded(string) bool {
	return false
}

func (n NoopWriter) SetScope(...string) Writer {
	return n
}
