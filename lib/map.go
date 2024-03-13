package lib

func GetMapKeys[M map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func GetMapValues[M map[K]V, K comparable, V any](m M) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
