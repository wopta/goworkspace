package consens

import (
	"context"
	"strings"
	"time"

	"github.com/wopta/goworkspace/models"
)

const (
	ruiSectionEStrategy = "rui_section_e"
	allNodesStrategy    = "all_nodes"
)

func newConsensStrategy(consens SystemConsens, node NodeWithConsens) (NeedConsensAlgorithm, error) {
	switch consens.Strategy {
	case ruiSectionEStrategy:
		return &RuisectionE{
			consens: consens,
			node:    node,
		}, nil
	case allNodesStrategy:
		return &AllNodes{
			consens: consens,
			node:    node,
		}, nil
	}
	return nil, errStrategyNotFound
}

type NeedConsensAlgorithm interface {
	Check(context.Context) (bool, error)
}

type RuisectionE struct {
	consens SystemConsens
	node    NodeWithConsens
}

func (w *RuisectionE) Check(ctx context.Context) (bool, error) {
	var ruiSection string
	switch w.node.Type {
	case models.AgentNetworkNodeType:
		ruiSection = w.node.Agent.RuiSection
	case models.AgencyNetworkNodeType:
		ruiSection = w.node.Agency.RuiSection
	case models.BrokerNetworkNodeType:
		ruiSection = w.node.Broker.RuiSection
	case models.AreaManagerNetworkNodeType:
		ruiSection = w.node.AreaManager.RuiSection
	case models.PartnershipNetworkNodeType:
		return false, errPartnershipNode
	}

	if !strings.EqualFold(ruiSection, ruiSectionE) {
		return true, nil
	}

	now := getTimestamp(ctx)

	if now.Before(w.consens.StartAt) {
		return true, nil
	}

	if now.Before(w.consens.ExpireAt) {
		return false, nil
	}

	return true, nil
}

type AllNodes struct {
	consens SystemConsens
	node    NodeWithConsens
}

func (w *AllNodes) Check(ctx context.Context) (bool, error) {
	now := getTimestamp(ctx)

	if now.Before(w.consens.StartAt) {
		return true, nil
	}

	if now.Before(w.consens.ExpireAt) {
		return false, nil
	}

	return true, nil
}

func getTimestamp(ctx context.Context) time.Time {
	if rawTime := ctx.Value(timestamp); rawTime != nil {
		return (rawTime).(time.Time)
	}
	return time.Time{}
}
