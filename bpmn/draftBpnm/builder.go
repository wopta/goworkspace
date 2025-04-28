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
	var err error
	if err = b.storage.mergeUnique(toMerge.storage); err != nil {
		return err
	}
	toMerge.storage = b.storage
	b.toInject, err = mergeUniqueMaps(b.toInject, toMerge.toInject)
	if err != nil {
		return fmt.Errorf("The merging process of injections went bad: %v", err)
	}
	b.Processes = append(b.Processes, toMerge.Processes...)
	b.handlers, err = mergeUniqueMaps(b.handlers, toMerge.handlers)
	if err != nil {
		return fmt.Errorf("The merging process of handlers went bad: %v", err)
	}
	return nil
}

func (b *BpnmBuilder) Build() (*FlowBpnm, error) {
	flow := new(FlowBpnm)
	flow.process = make(map[string]*processBpnm)

	var newProcess *processBpnm
	if b.storage == nil {
		return nil, errors.New("miss storage")
	}
	var builtActivities map[string]*activity
	var builtEndActivity map[string]*activity
	var err error
	var isInMap bool

	for _, p := range b.Processes {
		builtActivities = map[string]*activity{}
		if flow.process[p.Name] != nil {
			return nil, fmt.Errorf("Process %v's been already defined", newProcess.name)
		}
		newProcess = new(processBpnm)
		newProcess.description = p.Description
		newProcess.storageBpnm = b.storage
		newProcess.name = p.Name
		newProcess.requiredGlobalData = p.GlobalDataRequired
		builtActivities, err = b.buildActivities(p.Name, p.Activities...)
		if err != nil {
			return nil, err
		}
		if builtActivities[getNameEndActivity(p.Name)] != nil {
			return nil, fmt.Errorf("Cant use '%v' as name for an activity", getNameEndActivity(p.Name))
		}
		builtEndActivity, err = b.buildActivities(p.Name, getEndingActivityBuilder(p.Name)) //build the end activity
		if err != nil {
			return nil, err
		}
		builtActivities[getNameEndActivity(p.Name)] = builtEndActivity[getNameEndActivity(p.Name)]

		if _, isInMap = builtActivities[p.DefaultStart]; !isInMap {
			return nil, fmt.Errorf("Process '%v' has no activity named '%v' that can be used as default start", newProcess.name, p.DefaultStart)
		}
		newProcess.activities = builtActivities
		if err = newProcess.hydrateGateways(p.Activities); err != nil {
			return nil, err
		}
		newProcess.defaultStart = p.DefaultStart
		flow.process[newProcess.name] = newProcess
	}

	//Return error if some processes haven't been injected, when a process is injected it's removed from b.toInject
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
		b.handlers = make(map[string]activityHandler)
	}
	if b.toInject == nil {
		b.toInject = make(map[keyInjected]*processBpnm)
	}
	var order *order
	for i, p := range bpnmToInject.Processes { //to have a better error
		order = bpnmToInject.Processes[i].Order
		if order == nil {
			return fmt.Errorf("The 'order' field isn't filled")
		}
		if order.InWhatActivityInjected == "end" {
			order.InWhatActivityInjected = getNameEndActivity(order.InWhatProcessInjected)
		}
		if _, ok := b.toInject[getKeyInjectedProcess(order.InWhatProcessInjected, order.InWhatActivityInjected, order.Order)]; ok {
			return fmt.Errorf("Injection's been already done for: target process: '%v', process: injected '%v' with order '%v'", order.InWhatProcessInjected, p.Name, order.Order)
		}
	}
	process, err := bpnmToInject.Build()
	if err != nil {
		return fmt.Errorf("The building of injected process went bad: %v", err)
	}

	for i, p := range bpnmToInject.Processes {
		order = bpnmToInject.Processes[i].Order
		b.toInject[getKeyInjectedProcess(order.InWhatProcessInjected, order.InWhatActivityInjected, order.Order)] = process.process[p.Name]
	}

	if err = bpnmToInject.storage.setHigherStorage(b.storage); err != nil {
		return err
	}
	return nil
}

