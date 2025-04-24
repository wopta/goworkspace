package draftbpnm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/maja42/goval"
)

func (f *FlowBpnm) Run(processName string) error {
	process := f.process[processName]
	if process == nil {
		return fmt.Errorf("Process '%v' not founded", processName)
	}
	return f.RunAt(processName, process.defaultStart)
}

func (f *FlowBpnm) RunAt(processName, activityName string) error {
	log.Println("Run ", processName)
	process := f.process[processName]
	if process == nil {
		return fmt.Errorf("Process '%v' not founded", processName)
	}
	if e := checkGlobalResources(process.storageBpnm, process.requiredGlobalData); e != nil {
		return e
	}

	if e := process.loop(activityName); e != nil { //TODO: how to check if there is an infinite loop
		return e
	}
	log.Println("Stop ", processName)
	return nil
}

func (p *processBpnm) loop(nameActivity string) error {
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
			mapsMerged = mergeMaps(p.storageBpnm.getAllGlobal(), p.storageBpnm.getAllLocal())
			byte, err = json.Marshal(mapsMerged)
			if err != nil {
				return err
			}
			err = json.Unmarshal(byte, &mapsMerged)
			if err != nil {
				return err
			}

			listNewActivities, err = p.activeActivities[i].evaluateDecisions(p.name, p.storageBpnm, mapsMerged)
			if err != nil {
				return err
			}

			nextActivities = append(nextActivities, listNewActivities...)
		}
		if len(nextActivities) == 0 {
			if !callEndIfStop {
				return nil
			}
			return p.activities[getNameEndActivity(p.name)].runActivity(p.name, p.storageBpnm)
		}
		p.activeActivities = nextActivities
		if err = p.storageBpnm.clean(); err != nil {
			return err
		}
	}
}

func (act *activity) runActivity(nameProcess string, storage StorageData) error {
	log.Printf("Run process '%v', activity '%v'", nameProcess, act.name)
	if pre := act.preActivity; pre != nil {
		if err := pre.loop(pre.defaultStart); err != nil {
			return err
		}
	}
	if e := checkLocalResources(storage, act.requiredInputData); e != nil {
		return fmt.Errorf("Process '%v' with activity '%v' has an input error: %v", nameProcess, act.name, e.Error())
	}

	if act.handler != nil {
		var err error
		func() {
			defer func() {
				if act.recover == nil {
					return
				}
				if r := recover(); r != nil || err != nil {
					log.Printf("Run recorver process '%v', activity '%v'", nameProcess, act.name)
					err = act.recover(storage)
				}
				err = nil
			}()
			err = act.handler(storage)
		}()
		if err != nil {
			return err
		}
	}

	if post := act.postActivity; post != nil {
		if e := post.loop(post.defaultStart); e != nil {
			return e
		}
	}
	return nil
}

func (act *activity) evaluateDecisions(processName string, storage StorageData, date map[string]any) ([]*activity, error) {
	var res []*activity
	var resultEvaluation any
	var err error
	eval := goval.NewEvaluator()
	for _, ga := range act.gateway {
		if ga.decision == "" {
			if err = checkLocalResources(storage, act.requiredOutputData); err != nil {
				return nil, fmt.Errorf("Process '%v' with activity '%v' has an output error: %v", processName, act.name, err.Error())
			}
			storage.markWhatNeeded(act.requiredOutputData)
			return ga.nextActivities, nil
		}
		if len(ga.nextActivities) == 0 {
			log.Printf("Process '%v' has not activities", processName)
			return []*activity{}, nil
		}
		resultEvaluation, err = eval.Evaluate(ga.decision, date, nil)
		if err != nil {
			return nil, fmt.Errorf("Process '%v' with activity '%v' has an eval error: %v", processName, act.name, err.Error())
		}
		if resultEvaluation.(bool) {
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
	local := st.getAllLocal()
	for _, d := range req {
		v, ok := local[d.Name]
		if !ok {
			return fmt.Errorf("Required local resource is not found '%v'", d.Name)
		}
		if v.(DataBpnm).GetType() != d.Type {
			return fmt.Errorf("Local resource '%v' has a difference type, exp: '%v', got: '%v'", d.Name, d.Type, v.(DataBpnm).GetType())
		}
	}
	return nil
}

func checkGlobalResources(st StorageData, req []typeData) error {
	global := st.getAllGlobal()
	for _, d := range req {
		v, ok := global[d.Name]
		if !ok {
			return fmt.Errorf("Required global resource is not found '%v'", d.Name)
		}
		if v.(DataBpnm).GetType() != d.Type {
			return fmt.Errorf("Global resource '%v' has a difference type, exp: '%v', got: '%v'", d.Name, d.Type, v.(DataBpnm).GetType())
		}
	}
	return nil
}
