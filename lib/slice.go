package lib

func SliceRemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func SliceInsertIndex[T any](destination []T, element T, index int) []T {
	if len(destination) == index {
		return append(destination, element)
	}

	destination = append(destination[:index+1], destination[index:]...) // index < len(a)
	destination[index] = element

	return destination
}
func Make2D[T any](n, m int) [][]T {
	matrix := make([][]T, n)
	rows := make([]T, n*m)
	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}
	return matrix
}

func SliceContains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func SliceFilter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func SliceMap[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func SliceToSet[T comparable](s []T) []T {
	inResult := make(map[T]bool)
	var result []T
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func SliceFlatten[T any](s [][]T) []T {
	itemCount := 0
	for i := 0; i < len(s); i++ {
		itemCount = itemCount + len(s[i])
	}

	r := make([]T, itemCount)
	idx := 0
	for i := 0; i < len(s); i++ {
		for j := 0; j < len(s[i]); j++ {
			r[idx] = s[i][j]
			idx++
		}
	}
	return r
}

func SliceFlatMap[T, U any](s [][]T, f func(T) U) []U {
	return SliceFlatten(SliceMap(s, func(t []T) []U { return SliceMap(t, f) }))
}
