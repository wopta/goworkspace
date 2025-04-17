package draftbpnm

import (
	"errors"
	"fmt"
	"strings"
)

type Order struct {
	InWhatProcessBeInjected  string        `json:"inWhatProcessBeInjected"`
	InWhatActivityBeInjected string        `json:"inWhatActivityBeInjected"`
	Order                    OrderActivity `json:"order"`
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
	Order              Order             `json:"order"`
	Description        string            `json:"description"`
	Name               string            `json:"name"`
	Activities         []ActivityBuilder `json:"activities"`
	DefaultStart       string            `json:"defaultStart"`
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

func getKeyInjectedProcess(targetPro, targetAct string, order OrderActivity) KeyInject {
	order = OrderActivity(strings.ToLower(string(order)))
	return KeyInject{
		TargetProcess:  targetPro,
		TargetActivity: targetAct,
		OrderActivity:  order,
	}
}

func (b *BpnmBuilder) Inject(bpnmToInject *BpnmBuilder) error {
	if b.handlers == nil {
		b.handlers = make(map[string]ActivityHandler)
	}
	if b.toInject == nil {
		b.toInject = make(map[KeyInject]*ProcessBpnm)
	}
	process, err := bpnmToInject.Build()
	if err != nil {
		return err
	}
	var order Order
	for i, p := range bpnmToInject.Processes {
		order = bpnmToInject.Processes[i].Order
		if _, ok := b.toInject[getKeyInjectedProcess(order.InWhatProcessBeInjected, order.InWhatActivityBeInjected, order.Order)]; ok {
			return fmt.Errorf("Injection's been already done: target process: %v, process: injected %v", order.InWhatProcessBeInjected, p.Name)
		}
		b.toInject[getKeyInjectedProcess(order.InWhatProcessBeInjected, order.InWhatActivityBeInjected, order.Order)] = process.Process[p.Name]
	}
	return nil
}

func (b *BpnmBuilder) Build() (*FlowBpnm, error) {
	flow := new(FlowBpnm)
	flow.Process = make(map[string]*ProcessBpnm)

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
		process.DefaultStart = p.DefaultStart
		if err != nil {
			return nil, err
		}
		process.Activities = builtActivities
		if err := b.BuildGatewayBlock(p, process); err != nil {
			return nil, err
		}
		if flow.Process[process.Name] != nil {
			return nil, fmt.Errorf("Process %v's been already defined", process.Name)
		}

		flow.Process[process.Name] = process
	}
	//Return error if some processes isnt injected, when a process is injected it's removed from b.toInject
	if len(b.toInject) != 0 {
		var keyNoInjected string
		for i := range b.toInject {
			keyNoInjected += fmt.Sprintf("process: %v, activity: %v, order: %v\n", i.TargetProcess, i.TargetActivity, i.OrderActivity)
		}
		return nil, fmt.Errorf("Following injections went bad:\n  %v", keyNoInjected)
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
		newActivity.PreActivity = a.toInject[getKeyInjectedProcess(processName, activity.Name, PreActivity)]
		newActivity.PostActivity = a.toInject[getKeyInjectedProcess(processName, activity.Name, PostActivity)]
		//To check eventually if the some injection isnt possible
		delete(a.toInject, getKeyInjectedProcess(processName, activity.Name, PreActivity))
		delete(a.toInject, getKeyInjectedProcess(processName, activity.Name, PostActivity))

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
