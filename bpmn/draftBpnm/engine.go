package draftbpnm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/maja42/goval"
	"log"
)

func (f *FlowBpnm) Run(processName string) error {
	process := f.Process[processName]
	if process == nil {
		return fmt.Errorf("Process '%v' not founded", processName)
	}
	return f.RunAt(processName, process.DefaultStart)
}

func (f *FlowBpnm) RunAt(processName, activityName string) error {
	log.Println("Run ", processName)
	process := f.Process[processName]
	if process == nil {
		return fmt.Errorf("Process '%v' not founded", processName)
	}
	if e := checkValidityGlobalStorage(process.storageBpnm, process.RequiredGlobalData); e != nil {
		return e
	}

	if e := process.run(activityName); e != nil { //TODO: how to check if there is an infinite loop
		return e
	}
	log.Println("Stop ", processName)
	return nil
}

func (p *ProcessBpnm) run(nameActivity string) error {
	p.activeActivities = append(p.activeActivities, p.Activities[nameActivity])
	if p.storageBpnm == nil {
		return errors.New("miss storage")
	}
	if p.activeActivities == nil {
		return fmt.Errorf("Process '%v' has no activity '%v'", p.Name, nameActivity)
	}
	//TODO: implement at garbange collector with a counter, to remove old element not used in the next branch
	for {
		var nextActivities []*Activity
		for i := range p.activeActivities {
			if err := p.activeActivities[i].runActivity(p.Name, p.storageBpnm); err != nil {
				return err
			}
			//TODO: to improve

			m := mergeMaps(p.storageBpnm.GetAllGlobal(), p.storageBpnm.GetAllLocal())
			jsonMap := make(map[string]any)
			b, _ := json.Marshal(m)
			_ = json.Unmarshal(b, &jsonMap)

			list, e := p.activeActivities[i].evaluateDecisions(p.Name, p.storageBpnm, jsonMap)
			if e != nil {
				return e
			}
			nextActivities = append(nextActivities, list...)
		}
		if len(nextActivities) == 0 {
			return nil
		}
		p.activeActivities = nextActivities
		p.storageBpnm.clean()
	}
}

func (act *Activity) runActivity(nameProcess string, storage StorageData) error {
	log.Printf("Run process '%v', activity '%v'", nameProcess, act.Name)
	if pre := act.PreActivity; pre != nil {
		if err := pre.storageBpnm.Merge(storage); err != nil {
			return err
		}
		if err := pre.run(pre.DefaultStart); err != nil {
			return err
		}
	}

	if act.Branch != nil {
		if e := checkLocalStorage(storage, act.Branch.RequiredInputData); e != nil {
			return fmt.Errorf("Process '%v' with activity '%v' has an input error: %v", nameProcess, act.Name, e.Error())
		}
	}

	if act.handler != nil {
		if e := act.handler(storage); e != nil {
			return e
		}
	}

	if post := act.PostActivity; post != nil {
		post.storageBpnm.Merge(storage)
		if e := post.run(post.DefaultStart); e != nil {
			return e
		}
	}
	return nil
}

func (act *Activity) evaluateDecisions(processName string, storage StorageData, date map[string]any) ([]*Activity, error) {
	var res []*Activity
	if act.Branch == nil {
		return nil, nil
	}
	for _, ga := range act.Branch.Gateway { //é xor attualmente
		if ga.Decision == "" {
			return ga.NextActivities, nil
		}
		if len(ga.NextActivities) == 0 {
			log.Printf("Process '%v' has not activities", processName)
			return []*Activity{}, nil
		}
		eval := goval.NewEvaluator()
		result, e := eval.Evaluate(ga.Decision, date, nil)
		if e != nil {
			return nil, fmt.Errorf("Process '%v' with activity '%v' has an eval error: %v", processName, act.Name, e.Error())
		}
		if result.(bool) {
			if e := checkLocalStorage(storage, act.Branch.RequiredOutputData); e != nil {
				return nil, fmt.Errorf("Process '%v' with activity '%v' has error: %v", processName, act.Name, e.Error())
			}
			storage.markWhatNeeded(act.Branch.RequiredOutputData)
			res = append(res, ga.NextActivities...)
			break
		}
	}
	return res, nil
}

func checkLocalStorage(st StorageData, req []TypeData) error {
	for _, dR := range req {
		d, err := st.GetLocal(dR.Name)
		if err != nil {
			return fmt.Errorf("Resource required is not found '%v'", dR.Name)
		}
		if d.(DataBpnm).GetType() != dR.Type {
			return fmt.Errorf("Resource '%v' has a difference type, exp: '%v', got: '%v'", dR.Name, dR.Type, d.(DataBpnm).GetType())
		}
	}
	return nil
}

func checkValidityGlobalStorage(st StorageData, req []TypeData) error {
	global := st.GetAllGlobal()
	if len(global) != len(req) {
		return fmt.Errorf("Stored values (%v) and declared data (in globalData field) (%v) are different", len(global), len(req))
	}
	for _, d := range req {
		v, ok := global[d.Name]
		if !ok {
			return fmt.Errorf("Required global resource is not found '%v'", d.Name)
		}
		if v.(DataBpnm).GetType() != d.Type {
			return fmt.Errorf("Global sesource '%v' has a difference type, exp: '%v', got: '%v'", d.Name, d.Type, v.(DataBpnm).GetType())
		}
	}
	return nil
}
