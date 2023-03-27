package models

import (
	"encoding/json"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
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
	EmailVerified bool       `firestore:"emailVerified" json:"emailVerified,omitempty"`
	Uid           string     `firestore:"uid" json:"uid,omitempty"`
	BirthDate     string     `firestore:"birthDate" json:"birthDate,omitempty"`
	BirthCity     string     `firestore:"birthCity" json:"birthCity,omitempty"`
	BirthProvince string     `firestore:"birthProvince" json:"birthProvince,omitempty"`
	PictureUrl    string     `firestore:"pictureUrl" json:"pictureUrl,omitempty"`
	Location      Location   `firestore:"location" json:"location,omitempty"`
	Name          string     `firestore:"name" json:"name,omitempty"`
	Type          string     `firestore:"type" json:"type,omitempty"`
	Cluster       string     `firestore:"cluster" json:"cluster,omitempty"`
	Surname       string     `firestore:"surname" json:"surname,omitempty"`
	Address       string     `firestore:"address" json:"address,omitempty"`
	PostalCode    string     `firestore:"postalCode" json:"postalCode,omitempty"`
	City          string     `firestore:"city" json:"city,omitempty"`
	Locality      string     `firestore:"locality" json:"locality,omitempty"`
	StreetNumber  string     `firestore:"streetNumber,omitempty" json:"streetNumber,omitempty"`
	CityCode      string     `firestore:"cityCode" json:"cityCode,omitempty"`
	Role          string     `firestore:"role" json:"role,omitempty"`
	Work          string     `firestore:"work" json:"work,omitempty"`
	WorkType      string     `firestore:"workType" json:"workType,omitempty"`
	Mail          string     `firestore:"mail" json:"mail,omitempty"`
	Phone         string     `firestore:"phone" json:"phone,omitempty"`
	FiscalCode    string     `firestore:"fiscalCode" json:"fiscalCode,omitempty"`
	VatCode       string     `firestore:"vatCode" json:"vatCode,omitempty"`
	RiskClass     string     `firestore:"riskClass" json:"riskClass,omitempty"`
	CreationDate  string     `firestore:"creationDate,omitempty" json:"creationDate,omitempty"`
	UpdatedDate   string     `firestore:"updatedDate" json:"updatedDate,omitempty"`
	PoliciesUid   []string   `firestore:"policiesUid" json:"policiesUid,omitempty"`
	Claims        *[]Claim   `firestore:"claims" json:"claims,omitempty"`
	Consens       *[]Consens `firestore:"consens" json:"consens,omitempty"`
	IsAgent       bool       `firestore:"isAgent ,omitempty" json:"isAgent ,omitempty" bigquery:"isEmit"`
}
type Consens struct {
	Title   string `firestore:"title ,omitempty" json:"title,omitempty"`
	Consens string `firestore:"consens,omitempty" json:"consens,omitempty"`
	Key     int64  `firestore:"key,omitempty" json:"key,omitempty"`
	Answer  bool   `firestore:"answer,omitempty" json:"answer,omitempty"`
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
