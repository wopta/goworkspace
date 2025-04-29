package document

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
	"io"
	"net/http"
)

func ReservedFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err     error
		policy  *models.Policy
		product *models.Product
	)
	log.AddPrefix("ReservedFx")
	defer log.PopPrefix()
	log.Println("handler start ---------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("body: %s", string(body))

	err = json.Unmarshal(body, &policy)
	if err != nil {
		log.ErrorF("error unmarshaling request body: %s", err.Error())
		return "", nil, err
	}

	product = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	resp := Reserved(policy, product)

	respJson, err := json.Marshal(resp)

	log.Println("handler end -----------------------------------")

	return string(respJson), resp, err
}

func Reserved(policy *models.Policy, product *models.Product) *DocumentResponse {
	var (
		rawDoc []byte
		gsLink string
	)
	log.AddPrefix("Reserved")
	log.Println("function start ----------------------------------")

	switch policy.Name {
	case models.LifeProduct:
		log.Println("call lifeReserved...")
		gsLink, rawDoc = lifeReserved(policy, product)
	}

	return &DocumentResponse{
		LinkGcs: gsLink,
		Bytes:   base64.StdEncoding.EncodeToString(rawDoc),
	}
}
