package consens

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

type AcceptanceReq struct {
	Slug    string            `json:"slug"`
	Product string            `json:"product"`
	Answers map[string]string `json:"answers"`
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
			log.ErrorF("error: %v", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()
	log.AddPrefix("[AcceptanceFx] ")
	log.Println("Handler start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.ErrorF("error extracting authToken")
		return "", nil, err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	networkNode := network.GetNetworkNodeByUid(authToken.UserID)
	if networkNode == nil {
		log.ErrorF("error getting networkNode")
		err = errNetworkNodeNotFound
		return "", nil, err
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()
	if err != nil {
		log.ErrorF("error decoding request body")
		return "", nil, err
	}

	if request.Product == "" || request.Slug == "" || len(request.Answers) == 0 {
		err = errInvalidRequest
		log.Printf("invalid request body - should not be empty: product: '%s' slug: '%s' value: '%+v'",
			request.Product, request.Slug, request.Answers)
		return "", nil, err
	}

	filepath := folderPath + request.Product + "/" + request.Slug + ".json"
	consens, err = getConsensByPath(filepath)
	if err != nil {
		log.ErrorF("error getting consens")
		return "", nil, err
	}
	consens = enrichConsens(consens, networkNode)
	log.Printf("found consens: %+v", consens)

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
		Subtitle: consens.Subtitle,
		Content:  ContentToString(consens.Content, request.Answers, true),
		Answers:  request.Answers,
		GivenAt:  now,
	}
	log.Printf("constructing consent: %+v", nodeConsens)

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

	log.Println("saving networkNode to DBs...")
	if err := networkNode.SaveFirestore(); err != nil {
		log.ErrorF("error saving networkNode in firestore")
		return "", nil, err
	}
	if err := networkNode.SaveBigQuery(); err != nil {
		log.ErrorF("error saving networkNode in bigquery")
		return "", nil, err
	}

	audit := NodeConsensAudit{
		NetworkNodeUid:  networkNode.Uid,
		Name:            networkNode.GetName(),
		RuiCode:         networkNode.GetRuiCode(),
		RuiRegistration: networkNode.GetRuiRegistration(),
		FiscalCode:      networkNode.GetFiscalCode(),
		VatCode:         networkNode.GetVatCode(),
		Slug:            nodeConsens.Slug,
		Title:           nodeConsens.Title,
		Subtitle:        nodeConsens.Subtitle,
		Content:         nodeConsens.Content,
		Answers:         request.Answers,
		GivenAt:         now,
	}
	log.Println("saving consens audit to BigQuery...")
	if err := audit.Save(); err != nil {
		log.ErrorF("error saving consens audit")
		return "", nil, err
	}

	log.Println("sending consens mail to networkNode...")
	if err := sendConsensMail(networkNode, consens, nodeConsens); err != nil {
		log.ErrorF("error while sending mail: %v", err)
		log.Println("continuing acceptance process...")
		err = nil
	}

	log.Println("fetching undeclared consens...")
	if undeclaredConsens, err = getUndeclaredConsens(request.Product, networkNode); err != nil {
		log.ErrorF("error getting undeclared consens")
		return "", nil, err
	}

	response.Consens = make([]OutputConsens, 0, len(undeclaredConsens))
	for _, c := range undeclaredConsens {
		response.Consens = append(response.Consens, c.ToOutput())
	}

	if responseBytes, err = json.Marshal(response); err != nil {
		log.ErrorF("error marshalling response")
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

	availableConsents := make(map[string][]string)
	for _, content := range consens.Content {
		if content.InputValue != "" {
			if _, ok := availableConsents[content.InputName]; !ok {
				availableConsents[content.InputName] = make([]string, 0)
			}
			availableConsents[content.InputName] = append(availableConsents[content.InputName], content.InputValue)
		}
	}

	for key, val := range availableConsents {
		if v, ok := request.Answers[key]; !ok || !lib.SliceContains(val, v) {
			return errInvalidConsentValue
		}
	}

	return nil
}
