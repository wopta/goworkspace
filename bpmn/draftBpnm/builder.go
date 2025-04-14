package draftbpnm

import (
	"errors"
	"fmt"
)

type BpnmBuilder struct {
	handlers  map[string]ActivityHandler
	storage   StorageData
	Processes []*ProcessBuilder `json:"processes"`
	toInject  map[string]*ProcessBpnm
}

type TypeData struct {
	Name string
	Type string
}
type ProcessBuilder struct {
	GlobalDataRequired []TypeData        `json:"globalData"`
	Description        string            `json:"description"`
	Name               string            `json:"name"`
	Activities         []ActivityBuilder `json:"activities"`
}

type ActivityBuilder struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Branch      *BranchBuilder `json:"branch"`
}

type BranchBuilder struct {
	OutputDataRequired []TypeData     `json:"outputData,omitempty"`
	InputDataRequired  []TypeData     `json:"inputData,omitempty"`
	GatewayType        GatewayType    `json:"gatewayType,omitempty"`
	Gateways           []GatewayBlock `json:"gateways"`
	Recorver           string         `json:"recorver,omitempty"`
}

type GatewayBlock struct {
	NextActivities []string `json:"nextActivities"`
	Decision       string   `json:"decision,omitempty"`
}

type OrderActivity string

const (
	PreActivity  OrderActivity = "Pre"
	PostActivity OrderActivity = "Post"
)

type InjectActivity struct {
	Order              OrderActivity
	ActivitiesToInject []ActivityBuilder
}

func getNameInjectedProcess(targetPro, targetAct string, order OrderActivity) string {
	return targetPro + "_" + targetAct + "_" + string(order)
}

// injectedProcess: what process contains the activities to inject
// targetProcess: what process'll receive the new activities
// order: define when execute the activities, before or after targetActivity
func (b *BpnmBuilder) Inject(targetPro, targetAct, processToTake string, order OrderActivity, bpnmToInject *BpnmBuilder) error {
	if b.handlers == nil {
		b.handlers = make(map[string]ActivityHandler)
	}
	if b.toInject == nil {
		b.toInject = make(map[string]*ProcessBpnm)
	}
	process, err := bpnmToInject.Build()
	if err != nil {
		return fmt.Errorf("Injected process: target process %v, target activity %v, error: %v", targetPro, targetAct, err.Error())
	}
	for _, p := range process.Process {
		if p.Name == processToTake {
			b.toInject[getNameInjectedProcess(targetPro, targetAct, order)] = p
		}
	}
	return nil
}

func (b *BpnmBuilder) Build() (*FlowBpnm, error) {
	flow := new(FlowBpnm)
	var process *ProcessBpnm
	if b.storage == nil {
		return nil, errors.New("miss storage")
	}
	for _, p := range b.Processes {
		process = new(ProcessBpnm)
		process.Description = p.Description
		process.storageBpnm = b.storage
		process.Name = p.Name
		process.RequiredGlobalData = p.GlobalDataRequired
		builtActivities, err := b.BuildActivity(p.Activities, p.Name)
		if err != nil {
			return nil, err
		}
		process.Activities = builtActivities
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

func (a *BpnmBuilder) BuildActivity(activities []ActivityBuilder, processName string) (map[string]*Activity, error) {
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
		builtBranch, e := a.BuildBranchBuilder(activity.Branch)
		if e != nil {
			return nil, e
		}
		newActivity.Branch = builtBranch
		newActivity.PreActivity = a.toInject[getNameInjectedProcess(processName, activity.Name, PreActivity)]
		newActivity.PostActivity = a.toInject[getNameInjectedProcess(processName, activity.Name, PostActivity)]
		result[newActivity.Name] = newActivity
	}
	return result, nil
}

func (a *BpnmBuilder) BuildBranchBuilder(gatewayDto *BranchBuilder) (*Branch, error) {
	activity := new(Branch)
	activity.GatewayType = gatewayDto.GatewayType
	activity.RequiredInputData = gatewayDto.InputDataRequired
	activity.RequiredOutputData = gatewayDto.OutputDataRequired
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
