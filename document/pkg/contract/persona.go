package contract

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

type PersonaGenerator struct {
	*baseGenerator
}

func NewPersonaGenerator(engine *engine.Fpdf, policy *models.Policy, node *models.NetworkNode, product models.Product, isProposal bool) *PersonaGenerator {
	var worksForNode *models.NetworkNode
	if node != nil && node.WorksForUid != "" {
		worksForNode = network.GetNetworkNodeByUid(node.WorksForUid)
	}

	return &PersonaGenerator{
		baseGenerator: &baseGenerator{
			engine:       engine,
			isProposal:   isProposal,
			now:          time.Now(),
			signatureID:  0,
			networkNode:  node,
			worksForNode: worksForNode,
			policy:       policy,
		},
	}
}
