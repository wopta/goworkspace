package document

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

func ReservedFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy models.Policy
	)
	log.Println("[Reserved]")

	origin := r.Header.Get("Origin")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	err := json.Unmarshal(req, &policy)
	lib.CheckError(err)

	respObj := ReservedObj(origin, policy)
	respJson, err := json.Marshal(respObj)

	return string(respJson), respObj, err
}

func ReservedObj(origin string, policy models.Policy) DocumentResponse {
	var (
		out      []byte
		filename string
	)
	switch policy.Name {
	case models.LifeProduct:
		pdf := initFpdf()
		filename, out = LifeReserved(pdf, origin, policy)

	}

	return DocumentResponse{
		LinkGcs: filename,
		Bytes:   base64.StdEncoding.EncodeToString(out),
	}
}
