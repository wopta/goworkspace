package handlers

import (
	"log"
	"time"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

func savePolicy(state bpmnEngine.StorageData) error {
	var policy *flow.Policy
	var origin *flow.String
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
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
