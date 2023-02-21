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
