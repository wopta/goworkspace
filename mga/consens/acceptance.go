package consens

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type UpdateReq struct {
	Slug    string `json:"slug"`
	Product string `json:"product"`
	Value   string `json:"value"`
}

func AcceptanceFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err               error
		request           UpdateReq
		consens           SystemConsens
		undeclaredConsens []SystemConsens
		response          ConsensResp
		responseBytes     []byte
		now               = time.Now().UTC()
	)

	defer func() {
		if err != nil {
			log.Printf("error: %v", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()

	log.SetPrefix("[AcceptanceFx] ")
	log.Println("Handler start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.Println("error extracting authToken")
		return "", nil, err
	}

	networkNode := network.GetNetworkNodeByUid(authToken.UserID)
	if networkNode == nil {
		log.Println("error getting networkNode")
		err = errNetworkNodeNotFound
		return "", nil, err
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()
	if err != nil {
		log.Println("error decoding request body")
		return "", nil, err
	}

	filepath := folderPath + request.Product + "/" + request.Slug
	consens, err = getConsensByPath(filepath)
	if err != nil {
		log.Println("error getting consens")
		return "", nil, err
	}

	networkNode.Consens = append(networkNode.Consens, models.NodeConsens{
		Slug:     consens.Slug,
		ExpireAt: consens.ExpireAt,
		StartAt:  consens.StartAt,
		Title:    consens.Title,
		Content:  consens.ToString(),
		Value:    request.Value,
		GivenAt:  now,
	})

	if err := networkNode.SaveFirestore(); err != nil {
		log.Println("error saving networkNode in firestore")
		return "", nil, err
	}
	if err := networkNode.SaveBigQuery(""); err != nil {
		log.Println("error saving networkNode in bigquery")
		return "", nil, err
	}

	audit := NodeConsensAudit{
		Name:            networkNode.GetName(),
		RuiCode:         networkNode.GetRuiCode(),
		RuiRegistration: networkNode.GetRuiRegistration(),
		FiscalCode:      networkNode.GetFiscalCode(),
		VatCode:         networkNode.GetVatCode(),
		Slug:            consens.Slug,
		Title:           consens.Title,
		Content:         consens.ToString(),
		Answer:          request.Value,
		GivenAt:         now,
	}
	if err := audit.Save(); err != nil {
		log.Println("error saving consens audit")
		return "", nil, err
	}
	// send mail

	if undeclaredConsens, err = getUndeclaredConsens(request.Product, networkNode); err != nil {
		log.Println("error getting undeclared consens")
		return "", nil, err
	}

	response.Consens = undeclaredConsens

	if responseBytes, err = json.Marshal(response); err != nil {
		log.Println("error marshalling response")
		return "", nil, err
	}

	return string(responseBytes), response, err
}

func getConsensByPath(path string) (SystemConsens, error) {
	var (
		fileBytes []byte
		consens   SystemConsens
		err       error
	)

	if fileBytes, err = lib.GetFilesByEnvV2(path); err != nil {
		return SystemConsens{}, err
	}

	if err = json.Unmarshal(fileBytes, &consens); err != nil {
		return SystemConsens{}, err
	}

	return consens, err
}
