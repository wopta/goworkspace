package models

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.dev.wopta.it/goworkspace/lib"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	"google.golang.org/api/iterator"
)

type PolicyNotes struct {
	Notes []PolicyNote `json:"notes"`
}

type PolicyNote struct {
	Name               string    `json:"name" firestore:"name"`
	Surname            string    `json:"surname" firestore:"surname"`
	CreateDate         time.Time `json:"createDate" firestore:"createDate"`
	Type               string    `json:"type" firestore:"type"`
	ReadableByProducer bool      `json:"readableByProducer" firestore:"readableByProducer"`
	PolicyUid          string    `json:"policyUid" firestore:"policyUid"`
	Text               string    `json:"text" firestore:"text"`
	ExecutionId        string    `json:"executionId" firestore:"executionId"`
}

func AddNoteToPolicy(policyUid string, note PolicyNote) error {
	note.CreateDate = time.Now()
	note.PolicyUid = policyUid
	note.ExecutionId = env.GetExecutionId()
	if _, err := lib.GetFirestoreErr(PolicyCollection, policyUid); err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return errors.New("Policy not found")
		}
	}
	err := lib.SetFirestoreErr(lib.PolicyNoteCollection, lib.NewDoc(lib.PolicyNoteCollection), note)
	if err != nil {
		return err
	}
	return nil
}

func GetPolicyNotes(policyUid string) (PolicyNotes, error) {
	docsnap := lib.WhereFirestore(lib.PolicyNoteCollection, "policyUid", "==", policyUid)
	return policyNotesToListData(docsnap)
}

func policyNotesToListData(query *firestore.DocumentIterator) (PolicyNotes, error) {
	var result PolicyNotes
	for {
		d, err := query.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
		}
		var value PolicyNote
		e := d.DataTo(&value)
		if e != nil {
			return result, e
		}
		result.Notes = append(result.Notes, value)
	}
	return result, nil
}
func enrichNote(policyNote *PolicyNote, user string) {
	name, surname := os.Getenv("NameUser"), os.Getenv("SurnameUser")
	if name == "" || surname == "" {
		splitAuth := strings.SplitAfter(user, "id:")
		if len(splitAuth) == 2 {
			authId := splitAuth[1]
			userFirebase := lib.WhereLimitFirestore(UserCollection, "authId", "==", authId, 1)
			u, _ := FirestoreDocumentToUser(userFirebase)
			name = u.Name
			surname = u.Surname
		} else {
			return
		}
	}
	os.Setenv("NameUser", name)
	os.Setenv("SurnameUser", surname)
	policyNote.Name = name
	policyNote.Surname = surname
}
func (p *Policy) AddSystemNote(getterNote func(p *Policy) PolicyNote) error {
	note := getterNote(p)
	note.CreateDate = time.Now()
	note.PolicyUid = p.Uid
	note.Type = "System"
	note.ExecutionId = env.GetExecutionId()
	user := os.Getenv("User")
	if user != "" {
		enrichNote(&note, user)
	}
	if env.IsLocal() {
		note.Text += "(Operazione eseguita in ambiente locale, non veritiera)"
	}
	err := lib.SetFirestoreErr(lib.PolicyNoteCollection, lib.NewDoc(lib.PolicyNoteCollection), note)
	if err != nil {
		return err
	}
	return nil
}

func GetEmitNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "La polizza é stata emessa",
	}
}

func GetSignNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "La polizza é stata firmata",
	}
}

func GetPayNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "La polizza é stata pagata",
	}
}
func GetRenewNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "La polizza é stata rinnovata",
	}
}
func GetManualRenewNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "La polizza é stata rinnovata manualmente",
	}
}
func GetChangePaymentProviderNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "Cambio mandato",
	}
}
func GetApproveNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text:               "Polizza é stata approvata",
		ReadableByProducer: true,
	}
}
func GetRejectNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text:               "Polizza é stata rigettata",
		ReadableByProducer: true,
	}
}
func GetProposalNote(p *Policy) PolicyNote {
	var text string
	switch p.Status {
	case PolicyStatusProposal:
		text = "Salvataggio proposta"
	case PolicyStatusNeedsApproval:
		text = "Salvataggio proposta, Rapporto Visita Medica"
	}

	return PolicyNote{
		Text: text,
	}
}

func GetDeletePolicyNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "Polizza cancellata",
	}
}

func GetEmailNote(emailDetail string) func(*Policy) PolicyNote {
	return func(p *Policy) PolicyNote {
		return PolicyNote{
			Text: fmt.Sprintf("Email con oggetto '%v' inviata", emailDetail),
		}
	}
}

func GetErrorNote(processName string) func(*Policy) PolicyNote {
	return func(p *Policy) PolicyNote {
		return PolicyNote{
			Text: fmt.Sprintf("Operazione di '%v' é andata in errore!", processName),
		}
	}
}

func GetQuietanzamentoPolicyNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "Quietanzamento eseguito",
	}
}

func GetAddendumPolicyNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "Dati polizza modificati",
	}
}
