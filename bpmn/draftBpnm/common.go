package draftbpnm

type ActivityHandler func(StorageData) error
type DataBpnm interface {
	Type() string
}

// type date that you can use inside bpmn
type Error struct {
	Step        string
	Description string
	Result      bool
}

func (e *Error) Type() string {
	return "error"
}
