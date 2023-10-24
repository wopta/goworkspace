package document

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
	"io"
	"log"
	"net/http"
)

func ReservedFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err     error
		policy  *models.Policy
		product *models.Product
	)
	log.Println("[ReservedFx] handler start ---------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("[ReservedFx] body: %s", string(body))

	err = json.Unmarshal(body, &policy)
	if err != nil {
		log.Printf("[ReservedFx] error unmarshaling request body: %s", err.Error())
		return "", nil, err
	}

	product = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	resp := Reserved(policy, product)

	respJson, err := json.Marshal(resp)

	log.Println("[ReservedFx] handler end -----------------------------------")

	return string(respJson), resp, err
}

func Reserved(policy *models.Policy, product *models.Product) *DocumentResponse {
	var (
		rawDoc []byte
		gsLink string
	)
	log.Println("[Reserved] function start ----------------------------------")

	switch policy.Name {
	case models.LifeProduct:
		log.Println("[Reserved] call lifeReserved...")
		gsLink, rawDoc = lifeReserved(policy, product)
	}

	return &DocumentResponse{
		LinkGcs: gsLink,
		Bytes:   base64.StdEncoding.EncodeToString(rawDoc),
	}
}
