package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func proposalLifeAxaV2(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	lifeMainHeaderV2(pdf, policy, networkNode, true)

	return "", nil
}
