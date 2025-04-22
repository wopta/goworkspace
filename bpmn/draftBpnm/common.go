package draftbpnm

import "errors"

type ActivityHandler func(StorageData) error
type DataBpnm interface {
	GetType() string
}

type Error struct {
	Step        string
	Description string
	Result      bool
}

func (e *Error) GetType() string {
	return "error"
}

// mergeMaps merges two maps of the same key and value types.
// If both maps contain the same key, values from m2 will overwrite those from m1.
func mergeMaps[key comparable, out any](m1 map[key]out, m2 map[key]out) map[key]out {
	if m2 == nil {
		return m1
	}
	merged := make(map[key]out)
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
}

// mergeUniqueKeys merges two maps and returns a new map containing all key-value pairs.
// If the same key exists in both maps, an error is returned.
func mergeUniqueMaps[key comparable, out any](m1 map[key]out, m2 map[key]out) (map[key]out, error) {
	if m2 == nil {
		return m1, nil
	}
	merged := make(map[key]out)
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		if _, ok := merged[k]; ok {
			return nil, errors.New("The key '%v' is used in both maps")
		}
		merged[k] = v
	}
	return merged, nil
}
