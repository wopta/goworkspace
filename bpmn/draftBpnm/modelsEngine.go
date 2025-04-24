package draftbpnm

type FlowBpnm struct {
	Process map[string]*ProcessBpnm
}

type ProcessBpnm struct {
	Name               string
	activeActivities   []*Activity
	DefaultStart       string
	RequiredGlobalData []typeData
	Description        string
	Activities         map[string]*Activity
	storageBpnm        StorageData
}

type Activity struct {
	Name               string
	handler            ActivityHandler
	Description        string
	PreActivity        *ProcessBpnm
	PostActivity       *ProcessBpnm
	recover            ActivityHandler
	RequiredOutputData []typeData
	RequiredInputData  []typeData
	//	GatewayType        GatewayType
	Gateway []*Gateway
}

type GatewayType string

const (
	XOR GatewayType = "XOR"
	AND GatewayType = "AND"
)

type Gateway struct {
	NextActivities []*Activity
	Decision       string
}
