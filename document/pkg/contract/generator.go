package contract

import "github.com/wopta/goworkspace/document/internal/engine"

type baseGenerator struct {
	engine     *engine.Fpdf
	isProposal bool
}
