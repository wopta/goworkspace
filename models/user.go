package models

import (
	"encoding/json"
	"log"

	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
	"cloud.google.com/go/firestore"
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
	EmailVerified bool     `firestore:"emailVerified" json:"emailVerified,omitempty"`
	Uid           string   `firestore:"uid" json:"uid,omitempty"`
	BirthDate     string   `firestore:"birthDate" json:"birthDate,omitempty"`
	BirthCity     string   `firestore:"birthCity" json:"birthCity,omitempty"`
	BirthProvince string   `firestore:"birthProvince" json:"birthProvince,omitempty"`
	PictureUrl    string   `firestore:"pictureUrl" json:"pictureUrl,omitempty"`
	Name          string   `firestore:"name" json:"name,omitempty"`
	Type          string   `firestore:"type" json:"type,omitempty"`
	Cluster       string   `firestore:"cluster" json:"cluster,omitempty"`
	Surname       string   `firestore:"surname" json:"surname,omitempty"`
	Address       string   `firestore:"address" json:"address,omitempty"`
	PostalCode    string   `firestore:"postalCode" json:"postalCode,omitempty"`
	Role          string   `firestore:"role" json:"role,omitempty"`
	Work          string   `firestore:"work" json:"work,omitempty"`
	WorkType      string   `firestore:"workType" json:"workType,omitempty"`
	Mail          string   `firestore:"mail" json:"mail,omitempty"`
	Phone         string   `firestore:"phone" json:"phone,omitempty"`
	FiscalCode    string   `firestore:"fiscalCode" json:"fiscalCode,omitempty"`
	VatCode       string   `firestore:"vatCode" json:"vatCode,omitempty"`
	RiskClass     string   `firestore:"riskClass" json:"riskClass,omitempty"`
	CreationDate  string   `firestore:"creationDate,omitempty" json:"creationDate,omitempty"`
	UpdatedDate   string   `firestore:"updatedDate" json:"updatedDate,omitempty"`
	PoliciesUid   []string `firestore:"policiesUid" json:"policiesUid,omitempty"`
	Claims        []Claim  `firestore:"claims" json:"claims,omitempty"`
	IsAgent       bool     `firestore:"isEmit,omitempty" json:"isEmit,omitempty" bigquery:"isEmit"`
}

func FirestoreDocumentToUser(query *firestore.DocumentIterator) (User, error) {
	var result User
	userDocumentSnapshot, err := query.Next()

	if err != iterator.Done && err != nil {
		log.Println(`error happened while trying to get user`)
		log.Println(err)
		return result, err
	}

	e := userDocumentSnapshot.DataTo(&result)
	lib.CheckError(e)

	return result, e
}