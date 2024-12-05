package consens

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
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
		return "", nil, err
	}

	networkNode := network.GetNetworkNodeByUid(authToken.UserID)
	if networkNode == nil {
		err = errNetworkNodeNotFound
		return "", nil, err
	}
	tempBytes, _ := json.Marshal(networkNode)
	var tempNode NodeWithConsens
	json.Unmarshal(tempBytes, &tempNode)

	err = json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()
	if err != nil {
		log.Printf("error decoding request body: %s", err)
		return "", nil, err
	}

	filepath := folderPath + request.Product + "/" + request.Slug
	consens, err = getConsensByPath(filepath)
	if err != nil {
		return "", nil, err
	}

	tempNode.Consens = append(tempNode.Consens, NodeConsens{
		Slug:     consens.Slug,
		ExpireAt: consens.ExpireAt,
		StartAt:  consens.StartAt,
		Title:    consens.Title,
		Content:  consens.ToString(),
		Value:    request.Value,
		GivenAt:  now,
	})

	// save node
	if err := tempNode.SaveFirestore(); err != nil {
		return "", nil, err
	}
	if err := tempNode.SaveBigQuery(""); err != nil {
		return "", nil, err
	}
	// save consens
	audit := NodeConsensAudit{
		Name:            tempNode.GetName(),
		RuiCode:         tempNode.GetRuiCode(),
		RuiRegistration: tempNode.GetRuiRegistration(),
		FiscalCode:      tempNode.GetFiscalCode(),
		VatCode:         tempNode.GetVatCode(),
		Slug:            consens.Slug,
		Title:           consens.Title,
		Content:         consens.ToString(),
		Answer:          request.Value,
		GivenAt:         now,
	}
	if err := audit.Save(); err != nil {
		return "", nil, err
	}
	// send mail

	// get consens
	if undeclaredConsens, err = getUndeclaredConsens(request.Product, &tempNode); err != nil {
		return "", nil, err
	}

	response.Consens = undeclaredConsens

	if responseBytes, err = json.Marshal(response); err != nil {
		return "", nil, err
	}

	return string(responseBytes), response, err
}

func getConsensByPath(path string) (SystemConsens, error) {
	var consens SystemConsens
	fileBytes, err := lib.GetFilesByEnvV2(path)
	if err != nil {
		return SystemConsens{}, err
	}

	if err := json.Unmarshal(fileBytes, &consens); err != nil {
		return SystemConsens{}, err
	}

	return consens, nil
}
