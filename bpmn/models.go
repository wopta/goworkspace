package bpmn

import (
	"time"
)

type ProcessType string
type ProcessStatus string

const (
	Task         string = "TASK"
	Decision     string = "DECISION"
	Ready        string = "READY"
	Active       string = "ACTIVE"
	WithDrawn    string = "WITHDRAWN"
	Completing   string = "COMPLETING"
	Completed    string = "COMPLETED"
	Failing      string = "FAILING"
	Terminating  string = "TERMINATING"
	Compensating string = "COMPENSATING"
	Failed       string = "FAILED"
	Terminated   string = "TERMINATED"
	Compensated  string = "COMPENSATED"
)

type TimerState string

const TimerCreated TimerState = "CREATED"
const TimerTriggered TimerState = "TRIGGERED"
const TimerCancelled TimerState = "CANCELLED"

type Process struct {
	Id              int         `json:"id"`
	LayerId         int         `json:"layer"`
	Name            string      `json:"name"`
	Shape           string      `json:"shape"`
	Type            string      `json:"type"`
	Status          string      `json:"decision "`
	Decision        string      `json:"status"`
	Data            interface{} `json:"data"`
	X               float64     `json:"x"`
	Y               float64     `json:"y"`
	InProcess       []int       `json:"inProcess"`
	OutProcess      []int       `json:"outProcess"`
	OutTrueProcess  []int       `json:"outTrueProcess"`
	OutFalseProcess []int       `json:"outFalseProcess"`
	IsCompleted     bool        `json:"isCompleted"`
	IsFailed        bool        `json:"isFailed"`
	Error           string      `json:"error"`
}
type Timer struct {
	ElementId          string
	ElementInstanceKey int64
	ProcessKey         int64
	ProcessInstanceKey int64
	State              TimerState
	CreatedAt          time.Time
	DueAt              time.Time
	Duration           time.Duration
}
type job struct {
	ElementId          string
	ElementInstanceKey int64
	ProcessInstanceKey int64
	JobKey             int64

	CreatedAt time.Time
}
type State struct {
	name           string
	processes      []Process
	data           interface{}
	jobs           []*job
	timers         []*Timer
	scheduledFlows []string
	handlers       map[string]func(state *State) error
	IsFailed       bool
}
type activatedJob struct {
	completeHandler          func()
	failHandler              func(reason string)
	key                      int64
	processInstanceKey       int64
	bpmnProcessId            string
	processDefinitionVersion int32
	processDefinitionKey     int64
	elementId                string
	createdAt                time.Time
}
