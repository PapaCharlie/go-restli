package equals

func ObjectPointer[T Comparable[T]](left, right *T) bool {
	return GenericPointer(left, right, func(left, right T) bool {
		return left.Equals(right)
	})
}

func ObjectArray[T Comparable[T]](left, right []T) bool {
	return GenericArray(left, right, func(left, right T) bool {
		return left.Equals(right)
	})
}

func ObjectMap[T Comparable[T]](left, right map[string]T) bool {
	return GenericMap(left, right, func(left, right T) bool {
		return left.Equals(right)
	})
}

func ObjectArrayPointer[T Comparable[T]](left, right *[]T) bool {
	return GenericArrayPointer(left, right, func(left, right T) bool {
		return left.Equals(right)
	})
}

func ObjectMapPointer[T Comparable[T]](left, right *map[string]T) bool {
	return GenericMapPointer(left, right, func(left, right T) bool {
		return left.Equals(right)
	})
}
