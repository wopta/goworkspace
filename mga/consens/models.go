package consens

import (
	"errors"
	"time"

	"github.com/wopta/goworkspace/models"
)

var (
	errNetworkNodeNotFound = errors.New("network node not found")
	errPartnershipNode     = errors.New("partnership node does not have rui registration nor consens")
	errStrategyNotFound    = errors.New("strategy not found")
)

const (
	ruiSectionE         = "E"
	ruiSectionEStrategy = "rui_section_e"
	allNodesStrategy    = "all_nodes"
)

type NetworkConsens struct {
	Slug        string           `json:"slug"`
	ExpireAt    time.Time        `json:"expireAt"`
	StartAt     time.Time        `json:"startAt"`
	AvailableAt time.Time        `json:"availableAt"`
	Strategy    string           `json:"strategy"`
	Title       string           `json:"title"`
	Content     []ConsensContent `json:"content"`
}

type ConsensContent struct {
	Text       string `json:"text"`
	InputType  string `json:"inputType,omitempty"`
	InputName  string `json:"inputName,omitempty"`
	InputValue string `json:"inputValue,omitempty"`
}

// TODO: add the fields to the correct struct
type NodeWithConsens struct {
	models.NetworkNode
	Consens []NetworkConsens `json:"networkConsens"`
}
