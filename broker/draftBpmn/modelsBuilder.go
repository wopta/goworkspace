package draftbpmn

type BpnmBuilder struct {
	Processes []*processBuilder `json:"processes"`

	handlers map[string]activityHandler
	storage  StorageData
	toInject map[keyInject]*processBpnm
}

type processBuilder struct {
	GlobalDataRequired []typeData        `json:"globalData"`
	Order              *order            `json:"order"`
	Description        string            `json:"description"`
	Name               string            `json:"name"`
	Activities         []activityBuilder `json:"activities"`
	DefaultStart       string            `json:"defaultStart"`
}

type order struct {
	InWhatProcessInject  string        `json:"inWhatProcessInject"`
	InWhatActivityInject string        `json:"inWhatActivityInject"`
	Order                orderActivity `json:"order"`
}

type activityBuilder struct {
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	HandlerLess        bool           `json:"handlerless"`
	Recover            string         `json:"recover"`
	OutputDataRequired []typeData     `json:"outputData,omitempty"`
	InputDataRequired  []typeData     `json:"inputData,omitempty"`
	Gateways           []gatewayBlock `json:"gateways"`
	CallEndIfStop      *bool          `json:"callEndIfStop"`
}

type gatewayBlock struct {
	NextActivities []string `json:"nextActivities"`
	Decision       string   `json:"decision,omitempty"`
}

type orderActivity string

const (
	preActivity  orderActivity = "pre"
	postActivity orderActivity = "post"
)

type keyInject struct {
	targetProcess  string
	targetActivity string
	orderActivity
}

type typeData struct {
	Name string
	Type string
}
