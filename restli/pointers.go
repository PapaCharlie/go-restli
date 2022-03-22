package restli

// IntPointer returns a pointer to the given parameter, useful for inlining setting optional fields
func IntPointer(v int) *int {
	return &v
}

// Int32Pointer returns a pointer to the given parameter, useful for inlining setting optional fields
func Int32Pointer(v int32) *int32 {
	return &v
}

// Int64Pointer returns a pointer to the given parameter, useful for inlining setting optional fields
func Int64Pointer(v int64) *int64 {
	return &v
}

// Float32Pointer returns a pointer to the given parameter, useful for inlining setting optional fields
func Float32Pointer(v float32) *float32 {
	return &v
}

// Float64Pointer returns a pointer to the given parameter, useful for inlining setting optional fields
func Float64Pointer(v float64) *float64 {
	return &v
}

// BoolPointer returns a pointer to the given parameter, useful for inlining setting optional fields
func BoolPointer(v bool) *bool {
	return &v
}

// StringPointer returns a pointer to the given parameter, useful for inlining setting optional fields
func StringPointer(v string) *string {
	return &v
}

// BytesPointer returns a pointer to the given parameter, useful for inlining setting optional fields
func BytesPointer(v []byte) *[]byte {
	return &v
}
