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

func (f *FlowBpnm) RunAt(processName, startingActivity string) error {
	log.InfoF("Run %v", processName)
	process := f.process[processName]
	if process == nil {
		return fmt.Errorf("Process '%v' not found", processName)
	}

	process.storageBpnm.AddGlobal("statusFlow", &StatusFlow{CurrentProcess: process.name})
	if e := checkGlobalResources(process.storageBpnm, process.requiredGlobalData); e != nil {
		return e
	}

	var firstActivities []*activity
	if act := process.activities[startingActivity]; act != nil {
		firstActivities = []*activity{process.activities[startingActivity]}
	}
	if process.storageBpnm == nil {
		return errors.New("Miss storage")
	}
	if firstActivities == nil || len(firstActivities) == 0 {
		return fmt.Errorf("Process '%v' has no activity '%v'", process.name, startingActivity)
	}
	var err error

	if err = process.loop(process.storageBpnm, firstActivities...); err != nil {
		return err
	}
	return nil
}

func (p *processBpnm) loop(initialStorage StorageData, activities ...*activity) (err error) {
	var callEndIfStop bool = true
	for i := range activities {
		newStorage := NewStorageBpnm()
		err := newStorage.setHigherStorage(initialStorage)
		initialStorage = newStorage
		if err != nil {
			return err
		}

		if err = activities[i].runActivity(p.name, newStorage); err != nil {
			return err
		}
		callEndIfStop = callEndIfStop && activities[i].callEndIfStop
		//TODO: to improve
		mapsMerged := mergeMaps(newStorage.getAllGlobals(), newStorage.getAllLocals())
		byte, err := json.Marshal(mapsMerged)
		if err != nil {
			return err
		}
		err = json.Unmarshal(byte, &mapsMerged)
		if err != nil {
			return err
		}
		listNewActivities, err := activities[i].evaluateDecisions(p.name, newStorage, mapsMerged)
		lastActivity := activities[i].name
		if err != nil {
			return err
		}
		if len(listNewActivities) == 0 {
			if !callEndIfStop {
				log.InfoF("Finished %v", p.name)
				continue
			}
			if lastActivity != getNameEndActivity(p.name) {
				if err = p.activities[getNameEndActivity(p.name)].runActivity(p.name, newStorage); err != nil {
					return err
				}
			}
			log.InfoF("Finished %v", p.name)
			continue
		}
		newStorage.cleanNoMarkedResources()
		if err = p.loop(newStorage, listNewActivities...); err != nil {
			return err
		}
	}
	return err
}

func (act *activity) runActivity(nameProcess string, storage StorageData) (err error) {
	if pre := act.preActivity; pre != nil {
		pre.storageBpnm.setHigherStorage(storage)
		if err := pre.loop(pre.storageBpnm, pre.activities[pre.defaultStart]); err != nil {
			return err
		}
	}
	if e := checkLocalResources(storage, act.requiredInputData); e != nil {
		return fmt.Errorf("has an input error: %v", e.Error())
	}

	if err := callWithRecover(nameProcess, storage, act); err != nil {
		return err
	}

	if post := act.postActivity; post != nil {
		post.storageBpnm.setHigherStorage(storage)
		if err := post.loop(post.storageBpnm, post.activities[post.defaultStart]); err != nil {
			return err
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
	log.InfoF("Decision evaluation from '%v' ...", act.name)
	if len(act.gateway) == 0 {
		log.InfoF("Process '%v' with activity '%v' has not next activities", processName, act.name)
		return []*activity{}, nil
	}
	for _, ga := range act.gateway {
		//DEFAULT GATEWAY
		if len(ga.nextActivities) == 0 {
			log.InfoF("Process '%v' with activity '%v' has not next activities", processName, act.name)
			return []*activity{}, nil
		}
		if ga.decision == "" {
			if err = checkLocalResources(storage, act.requiredOutputData); err != nil {
				return nil, fmt.Errorf("Process '%v' with activity '%v' has an output error: %v", processName, act.name, err.Error())
			}
			log.InfoF("Decision evaluation: None")
			storage.markWhatNeeded(act.requiredOutputData)
			return ga.nextActivities, nil
		}

		//EVALUATE DECISION
		resultEvaluation, err = eval.Evaluate(ga.decision, date, nil)
		log.InfoF("Decision evaluation: ( %v )  => %+v", ga.decision, resultEvaluation)
		if err != nil {
			return nil, fmt.Errorf("Process '%v' with activity '%v' has an eval error: %v", processName, act.name, err.Error())
		}
		result, isBool := resultEvaluation.(bool)
		if !isBool {
			return nil, fmt.Errorf("Process '%v' with activity '%v' has an decision error: expected a 'bool' type, got a %v", processName, act.name, reflect.TypeOf(resultEvaluation).String())
		}
		if result {
			if err = checkLocalResources(storage, act.requiredOutputData); err != nil {
				return nil, fmt.Errorf("Process '%v' with activity '%v' has an output error: %v", processName, act.name, err.Error())
			}
			storage.markWhatNeeded(act.requiredOutputData)
			res = append(res, ga.nextActivities...)
			break
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
