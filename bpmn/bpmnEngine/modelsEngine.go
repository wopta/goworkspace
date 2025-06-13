package bpmnEngine

type FlowBpnm struct {
	process map[string]*processBpnm
}

type processBpnm struct {
	name               string
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
	gateway            []*gateway
	callEndIfStop      bool
}

type gateway struct {
	nextActivities []*activity
	decision       string
}

type StatusFlow struct {
	Parent          *StatusFlow
	CurrentProcess  string
	CurrentActivity string
}

func (*StatusFlow) GetType() string {
	return "_statusFlow"
}
