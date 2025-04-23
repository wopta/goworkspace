package draftbpnm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/maja42/goval"
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
	p.activeActivities = nil
	if act := p.Activities[nameActivity]; act != nil {
		p.activeActivities = append(p.activeActivities, p.Activities[nameActivity])
	}
	if p.storageBpnm == nil {
		return errors.New("Miss storage")
	}
	if p.activeActivities == nil || len(p.activeActivities) == 0 {
		return fmt.Errorf("Process '%v' has no activity '%v'", p.Name, nameActivity)
	}

	for {
		var nextActivities []*Activity
		for i := range p.activeActivities {
			if err := p.activeActivities[i].runActivity(p.Name, p.storageBpnm); err != nil {
				return err
			}
			//TODO: to improve
			m := mergeMaps(p.storageBpnm.getAllGlobal(), p.storageBpnm.getAllLocal())
			jsonMap := make(map[string]any)
			b, _ := json.Marshal(m)
			_ = json.Unmarshal(b, &jsonMap)

			listNewActivities, e := p.activeActivities[i].evaluateDecisions(p.Name, p.storageBpnm, jsonMap)

			if e != nil {
				return e
			}
			nextActivities = append(nextActivities, listNewActivities...)
		}
		if len(nextActivities) == 0 {
			p.Activities[fmt.Sprintf("%v_end", p.Name)].runActivity(p.Name, p.storageBpnm)
			return nil
		}
		p.activeActivities = nextActivities
		p.storageBpnm.clean()
	}
}

func (act *Activity) runActivity(nameProcess string, storage StorageData) error {
	log.Printf("Run process '%v', activity '%v'", nameProcess, act.Name)
	if pre := act.PreActivity; pre != nil {
		if err := pre.run(pre.DefaultStart); err != nil {
			return err
		}
	}
	if act.Branch != nil {
		fmt.Printf("process name %v, len %v\n", nameProcess, len(act.Branch.RequiredInputData))
		if e := checkLocalStorage(storage, act.Branch.RequiredInputData); e != nil {
			return fmt.Errorf("Process '%v' with activity '%v' has an input error: %v", nameProcess, act.Name, e.Error())
		}
	}

	if act.handler != nil {
		var err error
		func() {
			defer func() {
				if act.recover == nil {
					return
				}
				if r := recover(); r != nil || err != nil {
					log.Printf("Run recorver process '%v', activity '%v'", nameProcess, act.Name)
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

	if post := act.PostActivity; post != nil {
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
			if e := checkLocalStorage(storage, act.Branch.RequiredOutputData); e != nil {
				return nil, fmt.Errorf("Process '%v' with activity '%v' has an output error: %v", processName, act.Name, e.Error())
			}
			storage.markWhatNeeded(act.Branch.RequiredOutputData)
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
				return nil, fmt.Errorf("Process '%v' with activity '%v' has an output error: %v", processName, act.Name, e.Error())
			}
			storage.markWhatNeeded(act.Branch.RequiredOutputData)
			res = append(res, ga.NextActivities...)
			break
		}
	}
	return res, nil
}

func checkLocalStorage(st StorageData, req []TypeData) error {
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

func checkValidityGlobalStorage(st StorageData, req []TypeData) error {
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
