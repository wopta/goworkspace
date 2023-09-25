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
	Uid                      string                 `json:"uid" firestore:"uid" bigquery:"uid"`
	Code                     string                 `json:"code" firestore:"code" bigquery:"code"`
	AuthId                   string                 `json:"authId,omitempty" firestore:"authId,omitempty" bigquery:"-"`
	Name                     string                 `json:"name" firestore:"name" bigquery:"name"`
	Surname                  string                 `json:"surname" firestore:"surname" bigquery:"surname"`
	FiscalCode               string                 `json:"fiscalCode" firestore:"fiscalCode" bigquery:"fiscalCode"`
	VatCode                  string                 `json:"vatCode,omitempty" firestore:"vatCode,omitempty" bigquery:"vatCode"`
	BirthDate                string                 `json:"birthDate"                   firestore:"birthDate,omitempty"         bigquery:"-"`
	BigBirthDate             bigquery.NullDateTime  `json:"-"                           firestore:"-"                           bigquery:"birthDate"`
	BirthCity                string                 `json:"birthCity"                   firestore:"birthCity,omitempty"         bigquery:"birthCity"`
	BirthProvince            string                 `json:"birthProvince"               firestore:"birthProvince,omitempty"     bigquery:"birthProvince"`
	Residence                *Address               `json:"residence,omitempty"         firestore:"residence,omitempty"         bigquery:"-"`
	BigResidenceStreetName   string                 `json:"-"                           firestore:"-"                           bigquery:"residenceStreetName"`
	BigResidenceStreetNumber string                 `json:"-"                           firestore:"-"                           bigquery:"residenceStreetNumber"`
	BigResidenceCity         string                 `json:"-"                           firestore:"-"                           bigquery:"residenceCity"`
	BigResidencePostalCode   string                 `json:"-"                           firestore:"-"                           bigquery:"residencePostalCode"`
	BigResidenceLocality     string                 `json:"-"                           firestore:"-"                           bigquery:"residenceLocality"`
	BigResidenceCityCode     string                 `json:"-"                           firestore:"-"                           bigquery:"residenceCityCode"`
	Domicile                 *Address               `json:"domicile,omitempty"          firestore:"domicile,omitempty"          bigquery:"-"`
	BigDomicileStreetName    string                 `json:"-"                           firestore:"-"                           bigquery:"domicileStreetName"`
	BigDomicileStreetNumber  string                 `json:"-"                           firestore:"-"                           bigquery:"domicileStreetNumber"`
	BigDomicileCity          string                 `json:"-"                           firestore:"-"                           bigquery:"domicileCity"`
	BigDomicilePostalCode    string                 `json:"-"                           firestore:"-"                           bigquery:"domicilePostalCode"`
	BigDomicileLocality      string                 `json:"-"                           firestore:"-"                           bigquery:"domicileLocality"`
	BigDomicileCityCode      string                 `json:"-"                           firestore:"-"                           bigquery:"domicileCityCode"`
	Location                 Location               `json:"location"                    firestore:"location,omitempty"          bigquery:"-"`
	BigLocation              bigquery.NullGeography `json:"-"                           firestore:"-"                           bigquery:"location"`
	Mail                     string                 `json:"mail" firestore:"mail" bigquery:"mail"`
	Phone                    string                 `json:"phone,omitempty" firestore:"phone,omitempty" bigquery:"phone"`
	Role                     string                 `json:"role" firestore:"role" bigquery:"-"`
	ManagerUid               string                 `json:"managerUid,omitempty" firestore:"managerUid,omitempty" bigquery:"managerUid"`
	AgencyUid                string                 `json:"agencyUid"            firestore:"agencyUid"            bigquery:"agencyUid"`
	Agents                   []string               `json:"agents"               firestore:"agents"               bigquery:"-"`
	Users                    []string               `json:"users"                firestore:"users"                bigquery:"-"` // will contain users UIDs
	IsActive                 bool                   `json:"isActive"             firestore:"isActive"             bigquery:"isActive"`
	Products                 []Product              `json:"products"             firestore:"products"             bigquery:"-"`
	Policies                 []string               `json:"policies"             firestore:"policies"             bigquery:"-"` // will contain policies UIDs
	RuiCode                  string                 `json:"ruiCode"              firestore:"ruiCode"              bigquery:"ruiCode"`
	RuiSection               string                 `json:"ruiSection"           firestore:"ruiSection"           bigquery:"ruiSection"`
	RuiRegistration          time.Time              `json:"ruiRegistration"      firestore:"ruiRegistration"      bigquery:"-"`
	BigRuiRegistration       bigquery.NullDateTime  `json:"-"                    firestore:"-"                    bigquery:"ruiRegistration"`
	Data                     string                 `json:"-"                    firestore:"-"                    bigquery:"data"`
	CreationDate             time.Time              `json:"creationDate" firestore:"creationDate" bigquery:"-"`
	UpdatedDate              time.Time              `json:"updatedDate" firestore:"updatedDate" bigquery:"-"`
	BigCreationDate          bigquery.NullDateTime  `json:"-" firestore:"-" bigquery:"creationDate"`
	BigUpdatedDate           bigquery.NullDateTime  `json:"-" firestore:"-" bigquery:"updatedDate"`
}

func (agent *Agent) BigquerySave(origin string) error {
	agent.BigRuiRegistration = lib.GetBigQueryNullDateTime(agent.RuiRegistration)
	agent.BigCreationDate = lib.GetBigQueryNullDateTime(agent.CreationDate)
	agent.BigUpdatedDate = lib.GetBigQueryNullDateTime(agent.UpdatedDate)

	if agent.BirthDate != "" {
		birthDate, err := time.Parse(time.RFC3339, agent.BirthDate)
		if err != nil {
			return err
		}
		agent.BigBirthDate = lib.GetBigQueryNullDateTime(birthDate)
	}

	agent.BigLocation = bigquery.NullGeography{
		// TODO: Check if correct: Geography type uses the WKT format for geometry
		GeographyVal: fmt.Sprintf("POINT (%f %f)", agent.Location.Lng, agent.Location.Lat),
		Valid:        true,
	}

	data, err := json.Marshal(agent)
	if err != nil {
		return err
	}
	agent.Data = string(data) // includes agent.User data

	table := lib.GetDatasetByEnv(origin, AgentCollection)
	log.Println("agent save big query: " + agent.Uid)

	return lib.InsertRowsBigQuery(WoptaDataset, table, agent)
}

func GetAgentByAuthId(authId string) (*Agent, error) {
	agentFirebase := lib.WhereLimitFirestore(AgentCollection, "authId", "==", authId, 1)
	agent, err := FirestoreDocumentToAgent(agentFirebase)

	return agent, err
}

func GetAgentByUid(uid string) (*Agent, error) {
	agentFirebase := lib.WhereLimitFirestore(AgentCollection, "uid", "==", uid, 1)
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

	agent.UpdatedDate = time.Now().UTC()
	err = lib.SetFirestoreErr(fireAgent, agent.Uid, agent)
	if err != nil {
		log.Printf("[updateAgentPortfolio] ERROR saving agent: %s", err.Error())
		return err
	}

	err = agent.BigquerySave(origin)

	return err
}

func IsPolicyInAgentPortfolio(agentUid, policyUid string) bool {
	agent, err := GetAgentByUid(agentUid)
	if err != nil {
		log.Printf("[IsPolicyInAgentPortfolio] error retrieving agent: %s", err.Error())
		return false
	}

	return lib.SliceContains(agent.Policies, policyUid)
}
