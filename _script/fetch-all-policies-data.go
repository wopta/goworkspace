package _script

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/transaction"
	"google.golang.org/api/iterator"
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
			transactions := transaction.GetPolicyValidTransactions(p.Uid, nil)
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
		data []policyInfo
	)

	rawFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(rawFile, &data)
	if err != nil {
		log.Fatal(err)
	}

	dirAgent, _ := network.GetNetworkNodeByCode("W1.DIRAgent")
	if dirAgent == nil {
		log.Fatal("nil dirAgent")
	}

	var wg sync.WaitGroup

	for _, d := range data {
		wg.Add(1)

		go func(d policyInfo) {
			defer wg.Done()
			oldCodeCompany := d.Policy.CodeCompany

			d.Policy.CodeCompany = "PROD" + oldCodeCompany
			d.Policy.SignUrl = "www.wopta.it"
			d.Policy.PayUrl = "www.wopta.it"
			if d.Policy.Channel == "network" {
				d.Policy.ProducerUid = dirAgent.Uid
				d.Policy.ProducerCode = dirAgent.Code
			}
			d.Policy.Attachments = nil

			d.Policy.Contractor.Mail = "yousef.hammar+" + oldCodeCompany + "@wopta.it"
			d.Policy.Contractor.Phone = "+393334455667"
			d.Policy.Contractor.IdentityDocuments = nil

			for assetIndex, _ := range d.Policy.Assets {
				for guaranteeIndex, g := range d.Policy.Assets[assetIndex].Guarantees {
					if g.Beneficiaries != nil {
						for benIndex, ben := range *g.Beneficiaries {
							if ben.BeneficiaryType == models.BeneficiaryChosenBeneficiary {
								(*d.Policy.Assets[assetIndex].Guarantees[guaranteeIndex].Beneficiaries)[benIndex].Mail = "mail@wopta.it"
								(*d.Policy.Assets[assetIndex].Guarantees[guaranteeIndex].Beneficiaries)[benIndex].Phone = "+393334455667"
							}

						}
					}
				}

				d.Policy.Assets[assetIndex].Person.Mail = "yousef.hammar+" + oldCodeCompany + "@wopta.it"
				d.Policy.Assets[assetIndex].Person.Phone = "+393334455667"
				d.Policy.Assets[assetIndex].Person.IdentityDocuments = nil
			}

			if d.Policy.Contractors != nil {
				for contractorIndex, _ := range *d.Policy.Contractors {
					(*d.Policy.Contractors)[contractorIndex].Mail = "mail@wopta.it"
					(*d.Policy.Contractors)[contractorIndex].Phone = "+393334455667"
					(*d.Policy.Contractors)[contractorIndex].IdentityDocuments = nil
				}
			}

			toBeWritten := make(map[string]map[string]interface{})
			trMap := make(map[string]interface{})
			for trIndex, tr := range d.Transactions {
				d.Transactions[trIndex].BigQueryParse()
				d.Transactions[trIndex].PayUrl = "www.wopta.it"
				trMap[tr.Uid] = d.Transactions[trIndex]
			}

			d.Policy.BigQueryParse()

			toBeWritten = map[string]map[string]interface{}{
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
				log.Printf("error inserting policy rows %v", err)
				return
			}

			err = lib.InsertRowsBigQuery(lib.WoptaDataset, lib.TransactionsCollection, d.Transactions)
			if err != nil {
				log.Printf("error inserting transactions rows %v", err)
				return
			}
		}(d)
	}

	wg.Wait()

	log.Println("Done...")
}
