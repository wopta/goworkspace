package models

import (
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/civil"
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
	EmailVerified     bool                `firestore:"emailVerified" json:"emailVerified,omitempty" bigquery:"emailVerified"`
	Uid               string              `firestore:"uid" json:"uid,omitempty" bigquery:"uid" `
	Status            string              `firestore:"status,omitempty" json:"status,omitempty" bigquery:"status"`
	StatusHistory     []string            `firestore:"statusHistory,omitempty" json:"statusHistory,omitempty" bigquery:"-"`
	BirthDate         string              `firestore:"birthDate" json:"birthDate,omitempty" bigquery:"birthDate"`
	BirthCity         string              `firestore:"birthCity" json:"birthCity,omitempty" bigquery:"birthCity"`
	BirthProvince     string              `firestore:"birthProvince" json:"birthProvince,omitempty" bigquery:"birthProvince"`
	PictureUrl        string              `firestore:"pictureUrl" json:"pictureUrl,omitempty" bigquery:"-"`
	Location          Location            `firestore:"location" json:"location,omitempty" bigquery:"-"`
	Geo               latlng.LatLng       `firestore:"geo" json:"-" bigquery:"-"`
	Name              string              `firestore:"name" json:"name,omitempty" bigquery:"name"`
	Gender            string              `firestore:"gender" json:"gender,omitempty" bigquery:"gender"`
	Type              string              `firestore:"type" json:"type,omitempty" bigquery:"type"`
	Cluster           string              `firestore:"cluster" json:"cluster,omitempty" bigquery:"cluster"`
	Surname           string              `firestore:"surname" json:"surname,omitempty" bigquery:"surname"`
	Address           string              `firestore:"address" json:"address,omitempty" bigquery:"address"`
	PostalCode        string              `firestore:"postalCode" json:"postalCode,omitempty" bigquery:"postalCode"`
	City              string              `firestore:"city" json:"city,omitempty" bigquery:"city"`
	Locality          string              `firestore:"locality" json:"locality,omitempty" bigquery:"locality"`
	StreetNumber      string              `firestore:"streetNumber,omitempty" json:"streetNumber,omitempty" bigquery:"streetNumber"`
	CityCode          string              `firestore:"cityCode" json:"cityCode,omitempty" bigquery:"cityCode"`
	Role              string              `firestore:"role" json:"role,omitempty" bigquery:"role"`
	Work              string              `firestore:"work" json:"work,omitempty" bigquery:"work"`
	WorkType          string              `firestore:"workType" json:"workType,omitempty" bigquery:"workType"`
	Mail              string              `firestore:"mail" json:"mail,omitempty" bigquery:"mail"`
	Phone             string              `firestore:"phone" json:"phone,omitempty" bigquery:"phone"`
	FiscalCode        string              `firestore:"fiscalCode" json:"fiscalCode,omitempty" bigquery:"fiscalCode"`
	VatCode           string              `firestore:"vatCode" json:"vatCode" bigquery:"vatCode"`
	RiskClass         string              `firestore:"riskClass" json:"riskClass,omitempty" bigquery:"riskClass"`
	CreationDate      time.Time           `firestore:"creationDate,omitempty" json:"creationDate,omitempty" bigquery:"-"`
	UpdatedDate       time.Time           `firestore:"updatedDate,omitempty" json:"updatedDate,omitempty" bigquery:"-"`
	PoliciesUid       []string            `firestore:"policiesUid" json:"policiesUid,omitempty" bigquery:"-"`
	BigPoliciesUid    string              `firestore:"-" json:"-" bigquery:"policiesUid"`
	Claims            *[]Claim            `firestore:"claims" json:"claims,omitempty" bigquery:"-"`
	Consens           *[]Consens          `firestore:"consens" json:"consens,omitempty"`
	IsAgent           bool                `firestore:"isAgent ,omitempty" json:"isAgent,omitempty" bigquery:"isAgent"`
	Height            int                 `firestore:"height" json:"height" bigquery:"height"`
	Weight            int                 `firestore:"weight" json:"weight" bigquery:"weight"`
	Json              string              `firestore:"-" json:"-" bigquery:"json"`
	Residence         *Address            `json:"residence,omitempty" firestore:"residence,omitempty" bigquery:"-"`
	Domicile          *Address            `json:"domicile,omitempty" firestore:"domicile,omitempty" bigquery:"-"`
	IdentityDocuments []*IdentityDocument `json:"identityDocuments,omitempty" firestore:"identityDocuments,omitempty" bigquery:"-"`
}

