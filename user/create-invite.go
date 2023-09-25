package user

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
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
	Code            string           `json:"code"`
	Expiration      time.Time        `json:"expiration,omitempty" firestore:"expiration,omitempty"`
	Products        []models.Product `json:"products,omitempty" firestore:"products,omitempty"`
}

func CreateInviteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var createInviteRequest CreateInviteRequest

	reqBytes := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(reqBytes, &createInviteRequest)
	lib.CheckError(err)

	creatorUid, err := lib.GetUserIdFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Println("[CreateInvite] Invalid auth token")
		return `{"success": false}`, `{"success": false}`, nil
	}

	inviteUid, err := CreateInvite(createInviteRequest, r.Header.Get("Origin"), creatorUid)
	if err != nil {
		log.Printf("[CreateInvite]: %s", err.Error())
		return `{"success": false}`, `{"success": false}`, err
	}

	SendInviteMail(inviteUid, createInviteRequest.Email)
	return `{"success": true}`, `{"success": true}`, nil
}

func CreateInvite(inviteRequest CreateInviteRequest, origin, creatorUid string) (string, error) {
	log.Printf("[CreateInvite] Creating invite for user %s with role %s", inviteRequest.Email, inviteRequest.Role)

	collectionName := lib.GetDatasetByEnv(origin, invitesCollection)
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
		log.Println("[CreateInvite]: forbidden role")
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
	_, err := GetAuthUserByMail(origin, inviteRequest.Email)
	if err == nil {
		log.Printf("[CreateInvite]: user %s already exists", inviteRequest.Email)
		return "", errors.New("user already exists")
	}

	err = lib.SetFirestoreErr(collectionName, invite.Uid, invite)
	if err != nil {
		log.Printf("[CreateInvite]: could not create user %s", inviteRequest.Email)
		return "", errors.New("could not create user")
	}

	log.Printf("[CreateInvite] Created invite with uid %s", invite.Uid)
	return invite.Uid, nil
}

func SendInviteMail(inviteUid, email string) {
	var mailRequest mail.MailRequest

	mailRequest.From = "anna@wopta.it"
	mailRequest.To = []string{email}
	mailRequest.Subject = "Benvenuto in Wopta!"
	mailRequest.IsHtml = true

	lines := []string{
		"Ciao,",
		"Ecco il tuo invito al tuo account wopta.it.",
		"Accedi al link sottostante e crea la tua password.",
	}
	for _, line := range lines {
		mailRequest.Message = mailRequest.Message + `<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">` + line + `</p>`
	}

	mailRequest.Message = mailRequest.Message + ` 
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px"><br></p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">A presto,</p>
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#e50075;font-size:14px"><strong>Anna</strong> di Wopta Assicurazioni</p> `
	mailRequest.Title = "Invito a wopta.it"
	mailRequest.IsHtml = true
	mailRequest.IsLink = true
	mailRequest.Link = os.Getenv("WOPTA_CUSTOMER_AREA_BASE_URL") + "/login/inviteregistration?inviteUid=" + inviteUid
	mailRequest.LinkLabel = "Crea la tua password"

	mail.SendMail(mailRequest)
}
