package bpmn

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/maja42/goval"
	"github.com/wopta/goworkspace/models"
)

func (state *State) AddTaskHandler(name string, handler func(state *State) error) map[string]func(state *State) error {
	log.Printf("[AddTaskHandler] %s", name)
	if nil == state.Handlers {
		log.Println("[AddTaskHandler] state.Handlers == nil")
	}
	state.Handlers[name] = handler
	return state.Handlers
}

func (state *State) RunBpmn(processes []models.Process) {
	log.Println("[RunBpmn]")
	state.Processes = processes

	startProcesses := make([]models.Process, 0)

	for i, process := range processes {
		log.Printf("[RunBpmn] Index %d", i)
		if len(process.InProcess) == 0 {
			log.Printf("[RunBpmn] Adding process %s", process.Name)
			startProcesses = append(startProcesses, process)
		}
	}

	for _, process := range startProcesses {
		log.Printf("[RunBpmn] Running process %s", process.Name)
		state.runProcess(process)
	}
}

func (state *State) runNextProcess(process models.Process) {
	log.Println("[runNextProcess]")
	if !process.IsFailed {
		for _, x := range state.getProcesses(process.OutProcess) {
			state.runProcess(x)
		}
	}
}

func (state *State) runProcess(process models.Process) {
	log.Printf("[runProcess] %s", process.Name)
	id := process.Id
	state.Processes[id].Status = Active
	var (
		e error
		p models.Process
	)
	if process.Type == Task {
		e = state.Handlers[process.Name](state)
	}
	if process.Type == Decision {
		p, e = state.decisionStep(process)
		process = p
	}
	if e != nil {
		log.Printf("[runProcess] process %s FAILED", process.Name)
		state.Processes[id].Status = Failed
		state.IsFailed = true
	} else {
		log.Printf("[runProcess] process %s COMPLETED", process.Name)
		state.Processes[id].Status = Completed
		state.runNextProcess(process)
	}
}

func (state *State) getProcesses(ids []int) []models.Process {
	var processes []models.Process
	for _, id := range ids {
		for _, process := range state.Processes {
			if process.Id == id {
				processes = append(processes, process)
			}

		}
	}
	return processes
}

func (state *State) LoadProcesses(data string) ([]models.Process, error) {
	var processes []models.Process
	e := json.Unmarshal([]byte(data), &processes)
	state.Processes = processes
	return processes, e
}

func (state *State) decisionStep(process models.Process) (models.Process, error) {
	jsonMap := make(map[string]interface{})
	b, e := json.Marshal(state.Data)
	e = json.Unmarshal(b, &jsonMap)
	decision := strings.Replace(process.Decision, "\\", "\\", -1)
	log.Println(process.Decision)
	variables := jsonMap
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
