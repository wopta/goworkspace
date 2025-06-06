package flow

type String struct {
	String string
}

func (*String) GetType() string {
	return "string"
}

type BoolBpmn struct {
	Bool bool
}

func (*BoolBpmn) GetType() string {
	return "bool"
}
