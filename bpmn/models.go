package bpmn

import (
	"time"

	"github.com/wopta/goworkspace/models"
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
	Name           string
	Processes      []models.Process
	Data           *models.Policy
	DecisionData   *map[string]interface{}
	Jobs           []*job
	Timers         []*Timer
	ScheduledFlows []string
	Handlers       map[string]func(state *State) error
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
