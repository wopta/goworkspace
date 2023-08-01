package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

type Agent struct {
	User
	ManagerUid         string                `json:"managerUid,omitempty" firestore:"managerUid,omitempty" bigquery:"managerUid"`
	AgencyUid          string                `json:"agencyUid"            firestore:"agencyUid"            bigquery:"agencyUid"`
	Agents             []string              `json:"agents"               firestore:"agents"               bigquery:"-"`
	Users              []string              `json:"users"                firestore:"users"                bigquery:"-"` // will contain users UIDs
	IsActive           bool                  `json:"isActive"             firestore:"isActive"             bigquery:"isActive"`
	Products           []Product             `json:"products"             firestore:"products"             bigquery:"-"`
	Policies           []string              `json:"policies"             firestore:"policies"             bigquery:"-"` // will contain policies UIDs
	RuiCode            string                `json:"ruiCode"              firestore:"ruiCode"              bigquery:"ruiCode"`
	RuiSection         string                `json:"ruiSection"           firestore:"ruiSection"           bigquery:"ruiSection"`
	RuiRegistration    time.Time             `json:"ruiRegistration"      firestore:"ruiRegistration"      bigquery:"-"`
	BigRuiRegistration bigquery.NullDateTime `json:"-"                    firestore:"-"                    bigquery:"ruiRegistration"`
	Data               string                `json:"-"                    firestore:"-"                    bigquery:"data"`
}

func (agent *Agent) BigquerySave(origin string) error {
	err := agent.User.prepareForBigquerySave()
	if err != nil {
		return err
	}
	agent.User.Data = ""
	agent.BigRuiRegistration = lib.GetBigQueryNullDateTime(agent.RuiRegistration)

	data, err := json.Marshal(agent)
	if err != nil {
		return err
	}
	agent.Data = string(data) // includes agent.User data

	table := lib.GetDatasetByEnv(origin, AgentCollection)
	if err := agent.prepareForBigquerySave(); err != nil {
		return err
	}

	log.Println("agent save big query: " + agent.Uid)

	return lib.InsertRowsBigQuery(WoptaDataset, table, agent)
}

func GetAgentByAuthId(authId string) (*Agent, error) {
	agentFirebase := lib.WhereLimitFirestore(AgentCollection, "authId", "==", authId, 1)
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

func UpdateAgentPortfolio(policy *Policy, origin string) error {
	log.Printf("[updateAgentPortfolio] Policy %s", policy.Uid)
	if policy.AgentUid == "" {
		log.Printf("[updateAgentPortfolio] ERROR agent not set")
		return errors.New("agent not set")
	}

	var agent Agent
	fireAgent := lib.GetDatasetByEnv(origin, AgentCollection)
	docsnap, err := lib.GetFirestoreErr(fireAgent, policy.AgentUid)
	if err != nil {
		log.Printf("[updateAgentPortfolio] ERROR getting agent from firestore: %s", err.Error())
		return err
	}
	err = docsnap.DataTo(&agent)
	if err != nil {
		log.Printf("[updateAgentPortfolio] ERROR parsing agent: %s", err.Error())
		return err
	}
	agent.Policies = append(agent.Policies, policy.Uid)

	if !lib.SliceContains(agent.Users, policy.Contractor.Uid) {
		agent.Users = append(agent.Users, policy.Contractor.Uid)
	}

	err = lib.SetFirestoreErr(fireAgent, agent.Uid, agent)

	return err
}
