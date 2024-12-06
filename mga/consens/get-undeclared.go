package consens

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
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
			log.Printf("error: %v", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()

	log.SetPrefix("[GetUndeclaredConsensFx] ")
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

	product := r.URL.Query().Get("product")

	if consens, err = getUndeclaredConsens(product, networkNode); err != nil {
		log.Println("error getting undeclared consens")
		return "", nil, err
	}

	response.Consens = make([]OutputConsens, 0, len(consens))
	for _, c := range consens {
		response.Consens = append(response.Consens, c.ToOutput())
	}

	if responseBytes, err = json.Marshal(response); err != nil {
		log.Println("error marshalling response")
		return "", nil, err
	}

	return string(responseBytes), response, err
}

func getUndeclaredConsens(product string, networkNode *models.NetworkNode) ([]SystemConsens, error) {
	var (
		err               error
		fileList          = make([]string, 0)
		path              = folderPath + product
		allProductConsens = make([]SystemConsens, 0)
		undeclaredConsens = make([]SystemConsens, 0)
		now               = time.Now().UTC()
	)

	switch os.Getenv("env") {
	case "local":
		fileList, err = lib.ListLocalFolderContent(path)
	default:
		fileList, err = lib.ListGoogleStorageFolderContent(path)
	}

	if err != nil {
		return nil, err
	}

	log.Printf("found a total of %d consens", len(fileList))

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

		strategy, err := newConsensStrategy(c, *networkNode)
		if err != nil {
			return nil, err
		}
		valid, err := strategy.Check(ctx)
		if err != nil {
			return nil, err
		}

		if !valid {
			undeclaredConsens = append(undeclaredConsens, c)
		}
	}

	log.Printf("found a total of %d undeclared consens", len(undeclaredConsens))

	return undeclaredConsens, err
}
