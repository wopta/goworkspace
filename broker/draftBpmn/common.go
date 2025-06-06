package draftbpmn

import (
	"errors"
	"fmt"
	"strings"
)

type activityHandler func(StorageData) error
type DataBpnm interface {
	GetType() string
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

// Get a resource and return it.
// The priority is Local->Higher-Local..->Global->Higher-Global->Resource not found
func GetData[t DataBpnm](name string, storage StorageData) (t, error) {
	data, err := storage.GetLocal(name)
	var result t
	if err != nil {
		data, err = storage.GetGlobal(name)
	}
	if err != nil {
		return *new(t), err
	}

	result, ok := data.(t)
	if !ok {
		return *new(t), fmt.Errorf("Data '%v' has type %v, which differs from expected type '%v'", name, result.GetType(), data.GetType())
	}

	return result, nil
}

// Get a resource and assign it to 'data' parameter
// The priority is Local->Higher-Local..->Global->Higher-Global->Resource not found
func GetDataRef[t DataBpnm](name string, data *t, storage StorageData) (err error) {
	if data == nil {
		return errors.New("Reference can't be null")
	}
	*data, err = GetData[t](name, storage)
	return
}

// IsError gathers all errors(whether there are) and return them, otherwise return nil
func IsError(errs ...error) error {
	var res = strings.Builder{}
	for i := range errs {
		if errs[i] != nil {
			res.WriteString(errs[i].Error() + ",")
		}
	}
	if res.Len() == 0 {
		return nil
	}
	return errors.New(res.String())
}
