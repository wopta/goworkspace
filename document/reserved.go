package document

import (
	"encoding/json"
	"io"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	prd "gitlab.dev.wopta.it/goworkspace/product"
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

	resp, err := Reserved(policy, product)
	if err != nil {
		return "", nil, err
	}
	response, err := resp.Save()
	if err != nil {
		return "", nil, err
	}
	respJson, err := json.Marshal(response)

	log.Println("handler end -----------------------------------")

	return string(respJson), resp, err
}

func Reserved(policy *models.Policy, product *models.Product) (DocumentGenerated, error) {
	var (
		document DocumentGenerated
		err      error
	)
	log.AddPrefix("Reserved")
	log.Println("function start ----------------------------------")

	switch policy.Name {
	case models.LifeProduct:
		log.Println("call lifeReserved...")
		document, err = lifeReserved(policy, product)
	}

	return document, err
}