func (b *BpnmBuilder) AddHandler(nameHandler string, handler activityHandler) error {
	if b.handlers == nil {
		b.handlers = make(map[string]activityHandler)
	}
	if _, ok := b.handlers[nameHandler]; ok {
		return errors.New("Handler's been already defined")
	}
	b.handlers[nameHandler] = handler
	return nil
}

// only use it for test!!
func (b *BpnmBuilder) setHandler(nameHandler string, handler activityHandler) {
	if handler == nil {
		delete(b.handlers, nameHandler)
		return
	}
	b.handlers[nameHandler] = handler
}

func (b *BpnmBuilder) SetStorage(pool StorageData) {
	b.storage = pool
}

func (a *BpnmBuilder) buildActivities(processName string, activitiesToBuild ...activityBuilder) (map[string]*activity, error) {
	result := make(map[string]*activity)
	for _, activityToBuild := range activitiesToBuild {
		if _, ok := result[activityToBuild.Name]; ok {
			return nil, fmt.Errorf("Double event with same name '%v'", activityToBuild.Name)
		}
		newActivity := new(activity)

		handler, ok := a.handlers[activityToBuild.Name]
		if !activityToBuild.HandlerLess && !ok {
			return nil, fmt.Errorf("No handler registered for the activity: '%v'", activityToBuild.Name)
		}
		if pr := a.toInject[getKeyInjectedProcess(processName, activityToBuild.Name, preActivity)]; pr != nil {
			newActivity.preActivity = pr
			//To check eventually if the some injection isnt possible
			delete(a.toInject, getKeyInjectedProcess(processName, activityToBuild.Name, preActivity))
		}
		if pr := a.toInject[getKeyInjectedProcess(processName, activityToBuild.Name, postActivity)]; pr != nil {
			newActivity.postActivity = pr
			//To check eventually if the some injection isnt possible
			delete(a.toInject, getKeyInjectedProcess(processName, activityToBuild.Name, postActivity))
		}
		newActivity.name = activityToBuild.Name
		newActivity.description = activityToBuild.Description
		newActivity.handler = handler
		if activityToBuild.CallEndIfStop == nil {
			boolPtr := func(b bool) *bool {
				return &b
			}
			activityToBuild.CallEndIfStop = boolPtr(true)
		}
		newActivity.callEndIfStop = *activityToBuild.CallEndIfStop

		if activityToBuild.Recover != "" {
			rec, ok := a.handlers[activityToBuild.Recover]
			if !ok {
				return nil, fmt.Errorf("No handler registered for recovery '%v' in activity: '%v'", activityToBuild.Recover, activityToBuild.Name)
			}
			newActivity.recover = rec
		}
		newActivity.requiredInputData = activityToBuild.InputDataRequired
		newActivity.requiredOutputData = activityToBuild.OutputDataRequired

		result[newActivity.name] = newActivity
	}
	return result, nil
}

// hydrateGateways links each activity's gateways to their corresponding next activities.
// Returns an error if any referenced activity is missing.
func (p *processBpnm) hydrateGateways(activities []activityBuilder) error {
	for _, builderActivity := range activities {
		var gateways []*gateway = make([]*gateway, len(builderActivity.Gateways))
		for igat, builderGateway := range builderActivity.Gateways {
			gateway := &gateway{
				nextActivities: make([]*activity, len(builderGateway.NextActivities)),
				decision:       builderGateway.Decision,
			}
			for iact, nextJump := range builderGateway.NextActivities {
				if _, ok := p.activities[nextJump]; !ok {
					return fmt.Errorf("No event named %v", nextJump)
				}
				gateway.nextActivities[iact] = p.activities[nextJump]
			}
			gateways[igat] = gateway
		}
		p.activities[builderActivity.Name].gateway = gateways
	}
	return nil
}

func getEndingActivityBuilder(nameProcess string) activityBuilder {
	return activityBuilder{
		Name:        getNameEndActivity(nameProcess),
		Description: fmt.Sprint("end activity"),
		HandlerLess: true,
	}
}

func getNameEndActivity(nameProcess string) string {
	return "end_" + nameProcess
}
