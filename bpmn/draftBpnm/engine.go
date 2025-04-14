package draftbpnm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/maja42/goval"
)

func (f *FlowBpnm) Run(processName string) error {
	log.SetPrefix("Bpnm")
	defer log.SetPrefix("")
	log.Println("Run ", processName)
	for _, process := range f.Process {
		if process.Name == processName {
			if e := checkValidityGlobalStorage(process.storageBpnm, process.RequiredGlobalData); e != nil {
				return e
			}
		}

		if e := process.Run(); e != nil { //TODO: how to check if there is an infinite loop
			return e
		}
		log.Println("Stop ", processName)
		return nil
	}
	return errors.New("No process founded")
}

func (f *ProcessBpnm) Run() error {
	f.activeActivity = f.Activities["init"]
	for {
		if f.activeActivity.handler == nil {
			return fmt.Errorf("No handler defined for %v", f.activeActivity.Name)
		}

		if pre := f.activeActivity.PreActivity; pre != nil {
			pre.storageBpnm.Merge(f.storageBpnm)
			if e := pre.Run(); e != nil {
				return e
			}
		}

		if e := checkAndCleanLocalStorage(f.storageBpnm, f.activeActivity.Branch.RequiredInputData); e != nil {
			return fmt.Errorf("Activity: %v, input: %v", f.activeActivity.Name, e.Error())
		}

		if e := f.activeActivity.handler(f.storageBpnm); e != nil {
			return e
		}

		if post := f.activeActivity.PostActivity; post != nil {
			post.storageBpnm.Merge(f.storageBpnm)
			if e := post.Run(); e != nil {
				return e
			}
		}

		//TODO: to improve
		m := mergeMaps(f.storageBpnm.GetAllGlobal(), f.storageBpnm.GetAllLocal())
		jsonMap := make(map[string]any)
		b, _ := json.Marshal(m)
		_ = json.Unmarshal(b, &jsonMap)

		act := f.activeActivity
		f.activeActivity = nil
		if e := f.EvaluateDecisions(act, jsonMap); e != nil {
			return e
		}
		if f.activeActivity == nil {
			return nil
		}
	}
}

func (f *ProcessBpnm) EvaluateDecisions(act *Activity, date map[string]any) error {
	for _, ga := range act.Branch.Gateway { //é xor attualmente
		if ga.Decision == "" {
			f.activeActivity = ga.NextActivities[0]
			break
		}
		if len(ga.NextActivities) == 0 {
			log.Println("No activity")
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

func NewBpnmBuilder(path string) (*BpnmBuilder, error) {
	var Bpnm BpnmBuilder
	jsonProva, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(jsonProva, &Bpnm)
	return &Bpnm, nil
}

func checkAndCleanLocalStorage(st StorageData, req []TypeData) error {
	temp := st.GetAllLocal()
	st.ResetLocal()
	for _, dR := range req {
		d, ok := temp[dR.Name]
		if !ok {
			return fmt.Errorf("Resource required is not found %v", dR.Name)
		}
		if d.(DataBpnm).Type() != dR.Type {
			return fmt.Errorf("Resource %v has a difference type, exp:%v, got %v", dR.Name, dR.Type, d.(DataBpnm).Type())
		}
		if e := st.AddLocal(dR.Name, d.(DataBpnm)); e != nil {
			return e
		}
	}
	return nil
}

func checkValidityGlobalStorage(st StorageData, req []TypeData) error {
	for _, d := range req {
		v, err := st.GetGlobal(d.Name)
		if err != nil {
			return fmt.Errorf("Required resource is not found %v", d.Name)
		}
		if v.Type() != d.Type {
			return fmt.Errorf("Resource %v has a difference type, exp:%v, got %v", d.Name, d.Type, v.Type())
		}
	}
	return nil
}
