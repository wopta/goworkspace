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

type Agency struct {
	AuthId             string                `json:"authId"                 firestore:"authId"                 bigquery:"-"`
	Uid                string                `json:"uid"                    firestore:"uid"                    bigquery:"uid"`
	Email              string                `json:"email"                  firestore:"email"                  bigquery:"email"`
	VatCode            string                `json:"vatCode"                firestore:"vatCode"                bigquery:"vatCode"`
	Name               string                `json:"name"                   firestore:"name"                   bigquery:"name"`
	Manager            User                  `json:"manager"                firestore:"manager"                bigquery:"-"`
	BigManagerUid      string                `json:"-"                      firestore:"-"                      bigquery:"managerUid"`
	NodeSetting        NodeSetting           `json:"nodeSetting"            firestore:"nodeSetting"            bigquery:"-"`
	Users              []string              `json:"users"                  firestore:"users"                  bigquery:"-"`            // will contain users UIDs
	ParentAgency       string                `json:"parentAgency,omitempty" firestore:"parentAgency,omitempty" bigquery:"parentAgency"` // parent Agency UID
	Agencies           []string              `json:"agencies"               firestore:"agencies"               bigquery:"-"`            // will contain agencies UIDs
	Agents             []string              `json:"agents"                 firestore:"agents"                 bigquery:"-"`            // will contain agents UIDs
	IsActive           bool                  `json:"isActive"               firestore:"isActive"               bigquery:"isActive"`
	Products           []Product             `json:"products"               firestore:"products"               bigquery:"-"`
	Policies           []string              `json:"policies"               firestore:"policies"               bigquery:"-"` // will contain policies UIDs
	Steps              []Step                `json:"steps"                  firestore:"steps"                  bigquery:"-"`
	Skin               Skin                  `json:"skin"                   firestore:"skin"                   bigquery:"-"`
	RuiCode            string                `json:"ruiCode"                firestore:"ruiCode"                bigquery:"ruiCode"`
	RuiSection         string                `json:"ruiSection"             firestore:"ruiSection"             bigquery:"ruiSection"`
	RuiRegistration    time.Time             `json:"ruiRegistration"        firestore:"ruiRegistration"        bigquery:"-"`
	BigRuiRegistration bigquery.NullDateTime `json:"-"                      firestore:"-"                      bigquery:"ruiRegistration"`
	CreationDate       time.Time             `json:"creationDate"           firestore:"creationDate"           bigquery:"-"`
	BigCreationDate    bigquery.NullDateTime `json:"-"                      firestore:"-"                      bigquery:"creationDate"`
	UpdatedDate        time.Time             `json:"updatedDate"            firestore:"updatedDate"            bigquery:"-"`
	BigUpdatedDate     bigquery.NullDateTime `json:"-"                      firestore:"-"                      bigquery:"updatedDate"`
	Data               string                `json:"-"                      firestore:"-"                      bigquery:"data"`
}

type Skin struct {
	PrimaryColor   string `json:"primaryColor"   firestore:"primaryColor"   bigquery:"-"`
	SecondaryColor string `json:"secondaryColor" firestore:"secondaryColor" bigquery:"-"`
	LogoUrl        string `json:"logoUrl"        firestore:"logoUrl"        bigquery:"-"`
}

func (agency *Agency) BigquerySave(origin string) error {
	agency.BigManagerUid = agency.Manager.Uid
	agency.BigRuiRegistration = lib.GetBigQueryNullDateTime(agency.RuiRegistration)
	agency.BigCreationDate = lib.GetBigQueryNullDateTime(agency.CreationDate)
	agency.BigUpdatedDate = lib.GetBigQueryNullDateTime(agency.UpdatedDate)
	data, _ := json.Marshal(agency)
	agency.Data = string(data)

	table := lib.GetDatasetByEnv(origin, AgencyCollection)
	log.Println("[Agency] save big query: " + agency.Uid)

	return lib.InsertRowsBigQuery(WoptaDataset, table, agency)
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
	docsnap, err := lib.GetFirestoreErr(fireAgency, policy.AgencyUid)
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

	agency.UpdatedDate = time.Now().UTC()
	err = lib.SetFirestoreErr(fireAgency, agency.Uid, agency)
	if err != nil {
		log.Printf("[updateAgencyPortfolio] ERROR saving agency: %s", err.Error())
		return err
	}

	err = agency.BigquerySave(origin)

	return err
}
