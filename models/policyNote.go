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
