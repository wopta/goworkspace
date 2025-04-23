package draftbpnm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

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

func (b *BpnmBuilder) AddProcesses(toMerge *BpnmBuilder) error {
	var e error
	if e = b.storage.mergeUnique(toMerge.storage); e != nil {
		return e
	}
	toMerge.storage = b.storage
	b.toInject, e = mergeUniqueMaps(b.toInject, toMerge.toInject)
	if e != nil {
		return fmt.Errorf("The merging process of injections went bad: %v", e)
	}
	b.Processes = append(b.Processes, toMerge.Processes...)
	b.handlers, e = mergeUniqueMaps(b.handlers, toMerge.handlers)
	if e != nil {
		return fmt.Errorf("The merging process of handlers went bad: %v", e)
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
		p.Activities = append(p.Activities, buildEndingActivity(p.Name))
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

// Inject a processes that will be called before or after activities, it depends on the configuration Order
func (b *BpnmBuilder) Inject(bpnmToInject *BpnmBuilder) error {
	if b.storage == nil {
		return errors.New("No storage defined")
	}
	if b.handlers == nil {
		b.handlers = make(map[string]ActivityHandler)
	}
	if b.toInject == nil {
		b.toInject = make(map[KeyInject]*ProcessBpnm)
	}
	var order Order
	for i, p := range bpnmToInject.Processes { //to have a better error
		order = bpnmToInject.Processes[i].Order
		if order.InWhatActivityBeInjected == "end" {
			bpnmToInject.Processes[i].Order.InWhatActivityBeInjected = fmt.Sprintf("%v_end", order.InWhatProcessBeInjected)
		}
		if _, ok := b.toInject[getKeyInjectedProcess(order.InWhatProcessBeInjected, order.InWhatActivityBeInjected, order.Order)]; ok {
			return fmt.Errorf("Injection's been already done: target process: '%v', process: injected '%v'", order.InWhatProcessBeInjected, p.Name)
		}
	}
	process, err := bpnmToInject.Build()
	if err != nil {
		return err
	}

	for i, p := range bpnmToInject.Processes {
		order = bpnmToInject.Processes[i].Order
		b.toInject[getKeyInjectedProcess(order.InWhatProcessBeInjected, order.InWhatActivityBeInjected, order.Order)] = process.Process[p.Name]
	}

	if err = bpnmToInject.storage.setHigherStorage(b.storage); err != nil {
		return err
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
			return nil, fmt.Errorf("Double event with same name '%v'", activity.Name)
		}
		newActivity := new(Activity)

		handler, ok := a.handlers[activity.Name]
		if !activity.HandlerLess && !ok {
			return nil, fmt.Errorf("No handler registered for the activity: '%v'", activity.Name)
		}
		if pr := a.toInject[getKeyInjectedProcess(processName, activity.Name, PreActivity)]; pr != nil {
			newActivity.PreActivity = pr
		}
		if pr := a.toInject[getKeyInjectedProcess(processName, activity.Name, PostActivity)]; pr != nil {
			newActivity.PostActivity = pr
		}
		//To check eventually if the some injection isnt possible
		delete(a.toInject, getKeyInjectedProcess(processName, activity.Name, PreActivity))
		delete(a.toInject, getKeyInjectedProcess(processName, activity.Name, PostActivity))

		newActivity.Name = activity.Name
		newActivity.Description = activity.Description
		newActivity.handler = handler

		if activity.Recover != "" {
			rec, ok := a.handlers[activity.Recover]
			if !ok {
				return nil, fmt.Errorf("No handler registered for recovery '%v' in activity: '%v'", activity.Recover, activity.Name)
			}
			newActivity.recover = rec
		}

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
	//	activity.GatewayType = b.GatewayType
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

func buildEndingActivity(processName string) ActivityBuilder {
	return ActivityBuilder{
		Name:        fmt.Sprintf("%v_end", processName),
		Description: fmt.Sprint("end activity for ", processName),
		Branch:      nil,
		HandlerLess: true,
	}
}
