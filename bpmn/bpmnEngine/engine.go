package bpmnEngine

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/maja42/goval"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

// Run a process, it starts from defaultActivity defined in json
func (f *FlowBpnm) Run(processName string) error {
	process := f.process[processName]
	if process == nil {
		return fmt.Errorf("Process '%v' not found", processName)
	}
	return f.RunAt(processName, process.defaultStart)
}

// Run a process, it starts from 'startingActivity'
func (f *FlowBpnm) RunAt(processName, startingActivity string) error {
	log.InfoF("Run %v", processName)
	process := f.process[processName]
	if process == nil {
		return fmt.Errorf("Process '%v' not found", processName)
	}

	if e := process.loop(startingActivity); e != nil { //TODO: how to check if there is an infinite loop
		return e
	}
	log.InfoF("Finished %v", processName)
	return nil
}

func (p *processBpnm) loop(nameActivity string) error {
	p.storageBpnm.AddGlobal("statusFlow", &StatusFlow{CurrentProcess: p.name})
	if e := checkGlobalResources(p.storageBpnm, p.requiredGlobalData); e != nil {
		return e
	}

	p.activeActivities = nil
	if act := p.activities[nameActivity]; act != nil {
		p.activeActivities = []*activity{p.activities[nameActivity]}
	}
	if p.storageBpnm == nil {
		return errors.New("Miss storage")
	}
	if p.activeActivities == nil || len(p.activeActivities) == 0 {
		return fmt.Errorf("Process '%v' has no activity '%v'", p.name, nameActivity)
	}
	var err error
	var byte []byte
	var listNewActivities []*activity
	var nextActivities []*activity
	var callEndIfStop bool
	var lastActivity string
	var mapsMerged map[string]any
	for {
		nextActivities = make([]*activity, 0)
		callEndIfStop = true
		for i := range p.activeActivities {
			if err = p.activeActivities[i].runActivity(p.name, p.storageBpnm); err != nil {
				return err
			}
			callEndIfStop = callEndIfStop && p.activeActivities[i].callEndIfStop
			//TODO: to improve
			mapsMerged = mergeMaps(p.storageBpnm.getAllGlobals(), p.storageBpnm.getAllLocals())
			byte, err = json.Marshal(mapsMerged)
			if err != nil {
				return err
			}
			err = json.Unmarshal(byte, &mapsMerged)
			if err != nil {
				return err
			}
			listNewActivities, err = p.activeActivities[i].evaluateDecisions(p.name, p.storageBpnm, mapsMerged)
			lastActivity = p.activeActivities[i].name
			if err != nil {
				return err
			}

			nextActivities = append(nextActivities, listNewActivities...)
		}
		if len(nextActivities) == 0 {
			if !callEndIfStop {
				return nil
			}
			if lastActivity != getNameEndActivity(p.name) {
				return p.activities[getNameEndActivity(p.name)].runActivity(p.name, p.storageBpnm)
			}
			return nil
		}
		p.activeActivities = nextActivities
		if err = p.storageBpnm.cleanNoMarkedResources(); err != nil {
			return err
		}
	}
}

func (act *activity) runActivity(nameProcess string, storage StorageData) (err error) {
	if pre := act.preActivity; pre != nil {
		if e := pre.loop(pre.defaultStart); e != nil {
			return e
		}
	}
	if e := checkLocalResources(storage, act.requiredInputData); e != nil {
		return fmt.Errorf("has an input error: %v", e.Error())
	}

	if err := callWithRecover(nameProcess, storage, act); err != nil {
		return err
	}

	if post := act.postActivity; post != nil {
		if e := post.loop(post.defaultStart); e != nil {
			return e
		}
	}
	return nil
}

func callWithRecover(nameProcess string, storage StorageData, act *activity) (err error) {
	log.InfoF("Run process '%v', start activity '%v'", nameProcess, act.name)
	defer func() {
		if act.recover != nil {
			if r := recover(); r != nil || (err != nil) {
				log.InfoF("Run recorver process '%v', activity '%v'", nameProcess, act.name)
				err = act.recover(storage)
			}
		}
		status := ""
		if err == nil {
			status = "OK"
		} else {
			status = "Fail: " + err.Error()
		}
		log.InfoF("Run process '%v', finished activity '%v' with status: %v", nameProcess, act.name, status)
	}()
	status, err := storage.GetGlobal("statusFlow")
	if err != nil {
		return fmt.Errorf("Error setting status flow of Process '%v' with activity '%v'", nameProcess, act.name)
	}
	status.(*StatusFlow).CurrentActivity = act.name
	if act.handler == nil {
		return nil
	}
	return act.handler(storage)
}

func (act *activity) evaluateDecisions(processName string, storage StorageData, date map[string]any) ([]*activity, error) {
	var res []*activity
	var resultEvaluation any
	var err error
	eval := goval.NewEvaluator()
	if len(act.gateway) == 0 {
		log.InfoF("Process '%v' with activity '%v' has not next activities", processName, act.name)
		return []*activity{}, nil
	}
	for _, ga := range act.gateway {
		if len(ga.nextActivities) == 0 {
			log.InfoF("Process '%v' with activity '%v' has not next activities", processName, act.name)
			return []*activity{}, nil
		}
		if ga.decision == "" {
			if err = checkLocalResources(storage, act.requiredOutputData); err != nil {
				return nil, fmt.Errorf("Process '%v' with activity '%v' has an output error: %v", processName, act.name, err.Error())
			}
			storage.markWhatNeeded(act.requiredOutputData)
			return ga.nextActivities, nil
		}
		resultEvaluation, err = eval.Evaluate(ga.decision, date, nil)
		log.InfoF("Decision evaluation: ( %v )  => %+v", ga.decision, resultEvaluation)
		if err != nil {
			return nil, fmt.Errorf("Process '%v' with activity '%v' has an eval error: %v", processName, act.name, err.Error())
		}
		if ok, isBool := resultEvaluation.(bool); ok && isBool {
			if err = checkLocalResources(storage, act.requiredOutputData); err != nil {
				return nil, fmt.Errorf("Process '%v' with activity '%v' has an output error: %v", processName, act.name, err.Error())
			}
			storage.markWhatNeeded(act.requiredOutputData)
			res = append(res, ga.nextActivities...)
			break
		} else if !isBool {
			return nil, fmt.Errorf("Process '%v' with activity '%v' has an decision error: expected a 'bool' type, got a %v", processName, act.name, reflect.TypeOf(resultEvaluation).String())
		}
	}
	return res, nil
}

func checkLocalResources(st StorageData, req []typeData) error {
	local := st.getAllLocals()
	for _, requiredData := range req {
		storedData, exist := local[requiredData.Name]
		if !exist {
			return fmt.Errorf("Required local resource is not found '%v'", requiredData.Name)
		}
		if storedData.(DataBpnm).GetType() != requiredData.Type {
			return fmt.Errorf("Local resource '%v' has a difference type, exp: '%v', got: '%v'", requiredData.Name, requiredData.Type, storedData.(DataBpnm).GetType())
		}
	}
	return nil
}

func checkGlobalResources(st StorageData, req []typeData) error {
	global := st.getAllGlobals()
	for _, requiredData := range req {
		storedData, exist := global[requiredData.Name]
		if !exist {
			return fmt.Errorf("Required global resource is not found '%v'", requiredData.Name)
		}
		if storedData.(DataBpnm).GetType() != requiredData.Type {
			return fmt.Errorf("Global resource '%v' has a difference type, exp: '%v', got: '%v'", requiredData.Name, requiredData.Type, storedData.(DataBpnm).GetType())
		}
	}
	return nil
}
