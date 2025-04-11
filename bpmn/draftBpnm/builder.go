package draftbpnm

import (
	"errors"
	"fmt"
)

type BpnmBuilder struct {
	handlers  map[string]ActivityHandler
	storage   StorageData
	Processes []*ProcessBuilder `json:"processes"`
}

type TypeData struct {
	Name string
	Type string
}
type ProcessBuilder struct {
	GlobalData           []TypeData        `json:"globalData"`
	Description          string            `json:"description"`
	Name                 string            `json:"name"`
	Activities           []ActivityBuilder `json:"activities"`
	activitiesRegistered []string
	gatewayRegistered    []string
}

type ActivityBuilder struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Branch      *BranchBuilder `json:"branch"`
}

type BranchBuilder struct {
	OutputData  []TypeData     `json:"outputData,omitempty"`
	InputData   []TypeData     `json:"inputData,omitempty"`
	GatewayType string         `json:"gatewayType,omitempty"`
	Gateways    []GatewayBlock `json:"gateways"`
	Recorver    string         `json:"recorver,omitempty"`
}

type GatewayBlock struct {
	NextActivities []string `json:"nextActivities"`
	Decision       string   `json:"decision,omitempty"`
}

func (b *BpnmBuilder) Build() (*FlowBpnm, error) {
	flow := new(FlowBpnm)
	var process *ProcessBpnm
	for _, p := range b.Processes {
		process = new(ProcessBpnm)
		process.storageBpnm = b.storage
		process.Description = p.Description
		process.Name = p.Name
		process.RequiredGlobalData = p.GlobalData
		buildedActivities, err := b.BuildActivity(p.Activities)
		if err != nil {
			return nil, err
		}
		process.Activities = buildedActivities
		if err := b.BuildGatewayBlock(p, process); err != nil {
			return nil, err
		}
		flow.Process = append(flow.Process, process)
	}
	return flow, nil
}

func (b *BpnmBuilder) AddHandler(nameHandler string, handler ActivityHandler) error {
	if b.handlers == nil {
		b.handlers = make(map[string]ActivityHandler)
	}
	if _, ok := b.handlers[nameHandler]; ok {
		return errors.New("Handler's been already defined")
	}
	b.handlers[nameHandler] = handler
	return nil
}

func (b *BpnmBuilder) SetPoolDate(pool StorageData) {
	b.storage = pool
}

func (a *BpnmBuilder) BuildActivity(activities []ActivityBuilder) (map[string]*Activity, error) {
	result := make(map[string]*Activity)
	for _, activity := range activities {
		if _, ok := result[activity.Name]; ok {
			return nil, fmt.Errorf("Double event with same name %v", activity.Name)
		}
		newActivity := new(Activity)
		handler, ok := a.handlers[activity.Name]
		if !ok {
			return nil, fmt.Errorf("No handler registered for the activity: %v", activity.Name)
		}
		newActivity.Name = activity.Name
		newActivity.Description = activity.Description
		newActivity.handler = handler
		if activity.Branch == nil {
			return nil, fmt.Errorf("No branch founded for activity: %v", activity.Name)
		}
		buildedBranch, e := a.BuildBranchBuilder(activity.Branch)
		if e != nil {
			return nil, e
		}
		newActivity.Branch = buildedBranch
		result[newActivity.Name] = newActivity
	}
	return result, nil
}

func (a *BpnmBuilder) BuildBranchBuilder(gatewayDto *BranchBuilder) (*Branch, error) {
	activity := new(Branch)
	activity.GatewayType = GatewayType(gatewayDto.GatewayType)
	activity.RequiredInputData = gatewayDto.InputData
	activity.RequiredOutputData = gatewayDto.OutputData
	return activity, nil
}

func (a *BpnmBuilder) BuildGatewayBlock(processBuilder *ProcessBuilder, process *ProcessBpnm) error {
	for _, builderActivity := range processBuilder.Activities {
		var gateways []*Gateway = make([]*Gateway, 0)
		for _, builderGateway := range builderActivity.Branch.Gateways {
			gateway := &Gateway{
				NextActivities: make([]*Activity, 0),
				Decision:       builderGateway.Decision,
			}
			for _, nextJump := range builderGateway.NextActivities {
				if _, ok := process.Activities[nextJump]; !ok {
					return fmt.Errorf("No event named %v", nextJump)
				}
				gateway.NextActivities = append(gateway.NextActivities, process.Activities[nextJump])
			}
			gateways = append(gateways, gateway)
		}
		process.Activities[builderActivity.Name].Branch.Gateway = gateways
	}
	return nil
}
