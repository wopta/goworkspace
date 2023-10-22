package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
	"log"
)

func lifeProposal(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	var (
		rawDoc   []byte
		filename string
	)

	log.Println("[lifeProposal] function start ------------------------------")

	switch policy.ProductVersion {
	case models.ProductV1:
		filename, rawDoc = lifeAxaProposalV1(pdf, origin, policy, networkNode, product)
	case models.ProductV2:
		filename, rawDoc = lifeAxaProposalV2(pdf, origin, policy, networkNode, product)
	}

	log.Println("[lifeProposal] function end --------------------------------")

	return filename, rawDoc
}
