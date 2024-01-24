package models

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/type/latlng"
	"log"
	"time"
)

type Contractor struct {
	EmailVerified            bool                   `firestore:"emailVerified"               json:"emailVerified,omitempty"     bigquery:"-"`
	Uid                      string                 `firestore:"uid"                         json:"uid,omitempty"               bigquery:"uid"`
	Status                   string                 `firestore:"status,omitempty"            json:"status,omitempty"            bigquery:"-"`
	StatusHistory            []string               `firestore:"statusHistory,omitempty"     json:"statusHistory,omitempty"     bigquery:"-"`
	BirthDate                string                 `firestore:"birthDate"                   json:"birthDate,omitempty"         bigquery:"-"`
	BigBirthDate             bigquery.NullDateTime  `firestore:"-"                           json:"-"                           bigquery:"birthDate"`
	BirthCity                string                 `firestore:"birthCity"                   json:"birthCity,omitempty"         bigquery:"birthCity"`
	BirthProvince            string                 `firestore:"birthProvince"               json:"birthProvince,omitempty"     bigquery:"birthProvince"`
	PictureUrl               string                 `firestore:"pictureUrl"                  json:"pictureUrl,omitempty"        bigquery:"-"`
	Location                 Location               `firestore:"location"                    json:"location,omitempty"          bigquery:"-"`
	BigLocation              bigquery.NullGeography `firestore:"-"                           json:"-"                           bigquery:"location"`
	Geo                      latlng.LatLng          `firestore:"geo"                         json:"-"                           bigquery:"-"`
	Name                     string                 `firestore:"name"                        json:"name,omitempty"              bigquery:"name"`
	Gender                   string                 `firestore:"gender"                      json:"gender,omitempty"            bigquery:"gender"`
	Type                     string                 `firestore:"type"                        json:"type,omitempty"              bigquery:"-"`
	Cluster                  string                 `firestore:"cluster"                     json:"cluster,omitempty"           bigquery:"-"`
	Surname                  string                 `firestore:"surname"                     json:"surname,omitempty"           bigquery:"surname"`
	Address                  string                 `firestore:"address"                     json:"address,omitempty"           bigquery:"-"`
	PostalCode               string                 `firestore:"postalCode"                  json:"postalCode,omitempty"        bigquery:"-"`
	City                     string                 `firestore:"city"                        json:"city,omitempty"              bigquery:"-"`
	Locality                 string                 `firestore:"locality"                    json:"locality,omitempty"          bigquery:"-"`
	StreetNumber             string                 `firestore:"streetNumber,omitempty"      json:"streetNumber,omitempty"      bigquery:"-"`
	CityCode                 string                 `firestore:"cityCode"                    json:"cityCode,omitempty"          bigquery:"-"`
	Role                     string                 `firestore:"role"                        json:"role,omitempty"              bigquery:"role"`
	Work                     string                 `firestore:"work"                        json:"work,omitempty"              bigquery:"-"`
	WorkType                 string                 `firestore:"workType"                    json:"workType,omitempty"          bigquery:"-"`
	WorkStatus               string                 `firestore:"workStatus,omitempty"        json:"workStatus,omitempty"        bigquery:"-"`
	Mail                     string                 `firestore:"mail"                        json:"mail,omitempty"              bigquery:"mail"`
	Phone                    string                 `firestore:"phone"                       json:"phone,omitempty"             bigquery:"phone"`
	FiscalCode               string                 `firestore:"fiscalCode"                  json:"fiscalCode,omitempty"        bigquery:"fiscalCode"`
	VatCode                  string                 `firestore:"vatCode"                     json:"vatCode"                     bigquery:"vatCode"`
	RiskClass                string                 `firestore:"riskClass"                   json:"riskClass,omitempty"         bigquery:"-"`
	CreationDate             time.Time              `firestore:"creationDate,omitempty"      json:"creationDate,omitempty"      bigquery:"-"`
	BigCreationDate          bigquery.NullDateTime  `firestore:"-"                           json:"-"                           bigquery:"creationDate"`
	UpdatedDate              time.Time              `firestore:"updatedDate,omitempty"       json:"updatedDate,omitempty"       bigquery:"-"`
	BigUpdatedDate           bigquery.NullDateTime  `firestore:"-"                           json:"-"                           bigquery:"updatedDate"`
	PoliciesUid              []string               `firestore:"policiesUid"                 json:"policiesUid,omitempty"       bigquery:"-"`
	Claims                   *[]Claim               `firestore:"claims"                      json:"claims,omitempty"            bigquery:"-"`
	Consens                  *[]Consens             `firestore:"consens"                     json:"consens,omitempty"           bigquery:"-"`
	IsAgent                  bool                   `firestore:"isAgent,omitempty"           json:"isAgent,omitempty"           bigquery:"-"`
	Height                   int                    `firestore:"height"                      json:"height"                      bigquery:"-"`
	Weight                   int                    `firestore:"weight"                      json:"weight"                      bigquery:"-"`
	Json                     string                 `firestore:"-"                           json:"-"                           bigquery:"-"`
	Residence                *Address               `firestore:"residence,omitempty"         json:"residence,omitempty"         bigquery:"-"`
	BigResidenceStreetName   string                 `firestore:"-"                           json:"-"                           bigquery:"residenceStreetName"`
	BigResidenceStreetNumber string                 `firestore:"-"                           json:"-"                           bigquery:"residenceStreetNumber"`
	BigResidenceCity         string                 `firestore:"-"                           json:"-"                           bigquery:"residenceCity"`
	BigResidencePostalCode   string                 `firestore:"-"                           json:"-"                           bigquery:"residencePostalCode"`
	BigResidenceLocality     string                 `firestore:"-"                           json:"-"                           bigquery:"residenceLocality"`
	BigResidenceCityCode     string                 `firestore:"-"                           json:"-"                           bigquery:"residenceCityCode"`
	Domicile                 *Address               `firestore:"domicile,omitempty"          json:"domicile,omitempty"          bigquery:"-"`
	BigDomicileStreetName    string                 `firestore:"-"                           json:"-"                           bigquery:"domicileStreetName"`
	BigDomicileStreetNumber  string                 `firestore:"-"                           json:"-"                           bigquery:"domicileStreetNumber"`
	BigDomicileCity          string                 `firestore:"-"                           json:"-"                           bigquery:"domicileCity"`
	BigDomicilePostalCode    string                 `firestore:"-"                           json:"-"                           bigquery:"domicilePostalCode"`
	BigDomicileLocality      string                 `firestore:"-"                           json:"-"                           bigquery:"domicileLocality"`
	BigDomicileCityCode      string                 `firestore:"-"                           json:"-"                           bigquery:"domicileCityCode"`
	IdentityDocuments        []*IdentityDocument    `firestore:"identityDocuments,omitempty" json:"identityDocuments,omitempty" bigquery:"-"`
	AuthId                   string                 `firestore:"authId,omitempty"            json:"authId,omitempty"            bigquery:"-"`
	Statements               []*Statement           `firestore:"statements,omitempty"        json:"statements,omitempty"        bigquery:"-"`
	IsLegalPerson            bool                   `firestore:"isLegalPerson,omitempty"     json:"isLegalPerson,omitempty"     bigquery:"-"`
	CompanyAddress           *Address               `firestore:"companyAddress,omitempty" json:"companyAddress,omitempty" bigquery:"-"`
	Data                     string                 `firestore:"-"                           json:"-"                           bigquery:"data"`
}

