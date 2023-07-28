package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

type Agent struct {
	User
	ManagerUid      string    `json:"managerUid,omitempty" firestore:"managerUid,omitempty"`
	AgencyUid       string    `json:"agencyUid"            firestore:"agencyUid"`
	Agents          []string  `json:"agents"               firestore:"agents"`
	Users           []string  `json:"users"                firestore:"users"                bigquery:"-"` // will contain users UIDs
	IsActive        bool      `json:"isActive"             firestore:"isActive"`
	Products        []Product `json:"products"             firestore:"products"`
	Policies        []string  `json:"policies"             firestore:"policies"` // will contain policies UIDs
	RuiCode         string    `json:"ruiCode"              firestore:"ruiCode"`
	RuiSection      string    `json:"ruiSection"           firestore:"ruiSection"`
	RuiRegistration time.Time `json:"ruiRegistration"      firestore:"ruiRegistration"`
}

type AgentBigquery struct {
	Uid             string                `bigquery:"uid"`
	Name            string                `bigquery:"name"`
	Surname         string                `bigquery:"surname"`
	FiscalCode      string                `bigquery:"fiscalCode"`
	VatCode         string                `bigquery:"vatCode"`
	IsActive        bool                  `bigquery:"isActive"`
	RuiCode         string                `bigquery:"ruiCode"`
	RuiSection      string                `bigquery:"ruiSection"`
	RuiRegistration bigquery.NullDateTime `bigquery:"ruiRegistration"`
	AgencyUid       string                `bigquery:"agencyUid"`
	ManagerUid      string                `bigquery:"managerUid"`
	Agents          string                `bigquery:"agents"`
	CreationDate    bigquery.NullDateTime `bigquery:"creationDate"`
	UpdateDate      bigquery.NullDateTime `bigquery:"updateDate"`
	Data            string                `bigquery:"data"`
}

func (agent Agent) toBigquery() (AgentBigquery, error) {
	agentJson, err := json.Marshal(agent)
	if err != nil {
		return AgentBigquery{}, err
	}
	return AgentBigquery{
		Uid:             agent.Uid,
		Name:            agent.Name,
		Surname:         agent.Surname,
		FiscalCode:      agent.FiscalCode,
		VatCode:         agent.VatCode,
		IsActive:        agent.IsActive,
		RuiCode:         agent.RuiCode,
		RuiSection:      agent.RuiSection,
		RuiRegistration: lib.GetBigQueryNullDateTime(agent.RuiRegistration),
		AgencyUid:       agent.AgencyUid,
		ManagerUid:      agent.ManagerUid,
		Agents:          strings.Join(agent.Agents, ","),
		CreationDate:    lib.GetBigQueryNullDateTime(agent.CreationDate),
		UpdateDate:      lib.GetBigQueryNullDateTime(agent.UpdatedDate),
		Data:            string(agentJson),
	}, nil
}

func (agent Agent) BigquerySave(origin string) error {
	table := lib.GetDatasetByEnv(origin, "agent")
	agentBigquery, err := agent.toBigquery()
	if err != nil {
		return err
	}

	log.Println("agent save big query: " + agent.Uid)

	return lib.InsertRowsBigQuery("wopta", table, agentBigquery)
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
