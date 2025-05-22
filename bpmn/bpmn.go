package bpmn

import (
	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/models"
)

func NewBpmn(data models.Policy) *State {
	// Init workflow with a name, and max concurrent tasks
	log.Println("--------------------------NewBpmn-------------------------------------------")
	state := &State{
		Handlers: make(map[string]func(state *State) error),
		Data:     &data,
	}
	state.Handlers = make(map[string]func(state *State) error)
	return state
}
