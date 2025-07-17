package user

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type CreateInviteRequest struct {
	Role            string           `json:"role"`
	Email           string           `json:"email"`
	FiscalCode      string           `json:"fiscalCode,omitempty"`
	VatCode         string           `json:"vatCode,omitempty"`
	Name            string           `json:"name,omitempty"`
	Surname         string           `json:"surname,omitempty"`
	RuiCode         string           `json:"ruiCode,omitempty"`
	RuiSection      string           `json:"ruiSection,omitempty"`
	RuiRegistration time.Time        `json:"ruiRegistration"`
	Code            string           `json:"code"`
	Products        []models.Product `json:"products,omitempty"`
}

func (c *CreateInviteRequest) Normalize() {
	c.FiscalCode = lib.ToUpper(c.FiscalCode)
	c.VatCode = lib.TrimSpace(c.VatCode)
	c.Name = lib.ToUpper(c.Name)
	c.Surname = lib.ToUpper(c.Surname)
	c.Email = lib.ToUpper(c.Email)
}

type UserInvite struct {
	FiscalCode      string           `json:"fiscalCode,omitempty" firestore:"fiscalCode,omitempty"`
	VatCode         string           `json:"vatCode,omitempty" firestore:"vatCode,omitempty"`
	Name            string           `json:"name,omitempty" firestore:"name,omitempty"`
	Surname         string           `json:"surname,omitempty" firestore:"surname,omitempty"`
	Role            string           `json:"role,omitempty" firestore:"role,omitempty"`
	Email           string           `json:"email,omitempty" firestore:"email,omitempty"`
	Uid             string           `json:"uid,omitempty" firestore:"uid,omitempty"`
	CreatorUid      string           `json:"creatorUid,omitempty" firestore:"creatorUid,omitempty"`
	Consumed        bool             `json:"consumed" firestore:"consumed"`
	RuiCode         string           `json:"ruiCode,omitempty" firestore:"ruiCode,omitempty"`
	RuiSection      string           `json:"ruiSection,omitempty" firestore:"ruiSection,omitempty"`
	RuiRegistration time.Time        `json:"ruiRegistration" firestore:"ruiRegistration"`
	Code            string           `json:"code" firestore:"code"`
	Expiration      time.Time        `json:"expiration,omitempty" firestore:"expiration,omitempty"`
	Products        []models.Product `json:"products,omitempty" firestore:"products,omitempty"`
}

func CreateInviteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var createInviteRequest CreateInviteRequest

	log.AddPrefix("[CreateInviteFx] ")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	reqBytes := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(reqBytes, &createInviteRequest)
	lib.CheckError(err)

	createInviteRequest.Normalize()

	creatorUid, err := lib.GetUserIdFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Println("Invalid auth token")
		return "", nil, err
	}

	inviteUid, err := CreateInvite(createInviteRequest, creatorUid)
	if err != nil {
		log.ErrorF("error: %s", err.Error())
		return "", nil, err
	}

	mail.SendInviteMail(inviteUid, createInviteRequest.Email, false)

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}

func CreateInvite(inviteRequest CreateInviteRequest, creatorUid string) (string, error) {
	log.AddPrefix("CreateInvite")
	defer log.PopPrefix()
	log.Printf("Creating invite for user %s with role %s", inviteRequest.Email, inviteRequest.Role)

	collectionName := lib.InvitesCollection
	inviteUid := lib.NewDoc(collectionName)

	oneWeek := time.Hour * 168
	inviteExpiration := time.Now().UTC().Add(oneWeek)

	roles := models.GetAllRoles()
	var userRole string
	for _, role := range roles {
		if strings.EqualFold(inviteRequest.Role, role) {
			userRole = role
			break
		}
	}

	if userRole == "" {
		log.Println("forbidden role")
		return "", errors.New("forbidden role")
	}

	invite := UserInvite{
		Name:            inviteRequest.Name,
		Surname:         inviteRequest.Surname,
		VatCode:         inviteRequest.VatCode,
		FiscalCode:      inviteRequest.FiscalCode,
		Email:           inviteRequest.Email,
		Role:            userRole,
		Expiration:      inviteExpiration,
		Uid:             inviteUid,
		CreatorUid:      creatorUid,
		RuiCode:         inviteRequest.RuiCode,
		RuiSection:      inviteRequest.RuiSection,
		RuiRegistration: inviteRequest.RuiRegistration,
		Code:            inviteRequest.Code,
		Products:        inviteRequest.Products,
	}

	// check if user exists
	_, err := GetAuthUserByMail(inviteRequest.Email)
	if err == nil {
		log.ErrorF("user %s already exists", inviteRequest.Email)
		return "", errors.New("user already exists")
	}

	err = lib.SetFirestoreErr(collectionName, invite.Uid, invite)
	if err != nil {
		log.ErrorF("could not create user %s", inviteRequest.Email)
		return "", errors.New("could not create user")
	}

	log.Printf("Created invite with uid %s", invite.Uid)
	return invite.Uid, nil
}
