package draftbpnm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
	HandlerLess bool           `json:"handlerless"`
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

func NewBpnmBuilder(path string) (*BpnmBuilder, error) {
	var Bpnm BpnmBuilder
	jsonProva, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(jsonProva, &Bpnm)
	return &Bpnm, nil
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
		builtActivities, err := b.buildActivities(p.Activities, p.Name)
		if err != nil {
			return nil, err
		}
		process.Activities = builtActivities
		if err := process.hydrateGateways(p.Activities); err != nil {
			return nil, err
		}
		if flow.Process[process.Name] != nil {
			return nil, fmt.Errorf("Process %v's been already defined", process.Name)
		}

		if _, ok := builtActivities[p.DefaultStart]; !ok {
			return nil, fmt.Errorf("Process '%v' has no activity named '%v' that can be used as default start", process.Name, p.DefaultStart)
		}
		process.DefaultStart = p.DefaultStart
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

func (b *BpnmBuilder) SetStorage(pool StorageData) {
	b.storage = pool
}

func (a *BpnmBuilder) buildActivities(activities []ActivityBuilder, processName string) (map[string]*Activity, error) {
	result := make(map[string]*Activity)
	for _, activity := range activities {
		if _, ok := result[activity.Name]; ok {
			return nil, fmt.Errorf("Double event with same name %v", activity.Name)
		}
		newActivity := new(Activity)

		handler, ok := a.handlers[activity.Name]
		if !activity.HandlerLess && !ok {
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
		if activity.Branch != nil {
			builtBranch, e := activity.Branch.buildBranch()
			if e != nil {
				return nil, e
			}
			newActivity.Branch = builtBranch
		}
		result[newActivity.Name] = newActivity
	}
	return result, nil
}

func (b *BranchBuilder) buildBranch() (*Branch, error) {
	if b == nil {
		return nil, nil
	}
	activity := new(Branch)
	activity.GatewayType = b.GatewayType
	activity.RequiredInputData = b.InputDataRequired
	activity.RequiredOutputData = b.OutputDataRequired
	return activity, nil
}

func (p *ProcessBpnm) hydrateGateways(activities []ActivityBuilder) error {
	for _, builderActivity := range activities {
		if builderActivity.Branch == nil {
			continue
		}
		var gateways []*Gateway = make([]*Gateway, 0)
		for _, builderGateway := range builderActivity.Branch.Gateways {
			gateway := &Gateway{
				NextActivities: make([]*Activity, 0),
				Decision:       builderGateway.Decision,
			}
			for _, nextJump := range builderGateway.NextActivities {
				if _, ok := p.Activities[nextJump]; !ok {
					return fmt.Errorf("No event named %v", nextJump)
				}
				gateway.NextActivities = append(gateway.NextActivities, p.Activities[nextJump])
			}
			gateways = append(gateways, gateway)
		}
		p.Activities[builderActivity.Name].Branch.Gateway = gateways
	}
	return nil
}
