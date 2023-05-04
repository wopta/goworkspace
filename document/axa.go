package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func AxaContract(pdf *fpdf.Fpdf, policy models.Policy) (string, []byte) {
	var (
		filename string
		out      []byte
	)

	if policy.Name == "life" {
		filename, out = Life(pdf, policy)
	}

	return filename, out
}
