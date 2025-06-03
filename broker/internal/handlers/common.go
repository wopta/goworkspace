package handlers

import (
	"log"
	"time"

	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

func savePolicy(state bpmn.StorageData) error {
	var policy *flow.Policy
	var origin *flow.String
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
	)
	if err != nil {
		return err
	}

	policy.Updated = time.Now().UTC()
	log.Println("saving to firestore...")
	firePolicy := lib.GetDatasetByEnv(origin.String, lib.PolicyCollection)
	err = lib.SetFirestoreErr(firePolicy, policy.Uid, &policy)
	if err != nil {
		return err
	}
	log.Println("firestore saved!")

	policy.BigquerySave(origin.String)
	return nil
}
