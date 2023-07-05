package bpmn

import (
	"encoding/json"
	"log"
)

var (
	state *State
)

func AddTaskHandler(name string, handler func(state *State) error) map[string]func(state *State) error {
	if nil == state.handlers {
		state.handlers = make(map[string]func(state *State) error)
	}
	state.handlers[name] = handler

	return state.handlers
}

func RunBpmn(processes []Process, data interface{}) {
	state.processes = processes
	state.data = data

	for i, process := range processes {
		log.Println(i)
		if len(process.InProcess) == 0 {
			runProcess(process)

		} else if len(process.OutProcess) == 0 {
			runProcess(process)
			break

		} else {

		}

	}

}
func runNextProcess(process Process) {

	if !process.IsFailed {
		for _, x := range getProcesses(process.OutProcess) {
			runProcess(x)
			runNextProcess(x)
		}

	}

}
func runProcess(process Process) {
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
		runNextProcess(process)
	}
}
func getProcesses(ids []int) []Process {
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
func loadProcesses(data string) ([]Process, error) {
	var processes []Process
	e := json.Unmarshal([]byte(data), &processes)

	return processes, e
}
