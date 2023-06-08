package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
)

func CreateInviteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var createInviteRequest CreateInviteRequest

	reqBytes := lib.ErrorByte(ioutil.ReadAll(r.Body))
	json.Unmarshal(reqBytes, &createInviteRequest)

	creatorUid, err := lib.GetUserIdFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Println("[CreateInvite] Invalid auth token")
		return `{"success": false}`, `{"success": false}`, nil
	}

	if inviteUid, ok := CreateInvite(createInviteRequest.Email, createInviteRequest.Role, r.Header.Get("Origin"), creatorUid); ok {
		SendInviteMail(inviteUid, createInviteRequest.Email)
	}

	return `{"success": true}`, `{"success": true}`, nil
}

func CreateInvite(mail, role, origin, creatorUid string) (string, bool) {
	log.Printf("[CreateInvite] Creating invite for user %s with role %s", mail, role)

	collectionName := lib.GetDatasetByEnv(origin, invitesCollection)
	inviteUid := lib.NewDoc(collectionName)
	inviteExpiration := time.Now().UTC().Add(time.Hour * 168)

	invite := UserInvite{
		Email:      mail,
		Role:       role,
		Expiration: inviteExpiration,
		Uid:        inviteUid,
		CreatorUid: creatorUid,
	}

	lib.SetFirestore(collectionName, invite.Uid, invite)

	log.Printf("[CreateInvite] Created invite with uid %s", invite.Uid)
	return invite.Uid, true
}

func SendInviteMail(inviteUid, email string) {
	var mailRequest mail.MailRequest

	mailRequest.From = "anna@wopta.it"
	mailRequest.To = []string{email}
	mailRequest.Content = `<p>` + inviteUid + `</p>`
	mailRequest.Subject = "Benvenuto in Wopta!"
	mailRequest.IsHtml = true

	mail.SendMail(mailRequest)
}

type CreateInviteRequest struct {
	Role  string `json:"role"`
	Email string `json:"email"`
}

type UserInvite struct {
	Role       string    `json:"role,omitempty" firestore:"role,omitempty"`
	Email      string    `json:"email,omitempty" firestore:"email,omitempty"`
	Uid        string    `json:"uid,omitempty" firestore:"uid,omitempty"`
	CreatorUid string    `json:"creatorUid,omitempty" firestore:"creatorUid,omitempty"`
	Consumed   bool      `json:"consumed" firestore:"consumed"`
	Expiration time.Time `json:"expiration,omitempty" firestore:"expiration,omitempty"`
}
