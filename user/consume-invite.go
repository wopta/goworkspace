package user

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func ConsumeInviteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var ConsumeInviteRequest ConsumeInviteRequest

	reqBytes := lib.ErrorByte(ioutil.ReadAll(r.Body))
	json.Unmarshal(reqBytes, &ConsumeInviteRequest)

	if ok, _ := ConsumeInvite(ConsumeInviteRequest.InviteUid, ConsumeInviteRequest.Password, r.Header.Get("Origin")); ok {
		return `{"success": true}`, `{"success": true}`, nil
	}

	return `{"success": false}`, `{"success": false}`, nil
}

func ConsumeInvite(inviteUid, password, origin string) (bool, error) {
	log.Printf("[ConsumeInvite] Consuming invite %s", inviteUid)

	// Get the invite
	collection := lib.GetDatasetByEnv(origin, invitesCollection)
	docSnapshot, err := lib.GetFirestoreErr(collection, inviteUid)
	if err != nil {
		return false, err
	}

	var invite UserInvite
	err = docSnapshot.DataTo(&invite)
	if err != nil {
		return false, err
	}

	// Check if invite is not consumed nor expired
	if invite.Consumed || time.Now().UTC().After(invite.Expiration) {
		return false, errors.New("invite consumed or expired")
	}

	usersCollectionName := lib.GetDatasetByEnv(origin, usersCollection)

	// Create the user in auth with the invite data
	userRecord, err := lib.CreateUserWithEmailAndPassword(invite.Email, password, nil)
	if err != nil {
		return false, err
	}

	// create user in DB
	user := models.User{
		Mail:       invite.Email,
		Uid:        userRecord.UID,
		AuthId:     userRecord.UID,
		Role:       invite.Role,
		FiscalCode: invite.FiscalCode,
		Name:       invite.Name,
		Surname:    invite.Surname,
	}

	err = lib.SetFirestoreErr(usersCollectionName, user.Uid, user)
	if err != nil {
		return false, err
	}

	// update the user custom claim
	lib.SetCustomClaimForUser(user.AuthId, map[string]interface{}{
		"role": user.Role,
	})

	// update the invite to consumed
	invite.Consumed = true
	invitesCollectionName := lib.GetDatasetByEnv(origin, invitesCollection)
	lib.SetFirestore(invitesCollectionName, invite.Uid, invite)

	log.Printf("[ConsumeInvite] Consumed invite with uid %s", invite.Uid)
	return true, nil
}

type ConsumeInviteRequest struct {
	InviteUid string `json:"inviteUid"`
	Password  string `json:"password"`
}
