package draftbpnm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func getKeyInjectedProcess(targetPro, targetAct string, order orderActivity) keyInjected {
	order = orderActivity(strings.ToLower(string(order)))
	return keyInjected{
		targetProcess:  targetPro,
		targetActivity: targetAct,
		orderActivity:  order,
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

	var newProcess *ProcessBpnm
	if b.storage == nil {
		return nil, errors.New("miss storage")
	}
	for _, p := range b.Processes {
		if flow.Process[p.Name] != nil {
			return nil, fmt.Errorf("Process %v's been already defined", newProcess.Name)
		}
		newProcess = new(ProcessBpnm)
		newProcess.Description = p.Description
		newProcess.storageBpnm = b.storage
		newProcess.Name = p.Name
		newProcess.RequiredGlobalData = p.GlobalDataRequired
		builtActivities, err := b.buildActivities(p.Name, p.Activities...)
		if err != nil {
			return nil, err
		}
		if builtActivities["end"] != nil {
			return nil, errors.New("Cant use 'end' as name for an activity")
		}
		builtEndActivity, err := b.buildActivities(p.Name, getEndingActivityBuilder()) //build the end activity
		if err != nil {
			return nil, err
		}
		builtActivities["end"] = builtEndActivity["end"]

		if _, ok := builtActivities[p.DefaultStart]; !ok {
			return nil, fmt.Errorf("Process '%v' has no activity named '%v' that can be used as default start", newProcess.Name, p.DefaultStart)
		}
		newProcess.Activities = builtActivities
		if err := newProcess.hydrateGateways(p.Activities); err != nil {
			return nil, err
		}
		newProcess.DefaultStart = p.DefaultStart
		flow.Process[newProcess.Name] = newProcess
	}

	//Return error if some processes isnt injected, when a process is injected it's removed from b.toInject
	if len(b.toInject) != 0 {
		var keyNoInjected string
		for i := range b.toInject {
			keyNoInjected += fmt.Sprintf("process: %v, activity: %v, order: %v\n", i.targetProcess, i.targetActivity, i.orderActivity)
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
		b.toInject = make(map[keyInjected]*ProcessBpnm)
	}
	var order *Order
	for i, p := range bpnmToInject.Processes { //to have a better error
		order = bpnmToInject.Processes[i].Order
		if order == nil {
			return fmt.Errorf("No order defined, the 'order' field isnt filled")
		}
		if order.InWhatActivityInjected == "end" {
			order.InWhatActivityInjected = "end"
		}
		if _, ok := b.toInject[getKeyInjectedProcess(order.InWhatProcessInjected, order.InWhatActivityInjected, order.Order)]; ok {
			return fmt.Errorf("Injection's been already done: target process: '%v', process: injected '%v'", order.InWhatProcessInjected, p.Name)
		}
	}
	process, err := bpnmToInject.Build()
	if err != nil {
		return err
	}

	for i, p := range bpnmToInject.Processes {
		order = bpnmToInject.Processes[i].Order
		b.toInject[getKeyInjectedProcess(order.InWhatProcessInjected, order.InWhatActivityInjected, order.Order)] = process.Process[p.Name]
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

// only use it for test!!
func (b *BpnmBuilder) setHandler(nameHandler string, handler ActivityHandler) error {
	if b.handlers == nil {
		return errors.New("No handlers has been defined")
	}
	if _, ok := b.handlers[nameHandler]; !ok {
		return errors.New("Handler isn't defined")
	}
	b.handlers[nameHandler] = handler
	return nil
}

func (b *BpnmBuilder) SetStorage(pool StorageData) {
	b.storage = pool
}

func (a *BpnmBuilder) buildActivities(processName string, activities ...activityBuilder) (map[string]*Activity, error) {
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
		if pr := a.toInject[getKeyInjectedProcess(processName, activity.Name, preActivity)]; pr != nil {
			newActivity.PreActivity = pr
			//To check eventually if the some injection isnt possible
			delete(a.toInject, getKeyInjectedProcess(processName, activity.Name, preActivity))
		}
		if pr := a.toInject[getKeyInjectedProcess(processName, activity.Name, postActivity)]; pr != nil {
			newActivity.PostActivity = pr
			//To check eventually if the some injection isnt possible
			delete(a.toInject, getKeyInjectedProcess(processName, activity.Name, postActivity))
		}
		newActivity.Name = activity.Name
		newActivity.Description = activity.Description
		newActivity.handler = handler
		if activity.CallEndIfStop == nil {
			boolPtr := func(b bool) *bool {
				return &b
			}
			activity.CallEndIfStop = boolPtr(true)
		}
		newActivity.CallEndIfStop = *activity.CallEndIfStop

		if activity.Recover != "" {
			rec, ok := a.handlers[activity.Recover]
			if !ok {
				return nil, fmt.Errorf("No handler registered for recovery '%v' in activity: '%v'", activity.Recover, activity.Name)
			}
			newActivity.recover = rec
		}
		newActivity.RequiredInputData = activity.InputDataRequired
		newActivity.RequiredOutputData = activity.OutputDataRequired

		result[newActivity.Name] = newActivity
	}
	return result, nil
}

// hydrateGateways links each activity's gateways to their corresponding next activities.
// Returns an error if any referenced activity is missing.
func (p *ProcessBpnm) hydrateGateways(activities []activityBuilder) error {
	for _, builderActivity := range activities {
		var gateways []*Gateway = make([]*Gateway, 0)
		for _, builderGateway := range builderActivity.Gateways {
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
		p.Activities[builderActivity.Name].Gateway = gateways
	}
	return nil
}

func getEndingActivityBuilder() activityBuilder {
	return activityBuilder{
		Name:        "end",
		Description: fmt.Sprint("end activity"),
		HandlerLess: true,
	}
}
