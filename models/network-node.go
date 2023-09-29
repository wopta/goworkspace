package models

import (
	"log"
	"time"

	"cloud.google.com/go/bigquery"
)

type NetworkNode struct {
	Uid                      string                 `json:"uid" firestore:"uid" bigquery:"uid"`
	AuthId                   string                 `json:"authId,omitempty" firestore:"authId,omitempty" bigquery:"-"`
	Code                     string                 `json:"code" firestore:"code" bigquery:"code"`
	Type                     string                 `json:"type" firestore:"type" bigquery:"type"`
	Role                     string                 `json:"role" firestore:"role" bigquery:"role"`
	NetworkUid               string                 `json:"networkUid" firestore:"networkUid" bigquery:"networkUid"`
	NetworkCode              string                 `json:"networkCode" firestore:"networkCode" bigquery:"networkCode"`
	ParentUid                string                 `json:"parentUid,omitempty" firestore:"parentUid,omitempty" bigquery:"parentUid"`
	ManagerUid               string                 `json:"managerUid,omitempty" firestore:"managerUid,omitempty" bigquery:"managerUid"`
	Name                     string                 `json:"name" firestore:"name" bigquery:"name"`
	Surname                  string                 `json:"surname,omitempty" firestore:"surname,omitempty" bigquery:"surname"`
	FiscalCode               string                 `json:"fiscalCode,omitempty" firestore:"fiscalCode,omitempty" bigquery:"fiscalCode"`
	VatCode                  string                 `json:"vatCode,omitempty" firestore:"vatCode,omitempty" bigquery:"vatCode"`
	Mail                     string                 `json:"mail" firestore:"mail" bigquery:"mail"`
	Phone                    string                 `json:"phone,omitempty" firestore:"phone,omitempty" bigquery:"phone"`
	IsActive                 bool                   `json:"isActive" firestore:"isActive" bigquery:"isActive"`
	BirthDate                string                 `json:"birthDate,omitempty" firestore:"birthDate,omitempty" bigquery:"-"`
	BigBirthDate             bigquery.NullDateTime  `json:"-" firestore:"-" bigquery:"birthDate"`
	BirthCity                string                 `json:"birthCity,omitempty" firestore:"birthCity,omitempty" bigquery:"birthCity"`
	BirthProvince            string                 `json:"birthProvince,omitempty" firestore:"birthProvince,omitempty" bigquery:"birthProvince"`
	Residence                *Address               `json:"residence,omitempty" firestore:"residence,omitempty" bigquery:"-"`
	BigResidenceStreetName   string                 `json:"-" firestore:"-" bigquery:"residenceStreetName"`
	BigResidenceStreetNumber string                 `json:"-" firestore:"-" bigquery:"residenceStreetNumber"`
	BigResidenceCity         string                 `json:"-" firestore:"-" bigquery:"residenceCity"`
	BigResidencePostalCode   string                 `json:"-" firestore:"-" bigquery:"residencePostalCode"`
	BigResidenceLocality     string                 `json:"-" firestore:"-" bigquery:"residenceLocality"`
	BigResidenceCityCode     string                 `json:"-" firestore:"-" bigquery:"residenceCityCode"`
	Domicile                 *Address               `json:"domicile,omitempty" firestore:"domicile,omitempty" bigquery:"-"`
	BigDomicileStreetName    string                 `json:"-" firestore:"-" bigquery:"domicileStreetName"`
	BigDomicileStreetNumber  string                 `json:"-" firestore:"-" bigquery:"domicileStreetNumber"`
	BigDomicileCity          string                 `json:"-" firestore:"-" bigquery:"domicileCity"`
	BigDomicilePostalCode    string                 `json:"-" firestore:"-" bigquery:"domicilePostalCode"`
	BigDomicileLocality      string                 `json:"-" firestore:"-" bigquery:"domicileLocality"`
	BigDomicileCityCode      string                 `json:"-" firestore:"-" bigquery:"domicileCityCode"`
	Location                 Location               `json:"location,omitempty" firestore:"location,omitempty" bigquery:"-"`
	BigLocation              bigquery.NullGeography `json:"-" firestore:"-" bigquery:"location"`
	Users                    []string               `json:"users" firestore:"users" bigquery:"-"`
	Products                 []Product              `json:"products" firestore:"products" bigquery:"-"`
	Policies                 []string               `json:"policies" firestore:"policies" bigquery:"-"`
	RuiCode                  string                 `json:"ruiCode" firestore:"ruiCode" bigquery:"ruiCode"`
	RuiSection               string                 `json:"ruiSection" firestore:"ruiSection" bigquery:"ruiSection"`
	RuiRegistration          time.Time              `json:"ruiRegistration" firestore:"ruiRegistration" bigquery:"-"`
	BigRuiRegistration       bigquery.NullDateTime  `json:"-" firestore:"-" bigquery:"ruiRegistration"`
	CreationDate             time.Time              `json:"creationDate" firestore:"creationDate" bigquery:"-"`
	BigCreationDate          bigquery.NullDateTime  `json:"-" firestore:"-" bigquery:"creationDate"`
	UpdatedDate              time.Time              `json:"updatedDate" firestore:"updatedDate" bigquery:"-"`
	BigUpdatedDate           bigquery.NullDateTime  `json:"-" firestore:"-" bigquery:"updatedDate"`
	NodeSetting              NodeSetting            `json:"nodeSetting,omitempty" firestore:"nodeSetting,omitempty" bigquery:"-"`
	Steps                    []Step                 `json:"steps,omitempty" firestore:"steps,omitempty" bigquery:"-"`
	Skin                     Skin                   `json:"skin,omitempty" firestore:"skin,omitempty" bigquery:"-"`
	Data                     string                 `json:"-" firestore:"-" bigquery:"data"`
}

func (nn *NetworkNode) SaveBigQuery(origin string) error {
	log.Println("[NetworkNode.SaveBigQuery] TO BE IMPLEMENTED")
	return nil
}
