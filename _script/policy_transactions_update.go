package _script

import (
	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
	"log"
	"time"
)

func getAllPolicies() ([]models.Policy, error) {
	var policies = make([]models.Policy, 0)
	docIterator := lib.OrderFirestore(lib.PolicyCollection, "uid", firestore.Asc)

	snapshots, err := docIterator.GetAll()
	if err != nil {
		log.Printf("error getting polcies from Firestore: %s", err.Error())
		return policies, err
	}

	for _, snapshot := range snapshots {
		var policy models.Policy
		err = snapshot.DataTo(&policy)
		if err != nil {
			log.Printf("error parsing policy %s: %s", snapshot.Ref.ID, err.Error())
		} else {
			policies = append(policies, policy)
		}
	}

	return policies, nil
}

func PolicyTransactionsUpdate() {
	policies, err := getAllPolicies()
	if err != nil {
		return
	}

	for _, p := range policies {
		m := map[string]map[string]interface{}{
			models.PolicyCollection:       make(map[string]interface{}),
			models.TransactionsCollection: make(map[string]interface{}),
		}

		p.Annuity = 0
		if p.Name == models.LifeProduct {
			p.IsRenewable = true
		} else {
			p.IsRenewable = false
		}
		p.Updated = time.Now().UTC()
		m[models.PolicyCollection][p.Uid] = p

		transactions := transaction.GetPolicyTransactions("", p.Uid)
		for _, t := range transactions {
			t.Annuity = 0
			t.UpdateDate = time.Now()
			m[models.TransactionsCollection][t.Uid] = t
		}

		log.Printf("%v", m)

		err = lib.SetBatchFirestoreErr(m)
		if err != nil {
			log.Printf("error saving policy with uid %s into firestore: %s", p.Uid, err.Error())
			continue
		}

		p.BigquerySave("")

		for _, t := range m[models.TransactionsCollection] {
			tr := t.(models.Transaction)
			tr.BigQuerySave("")
		}
	}

}
