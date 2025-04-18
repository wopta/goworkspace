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

func (f *ProcessBpnm) run(nameActivity string) error {
	f.activeActivities = append(f.activeActivities, f.Activities[nameActivity])
	if f.storageBpnm == nil {
		return errors.New("miss storage")
	}
	if f.activeActivities == nil {
		return fmt.Errorf("Process '%v' has no activity '%v'", f.Name, nameActivity)
	}
	//TODO: implement at garbange collector with a counter, to remove old element not used in the next branch
	for {
		var nextActivities []*Activity
		for i := range f.activeActivities {
			if err := f.runActivity(f.activeActivities[i], f.storageBpnm); err != nil {
				return err
			}
			//TODO: to improve

			m := mergeMaps(f.storageBpnm.GetAllGlobal(), f.storageBpnm.GetAllLocal())
			jsonMap := make(map[string]any)
			b, _ := json.Marshal(m)
			_ = json.Unmarshal(b, &jsonMap)

			list, e := f.evaluateDecisions(f.activeActivities[i], jsonMap)
			if e != nil {
				return e
			}
			nextActivities = append(nextActivities, list...)
		}
		if len(nextActivities) == 0 {
			return nil
		}
		f.activeActivities = nextActivities
	}
}

func (f *ProcessBpnm) runActivity(act *Activity, storage StorageData) error {
	if act.handler == nil {
		return fmt.Errorf("Process '%v' has no handler defined for activity '%v'", f.Name, act.Name)
	}
	log.Printf("Run process '%v', activity '%v'", f.Name, act.Name)
	if pre := act.PreActivity; pre != nil {
		pre.storageBpnm.Merge(storage)
		if e := pre.run(pre.DefaultStart); e != nil {
			return e
		}
	}

	if act.Branch != nil {
		if e := checkLocalStorage(storage, act.Branch.RequiredInputData); e != nil {
			return fmt.Errorf("Process '%v' with activity '%v' has an input error: %v", f.Name, act.Name, e.Error())
		}
	}

	if e := act.handler(f.storageBpnm); e != nil {
		return e
	}

	if post := act.PostActivity; post != nil {
		post.storageBpnm.Merge(storage)
		if e := post.run(post.DefaultStart); e != nil {
			return e
		}
	}
	return nil
}

func (f *ProcessBpnm) evaluateDecisions(act *Activity, date map[string]any) ([]*Activity, error) {
	var res []*Activity
	if act.Branch == nil {
		return nil, nil
	}
	for _, ga := range act.Branch.Gateway { //é xor attualmente
		if ga.Decision == "" {
			return ga.NextActivities, nil
		}
		if len(ga.NextActivities) == 0 {
			log.Printf("Process '%v' has not activities", f.Name)
			return []*Activity{}, nil
		}
		eval := goval.NewEvaluator()
		result, e := eval.Evaluate(ga.Decision, date, nil)
		if e != nil {
			return nil, fmt.Errorf("Process '%v' with activity '%v' has an eval error: %v", f.Name, act.Name, e.Error())
		}
		if result.(bool) {
			if e := checkLocalStorage(f.storageBpnm, act.Branch.RequiredOutputData); e != nil {
				return nil, fmt.Errorf("Process '%v' with activity '%v' has error: %v", f.Name, act.Name, e.Error())
			}
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
	for _, d := range req {
		v, err := st.GetGlobal(d.Name)
		if err != nil {
			return fmt.Errorf("Required global resource is not found '%v'", d.Name)
		}
		if v.GetType() != d.Type {
			return fmt.Errorf("Global sesource '%v' has a difference type, exp: '%v', got: '%v'", d.Name, d.Type, v.GetType())
		}
	}
	return nil
}
