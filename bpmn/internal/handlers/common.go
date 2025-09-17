package handlers

import (
	"log"
	"time"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

func savePolicy(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
	)
	if err != nil {
		return err
	}

	policy.Updated = time.Now().UTC()
	log.Println("saving to firestore...")
	err = lib.SetFirestoreErr(lib.PolicyCollection, policy.Uid, &policy)
	if err != nil {
		return err
	}
	log.Println("firestore saved!")

	policy.BigquerySave()
	return nil
}
