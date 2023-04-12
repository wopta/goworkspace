package models

import (
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
	latlng "google.golang.org/genproto/googleapis/type/latlng"
)

func UnmarshalUser(data []byte) (Claim, error) {
	var r Claim
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *User) Marshal() ([]byte, error) {

	return json.Marshal(r)
}

type User struct {
	EmailVerified  bool          `firestore:"emailVerified" json:"emailVerified,omitempty" bigquery:"emailVerified"`
	Uid            string        `firestore:"uid" json:"uid,omitempty" bigquery:"uid" `
	BirthDate      string        `firestore:"birthDate" json:"birthDate,omitempty" bigquery:"birthDate"`
	BirthCity      string        `firestore:"birthCity" json:"birthCity,omitempty" bigquery:"birthCity"`
	BirthProvince  string        `firestore:"birthProvince" json:"birthProvince,omitempty" bigquery:"birthProvince"`
	PictureUrl     string        `firestore:"pictureUrl" json:"pictureUrl,omitempty" bigquery:"-"`
	Location       Location      `firestore:"location" json:"location,omitempty" bigquery:"-"`
	Geo            latlng.LatLng `firestore:"geo" json:"-" bigquery:"-"`
	Name           string        `firestore:"name" json:"name,omitempty" bigquery:"name"`
	Type           string        `firestore:"type" json:"type,omitempty" bigquery:"type"`
	Cluster        string        `firestore:"cluster" json:"cluster,omitempty" bigquery:"cluster"`
	Surname        string        `firestore:"surname" json:"surname,omitempty" bigquery:"surname"`
	Address        string        `firestore:"address" json:"address,omitempty" bigquery:"address"`
	PostalCode     string        `firestore:"postalCode" json:"postalCode,omitempty" bigquery:"postalCode"`
	City           string        `firestore:"city" json:"city,omitempty" bigquery:"city"`
	Locality       string        `firestore:"locality" json:"locality,omitempty" bigquery:"locality"`
	StreetNumber   string        `firestore:"streetNumber,omitempty" json:"streetNumber,omitempty" bigquery:"streetNumber"`
	CityCode       string        `firestore:"cityCode" json:"cityCode,omitempty" bigquery:"cityCode"`
	Role           string        `firestore:"role" json:"role,omitempty" bigquery:"role"`
	Work           string        `firestore:"work" json:"work,omitempty" bigquery:"work"`
	WorkType       string        `firestore:"workType" json:"workType,omitempty" bigquery:"workType"`
	Mail           string        `firestore:"mail" json:"mail,omitempty" bigquery:"mail"`
	Phone          string        `firestore:"phone" json:"phone,omitempty" bigquery:"phone"`
	FiscalCode     string        `firestore:"fiscalCode" json:"fiscalCode,omitempty" bigquery:"fiscalCode"`
	VatCode        string        `firestore:"vatCode" json:"vatCode" bigquery:"vatCode"`
	RiskClass      string        `firestore:"riskClass" json:"riskClass,omitempty" bigquery:"riskClass"`
	CreationDate   time.Time     `firestore:"creationDate,omitempty" json:"creationDate,omitempty" bigquery:"-"`
	UpdatedDate    time.Time     `firestore:"updatedDate" json:"updatedDate,omitempty" bigquery:"-"`
	PoliciesUid    []string      `firestore:"policiesUid" json:"policiesUid,omitempty" bigquery:"-"`
	BigPoliciesUid string        `firestore:"-" json:"-" bigquery:"policiesUid"`
	Claims         *[]Claim      `firestore:"claims" json:"claims,omitempty" bigquery:"-"`
	Consens        *[]Consens    `firestore:"consens" json:"consens,omitempty"`
	IsAgent        bool          `firestore:"isAgent ,omitempty" json:"isAgent,omitempty" bigquery:"isAgent"`
	Height         int           `firestore:"height" json:"height" bigquery:"height"`
	Weight         int           `firestore:"weight" json:"weight" bigquery:"weight"`
	Json           string        `firestore:"-" json:"-" bigquery:"json"`
}
type Consens struct {
	UserUid      string    `firestore:"useruid" json:"useruid,omitempty" bigquery:"useruid" `
	Title        string    `firestore:"title ,omitempty" json:"title,omitempty"`
	Consens      string    `firestore:"consens,omitempty" json:"consens,omitempty"`
	Key          int64     `firestore:"key,omitempty" json:"key,omitempty"`
	Answer       bool      `firestore:"answer,omitempty" json:"answer,omitempty"`
	CreationDate time.Time `firestore:"creationDate,omitempty" json:"creationDate,omitempty" bigquery:"-"`
}

func FirestoreDocumentToUser(query *firestore.DocumentIterator) (User, error) {
	var result User
	userDocumentSnapshot, err := query.Next()

	if (err != iterator.Done && err != nil) || userDocumentSnapshot == nil {
		log.Println(`error happened while trying to get user`)
		log.Println(err)
		return result, err
	}

	e := userDocumentSnapshot.DataTo(&result)
	if len(result.Uid) == 0 {
		result.Uid = userDocumentSnapshot.Ref.ID
	}
	lib.CheckError(e)

	return result, e
}
