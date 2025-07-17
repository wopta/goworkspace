package user

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"firebase.google.com/go/v4/auth"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type ConsumeInviteReq struct {
	InviteUid string `json:"inviteUid"`
	Password  string `json:"password"`
}

func ConsumeInviteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var ConsumeInviteRequest ConsumeInviteReq

	log.AddPrefix("ConsumeInviteFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	reqBytes := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(reqBytes, &ConsumeInviteRequest)
	if err != nil {
		log.ErrorF("error unmarshaling request: %s", err.Error())
		return "", nil, err
	}

	_, err = ConsumeInvite(ConsumeInviteRequest.InviteUid, ConsumeInviteRequest.Password)
	if err != nil {
		log.ErrorF("error consuming invite: %s", err.Error())
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}

func ConsumeInvite(inviteUid, password string) (bool, error) {
	log.AddPrefix("ConsumeInvite")
	defer log.PopPrefix()
	log.Printf("Consuming invite %s", inviteUid)

	// Get the invite
	collection := lib.InvitesCollection
	docSnapshot, err := lib.GetFirestoreErr(collection, inviteUid)
	if err != nil {
		log.ErrorF("error retrieving data from firestore: %s", err.Error())
		return false, err
	}

	var invite UserInvite
	err = docSnapshot.DataTo(&invite)
	if err != nil {
		log.ErrorF("error unmarshaling data: %s", err.Error())
		return false, err
	}

	// Check if invite is not consumed nor expired
	if invite.Consumed || time.Now().UTC().After(invite.Expiration) {
		log.ErrorF("cannot consume invite with Consumed %t and Expiration %s", invite.Consumed, invite.Expiration.String())
		return false, errors.New("invite consumed or expired")
	}

	docUid := ""
	collection = ""
	switch invite.Role {
	case models.UserRoleAgent:
		collection = models.AgentCollection
		docUid = lib.NewDoc(collection) + "_agent"
		agentRecord, err := lib.CreateUserWithEmailAndPassword(invite.Email, password, &docUid)
		if err != nil {
			return false, err
		}
		createAgent(collection, agentRecord, invite)
	case models.UserRoleAgency:
		collection = models.AgencyCollection
		docUid = lib.NewDoc(collection) + "_agency"
		agencyRecord, err := lib.CreateUserWithEmailAndPassword(invite.Email, password, &docUid)
		if err != nil {
			return false, err
		}
		createAgency(collection, agencyRecord, invite)
	default:
		collection = lib.UserCollection
		userRecord, err := lib.CreateUserWithEmailAndPassword(invite.Email, password, &docUid)
		if err != nil {
			return false, err
		}
		createUser(collection, userRecord, invite)
	}

	// update the invite to consumed
	invite.Consumed = true
	invitesCollectionName := lib.InvitesCollection
	lib.SetFirestore(invitesCollectionName, invite.Uid, invite)

	log.Printf("Consumed invite with uid %s", invite.Uid)
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
		VatCode:      invite.VatCode,
		Name:         invite.Name,
		Surname:      invite.Surname,
		CreationDate: time.Now().UTC(),
		UpdatedDate:  time.Now().UTC(),
	}

	err := lib.SetFirestoreErr(collection, user.Uid, user)
	lib.CheckError(err)

	err = user.BigquerySave()
	lib.CheckError(err)

	// update the user custom claim
	lib.SetCustomClaimForUser(user.AuthId, map[string]interface{}{
		"role": user.Role,
	})
}

func createAgent(collection string, userRecord *auth.UserRecord, invite UserInvite) {
	// create user in DB
	for productIndex, _ := range invite.Products {
		invite.Products[productIndex].IsAgentActive = true
	}

	agent := models.Agent{
		Mail:            invite.Email,
		Uid:             userRecord.UID,
		AuthId:          userRecord.UID,
		Code:            invite.Code,
		Role:            invite.Role,
		FiscalCode:      invite.FiscalCode,
		VatCode:         invite.VatCode,
		Name:            invite.Name,
		Surname:         invite.Surname,
		CreationDate:    time.Now().UTC(),
		UpdatedDate:     time.Now().UTC(),
		RuiCode:         invite.RuiCode,
		RuiSection:      invite.RuiSection,
		RuiRegistration: invite.RuiRegistration,
		Products:        invite.Products,
	}

	err := lib.SetFirestoreErr(collection, agent.Uid, agent)
	lib.CheckError(err)

	err = agent.BigquerySave()
	lib.CheckError(err)

	// update the user custom claim
	lib.SetCustomClaimForUser(agent.AuthId, map[string]interface{}{
		"role": agent.Role,
	})
}

func createAgency(collection string, userRecord *auth.UserRecord, invite UserInvite) {
	// create user in DB
	for productIndex, _ := range invite.Products {
		invite.Products[productIndex].IsAgencyActive = true
	}

	agency := models.Agency{
		AuthId:          userRecord.UID,
		Uid:             userRecord.UID,
		Code:            invite.Code,
		Name:            invite.Name,
		Email:           invite.Email,
		VatCode:         invite.VatCode,
		RuiCode:         invite.RuiCode,
		RuiSection:      invite.RuiSection,
		RuiRegistration: invite.RuiRegistration,
		CreationDate:    time.Now().UTC(),
		UpdatedDate:     time.Now().UTC(),
		Products:        invite.Products,
	}

	err := lib.SetFirestoreErr(collection, agency.Uid, agency)
	lib.CheckError(err)

	err = agency.BigquerySave()
	lib.CheckError(err)

	// update the user custom claim
	lib.SetCustomClaimForUser(agency.AuthId, map[string]interface{}{
		"role": invite.Role,
	})
}
