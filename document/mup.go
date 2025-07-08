package document

import (
	"bytes"

	"gitlab.dev.wopta.it/goworkspace/document/pkg/contract"
)

func GenerateMup(companyName string, consultancyPrice float64, channel string) (out bytes.Buffer, err error) {
	return contract.GenerateMup(companyName, consultancyPrice, channel)
}
