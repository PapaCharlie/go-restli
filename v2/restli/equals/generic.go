package equals

type Comparable[T any] interface {
	Equals(T) bool
}

func GenericPointer[T any](left, right *T, equals func(left, right T) bool) bool {
	if left != right {
		if left == nil || right == nil {
			return false
		}

		if !equals(*left, *right) {
			return false
		}
	}

	return true
}

func GenericArray[T any](left, right []T, equals func(left, right T) bool) bool {
	if len(left) != len(right) {
		return false
	}

	for i, l := range left {
		if !equals(l, right[i]) {
			return false
		}
	}

	return true
}

func GenericMap[T any](left, right map[string]T, equals func(left, right T) bool) bool {
	if len(left) != len(right) {
		return false
	}

	for k, lv := range left {
		rv, ok := right[k]
		if !ok || !equals(lv, rv) {
			return false
		}
	}

	return true
}

func GenericArrayPointer[T any](left, right *[]T, equals func(left, right T) bool) bool {
	return GenericPointer(left, right, func(left, right []T) bool {
		return GenericArray(left, right, equals)
	})
}

func GenericMapPointer[T any](left, right *map[string]T, equals func(left, right T) bool) bool {
	return GenericPointer(left, right, func(left, right map[string]T) bool {
		return GenericMap(left, right, equals)
	})
}
