package draftbpnm

type FlowBpnm struct {
	process map[string]*processBpnm
}

type processBpnm struct {
	name               string
	activeActivities   []*activity
	defaultStart       string
	requiredGlobalData []typeData
	description        string
	activities         map[string]*activity
	storageBpnm        StorageData
}

type activity struct {
	name               string
	handler            activityHandler
	description        string
	preActivity        *processBpnm
	postActivity       *processBpnm
	recover            activityHandler
	requiredOutputData []typeData
	requiredInputData  []typeData
	//	GatewayType        GatewayType
	gateway       []*gateway
	callEndIfStop bool
}

type gateway struct {
	nextActivities []*activity
	decision       string
}