func (c *Contractor) ToUser() *User {
	var user User

	rawUser, err := json.Marshal(c)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(rawUser, &user)
	if err != nil {
		return nil
	}

	return &user
}

func (c *Contractor) initBigqueryData() error {
	rawContractor, err := json.Marshal(c)
	if err != nil {
		return err
	}
	c.Data = string(rawContractor)

	if c.BirthDate != "" {
		birthDate, err := time.Parse(time.RFC3339, c.BirthDate)
		if err != nil {
			return err
		}
		c.BigBirthDate = lib.GetBigQueryNullDateTime(birthDate)
	}

	if c.Residence != nil {
		c.BigResidenceStreetName = c.Residence.StreetName
		c.BigResidenceStreetNumber = c.Residence.StreetNumber
		c.BigResidenceCity = c.Residence.City
		c.BigResidencePostalCode = c.Residence.PostalCode
		c.BigResidenceLocality = c.Residence.Locality
		c.BigResidenceCityCode = c.Residence.CityCode
	}

	if c.Domicile != nil {
		c.BigDomicileStreetName = c.Domicile.StreetName
		c.BigDomicileStreetNumber = c.Domicile.StreetNumber
		c.BigDomicileCity = c.Domicile.City
		c.BigDomicilePostalCode = c.Domicile.PostalCode
		c.BigDomicileLocality = c.Domicile.Locality
		c.BigDomicileCityCode = c.Domicile.CityCode
	}

	c.BigLocation = bigquery.NullGeography{
		// TODO: Check if correct: Geography type uses the WKT format for geometry
		GeographyVal: fmt.Sprintf("POINT (%f %f)", c.Location.Lng, c.Location.Lat),
		Valid:        true,
	}
	c.BigCreationDate = lib.GetBigQueryNullDateTime(c.CreationDate)
	c.BigUpdatedDate = lib.GetBigQueryNullDateTime(c.UpdatedDate)

	return nil
}

func (c *Contractor) BigquerySave(origin string) error {
	table := lib.GetDatasetByEnv(origin, UserCollection)

	if err := c.initBigqueryData(); err != nil {
		return err
	}

	log.Println("user save big query: " + c.Uid)

	return lib.InsertRowsBigQuery(WoptaDataset, table, c)
}

func (c *Contractor) GetIdentityDocument() *IdentityDocument {
	var (
		lastUpdate       time.Time
		selectedDocument *IdentityDocument
	)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for _, identityDocument := range c.IdentityDocuments {
		if identityDocument.LastUpdate.After(lastUpdate) && identityDocument.ExpiryDate.After(today) {
			selectedDocument = identityDocument
			lastUpdate = identityDocument.LastUpdate
		}
	}
	return selectedDocument
}

func FirestoreDocumentToContractor(query *firestore.DocumentIterator) (Contractor, error) {
	var result Contractor
	userDocumentSnapshot, err := query.Next()

	if err == iterator.Done && userDocumentSnapshot == nil {
		log.Println("user not found in firebase DB")
		return result, fmt.Errorf("no user found")
	}

	if err != iterator.Done && err != nil {
		log.Println(`error happened while trying to get user`)
		return result, err
	}

	e := userDocumentSnapshot.DataTo(&result)
	if len(result.Uid) == 0 {
		result.Uid = userDocumentSnapshot.Ref.ID
	}

	return result, e
}
