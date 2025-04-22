package draftbpnm

type FlowBpnm struct {
	Process map[string]*ProcessBpnm
}

type ProcessBpnm struct {
	Name               string
	activeActivities   []*Activity
	DefaultStart       string
	RequiredGlobalData []TypeData
	Description        string
	Activities         map[string]*Activity
	storageBpnm        StorageData
}

type Activity struct {
	Name         string
	handler      ActivityHandler
	Description  string
	PreActivity  *ProcessBpnm
	PostActivity *ProcessBpnm
	Branch       *Branch
}

type GatewayType string

const (
	XOR GatewayType = "XOR"
	AND GatewayType = "AND"
)

type Branch struct {
	RequiredOutputData []TypeData
	RequiredInputData  []TypeData
	//	GatewayType        GatewayType
	Gateway []*Gateway
}

type Gateway struct {
	NextActivities []*Activity
	Decision       string
}
