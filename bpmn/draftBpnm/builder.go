package draftbpnm

import (
	"errors"
	"fmt"
	"log"
)

type BpnmBuilder struct {
	handlers  map[string]ActivityHandler
	storage   StorageData
	Processes []*ProcessBuilder `json:"processes"`
}

func (b *BpnmBuilder) Build() (*FlowBpnm, error) {
	flow := new(FlowBpnm)
	var process *ProcessBpnm
	b.AddHandler("init", func(st StorageData) error {
		log.Print("init bello")
		return nil
	})
	for _, p := range b.Processes {
		process = new(ProcessBpnm)
		process.storageBpnm = b.storage
		process.Description = p.Description
		process.Name = p.Name
		process.RequiredGlobalData = p.GlobalData
		if ac, activityRegistered, err := b.BuildActivity(p.Activities); err != nil {
			return nil, err
		} else {
			process.Activities = ac
			if err := b.BuildGatewayBlock(p, process, activityRegistered); err != nil {
				return nil, err
			}
			flow.Process = append(flow.Process, process)
		}
	}
	//need to add gateway

	return flow, nil
}

func (b *BpnmBuilder) AddHandler(nameHandler string, handler ActivityHandler) error {
	if b.handlers == nil {
		b.handlers = make(map[string]ActivityHandler)
	}
	if _, ok := b.handlers[nameHandler]; ok {
		return errors.New("handler's been already defined")
	}
	b.handlers[nameHandler] = handler
	return nil
}

func (b *BpnmBuilder) SetPoolDate(pool StorageData) {
	b.storage = pool
}

type TypeData struct {
	Name string
	Type string
}
type ProcessBuilder struct {
	GlobalData           []TypeData `json:"globalData"`
	Description          string     `json:"description"`
	Name                 string
	Activities           []ActivityBuilder `json:"activities"`
	activitiesRegistered []string
	gatewayRegistered    []string
}

type ActivityBuilder struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Branch      *BranchBuilder `json:"branch"`
}

func (a *BpnmBuilder) BuildActivity(activities []ActivityBuilder) (map[string]*Activity, []string, error) {
	result := make(map[string]*Activity)
	registeredActivities := make([]string, len(activities))
	for _, activityB := range activities {
		activity := new(Activity)
		if handler, ok := a.handlers[activityB.Name]; !ok {
			return nil, nil, fmt.Errorf("no function registered with name: %v", activityB.Name)
		} else {
			activity.Name = activityB.Name
			activity.Description = activityB.Description
			activity.handler = handler
		}
		if g, e := a.BuildBranchBuilder(activityB.Branch); e != nil {
			return nil, nil, e
		} else {
			activity.Branch = g
		}
		result[activity.Name] = activity
		registeredActivities = append(registeredActivities, activity.Name)

	}
	return result, registeredActivities, nil
}

type BranchBuilder struct {
	OutputData  []TypeData     `json:"outputData,omitempty"`
	InputData   []TypeData     `json:"inputData,omitempty"`
	GatewayType string         `json:"gatewayType,omitempty"`
	Gateways    []GatewayBlock `json:"gateways"`
	Recorver    string         `json:"recorver,omitempty"`
}

func (a *BpnmBuilder) BuildBranchBuilder(gatewayDto *BranchBuilder) (*Branch, error) {
	if gatewayDto == nil {
		return nil, errors.New("no branch founded")
	}
	activity := new(Branch)
	activity.GatewayType = GatewayType(gatewayDto.GatewayType)
	activity.OutputData = make(map[string]*DataBpnm)
	activity.InputData = make(map[string]*DataBpnm)
	activity.RequiredInputData = gatewayDto.InputData
	activity.RequiredOutputData = gatewayDto.OutputData
	return activity, nil
}

type GatewayBlock struct {
	NextActivities []string `json:"nextActivities"`
	Decision       string   `json:"decision,omitempty"`
}

func (a *BpnmBuilder) BuildGatewayBlock(pb *ProcessBuilder, p *ProcessBpnm, registeredActivities []string) error {
	for _, act := range pb.Activities {
		var result []*Gateway = make([]*Gateway, 0)
		for _, gat := range act.Branch.Gateways {
			gateway := &Gateway{
				NextActivities: make([]*Activity, 0),
				Decision:       gat.Decision,
			}
			for _, nextJump := range gat.NextActivities {
				if _, ok := p.Activities[nextJump]; !ok {
					return errors.New("no activity registered")
				}
				gateway.NextActivities = append(gateway.NextActivities, p.Activities[nextJump])
			}
			result = append(result, gateway)
		}
		p.Activities[act.Name].Branch.Gateway = result
	}
	return nil
}