type Consens struct {
	UserUid         string         `firestore:"useruid" json:"useruid,omitempty" bigquery:"useruid" `
	Title           string         `firestore:"title ,omitempty" json:"title,omitempty" bigquery:"title"`
	Consens         string         `firestore:"consens,omitempty" json:"consens,omitempty" bigquery:"consens"`
	Key             int64          `firestore:"key,omitempty" json:"key,omitempty" bigquery:"key"`
	Answer          bool           `firestore:"answer,omitempty" json:"answer,omitempty" bigquery:"answer"`
	CreationDate    time.Time      `firestore:"creationDate,omitempty" json:"creationDate,omitempty" bigquery:"-" bigquery:"-"`
	BigCreationDate civil.DateTime `bigquery:"creationDate" firestore:"-" json:"-"`
	Mail            string         `firestore:"-" json:"-" bigquery:"mail"`
}

type Address struct {
	StreetName   string `firestore:"streetName" json:"streetName,omitempty" bigquery:"streetName"`
	StreetNumber string `firestore:"streetNumber,omitempty" json:"streetNumber,omitempty" bigquery:"streetNumber"`
	City         string `firestore:"city" json:"city,omitempty" bigquery:"city"`
	PostalCode   string `firestore:"postalCode" json:"postalCode,omitempty" bigquery:"postalCode"`
	Locality     string `firestore:"locality" json:"locality,omitempty" bigquery:"locality"`
	CityCode     string `firestore:"cityCode" json:"cityCode,omitempty" bigquery:"cityCode"`
}

type IdentityDocument struct {
	Code             int       `json:"code" firestore:"code" bigquery:"-"`
	Type             string    `json:"type" firestore:"type" bigquery:"-"`
	Number           string    `json:"number" firestore:"number" bigquery:"-"`
	IssuingAuthority string    `json:"issuingAuthority" firestore:"issuingAuthority" bigquery:"-"`
	PlaceOfIssue     string    `json:"placeOfIssue" firestore:"placeOfIssue" bigquery:"-"`
	DateOfIssue      time.Time `json:"dateOfIssue" firestore:"dateOfIssue" bigquery:"-"`
	ExpiryDate       time.Time `json:"expiryDate" firestore:"expiryDate" bigquery:"-"`
	Link             string    `json:"link,omitempty" firestore:"link,omitempty" bigquery:"-"`
	MimeType         string    `json:"mimeType,omitempty" firestore:"mimeType,omitempty" bigquery:"-"`
	Base64Encoding   string    `json:"base64Encoding,omitempty" firestore:"-" bigquery:"-"`
}

func FirestoreDocumentToUser(query *firestore.DocumentIterator) (User, error) {
	var result User
	userDocumentSnapshot, err := query.Next()

	if err == iterator.Done && userDocumentSnapshot == nil {
		log.Println("user not found in firebase DB")
		return result, err
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

func UserUpdateByFiscalcode(user User) (string, error) {
	var (
		useruid string
		e       error
	)
	docsnap := lib.WhereFirestore("users", "fiscalCode", "==", user.FiscalCode)
	userL, e := FirestoreDocumentToUser(docsnap)
	usersFire := lib.GetDatasetByEnv(user.Name, "users")
	if len(user.Uid) == 0 {
		user.CreationDate = time.Now()
		user.UpdatedDate = time.Now()
		ref2, _ := lib.PutFirestore(usersFire, user)
		log.Println("Proposal User uid", ref2)
		useruid = ref2.ID
	} else {
		useruid = user.Uid
		userL.UpdatedDate = time.Now()
		_, e = lib.FireUpdate(usersFire, useruid, userL)
	}
	return useruid, e
}
