package _script

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
	"google.golang.org/api/iterator"
	"log"
	"os"
	"sync"
)

type policiesData struct {
	Data []policyInfo `json:"data"`
	mux  sync.Mutex
}

type policyInfo struct {
	Policy       models.Policy        `json:"policy"`
	Transactions []models.Transaction `json:"transactions"`
}

func FetchAllPoliciesData() {
	var (
		policies = make([]models.Policy, 0)
		data     = policiesData{Data: make([]policyInfo, 0)}
	)

	queries := lib.Firequeries{
		Queries: []lib.Firequery{
			{Field: "isPay", Operator: "==", QueryValue: true},
			{Field: "isDeleted", Operator: "==", QueryValue: false},
		},
	}

	iter, err := queries.FirestoreWherefields(lib.PolicyCollection)
	if err != nil {
		panic(err)
	}
	defer iter.Stop() // add this line to ensure resources cleaned up
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		var policy models.Policy
		if err := doc.DataTo(&policy); err == nil {
			policies = append(policies, policy)
		}
	}

	var wg sync.WaitGroup

	for _, p := range policies {
		wg.Add(1)
		go func(p models.Policy) {
			defer wg.Done()
			transactions := transaction.GetPolicyActiveTransactions("", p.Uid)
			data.mux.Lock()
			data.Data = append(data.Data, policyInfo{p, transactions})
			data.mux.Unlock()
		}(p)

	}

	wg.Wait()

	rawOut, err := json.Marshal(data.Data)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./_script/prod_data.json", rawOut, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
