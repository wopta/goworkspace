package mga

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type ConsumeNetworkNodeInviteReq struct {
	InviteUid string `json:"inviteUid"`
	Password  string `json:"password"`
}

func ConsumeNetworkNodeInviteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req ConsumeNetworkNodeInviteReq
	)

	log.Println("[ConsumeNetworkNodeInviteFx] handler start -----------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	origin := r.Header.Get("Origin")

	err := json.Unmarshal(body, &req)
	if err != nil {
		log.Println("[ConsumeNetworkNodeInviteFx] error unmarshaling request body")
		return "", nil, err
	}

	log.Printf("[consumeNetworkNodeInvite] Consuming invite %s...", req.InviteUid)

	err = consumeNetworkNodeInvite(origin, req.InviteUid, req.Password)
	if err != nil {
		log.Printf("[ConsumeNetworkNodeInviteFx] error consuming invite %s: %s", req.InviteUid, err.Error())
		return "", nil, err
	}

	return "", nil, nil
}

func consumeNetworkNodeInvite(origin, inviteUid, password string) error {
	var (
		invite      NetworkNodeInvite
		networkNode *models.NetworkNode
	)

	log.Printf("consumeNetworkNodeInvite] getting invite %s from Firestore...", inviteUid)

	fireInvites := lib.GetDatasetByEnv(origin, models.InvitesCollection)
	docsnap, err := lib.GetFirestoreErr(fireInvites, inviteUid)
	if err != nil {
		log.Printf("consumeNetworkNodeInvite] error getting invite %s from Firestore", inviteUid)
		return err
	}

	err = docsnap.DataTo(&invite)
	if err != nil {
		log.Printf("consumeNetworkNodeInvite] error unmarshaling invite %s", inviteUid)
		return err
	}

	if invite.Consumed || time.Now().UTC().After(invite.Expiration) {
		log.Printf("[consumeNetworkNodeInvite] cannot consume invite with Consumed %t and Expiration %s", invite.Consumed, invite.Expiration.String())
		return errors.New("invite consumed or expired")
	}

	log.Printf("consumeNetworkNodeInvite] getting network node %s from Firestore...", invite.NetworkNodeUid)

	networkNode, err = network.GetNodeByUid(invite.NetworkNodeUid)
	if err != nil {
		log.Printf("[consumeNetworkNodeInvite] error getting network node %s from Firestore", invite.NetworkNodeUid)
		return err
	}

	userRecord, err := lib.CreateUserWithEmailAndPassword(networkNode.Mail, password, &networkNode.Uid)
	if err != nil {
		log.Printf("[consumeNetworkNodeInvite] error creating network node %s auth account", networkNode.Uid)
		return err
	}

	networkNode.AuthId = userRecord.UID
	networkNode.UpdatedDate = time.Now().UTC()

	log.Printf("[consumeNetworkInvite] updating network node %s in Firestore...", networkNode.Uid)

	fireNetworkNode := lib.GetDatasetByEnv(origin, models.NetworkNodesCollection)
	err = lib.SetFirestoreErr(fireNetworkNode, networkNode.Uid, networkNode)
	if err != nil {
		log.Printf("[consumeNetworkInvite] error update network node %s in Firestore", networkNode.Uid)
		return err
	}

	log.Printf("[consumeNetworkInvite] updating network node %s in BigQuery...", networkNode.Uid)

	networkNode.SaveBigQuery(origin)

	invite.Consumed = true
	invite.ConsumeDate = time.Now().UTC()

	log.Printf("[consumeNetworkInvite] updating invite %s in Firestore...", inviteUid)

	err = lib.SetFirestoreErr(fireInvites, inviteUid, invite)
	if err != nil {
		log.Printf("[consumeNetworkInvite] error updating invite %s in Firestore...", inviteUid)
		return err
	}

	return nil
}
