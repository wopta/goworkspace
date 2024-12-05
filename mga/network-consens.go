package mga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type key string

const timestamp = key("timestamp")

var (
	errNetworkNodeNotFound = errors.New("network node not found")
	errPartnershipNode     = errors.New("partnership node does not have rui registration nor consens")
	errStrategyNotFound    = errors.New("strategy not found")
)

const (
	ruiSectionE         = "E"
	ruiSectionEStrategy = "rui_section_e"
	allNodesStrategy    = "all_nodes"
)

type NetworkConsens struct {
	Slug        string    `json:"slug"`
	ExpireAt    time.Time `json:"expireAt"`
	StartAt     time.Time `json:"startAt"`
	AvailableAt time.Time `json:"availableAt"`
	Strategy    string    `json:"strategy"`
}

// TODO: add the fields to the correct struct
type NodeWithConsens struct {
	models.NetworkNode
	Consens []NetworkConsens `json:"networkConsens"`
}

func GetUndeclaredConsensFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err           error
		consens       []NetworkConsens
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

	product := chi.URLParam(r, "product")

	consens, err = getUndeclaredConsens(product, &tempNode)

	responseBytes, err = json.Marshal(consens)
	if err != nil {
		return "", nil, err
	}

	return string(responseBytes), consens, err
}

func getUndeclaredConsens(product string, networkNode *NodeWithConsens) ([]NetworkConsens, error) {
	var (
		err               error
		fileList          = make([]string, 0)
		folderPath        = fmt.Sprintf("consens/network/%s", product)
		allProductConsens = make([]NetworkConsens, 0)
		undeclaredConsens = make([]NetworkConsens, 0)
		now               = time.Now().UTC()
	)

	switch os.Getenv("env") {
	case "local":
		fileList, err = lib.ListLocalFolderContent(folderPath)
	default:
		fileList, err = lib.ListGoogleStorageFolderContent(folderPath)
	}

	if err != nil {
		return nil, err
	}

	for _, file := range fileList {
		var (
			fileBytes []byte
			c         NetworkConsens
		)
		fileBytes, err = lib.GetFilesByEnvV2(file)
		if err != nil {
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

type NeedConsensAlgorithm interface {
	Check(context.Context) (bool, error)
}

type RuisectionE struct {
	consens NetworkConsens
	node    NodeWithConsens
}

func (w *RuisectionE) Check(ctx context.Context) (bool, error) {
	var ruiSection string
	switch w.node.Type {
	case models.AgentNetworkNodeType:
		ruiSection = w.node.Agent.RuiSection
	case models.AgencyNetworkNodeType:
		ruiSection = w.node.Agency.RuiSection
	case models.BrokerNetworkNodeType:
		ruiSection = w.node.Broker.RuiSection
	case models.AreaManagerNetworkNodeType:
		ruiSection = w.node.AreaManager.RuiSection
	case models.PartnershipNetworkNodeType:
		return false, errPartnershipNode
	}

	if !strings.EqualFold(ruiSection, ruiSectionE) {
		return true, nil
	}

	now := getTimestamp(ctx)

	if now.Before(w.consens.StartAt) {
		return true, nil
	}

	if now.Before(w.consens.ExpireAt) {
		return false, nil
	}

	return true, nil
}

type AllNodes struct {
	consens NetworkConsens
	node    NodeWithConsens
}

func (w *AllNodes) Check(ctx context.Context) (bool, error) {
	now := getTimestamp(ctx)

	if now.Before(w.consens.StartAt) {
		return true, nil
	}

	if now.Before(w.consens.ExpireAt) {
		return false, nil
	}

	return true, nil
}

func getTimestamp(ctx context.Context) time.Time {
	if rawTime := ctx.Value(timestamp); rawTime != nil {
		return (rawTime).(time.Time)
	}
	return time.Time{}
}

func newConsensStrategy(consens NetworkConsens, node NodeWithConsens) (NeedConsensAlgorithm, error) {
	switch consens.Strategy {
	case ruiSectionEStrategy:
		return &RuisectionE{
			consens: consens,
			node:    node,
		}, nil
	case allNodesStrategy:
		return &AllNodes{
			consens: consens,
			node:    node,
		}, nil
	}
	return nil, errStrategyNotFound
}
