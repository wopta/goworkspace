package models

import (
	"errors"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"google.golang.org/api/iterator"
)

type PolicyNotes struct {
	Notes []PolicyNote `json:"notes"`
}

type PolicyNote struct {
	Name         string    `json:"name" firestore:"name"`
	Username     string    `json:"username" firestore:"username"`
	CreateDate   time.Time `json:"createDate" firestore:"createDate"`
	Type         string    `json:"type" firestore:"type"`
	OnlyProducer bool      `json:"onlyProducer" firestore:"onlyProducer"`
	CreatedBy    string    `json:"createdBy" firestore:"createdBy"`
	PolicyUid    string    `json:"policyUid" firestore:"policyUid"`
	Text         string    `json:"text" firestore:"text"`
}

func AddNoteToPolicy(policyUid string, note PolicyNote) error {
	note.CreateDate = time.Now()
	note.PolicyUid = policyUid
	if note.CreatedBy == "" {
		return errors.New("Need to specify the creator")
	}
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
	if note.CreatedBy == "" {
		return errors.New("Need to specify the creator")
	}
	err := lib.SetFirestoreErr(lib.PolicyNoteCollection, lib.NewDoc(lib.PolicyNoteCollection), note)
	if err != nil {
		return err
	}
	return nil
}

func GetEmitNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text:         "La polizza é stata emessa",
		Type:         "System",
		OnlyProducer: false,
	}
}

func GetSignNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text:         "La polizza é stata firmata",
		Type:         "System",
		OnlyProducer: false,
	}
}

func GetPayNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text:         "La polizza é stata pagata",
		Type:         "System",
		OnlyProducer: false,
	}
}
func GetRenewNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text:         "La polizza é stata rinnovata",
		Type:         "System",
		OnlyProducer: false,
	}
}
func GetManualRenewNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text:         "La polizza é stata rinnovata manualmente",
		Type:         "System",
		OnlyProducer: false,
	}
}
func GetChangePaymentProviderNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text:         "Cambio mandato",
		Type:         "System",
		OnlyProducer: false,
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
		Text:         text,
		Type:         "System",
		OnlyProducer: false,
	}
}
func GetChangeAppendiceNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text:         "Appendice modificata",
		Type:         "System",
		OnlyProducer: false,
	}
}
func GetDeletePolicyNote(_ *Policy) PolicyNote {
	return PolicyNote{
		Text:         "Appendice modificata",
		Type:         "System",
		OnlyProducer: false,
	}
}
