package models

import (
	"errors"
	"fmt"
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
	Username           string    `json:"surname" firestore:"surname"`
	CreateDate         time.Time `json:"createDate" firestore:"createDate"`
	Type               string    `json:"type" firestore:"type"`
	ReadableByProducer bool      `json:"readableByProducer" firestore:"readableByProducer"`
	PolicyUid          string    `json:"policyUid" firestore:"policyUid"`
	Text               string    `json:"text" firestore:"text"`
}

func AddNoteToPolicy(policyUid string, note PolicyNote) error {
	note.CreateDate = time.Now()
	note.PolicyUid = policyUid
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

func (p *Policy) AddSystemNote(getterNote func(p *Policy) PolicyNote) error {
	note := getterNote(p)
	note.CreateDate = time.Now()
	note.PolicyUid = p.Uid
	note.Type = "System"
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
func GetChangeAppendiceNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "Appendice modificata",
	}
}
func GetDeletePolicyNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text: "Appendice modificata",
	}
}

func GetEmailNote(emailDetail string) func(*Policy) PolicyNote {
	return func(p *Policy) PolicyNote {
		return PolicyNote{
			Text: fmt.Sprintf("Email con oggetto '%v' é stata inviata", emailDetail),
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
