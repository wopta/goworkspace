package document

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/document/pkg/contract"
	lib "gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type requestGenerateMup struct {
	CompanyName      string  `json:"companyName"`
	ConsultancyPrice float64 `json:"consultancyPrice"`
}

func GenerateMupFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req  requestGenerateMup
		resp struct {
			bytes []byte
		}
	)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(body, &req)
	if err != nil {
		log.ErrorF("error unmarshaling request body")
		return "", nil, err
	}
	nodeUid := chi.URLParam(r, "nodeUid")

	bytes, err := contract.GenerateMup(req.CompanyName, req.ConsultancyPrice, models.NetworkChannel, nodeUid)

	resp.bytes = bytes.Bytes()
	respBytes, err := json.Marshal(resp)
	return string(respBytes), bytes, err
}
