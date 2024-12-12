package consens

import (
	"context"
	"strings"

	"github.com/wopta/goworkspace/models"
)

const (
	ruiSectionEStrategy = "rui_section_e_mga_proponent"
	allNodesStrategy    = "all_nodes"
)

func newConsensStrategy(consens SystemConsens, node models.NetworkNode) (NeedConsentAlgorithm, error) {
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

type NeedConsentAlgorithm interface {
	Check(context.Context) (bool, error)
}

type RuisectionE struct {
	consens SystemConsens
	node    models.NetworkNode
}

func (w *RuisectionE) Check(ctx context.Context) (bool, error) {
	ruiSection := w.node.GetRuiSection()
	if ruiSection == "" {
		return false, errRuiSectionNotSet
	}

	if !strings.EqualFold(ruiSection, ruiSectionE) {
		return true, nil
	}

	if !w.node.IsMgaProponent {
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
	node    models.NetworkNode
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
