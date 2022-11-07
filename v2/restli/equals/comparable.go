package equals

func ComparablePointer[T comparable](left, right *T) bool {
	return GenericPointer(left, right, func(left, right T) bool {
		return left == right
	})
}

func ComparableArray[T comparable](left, right []T) bool {
	return GenericArray(left, right, func(left, right T) bool {
		return left == right
	})
}

func ComparableMap[T comparable](left, right map[string]T) bool {
	return GenericMap(left, right, func(left, right T) bool {
		return left == right
	})
}

func ComparableArrayPointer[T comparable](left, right *[]T) bool {
	return GenericArrayPointer(left, right, func(left, right T) bool {
		return left == right
	})
}

func ComparableMapPointer[T comparable](left, right *map[string]T) bool {
	return GenericMapPointer(left, right, func(left, right T) bool {
		return left == right
	})
}
