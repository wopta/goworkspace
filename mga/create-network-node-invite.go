package mga

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib/log"
	"io"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type CreateNetworkNodeInviteRequest struct {
	NetworkNodeUid string `json:"networkNodeUid"`
}

type NetworkNodeInvite struct {
	Uid            string    `json:"uid,omitempty" firestore:"uid,omitempty"`
	CreatorUid     string    `json:"creatorUid,omitempty" firestore:"creatorUid,omitempty"`
	Consumed       bool      `json:"consumed,omitempty" firestore:"consumed,omitempty"`
	NetworkNodeUid string    `json:"networkNodeUid,omitempty" firestore:"networkNodeUid,omitempty"`
	CreationDate   time.Time `json:"creationDate,omitempty" firestore:"creationDate,omitempty"`
	ConsumeDate    time.Time `json:"consumeDate,omitempty" firestore:"consumeDate,omitempty"`
	Expiration     time.Time `json:"expiration,omitempty" firestore:"expiration,omitempty"`
}

func CreateNetworkNodeInviteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req         CreateNetworkNodeInviteRequest
		networkNode *models.NetworkNode
	)

	log.AddPrefix("[CreateNetworkNodeInviteFx] ")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	origin := r.Header.Get("Origin")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(body, &req)
	if err != nil {
		log.ErrorF("error unmarshalling body")
		return "", "", err
	}

	token := r.Header.Get("Authorization")

	log.Printf("getting creatorUid from token %s", token)

	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("invalid JWT %s", token)
		return "", "", err
	}

	log.Printf("getting network node %s from Firestore...", req.NetworkNodeUid)

	networkNode, err = network.GetNodeByUid(req.NetworkNodeUid)
	if err != nil {
		log.ErrorF("error getting network node %s from Firestore...", req.NetworkNodeUid)
		return "", "", err
	}

	log.Printf("generating invite for network node %s", req.NetworkNodeUid)

	inviteUid, err := createNetworkNodeInvite(origin, networkNode.Uid, authToken.UserID)
	if err != nil {
		log.ErrorF("error generating invite for network node %s", req.NetworkNodeUid)
		return "", "", err
	}

	log.Printf("sending network node invite mail to %s", networkNode.Mail)

	mail.SendInviteMail(inviteUid, networkNode.Mail, true)

	log.Printf("network node invite mail sent to %s", networkNode.Mail)
	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}

func createNetworkNodeInvite(origin, networkNodeUid, creatorUid string) (string, error) {
	fireInvite := lib.GetDatasetByEnv(origin, models.InvitesCollection)
	inviteUid := lib.NewDoc(fireInvite)

	invite := NetworkNodeInvite{
		Uid:            inviteUid,
		CreatorUid:     creatorUid,
		Consumed:       false,
		NetworkNodeUid: networkNodeUid,
		CreationDate:   time.Now().UTC(),
		Expiration:     time.Now().UTC().Add(time.Hour * 168),
	}
	log.AddPrefix("CreateNetworkNodeInvite")
	defer log.PopPrefix()
	log.Printf("saving network node invite %s to Firestore...", inviteUid)

	err := lib.SetFirestoreErr(fireInvite, inviteUid, invite)
	if err != nil {
		log.ErrorF("error saving network node invite %s to Firestore", inviteUid)
		return "", err
	}

	log.Printf("network node invite %s saved", inviteUid)
	return inviteUid, nil
}
