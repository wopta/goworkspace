package draftbpnm

type Order struct {
	InWhatProcessInjected  string        `json:"inWhatProcessInjected"`
	InWhatActivityInjected string        `json:"inWhatActivityInjected"`
	Order                  OrderActivity `json:"order"`
}
type BpnmBuilder struct {
	handlers  map[string]ActivityHandler
	storage   StorageData
	Processes []*ProcessBuilder `json:"processes"`
	toInject  map[KeyInject]*ProcessBpnm
}

type TypeData struct {
	Name string
	Type string
}

type ProcessBuilder struct {
	GlobalDataRequired []TypeData        `json:"globalData"`
	Order              *Order            `json:"order"`
	Description        string            `json:"description"`
	Name               string            `json:"name"`
	Activities         []ActivityBuilder `json:"activities"`
	DefaultStart       string            `json:"defaultStart"`
}

type ActivityBuilder struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Branch      *BranchBuilder `json:"branch"`
	HandlerLess bool           `json:"handlerless"`
	Recover     string         `json:"recover"`
}

type BranchBuilder struct {
	OutputDataRequired []TypeData     `json:"outputData,omitempty"`
	InputDataRequired  []TypeData     `json:"inputData,omitempty"`
	Gateways           []GatewayBlock `json:"gateways"`
	//unused field
	//GatewayType        GatewayType    `json:"gatewayType,omitempty"`
	//Recorver           string         `json:"recorver,omitempty"`
}

type GatewayBlock struct {
	NextActivities []string `json:"nextActivities"`
	Decision       string   `json:"decision,omitempty"`
}

type OrderActivity string

const (
	PreActivity  OrderActivity = "pre"
	PostActivity OrderActivity = "post"
)

type InjectActivity struct {
	Order              OrderActivity
	ActivitiesToInject []ActivityBuilder
}
type KeyInject struct {
	TargetProcess  string
	TargetActivity string
	OrderActivity
}
