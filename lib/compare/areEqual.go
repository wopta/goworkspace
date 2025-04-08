package compare

type Comparable[T any] interface {
	IsEqual(T) bool
}

// Check if 'a' and 'b' are equal, it uses 'IsEqual' method.
func AreEqual[T Comparable[T]](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil && b != nil {
		return false
	}
	return (*a).IsEqual(*b)
}

// Check if 'a' and 'b' are equal, it uses isEqual parameter.
func AreEqualFunc[T any](a, b *T, isEqual func(a, b *T) bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil && b != nil {
		return false
	}
	return isEqual(a, b)
}

// Check if the elements of 'a' and 'b' are equal, it uses IsEqual method.
func AreSlicesEqual[T Comparable[T]](a, b []*T) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil && b != nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !(*a[i]).IsEqual(*b[i]) {
			return false
		}
	}
	return true
}

// Check if the eleemnts of 'a' and 'b' are equal, it uses isEqual parameter.
func AreSlicesEqualFunc[T any](a, b *[]T, isEqual func(T, T) bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil && b != nil {
		return false
	}
	if len(*a) != len(*b) {
		return false
	}
	for i := range *a {
		if !isEqual((*a)[i], (*b)[i]) {
			return false
		}
	}
	return true
}
