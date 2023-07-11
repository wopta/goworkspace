package bpmn

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/maja42/goval"
)

func (state *State) AddTaskHandler(name string, handler func(state *State) error) map[string]func(state *State) error {
	log.Println("AddTaskHand")
	if nil == state.Handlers {
		log.Println("nil")
	}
	state.Handlers[name] = handler
	return state.Handlers
}

func (state *State) RunBpmn(processes []Process, data interface{}) {
	state.Processes = processes
	state.Data = data

	for i, process := range processes {
		log.Println(i)
		if len(process.InProcess) == 0 {
			state.runProcess(process)

		} else if len(process.OutProcess) == 0 {
			state.runProcess(process)
			break

		} else {

		}

	}

}
func (state *State) runNextProcess(process Process) {
	log.Println("runNextProcess")
	if !process.IsFailed {
		for _, x := range state.getProcesses(process.OutProcess) {
			state.runProcess(x)
			state.runNextProcess(x)
		}

	}

}
func (state *State) runProcess(process Process) {
	log.Println("runProcess")
	id := process.Id
	state.Processes[id].Status = Active
	var (
		e error
		p Process
	)
	if process.Type == Task {
		e = state.Handlers[process.Name](state)
	}
	if process.Type == Decision {
		p, e = state.decisionStep(process)
		process = p
	}
	if e != nil {
		state.Processes[id].Status = Failed
		state.IsFailed = true
	} else {
		state.Processes[id].Status = Completed
		state.runNextProcess(process)
	}
}
func (state *State) getProcesses(ids []int) []Process {
	var processes []Process
	for _, id := range ids {
		for _, process := range state.Processes {
			if process.Id == id {
				processes = append(processes, process)
			}

		}
	}
	return processes
}
func (state *State) loadProcesses(data string) ([]Process, error) {
	var processes []Process
	e := json.Unmarshal([]byte(data), &processes)

	return processes, e
}
func (state *State) decisionStep(process Process) (Process, error) {

	decision := strings.Replace(process.Decision, "\\", "\\", -1)
	log.Println(process.Decision)
	variables := state.DecisionData
	eval := goval.NewEvaluator()
	result, e := eval.Evaluate(decision, variables, nil) // Returns <true, nil>
	log.Println(result)
	if result.(bool) {
		process.OutProcess = process.OutTrueProcess
	} else {
		process.OutProcess = process.OutFalseProcess
	}

	return process, e
}
