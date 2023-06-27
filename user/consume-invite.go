package user

import (
	"encoding/json"
	"errors"
	"firebase.google.com/go/v4/auth"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type ConsumeInviteReq struct {
	InviteUid string `json:"inviteUid"`
	Password  string `json:"password"`
}

func ConsumeInviteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var ConsumeInviteRequest ConsumeInviteReq

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

	docUid := ""
	collection = ""
	switch invite.Role {
	case models.UserRoleAgent:
		collection = lib.GetDatasetByEnv(origin, models.AgentCollection)
		docUid = lib.NewDoc(collection) + "_agent"
		agentRecord, err := lib.CreateUserWithEmailAndPassword(invite.Email, password, &docUid)
		if err != nil {
			return false, err
		}
		createAgent(collection, agentRecord, invite)
	case models.UserRoleAgency:
		collection = lib.GetDatasetByEnv(origin, models.AgencyCollection)
		docUid = lib.NewDoc(collection) + "_agency"
		agencyRecord, err := lib.CreateUserWithEmailAndPassword(invite.Email, password, &docUid)
		if err != nil {
			return false, err
		}
		createAgency(collection, agencyRecord, invite)
	default:
		collection = lib.GetDatasetByEnv(origin, usersCollection)
		userRecord, err := lib.CreateUserWithEmailAndPassword(invite.Email, password, &docUid)
		if err != nil {
			return false, err
		}
		createUser(collection, userRecord, invite)
	}

	// update the invite to consumed
	invite.Consumed = true
	invitesCollectionName := lib.GetDatasetByEnv(origin, invitesCollection)
	lib.SetFirestore(invitesCollectionName, invite.Uid, invite)

	log.Printf("[ConsumeInvite] Consumed invite with uid %s", invite.Uid)
	return true, nil
}

func createUser(collection string, userRecord *auth.UserRecord, invite UserInvite) {
	// create user in DB
	user := models.User{
		Mail:         invite.Email,
		Uid:          userRecord.UID,
		AuthId:       userRecord.UID,
		Role:         invite.Role,
		FiscalCode:   invite.FiscalCode,
		Name:         invite.Name,
		Surname:      invite.Surname,
		CreationDate: time.Now().UTC(),
		UpdatedDate:  time.Now().UTC(),
	}

	err := lib.SetFirestoreErr(collection, user.Uid, user)
	lib.CheckError(err)

	// update the user custom claim
	lib.SetCustomClaimForUser(user.AuthId, map[string]interface{}{
		"role": user.Role,
	})
}

func createAgent(collection string, userRecord *auth.UserRecord, invite UserInvite) {
	// create user in DB
	agent := models.Agent{
		User: models.User{
			Mail:         invite.Email,
			Uid:          userRecord.UID,
			AuthId:       userRecord.UID,
			Role:         invite.Role,
			FiscalCode:   invite.FiscalCode,
			VatCode:      invite.VatCode,
			Name:         invite.Name,
			Surname:      invite.Surname,
			CreationDate: time.Now().UTC(),
			UpdatedDate:  time.Now().UTC(),
		},
		RuiCode:         invite.RuiCode,
		RuiRegistration: invite.RuiRegistration,
	}

	err := lib.SetFirestoreErr(collection, agent.Uid, agent)
	lib.CheckError(err)

	// update the user custom claim
	lib.SetCustomClaimForUser(agent.AuthId, map[string]interface{}{
		"role": agent.Role,
	})
}

func createAgency(collection string, userRecord *auth.UserRecord, invite UserInvite) {
	// create user in DB
	agency := models.Agency{
		AuthId:          userRecord.UID,
		Uid:             userRecord.UID,
		Name:            invite.Name,
		Email:           invite.Email,
		VatCode:         invite.VatCode,
		RuiCode:         invite.RuiCode,
		RuiRegistration: invite.RuiRegistration,
		CreationDate:    time.Now().UTC(),
		UpdatedDate:     time.Now().UTC(),
	}

	err := lib.SetFirestoreErr(collection, agency.Uid, agency)
	lib.CheckError(err)

	// update the user custom claim
	lib.SetCustomClaimForUser(agency.AuthId, map[string]interface{}{
		"role": invite.Role,
	})
}
