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

type Agency struct {
	AuthId          string      `json:"authId"                 firestore:"authId"`
	Uid             string      `json:"uid"                    firestore:"uid"`
	Email           string      `json:"email"                  firestore:"email"`
	VatCode         string      `json:"vatCode"                firestore:"vatCode"`
	Name            string      `json:"name"                   firestore:"name"`
	Manager         User        `json:"manager"                firestore:"manager"`
	NodeSetting     NodeSetting `json:"nodeSetting"            firestore:"nodeSetting"`
	Users           []string    `json:"users"                  firestore:"users"                  bigquery:"-"` // will contain users UIDs
	ParentAgency    string      `json:"parentAgency,omitempty" firestore:"parentAgency,omitempty"`              // parent Agency UID
	Agencies        []string    `json:"agencies"               firestore:"agencies"`                            // will contain agencies UIDs
	Agents          []string    `json:"agents"                 firestore:"agents"`                              // will contain agents UIDs
	IsActive        bool        `json:"isActive"               firestore:"isActive"`
	Products        []Product   `json:"products"               firestore:"products"`
	Policies        []string    `json:"policies"               firestore:"policies"` // will contain policies UIDs
	Steps           []Step      `json:"steps"                  firestore:"steps"`
	Skin            Skin        `json:"skin"                   firestore:"skin"`
	RuiCode         string      `json:"ruiCode"                firestore:"ruiCode"`
	RuiSection      string      `json:"ruiSection"             firestore:"ruiSection"`
	RuiRegistration time.Time   `json:"ruiRegistration"        firestore:"ruiRegistration"`
	CreationDate    time.Time   `json:"creationDate"           firestore:"creationDate"`
	UpdatedDate     time.Time   `json:"updatedDate"            firestore:"updatedDate"`
}

// TODO: missing Surname in Agency
// TODO: missing VatCode in Agency
type AgencyBigquery struct {
	Uid             string                `bigquery:"uid"`
	Name            string                `bigquery:"name"`
	Surname         string                `bigquery:"surname"`
	VatCode         string                `bigquery:"vatCode"`
	IsActive        bool                  `bigquery:"isActive"`
	RuiCode         string                `bigquery:"ruiCode"`
	RuiSection      string                `bigquery:"ruiSection"`
	RuiRegistration bigquery.NullDateTime `bigquery:"ruiRegistration"`
	ManagerUid      string                `bigquery:"managerUid"`
	ParentAgency    string                `bigquery:"parentAgency"` // parent Agency UID
	Agencies        string                `bigquery:"agencies"`     // will contain agencies UIDs
	Agents          string                `bigquery:"agents"`       // will contain agents UIDs
	CreationDate    bigquery.NullDateTime `bigquery:"creationDate"`
	UpdateDate      bigquery.NullDateTime `bigquery:"updateDate"`
	Data            string                `bigquery:"data"`
}

func (agency Agency) toBigquery() (AgencyBigquery, error) {
	agencyJson, err := json.Marshal(agency)
	if err != nil {
		return AgencyBigquery{}, err
	}
	return AgencyBigquery{
		Uid:             agency.Uid,
		Name:            agency.Name,
		Surname:         "",
		VatCode:         "",
		IsActive:        agency.IsActive,
		RuiCode:         agency.RuiCode,
		RuiSection:      agency.RuiSection,
		RuiRegistration: lib.GetBigQueryNullDateTime(agency.RuiRegistration),
		ManagerUid:      agency.Manager.Uid,
		Agents:          strings.Join(agency.Agents, ","),
		CreationDate:    lib.GetBigQueryNullDateTime(agency.CreationDate),
		UpdateDate:      lib.GetBigQueryNullDateTime(agency.UpdatedDate),
		Data:            string(agencyJson),
	}, nil
}

func (agency Agency) BigquerySave(origin string) error {
	table := lib.GetDatasetByEnv(origin, "agency")
	agencyBigquery, err := agency.toBigquery()
	if err != nil {
		return err
	}

	log.Println("agency save big query: " + agency.Uid)

	return lib.InsertRowsBigQuery("wopta", table, agencyBigquery)
}

type Skin struct {
	PrimaryColor   string `json:"primaryColor"   firestore:"primaryColor"   bigquery:"-"`
	SecondaryColor string `json:"secondaryColor" firestore:"secondaryColor" bigquery:"-"`
	LogoUrl        string `json:"logoUrl"        firestore:"logoUrl"        bigquery:"-"`
}

func GetAgencyByAuthId(authId string) (*Agency, error) {
	agencyFirebase := lib.WhereLimitFirestore(AgencyCollection, "uid", "==", authId, 1)
	agency, err := FirestoreDocumentToAgency(agencyFirebase)

	return agency, err
}

func FirestoreDocumentToAgency(query *firestore.DocumentIterator) (*Agency, error) {
	var result Agency
	agencyDocumentSnapshot, err := query.Next()

	if err == iterator.Done && agencyDocumentSnapshot == nil {
		log.Println("agency not found in firebase DB")
		return &result, fmt.Errorf("no agent found")
	}

	if err != iterator.Done && err != nil {
		log.Println(`error happened while trying to get agency`)
		return &result, err
	}

	e := agencyDocumentSnapshot.DataTo(&result)
	if len(result.Uid) == 0 {
		result.Uid = agencyDocumentSnapshot.Ref.ID
	}

	return &result, e
}

func UpdateAgencyPortfolio(policy *Policy, origin string) error {
	log.Printf("[updateAgencyPortfolio] Policy %s", policy.Uid)
	if policy.AgencyUid == "" {
		log.Printf("[updateAgencyPortfolio] ERROR agency not set")
		return errors.New("agency not set")
	}

	var agency Agency
	fireAgency := lib.GetDatasetByEnv(origin, AgencyCollection)
	docsnap, err := lib.GetFirestoreErr(fireAgency, policy.AgentUid)
	if err != nil {
		log.Printf("[updateAgencyPortfolio] ERROR getting agency from firestore: %s", err.Error())
		return err
	}
	err = docsnap.DataTo(&agency)
	if err != nil {
		log.Printf("[updateAgencyPortfolio] ERROR parsing agency: %s", err.Error())
		return err
	}
	agency.Policies = append(agency.Policies, policy.Uid)

	if !lib.SliceContains(agency.Users, policy.Contractor.Uid) {
		agency.Users = append(agency.Users, policy.Contractor.Uid)
	}

	err = lib.SetFirestoreErr(fireAgency, agency.Uid, agency)

	return err
}
