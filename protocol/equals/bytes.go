package equals

import "bytes"

func Bytes(left, right []byte) bool {
	return bytes.Equal(left, right)
}

func BytesPointer(left, right *[]byte) bool {
	return GenericPointer(left, right, Bytes)
}

func BytesArray(left, right [][]byte) bool {
	return GenericArray(left, right, Bytes)
}

func BytesMap(left, right map[string][]byte) bool {
	return GenericMap(left, right, Bytes)
}

func BytesArrayPointer(left, right *[][]byte) bool {
	return GenericArrayPointer(left, right, Bytes)
}

func BytesMapPointer(left, right *map[string][]byte) bool {
	return GenericMapPointer(left, right, Bytes)
}
