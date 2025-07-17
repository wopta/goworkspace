package _script

import (
	"log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

func UpdateContractorName(policyUid, contractorName string) {
	var (
		err    error
		policy models.Policy
	)

	policy, err = plc.GetPolicy(policyUid)
	if err != nil {
		log.Printf("error fetching policy %s from Firestore: %s", policyUid, err.Error())
		return
	}

	policy.Contractor.Name = contractorName
	policy.Assets[0].Person.Name = contractorName

	err = lib.SetFirestoreErr(models.PolicyCollection, policyUid, policy)
	if err != nil {
		log.Printf("error saving policy %s in Firestore: %s", policyUid, err.Error())
		return
	}

	policy.BigquerySave()

	return
}
