package models

import (
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

type Agent struct {
	User
	ManagerUid string    `json:"managerUid,omitempty" firestore:"managerUid,omitempty" bigquery:"-"`
	AgencyUid  string    `json:"agencyUid" firestore:"agencyUid" bigquery:"-"`
	Agents     []string  `json:"agents" firestore:"agents" bigquery:"-"`
	Portfolio  []string  `json:"portfolio" firestore:"portfolio" bigquery:"-"` // will contain users UIDs
	IsActive   bool      `json:"isActive" firestore:"isActive" bigquery:"-"`
	Products   []Product `json:"products" firestore:"products" bigquery:"-"`
	Policies   []string  `json:"policies" firestore:"policies" bigquery:"-"` // will contain policies UIDs
	RuiCode    string    `json:"ruiCode" firestore:"ruiCode" bigquery:"-"`
}

func GetAgentByAuthId(authId string) (*Agent, error) {
	agentFirebase := lib.WhereLimitFirestore(AgentCollection, "user.authId", "==", authId, 1)
	agent, err := FirestoreDocumentToAgent(agentFirebase)

	return agent, err
}

func FirestoreDocumentToAgent(query *firestore.DocumentIterator) (*Agent, error) {
	var result Agent
	agentDocumentSnapshot, err := query.Next()

	if err == iterator.Done && agentDocumentSnapshot == nil {
		log.Println("agent not found in firebase DB")
		return &result, fmt.Errorf("no agent found")
	}

	if err != iterator.Done && err != nil {
		log.Println(`error happened while trying to get agent`)
		return &result, err
	}

	e := agentDocumentSnapshot.DataTo(&result)
	if len(result.Uid) == 0 {
		result.Uid = agentDocumentSnapshot.Ref.ID
	}

	return &result, e
}
