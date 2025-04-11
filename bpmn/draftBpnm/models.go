package draftbpnm

type FlowBpnm struct {
	Process []*ProcessBpnm
}

type ProcessBpnm struct {
	Name               string
	activeActivity     *Activity
	RequiredGlobalData []TypeData
	Description        string
	Activities         map[string]*Activity
	storageBpnm        StorageData
}

type Activity struct {
	Name        string
	handler     ActivityHandler
	Description string
	Branch      *Branch
}
type GatewayType string

const (
	XOR GatewayType = "XOR"
	AND GatewayType = "AND"
)

type Branch struct {
	RequiredOutputData []TypeData
	RequiredInputData  []TypeData
	GatewayType        GatewayType
	Gateway            []*Gateway
}
type Gateway struct {
	NextActivities []*Activity
	Decision       string
}
