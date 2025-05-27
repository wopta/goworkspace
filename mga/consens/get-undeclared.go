package consens

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"slices"
	"sort"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

func GetUndeclaredConsensFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err           error
		response      ConsensResp
		consens       []SystemConsens
		responseBytes []byte
	)

	defer func() {
		if err != nil {
			log.ErrorF("error: %v", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()

	log.AddPrefix("[GetUndeclaredConsensFx] ")
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

	product := r.URL.Query().Get("product")

	log.Println("fetching undeclared consens...")
	if consens, err = getUndeclaredConsens(product, networkNode); err != nil {
		log.ErrorF("error getting undeclared consens")
		return "", nil, err
	}

	response.Consens = make([]OutputConsens, 0, len(consens))
	for _, c := range consens {
		response.Consens = append(response.Consens, c.ToOutput())
	}

	if responseBytes, err = json.Marshal(response); err != nil {
		log.ErrorF("error marshalling response")
		return "", nil, err
	}

	return string(responseBytes), response, err
}

func getUndeclaredConsens(product string, networkNode *models.NetworkNode) ([]SystemConsens, error) {
	var (
		err               error
		allProductConsens []SystemConsens
		undeclaredConsens = make([]SystemConsens, 0)
		now               = time.Now().UTC()
	)

	log.Printf("retrieving product %s consens...", product)
	if allProductConsens, err = getProductConsens(product); err != nil {
		return nil, err
	}

	for i, c := range allProductConsens {
		allProductConsens[i] = enrichConsens(c, networkNode)
	}

	log.Printf("found a total of %d consens for product %s", len(allProductConsens), product)

	nodeConsensList := make([]string, 0, len(networkNode.Consens))
	for _, c := range networkNode.Consens {
		nodeConsensList = append(nodeConsensList, c.Slug)
	}

	ctx := context.WithValue(context.Background(), timestamp, now)

	for _, c := range allProductConsens {
		if lib.SliceContains(nodeConsensList, c.Slug) {
			log.Printf("consent given for consens %s", c.Slug)
			continue
		}

		log.Println("checking consens configuration...")
		strategy, err := newConsensStrategy(c, *networkNode)
		if err != nil {
			return nil, err
		}
		log.Printf("executing strategy '%s'...", c.Strategy)
		valid, err := strategy.Check(ctx)
		if err != nil {
			return nil, err
		}

		if !valid {
			log.Printf("adding consens '%s' to undeclared list", c.Slug)
			undeclaredConsens = append(undeclaredConsens, c)
		}
	}

	log.Printf("found a total of %d undeclared consens", len(undeclaredConsens))

	return undeclaredConsens, err
}

func getProductConsens(product string) ([]SystemConsens, error) {
	var (
		path              = folderPath + product
		fileList          []string
		allProductConsens []SystemConsens
		err               error
	)
	switch os.Getenv("env") {
	case env.Local:
		fileList, err = lib.ListLocalFolderContent(path)
	default:
		fileList, err = lib.ListGoogleStorageFolderContent(path)
	}

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	log.Printf("found a total of %d files", len(fileList))

	if len(fileList) == 0 {
		return nil, nil
	}

	switch product {
	case allProducts:
		allProductConsens, err = getAllAvailableConsens(fileList)
	default:
		allProductConsens, err = getLastestNeededConsens(fileList)
	}
	if err != nil {
		return nil, err
	}

	return allProductConsens, nil
}

func getAllAvailableConsens(fileList []string) ([]SystemConsens, error) {
	var (
		allProductConsens = make([]SystemConsens, 0, len(fileList))
		err               error
	)

	for _, file := range fileList {
		var (
			fileBytes []byte
			c         SystemConsens
		)
		if fileBytes, err = lib.GetFilesByEnvV2(file); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(fileBytes, &c); err != nil {
			return nil, err
		}
		allProductConsens = append(allProductConsens, c)
	}

	return allProductConsens, err
}

func getLastestNeededConsens(fileList []string) ([]SystemConsens, error) {
	var (
		latestNeededConsens = make([]SystemConsens, 0, 1)
		allProductConsens   MultipleConsens
		temp                []SystemConsens
		err                 error
	)

	if temp, err = getAllAvailableConsens(fileList); err != nil {
		return nil, err
	}

	allProductConsens = append(allProductConsens, temp...)
	sort.Sort(allProductConsens)

	latestNeededConsens = append(latestNeededConsens, allProductConsens[0])

	return latestNeededConsens, err
}

type MultipleConsens []SystemConsens

func (c MultipleConsens) Len() int           { return len(c) }
func (c MultipleConsens) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c MultipleConsens) Less(i, j int) bool { return c[i].StartAt.After(c[j].StartAt) }

func enrichConsens(consens SystemConsens, networkNode *models.NetworkNode) SystemConsens {
	regexNodeName := regexp.MustCompile("{{NODE_NAME}}")
	regexNodeFiscalCode := regexp.MustCompile("{{NODE_FISCALCODE}}")
	regexNodeRuiSection := regexp.MustCompile("{{NODE_RUI_SECTION}}")
	regexNodeRuiCode := regexp.MustCompile("{{NODE_RUI_CODE}}")
	regexNodeRuiRegistrationDate := regexp.MustCompile("{{NODE_RUI_REGISTRATION_DATE}}")
	regexNodeDesignation := regexp.MustCompile("{{NODE_DESIGNATION}}")

	name := networkNode.GetName()
	fiscalCode := networkNode.GetFiscalCode()
	ruiSection := networkNode.GetRuiSection()
	ruiCode := networkNode.GetRuiCode()
	ruiRegistration := networkNode.GetRuiRegistration()
	designation := networkNode.Designation
	if slices.Contains([]string{models.AgencyNetworkNodeType, models.BrokerNetworkNodeType}, networkNode.Type) {
		name = networkNode.GetManagerName()
		fiscalCode = networkNode.GetManagerFiscalCode()
		ruiSection = networkNode.GetManagerRuiSection()
		ruiCode = networkNode.GetManagerRuiCode()
		ruiRegistration = networkNode.GetManagerRuiResgistration()
	}

	for j, cont := range consens.Content {
		text := cont.Text
		text = regexNodeName.ReplaceAllString(text, name)
		text = regexNodeFiscalCode.ReplaceAllString(text, fiscalCode)
		text = regexNodeRuiSection.ReplaceAllString(text, ruiSection)
		text = regexNodeRuiCode.ReplaceAllString(text, ruiCode)
		text = regexNodeRuiRegistrationDate.ReplaceAllString(text, ruiRegistration.Format("02/01/2006"))
		text = regexNodeDesignation.ReplaceAllString(text, designation)
		consens.Content[j].Text = text
	}

	return consens
}
