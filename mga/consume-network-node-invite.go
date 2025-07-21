package mga

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

type ConsumeNetworkNodeInviteReq struct {
	InviteUid string `json:"inviteUid"`
	Password  string `json:"password"`
}

func ConsumeNetworkNodeInviteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req ConsumeNetworkNodeInviteReq
	)

	log.AddPrefix("[ConsumeNetworkNodeInviteFx] ")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(body, &req)
	if err != nil {
		log.ErrorF("error unmarshaling request body")
		return "", "", err
	}

	log.Printf("Consuming invite %s...", req.InviteUid)

	err = consumeNetworkNodeInvite(req.InviteUid, req.Password)
	if err != nil {
		log.ErrorF("error consuming invite %s: %s", req.InviteUid, err.Error())
		return "", "", err
	}

	log.Println("Handler end -------------------------------------------------")

	return "{}", "", nil
}

func consumeNetworkNodeInvite(inviteUid, password string) error {
	var (
		invite      NetworkNodeInvite
		networkNode *models.NetworkNode
	)
	log.AddPrefix("ConsumeNetworkNodeInvite")
	defer log.PopPrefix()
	log.Printf("getting invite %s from Firestore...", inviteUid)

	fireInvites := models.InvitesCollection
	docsnap, err := lib.GetFirestoreErr(fireInvites, inviteUid)
	if err != nil {
		log.ErrorF("error getting invite %s from Firestore", inviteUid)
		return err
	}

	err = docsnap.DataTo(&invite)
	if err != nil {
		log.ErrorF("error unmarshaling invite %s", inviteUid)
		return err
	}

	if invite.Consumed || time.Now().UTC().After(invite.Expiration) {
		log.Printf("cannot consume invite with Consumed %t and Expiration %s", invite.Consumed, invite.Expiration.String())
		return errors.New("invite consumed or expired")
	}

	log.Printf("getting network node %s from Firestore...", invite.NetworkNodeUid)

	networkNode, err = network.GetNodeByUidErr(invite.NetworkNodeUid)
	if err != nil {
		log.ErrorF("error getting network node %s from Firestore", invite.NetworkNodeUid)
		return err
	}
	if networkNode == nil {
		return fmt.Errorf("error no node found: %v", invite.NetworkNodeUid)
	}

	userRecord, err := lib.CreateUserWithEmailAndPassword(networkNode.Mail, password, &networkNode.Uid)
	if err != nil {
		log.ErrorF("error creating network node %s auth account", networkNode.Uid)
		return err
	}

	networkNode.AuthId = userRecord.UID
	networkNode.UpdatedDate = time.Now().UTC()

	log.Printf("updating network node %s in Firestore...", networkNode.Uid)

	fireNetworkNode := models.NetworkNodesCollection
	err = lib.SetFirestoreErr(fireNetworkNode, networkNode.Uid, networkNode)
	if err != nil {
		log.ErrorF("error update network node %s in Firestore", networkNode.Uid)
		return err
	}

	log.Printf("updating network node %s in BigQuery...", networkNode.Uid)

	networkNode.SaveBigQuery()

	invite.Consumed = true
	invite.ConsumeDate = time.Now().UTC()

	log.Printf("updating invite %s in Firestore...", inviteUid)

	err = lib.SetFirestoreErr(fireInvites, inviteUid, invite)
	if err != nil {
		log.ErrorF("error updating invite %s in Firestore...", inviteUid)
		return err
	}

	lib.SetCustomClaimForUser(networkNode.AuthId, map[string]interface{}{
		"isNetworkNode": true,
		"role":          networkNode.Role,
		"type":          networkNode.Type,
	})

	return nil
}
