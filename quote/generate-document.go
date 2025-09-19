package quote

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/document"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/product"
)

type DraftResponse struct {
	RawDoc string `json:"rawDoc"`
}

func generateDocumentFx(_ http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err           error
		policy        models.Policy
		nn            *models.NetworkNode
		w             *models.Warrant
		docBytes      []byte
		responseBytes []byte
	)

	log.AddPrefix("GenerateQuoteDocumentFx")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.ErrorF("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start -----------------------------------------------")

	if err = json.NewDecoder(r.Body).Decode(&policy); err != nil {
		log.ErrorF("error decoding request body")
		return "", nil, err
	}

	policy.Normalize()

	if nn = network.GetNetworkNodeByUid(policy.ProducerUid); nn != nil {
		w = nn.GetWarrant()
	}

	prd := product.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, nn, w)

	if docBytes, err = document.Quote(&policy, prd); err != nil {
		log.ErrorF("error generating document")
		return "", nil, err
	}

	response := DraftResponse{
		RawDoc: base64.StdEncoding.EncodeToString(docBytes),
	}

	if responseBytes, err = json.Marshal(response); err != nil {
		log.ErrorF("error marshaling response")
		return "", nil, err
	}

	return string(responseBytes), response, err
}
