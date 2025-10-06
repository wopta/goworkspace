package bpmnEngine

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/maja42/goval"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

type BpmnFlow string

const (
	Emit            BpmnFlow = "emit"
	Lead            BpmnFlow = "lead"
	Proposal        BpmnFlow = "proposal"
	RequestApproval BpmnFlow = "requestApproval"
	Acceptance      BpmnFlow = "acceptance"
	Pay             BpmnFlow = "pay"
	Sign            BpmnFlow = "sign"
)

// Run a process, it starts from defaultActivity defined in json
func (f *FlowBpnm) Run(flow BpmnFlow) error {
	process := f.process[string(flow)]
	if process == nil {
		return fmt.Errorf("Process '%v' not found", string(flow))
	}
	return f.RunAt(flow, process.defaultStart)
}

func (f *FlowBpnm) RunAt(flow BpmnFlow, startingActivity string) error {
	log.InfoF("Run %v", flow)
	process := f.process[string(flow)]
	if process == nil {
		return fmt.Errorf("Process '%v' not found", string(flow))
	}

	if err := process.storageBpnm.AddGlobal("statusFlow", &StatusFlow{CurrentProcess: BpmnFlow(process.name), CurrentActivity: "Starting process"}); err != nil {
		return err
	}

	if err := checkGlobalResources(process.storageBpnm, process.requiredGlobalData); err != nil {
		return err
	}

	var firstActivities []*activity
	if act := process.activities[startingActivity]; act != nil {
		firstActivities = []*activity{process.activities[startingActivity]}
	}
	if process.storageBpnm == nil {
		return errors.New("Miss storage")
	}
	if len(firstActivities) == 0 {
		return fmt.Errorf("Process '%v' has no activity '%v'", process.name, startingActivity)
	}
	var err error

	if err = process.loop(process.storageBpnm, firstActivities...); err != nil {
		return err
	}
	if process.lastActivity != nil && !process.lastActivity.callEndIfStop {
		return nil
	}
	if err = process.activities[getNameEndActivity(process.name)].runActivity(process, process.storageBpnm); err != nil {
		return err
	}

	return nil
}

func (p *processBpnm) loop(initialStorage *StorageBpnm, activities ...*activity) (err error) {
	var newStorage *StorageBpnm
	for i := range activities {
		newStorage = NewStorageBpnm()
		err = newStorage.setHigherStorage(initialStorage)
		if err != nil {
			return err
		}
		if err := newStorage.AddGlobal("statusFlow", &StatusFlow{CurrentProcess: BpmnFlow(p.name)}); err != nil {
			return err
		}

		initialStorage = newStorage

		if err = activities[i].runActivity(p, newStorage); err != nil {
			return err
		}
		p.lastActivity = activities[i]

		mapsMerged := newStorage.GetMap()
		listNewActivities, err := activities[i].evaluateDecisions(p.name, newStorage, mapsMerged)
		if err != nil {
			return err
		}
		if len(listNewActivities) == 0 {
			continue
		}
		newStorage.cleanNoMarkedResources()
		if err = p.loop(newStorage, listNewActivities...); err != nil {
			return err
		}
	}
	return err
}

func (act *activity) runActivity(process *processBpnm, storage *StorageBpnm) (err error) {
	defer func() {
		if process.recover != nil {
			if r := recover(); r != nil || (err != nil) {
				log.InfoF("Run recorver process '%v', for activity '%v'", process.name, act.name)
				for i := range process.recover {
					err = errors.Join(err, process.recover[i](storage))
				}
				if r != nil {
					err = errors.Join(err, errors.New(r.(string)))
				}
			}
		}
		status := ""
		if err == nil {
			status = "OK"
		} else {
			status = "Fail: " + err.Error()
		}
		log.InfoF("Run process '%v', finished activity '%v' with status: %v", process.name, act.name, status)
	}()
	for i := range act.preActivity {
		if pre := act.preActivity[i]; pre != nil {
			err = pre.storageBpnm.setHigherStorage(storage)
			if err != nil {
				return err
			}
			if err = pre.loop(pre.storageBpnm, pre.activities[pre.defaultStart]); err != nil {
				return err
			}
		}
	}
	if err = checkResources(storage, act.requiredInputData); err != nil {
		return fmt.Errorf("Process '%v' with Activity  '%v' has an input error: %v", process.name, act.name, err.Error())
	}

	if err = callActivity(process, storage, act); err != nil {
		return err
	}
	for i := range act.postActivity {
		if post := act.postActivity[i]; post != nil {
			err = post.storageBpnm.setHigherStorage(storage)
			if err != nil {
				return err
			}
			if err = post.loop(post.storageBpnm, post.activities[post.defaultStart]); err != nil {
				return err
			}
		}
	}
	return nil
}

func callActivity(process *processBpnm, storage *StorageBpnm, act *activity) (err error) {
	log.InfoF("Run process '%v', start activity '%v'", process.name, act.name)
	var status DataBpnm
	status, err = storage.GetGlobal("statusFlow")
	if err != nil {
		return fmt.Errorf("Error setting status flow of Process '%v' with activity '%v'", process.name, act.name)
	}
	status.(*StatusFlow).CurrentActivity = act.name
	if act.handler == nil {
		return nil
	}
	if act.handler == nil {
		return
	}
	err = act.handler(storage)
	return err
}

func (act *activity) evaluateDecisions(processName string, storage *StorageBpnm, date map[string]any) ([]*activity, error) {
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
			if err = checkResources(storage, act.requiredOutputData); err != nil {
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
			if err = checkResources(storage, act.requiredOutputData); err != nil {
				return nil, fmt.Errorf("Process '%v' with activity '%v' has an output error: %v", processName, act.name, err.Error())
			}
			storage.markWhatNeeded(act.requiredOutputData)
			res = append(res, ga.nextActivities...)
			break
		}
	}
	return res, nil
}

func checkResources(st *StorageBpnm, req []typeData) error {
	local := st.getAllGlobals()
	for name, value := range st.getAllLocals() {
		local[name] = value
	}
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

func checkGlobalResources(st *StorageBpnm, req []typeData) error {
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
