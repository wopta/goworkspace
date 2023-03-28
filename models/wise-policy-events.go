package models

import "time"

type WisePolicyEvents struct {
	EventType        string    `json:"cdTipoEvento"`
	EventDescription string    `json:"txDescrizioneEvento"`
	IssueDate        time.Time `json:"dtCompetenzaEmesso"`
	SystemDate       string    `json:"dtSysEvento"`
}
