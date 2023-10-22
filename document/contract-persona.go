package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func personaContract(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	var (
		filename string
		out      []byte
	)

	filename, out = personaGlobalContractV1(pdf, policy, networkNode, product)

	return filename, out
}
