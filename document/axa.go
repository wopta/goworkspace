package document

import (
	"github.com/wopta/goworkspace/models"
)

func AxaContract(policy models.Policy) (string, []byte) {
	var (
		filename string
		out      []byte
	)

	pdf := initFpdf()

	if policy.Name == "life" {
		filename, out = Life(pdf, policy)
	}

	return filename, out
}
