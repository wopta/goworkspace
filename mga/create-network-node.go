package mga

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type CreateNetworkNodeRequest struct {
	Code        string                  `json:"code"`
	Type        string                  `json:"type"`
	Role        string                  `json:"role"`
	NetworkCode string                  `json:"networkCode"`
	ManagerUid  string                  `json:"managerUid,omitempty"`
	ParentUid   string                  `json:"parentUid,omitempty"`
	Agent       *models.AgentNode       `json:"agent,omitempty"`
	Agency      *models.AgencyNode      `json:"agency,omitempty"`
	Broker      *models.AgencyNode      `json:"broker,omitempty"`
	Partnership *models.PartnershipNode `json:"partnership,omitempty"`
	Products    []models.Product        `json:"products,omitempty"`
}

func CreateNetworkNodeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		request CreateNetworkNodeRequest
		err     error
	)

	log.Println("[CreateNetworkNodeFx] Handler start -------------------------")

	origin := r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("[CreateNetworkNodeFx] request body: %s", string(body))
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("[CreateNetworkNodeFx] error unmarshaling request: %s", err.Error())
		return "", "", err
	}

	node := createNetworkNode(request, origin)

	fireNetwork := lib.GetDatasetByEnv(origin, models.NetworkNodesCollection)
	err = lib.SetFirestoreErr(fireNetwork, node.Uid, node)
	if err != nil {
		log.Printf("[CreateNetworkNodeFx] error saving node to firestore: %s", err.Error())
		return "", "", err
	}

	node.SaveBigQuery(origin)

	log.Println("[CreateNetworkNodeFx] network node successfully created!")

	return "", "", err
}

// TODO: mode to network domain
func createNetworkNode(request CreateNetworkNodeRequest, origin string) *models.NetworkNode {
	uid := lib.NewDoc(models.NetworkNodesCollection)
	now := time.Now().UTC()

	log.Printf("[createNetworkNode] creating node with uid %s", uid)

	return &models.NetworkNode{
		Uid:          uid,
		Code:         request.Code,
		Type:         request.Type,
		Role:         request.Role,
		NetworkCode:  request.NetworkCode,
		NetworkUid:   request.NetworkCode,
		Agent:        request.Agent,
		Agency:       request.Agency,
		Broker:       request.Broker,
		Partnership:  request.Partnership,
		Products:     request.Products,
		IsActive:     true,
		CreationDate: now,
		UpdatedDate:  now,
	}
}
