package _script

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/transaction"
	"google.golang.org/api/iterator"
	"log"
	"os"
	"sync"
)

const filePath = "./_script/prod_data.json"

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

	err = os.WriteFile(filePath, rawOut, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func ImportPoliciesData() {
	var (
		data        []policyInfo
		toBeWritten = make(map[string]map[string]map[string]interface{})
	)

	rawFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(rawFile, &data)
	if err != nil {
		log.Fatal(err)
	}

	dirAgent := network.GetNetworkNodeByCode("W1.DIRAgent")

	for _, d := range data {
		go func(d policyInfo) {
			oldCodeCompany := d.Policy.CodeCompany

			d.Policy.CodeCompany = "PROD" + d.Policy.CodeCompany
			d.Policy.SignUrl = "www.wopta.it"
			d.Policy.PayUrl = "www.wopta.it"
			if d.Policy.Channel == "network" {
				d.Policy.ProducerUid = dirAgent.Uid
				d.Policy.ProducerCode = dirAgent.Code
			}

			d.Policy.Contractor.Mail = "yousef.hammar+" + oldCodeCompany + "@wopta.it"
			d.Policy.Contractor.Phone = "+393334455667"
			d.Policy.Contractor.IdentityDocuments = nil

			for _, g := range d.Policy.Assets[0].Guarantees {
				if g.Slug == "death" && g.Beneficiaries != nil {
					for benIndex, _ := range *g.Beneficiaries {
						(*g.Beneficiaries)[benIndex].Mail = "mail@wopta.it"
						(*g.Beneficiaries)[benIndex].Phone = "+393334455667"
					}
				}
			}

			d.Policy.Assets[0].Person.Mail = "yousef.hammar+" + oldCodeCompany + "@wopta.it"
			d.Policy.Assets[0].Person.Phone = "+393334455667"
			d.Policy.Assets[0].Person.IdentityDocuments = nil

			if d.Policy.Contractors != nil {
				for contractorIndex, _ := range *d.Policy.Contractors {
					(*d.Policy.Contractors)[contractorIndex].Mail = "mail@wopta.it"
					(*d.Policy.Contractors)[contractorIndex].Phone = "+393334455667"
					(*d.Policy.Contractors)[contractorIndex].IdentityDocuments = nil
				}
			}

			trMap := make(map[string]interface{})
			for trIndex, tr := range d.Transactions {
				d.Transactions[trIndex].PayUrl = "www.wopta.it"
				trMap[tr.Uid] = d.Transactions[trIndex]
			}

			toBeWritten[d.Policy.Uid] = map[string]map[string]interface{}{
				lib.PolicyCollection: {
					d.Policy.Uid: d.Policy,
				},
				lib.TransactionsCollection: trMap,
			}

			err = lib.SetBatchFirestoreErr(toBeWritten)
			if err != nil {
				log.Printf("error setting batch err %v", err)
				return
			}

			err = lib.InsertRowsBigQuery(lib.WoptaDataset, lib.PolicyCollection, d.Policy)
			if err != nil {
				log.Printf("error inserting rows %v", err)
				return
			}

			err = lib.InsertRowsBigQuery(lib.WoptaDataset, lib.TransactionsCollection, d.Transactions)
			if err != nil {
				log.Printf("error inserting rows %v", err)
				return
			}
		}(d)
	}
}
