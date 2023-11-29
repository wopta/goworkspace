package _script

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	"log"
)

func UpdateContractorAddress(policyUid, city, cityCode, locality string) {
	var (
		err    error
		policy models.Policy
	)

	policy, err = plc.GetPolicy(policyUid, "")
	if err != nil {
		log.Printf("error fetching policy %s from Firestore: %s", policyUid, err.Error())
		return
	}

	policy.Contractor.Residence.City = city
	policy.Contractor.Residence.CityCode = cityCode
	policy.Contractor.Residence.Locality = locality

	policy.Assets[0].Person.Residence.City = city
	policy.Assets[0].Person.Residence.CityCode = cityCode
	policy.Assets[0].Person.Residence.Locality = locality

	err = lib.SetFirestoreErr(models.PolicyCollection, policyUid, policy)
	if err != nil {
		log.Printf("error saving policy %s in Firestore: %s", policyUid, err.Error())
		return
	}

	policy.BigquerySave("")

	return
}
