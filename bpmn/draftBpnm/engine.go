package draftbpnm

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"

	"github.com/maja42/goval"
)

func (f *FlowBpnm) Run(processName string) error {
	for _, process := range f.Process {
		if process.Name == processName {
			if e := checkValidityGlobalStorage(process.storageBpnm, process.RequiredGlobalData); e != nil {
				return e
			}
		}

		if e := process.Run(); e != nil {
			return e
		}
		return nil
	}
	return errors.New("no process founded")
}

func (f *ProcessBpnm) Run() error {
	f.activeActivity = f.Activities["init"]
	for {
		if f.activeActivity.handler == nil {
			return fmt.Errorf("no handler defined for %v", f.activeActivity.Name)
		}
		if e := checkAndCleanLocalStorage(f.storageBpnm, f.activeActivity.Branch.RequiredInputData); e != nil {
			return fmt.Errorf("Activity: %v, input: %v", f.activeActivity.Name, e.Error())
		}
		if e := f.activeActivity.handler(f.storageBpnm); e != nil {
			return e
		}

		//TODO: to improve
		m := mergeMaps(f.storageBpnm.GetAllGlobal(), f.storageBpnm.GetAllLocal())
		jsonMap := make(map[string]any)
		b, _ := json.Marshal(m)
		_ = json.Unmarshal(b, &jsonMap)

		act := f.activeActivity
		f.activeActivity = nil
		if e := f.RunDecision(act, jsonMap); e != nil {
			return e
		}
		if f.activeActivity == nil {
			return nil
		}
	}
}

func (f *ProcessBpnm) RunDecision(act *Activity, date map[string]any) error {
	for _, ga := range act.Branch.Gateway { //é xor attualmente
		if ga.Decision == "" {
			f.activeActivity = ga.NextActivities[0]
			break
		}
		if len(ga.NextActivities) == 0 {
			return nil
		}
		eval := goval.NewEvaluator()
		result, e := eval.Evaluate(ga.Decision, date, nil)
		if e != nil {
			return fmt.Errorf("Activity: %v, error eval:%v ", act.Name, e.Error())
		}
		if result.(bool) {
			if e := checkAndCleanLocalStorage(f.storageBpnm, act.Branch.RequiredOutputData); e != nil {
				return fmt.Errorf("Activity: %v, output: %v", act.Name, e.Error())
			}
			f.activeActivity = ga.NextActivities[0]
			break
		}
	}
	return nil
}

func NewBpnmBuilder() (*BpnmBuilder, error) {
	var Bpnm BpnmBuilder
	jsonProva, err := os.ReadFile("prova.json")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(jsonProva, &Bpnm)
	return &Bpnm, nil
}

func checkAndCleanLocalStorage(st StorageData, req []TypeData) error {
	temp := maps.Clone(st.GetAllLocal())
	st.resetLocal()
	for _, dR := range req {
		d, ok := temp[dR.Name]
		if !ok {
			return fmt.Errorf("resource required is not found %v", dR.Name)
		}
		if d.(DataBpnm).Type() != dR.Type {
			return fmt.Errorf("resource %v has a differente type, exp:%v, got %v", dR.Name, dR.Type, d.(DataBpnm).Type())
		}
		st.AddLocal(dR.Name, d.(DataBpnm))
	}
	return nil
}

func checkValidityGlobalStorage(st StorageData, req []TypeData) error {
	for _, d := range req {
		v, err := st.GetGlobal(d.Name)
		if err != nil {
			return fmt.Errorf("required resource is not present %v", d.Name)
		}
		if v.Type() != d.Type {
			return fmt.Errorf("resource %v has a differente type, exp:%v, got %v", d.Name, d.Type, v.Type())
		}
	}
	return nil
}

func mergeMaps(m1 map[string]any, m2 map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
}
