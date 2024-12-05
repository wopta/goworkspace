package consens

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/network"
)

type key string

const timestamp = key("timestamp")

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

	product := r.URL.Query().Get("product")

	if consens, err = getUndeclaredConsens(product, &tempNode); err != nil {
		return "", nil, err
	}

	response.Consens = consens

	if responseBytes, err = json.Marshal(response); err != nil {
		return "", nil, err
	}

	return string(responseBytes), response, err
}

func getUndeclaredConsens(product string, networkNode *NodeWithConsens) ([]SystemConsens, error) {
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

	return undeclaredConsens, err
}
