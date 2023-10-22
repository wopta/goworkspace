package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func gapProposal(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode) (string, []byte) {
	var (
		filename string
		out      []byte
	)

	filename, out = gapSogessurProposalV1(pdf, origin, policy, networkNode)

	return filename, out
}
