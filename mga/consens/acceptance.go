package consens

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type AcceptanceReq struct {
	Slug    string `json:"slug"`
	Product string `json:"product"`
	Value   string `json:"value"`
}

func AcceptanceFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err               error
		request           AcceptanceReq
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

	if request.Product == "" || request.Slug == "" || request.Value == "" {
		err = errInvalidRequest
		log.Printf("invalid request body - should not be empty: product: '%s' slug: '%s' value: '%s'",
			request.Product, request.Slug, request.Value)
		return "", nil, err
	}

	filepath := folderPath + request.Product + "/" + request.Slug + ".json"
	consens, err = getConsensByPath(filepath)
	if err != nil {
		log.Println("error getting consens")
		return "", nil, err
	}

	ctx := context.WithValue(context.Background(), timestamp, now)
	if err = consentMayBeGiven(ctx, consens, request, networkNode); err != nil {
		log.Println("invalid consent")
		return "", nil, err
	}

	nodeConsens := models.NodeConsens{
		Slug:     consens.Slug,
		ExpireAt: consens.ExpireAt,
		StartAt:  consens.StartAt,
		Title:    consens.Title,
		Content:  consens.ToString(),
		Value:    request.Value,
		GivenAt:  now,
	}

	nodeConsensIndex := slices.IndexFunc(networkNode.Consens, func(c models.NodeConsens) bool {
		return c.Slug == consens.Slug
	})

	if nodeConsensIndex == -1 {
		log.Println("appending new consent")
		networkNode.Consens = append(networkNode.Consens, nodeConsens)
	} else {
		log.Println("updating given consent")
		networkNode.Consens[nodeConsensIndex] = nodeConsens
	}

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
	
	if err := sendConsensMail(networkNode, nodeConsens); err != nil {
		log.Printf("error while sending mail: %v", err)
		log.Println("continuing acceptance process...")
		err = nil
	}

	if undeclaredConsens, err = getUndeclaredConsens(request.Product, networkNode); err != nil {
		log.Println("error getting undeclared consens")
		return "", nil, err
	}

	response.Consens = make([]OutputConsens, 0, len(undeclaredConsens))
	for _, c := range undeclaredConsens {
		response.Consens = append(response.Consens, c.ToOutput())
	}

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

func consentMayBeGiven(ctx context.Context, consens SystemConsens, request AcceptanceReq, networkNode *models.NetworkNode) error {
	var (
		availableConsens []SystemConsens
		err              error
	)
	now := getTimestamp(ctx)

	if availableConsens, err = getUndeclaredConsens(request.Product, networkNode); err != nil {
		return err
	}
	availableConsensSlugs := lib.SliceMap(availableConsens, func(c SystemConsens) string { return c.Slug })

	if !lib.SliceContains(availableConsensSlugs, request.Slug) {
		return errInvalidConsensToBeGiven
	}

	if now.After(consens.ExpireAt) {
		return errConsensExpired
	}

	availableConsents := make([]string, 0)
	for _, content := range consens.Content {
		if content.InputValue != "" {
			availableConsents = append(availableConsents, content.InputValue)
		}
	}

	if !lib.SliceContains(availableConsents, request.Value) {
		return errInvalidConsentValue
	}

	return nil
}
