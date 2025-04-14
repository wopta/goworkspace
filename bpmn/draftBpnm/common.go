package draftbpnm

type ActivityHandler func(StorageData) error
type DataBpnm interface {
	Type() string
}

type Error struct {
	Step        string
	Description string
	Result      bool
}

func (e *Error) Type() string {
	return "error"
}

func mergeMaps(m1 map[string]any, m2 map[string]any) map[string]any {
	if m2 == nil {
		return m1
	}
	merged := make(map[string]any)
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
}
