package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

type Agency struct {
	AuthId          string      `json:"authId" firestore:"authId" bigquery:"-"`
	Uid             string      `json:"uid" firestore:"uid" bigquery:"-"`
	Email           string      `json:"email" firestore:"email" bigquery:"-"`
	VatCode         string      `json:"vatCode" firestore:"vatCode" bigquery:"-"`
	Name            string      `json:"name" firestore:"name" bigquery:"-"`
	Manager         User        `json:"manager" firestore:"manager" bigquery:"-"`
	NodeSetting     NodeSetting `json:"nodeSetting" firestore:"nodeSetting" bigquery:"-"`
	Users           []string    `json:"users" firestore:"users" bigquery:"-"`                                   // will contain users UIDs
	ParentAgency    string      `json:"parentAgency,omitempty" firestore:"parentAgency,omitempty" bigquery:"-"` // parent Agency UID
	Agencies        []string    `json:"agencies" firestore:"agencies" bigquery:"-"`                             // will contain agencies UIDs
	Agents          []string    `json:"agents" firestore:"agents" bigquery:"-"`                                 // will contain agents UIDs
	IsActive        bool        `json:"isActive" firestore:"isActive" bigquery:"-"`
	Products        []Product   `json:"products" firestore:"products" bigquery:"-"`
	Policies        []string    `json:"policies" firestore:"policies" bigquery:"-"` // will contain policies UIDs
	Steps           []Step      `json:"steps" firestore:"steps" bigquery:"-"`
	Skin            Skin        `json:"skin" firestore:"skin" bigquery:"-"`
	RuiCode         string      `json:"ruiCode" firestore:"ruiCode" bigquery:"-"`
	RuiSection      string      `json:"ruiSection" firestore:"ruiSection" bigquery:"-"`
	RuiRegistration time.Time   `json:"ruiRegistration" firestore:"ruiRegistration" bigquery:"-"`
	CreationDate    time.Time   `json:"creationDate" firestore:"creationDate" bigquery:"-"`
	UpdatedDate     time.Time   `json:"updatedDate" firestore:"updatedDate" bigquery:"-"`
}

type Skin struct {
	PrimaryColor   string `json:"primaryColor" firestore:"primaryColor" bigquery:"-"`
	SecondaryColor string `json:"secondaryColor" firestore:"secondaryColor" bigquery:"-"`
	LogoUrl        string `json:"logoUrl" firestore:"logoUrl" bigquery:"-"`
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
