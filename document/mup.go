package document

import (
	"bytes"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/document/pkg/contract"
)

func GenerateMupFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	return contract.GenerateMup(companyName, consultancyPrice, channel)
}
