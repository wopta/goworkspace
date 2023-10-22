package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func gapContract(pdf *fpdf.Fpdf, origin string, policy *models.Policy) (string, []byte) {
	var (
		filename string
		out      []byte
	)

	filename, out = GapSogessur(pdf, origin, policy)

	return filename, out
}
