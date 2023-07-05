package bpmn

import (
	"encoding/json"
	"log"
)

func (state *State) AddTaskHandler(name string, handler func(state *State) error) map[string]func(state *State) error {
	log.Println("AddTaskHand")
	if nil == state.handlers {
		log.Println("nil")

	}
	state.handlers[name] = handler

	return state.handlers
}

func (state *State) RunBpmn(processes []Process, data interface{}) {
	state.processes = processes
	state.data = data

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
	state.processes[id].Status = Active
	var (
		e error
	)
	if process.Type == Task {
		e = state.handlers[process.Name](state)
	}
	if process.Type == Decision {
		e = state.handlers[process.Name](state)
	}
	if e != nil {
		state.processes[id].Status = Failed
		state.IsFailed = true
	} else {
		state.processes[id].Status = Completed
		state.runNextProcess(process)
	}
}
func (state *State) getProcesses(ids []int) []Process {
	var processes []Process
	for _, id := range ids {
		for _, process := range state.processes {
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
