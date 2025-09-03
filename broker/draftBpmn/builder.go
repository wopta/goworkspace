package draftbpmn

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib"
)

func getKeyInjectProcess(targetPro, targetAct string, order activityOrder) injectionKey {
	order = activityOrder(strings.ToLower(string(order)))
	return injectionKey{
		targetProcess:  targetPro,
		targetActivity: targetAct,
		activityOrder:  order,
	}
}

func NewBpnmBuilderRawPath(path string) (*BpnmBuilder, error) {
	var Bpnm BpnmBuilder
	jsonProva, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(jsonProva, &Bpnm)
	return &Bpnm, nil
}
func NewBpnmBuilder(path string) (*BpnmBuilder, error) {
	var Bpnm BpnmBuilder
	jsonProva, err := lib.GetFilesByEnvV2(path)
	if err != nil {
		return nil, err
	}
	if len(jsonProva) == 0 {
		return nil, errors.New("Json not found: " + path)
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
			return nil, fmt.Errorf("Cant use '%v' as name for an activity, since it's a builtin activity", getNameEndActivity(p.Name))
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
			return nil, fmt.Errorf("Process '%v' with %v", p.Name, err)
		}
		newProcess.defaultStart = p.DefaultStart
		flow.process[newProcess.name] = newProcess
	}

	//Return error if some processes haven't been injected, when a process is injected it's removed from b.toInject
	if len(b.toInject) != 0 {
		var keyNoInject string
		for i := range b.toInject {
			keyNoInject += fmt.Sprintf("process: %v, activity: %v, order: %v\n", i.targetProcess, i.targetActivity, i.activityOrder)
		}
		return nil, fmt.Errorf("Following injections haven't been done:\n%v", keyNoInject)
	}
	return flow, nil
}

// Inject a processes that will be called before or after activity's handler, it depends on the configuration Order
func (b *BpnmBuilder) Inject(bpnmToInject *BpnmBuilder) error {
	if b.storage == nil {
		return errors.New("No storage defined")
	}
	if b.handlers == nil {
		b.handlers = make(map[string]activityHandler)
	}
	if b.toInject == nil {
		b.toInject = make(map[injectionKey]*processBpnm)
	}
	var order *order
	for i, p := range bpnmToInject.Processes { //to have a better error
		order = bpnmToInject.Processes[i].Order
		if order == nil {
			return fmt.Errorf("The 'order' field isn't filled")
		}
		if order.InWhatActivityInject == "end" {
			order.InWhatActivityInject = getNameEndActivity(order.InWhatProcessInject)
		}
		if _, exist := b.toInject[getKeyInjectProcess(order.InWhatProcessInject, order.InWhatActivityInject, order.Order)]; exist {
			return fmt.Errorf("Injection's been already done for: target process: '%v', process: injected '%v' with order '%v'", order.InWhatProcessInject, p.Name, order.Order)
		}
	}
	process, err := bpnmToInject.Build()
	if err != nil {
		return fmt.Errorf("The building of injected process went bad: %v", err)
	}

	for i, p := range bpnmToInject.Processes {
		order = bpnmToInject.Processes[i].Order
		b.toInject[getKeyInjectProcess(order.InWhatProcessInject, order.InWhatActivityInject, order.Order)] = process.process[p.Name]
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
	if _, exist := b.handlers[nameHandler]; exist {
		return fmt.Errorf("Handler %v has been already defined", nameHandler)
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

// Build a list of activity, association the handlers and injected processes to each activities
// The matching between gateways and activities has to been done yet, with 'hydrateGateways'
func (a *BpnmBuilder) buildActivities(processName string, activitiesToBuild ...activityBuilder) (map[string]*activity, error) {
	result := make(map[string]*activity)
	for _, activityToBuild := range activitiesToBuild {
		if _, exist := result[activityToBuild.Name]; exist {
			return nil, fmt.Errorf("Double event with same name '%v'", activityToBuild.Name)
		}
		newActivity := new(activity)

		handler, exist := a.handlers[activityToBuild.Name]
		if !activityToBuild.HandlerLess && !exist {
			return nil, fmt.Errorf("No handler registered for the activity: '%v'", activityToBuild.Name)
		}
		if pr := a.toInject[getKeyInjectProcess(processName, activityToBuild.Name, preActivity)]; pr != nil {
			newActivity.preActivity = pr
			//To check eventually if the some injection isnt possible
			delete(a.toInject, getKeyInjectProcess(processName, activityToBuild.Name, preActivity))
		}
		if pr := a.toInject[getKeyInjectProcess(processName, activityToBuild.Name, postActivity)]; pr != nil {
			newActivity.postActivity = pr
			//To check eventually if the some injection isnt possible
			delete(a.toInject, getKeyInjectProcess(processName, activityToBuild.Name, postActivity))
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
		//Check if a recover handler has been specified
		if activityToBuild.Recover != "" {
			rec, exist := a.handlers[activityToBuild.Recover]
			if !exist {
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
				if nextJump == "end" {
					nextJump = getNameEndActivity(p.name)
				}
				if _, exist := p.activities[nextJump]; !exist {
					return fmt.Errorf("No event named %v", nextJump)
				}
				gateway.nextActivities[iact] = p.activities[nextJump]
				if e := isInputProvidedByOutput(gateway.nextActivities[iact].requiredInputData, builderActivity.OutputDataRequired); e != nil {
					prefix := fmt.Sprintf("input activity: '%v' and output activity: '%v'", gateway.nextActivities[iact].name, builderActivity.Name)
					return fmt.Errorf(prefix+", has error: %v", e.Error())
				}
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

// isInputProvidedByOutput: check if every inputs is correctly provided by outputs, otherwise return error
func isInputProvidedByOutput(inputs []typeData, outputs []typeData) error {
	//check if input is equal to output
	checkData := func(input typeData, output typeData) (isFounded bool, err error) {
		if input.Name == output.Name {
			if input.Type == output.Type {
				return true, nil
			}
			isFounded = true
			err = fmt.Errorf("The type of output data '%v' differ from the input one, '%v'!='%v'", output.Name, output.Type, input.Type)
			return
		}
		return false, nil
	}
	if len(inputs) > len(outputs) {
		return errors.New("Insufficient number of output data")
	}
	for _, input := range inputs {
		err := fmt.Errorf("The input %v isn't provided by output", input)
		for _, output := range outputs {
			isFounded, errComparison := checkData(input, output)
			if !isFounded {
				continue
			}
			err = nil
			if errComparison != nil {
				return errComparison
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
