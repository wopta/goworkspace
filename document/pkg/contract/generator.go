package contract

import (
	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/lib"
)

type baseGenerator struct {
	engine     *engine.Fpdf
	isProposal bool
}

func (bg *baseGenerator) woptaHeader() {
	bg.engine.SetHeader(func() {
		bg.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 10, 6, 0, 10)
		bg.engine.NewLine(10)

		if bg.isProposal {
			bg.engine.DrawWatermark(constants.Proposal)
		}
	})
}
