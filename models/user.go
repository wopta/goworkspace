package models

import (
	"encoding/json"
	"fmt"
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
	WorkStatus        string              `json:"workStatus,omitempty" firestore:"workStatus,omitempty" bigquery:"-"`
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
	AuthId            string              `json:"authId,omitempty" firestore:"authId,omitempty" bigquery:"-"`
	Statements        []*Statement        `json:"statements,omitempty" firestore:"statements,omitempty" bigquery:"-"`
}

type Consens struct {
	UserUid         string         `firestore:"useruid" json:"useruid,omitempty" bigquery:"useruid" `
	Title           string         `firestore:"title ,omitempty" json:"title,omitempty" bigquery:"title"`
	Consens         string         `firestore:"consens,omitempty" json:"consens,omitempty" bigquery:"consens"`
	Key             int64          `firestore:"key,omitempty" json:"key,omitempty" bigquery:"key"`
	Answer          bool           `firestore:"answer" json:"answer" bigquery:"answer"`
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

func (u *User) GetIdentityDocument() *IdentityDocument {
	var (
		lastUpdate       time.Time
		selectedDocument *IdentityDocument
	)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for _, identityDocument := range u.IdentityDocuments {
		if identityDocument.LastUpdate.After(lastUpdate) && identityDocument.ExpiryDate.After(today) {
			selectedDocument = identityDocument
			lastUpdate = identityDocument.LastUpdate
		}
	}
	return selectedDocument
}

func FirestoreDocumentToUser(query *firestore.DocumentIterator) (User, error) {
	var result User
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

func GetUserUIDByFiscalCode(origin string, fiscalCode string) (string, bool, error) {
	usersFire := lib.GetDatasetByEnv(origin, "users")
	docSnap := lib.WhereFirestore(usersFire, "fiscalCode", "==", fiscalCode)
	retrievedUser, err := FirestoreDocumentToUser(docSnap)
	if err != nil && err.Error() != "no user found" {
		return "", false, err
	}
	if retrievedUser.Uid != "" {
		return retrievedUser.Uid, false, nil
	}
	return lib.NewDoc(usersFire), true, nil
}

func UpdateUserByFiscalCode(origin string, user User) (string, error) {
	var (
		err error
	)

	usersFire := lib.GetDatasetByEnv(origin, "users")
	docSnap := lib.WhereFirestore(usersFire, "fiscalCode", "==", user.FiscalCode)
	retrievedUser, err := FirestoreDocumentToUser(docSnap)
	if retrievedUser.Uid != "" {
		for _, identityDocument := range user.IdentityDocuments {
			retrievedUser.IdentityDocuments = append(retrievedUser.IdentityDocuments, identityDocument)
		}

		retrievedUser.Consens = updateUserConsens(retrievedUser.Consens, user.Consens)

		updatedUser := map[string]interface{}{
			"address":           user.Address,
			"postalCode":        user.PostalCode,
			"city":              user.City,
			"locality":          user.Locality,
			"cityCode":          user.CityCode,
			"streetNumber":      user.StreetNumber,
			"location":          user.Location,
			"identityDocuments": retrievedUser.IdentityDocuments,
			"consens":           retrievedUser.Consens,
			"residence":         user.Residence,
			"domicile":          user.Domicile,
			"updatedDate":       time.Now().UTC(),
		}
		if user.Height != 0 {
			updatedUser["height"] = user.Height
		}
		if user.Weight != 0 {
			updatedUser["weight"] = user.Weight
		}

		_, err = lib.FireUpdate(usersFire, retrievedUser.Uid, updatedUser)
		return retrievedUser.Uid, err
	}

	return "", fmt.Errorf("no user found with this fiscal code")
}

func updateUserConsens(oldConsens *[]Consens, newConsens *[]Consens) *[]Consens {
	if newConsens == nil {
		return oldConsens
	}
	if oldConsens == nil {
		return newConsens
	}
	for _, consens := range *newConsens {
		found := false
		for index, savedConsens := range *oldConsens {
			if consens.Key == savedConsens.Key {
				(*oldConsens)[index].Answer = consens.Answer
				(*oldConsens)[index].Title = consens.Title
				found = true
			}
		}
		if !found {
			*oldConsens = append(*oldConsens, consens)
		}
	}
	return oldConsens
}

func UsersToListData(query *firestore.DocumentIterator) []User {
	result := make([]User, 0)
	for {
		d, err := query.Next()
		if err != nil {
			log.Println("error")
			if err == iterator.Done {
				log.Println("iterator.Done")
				break
			}
			break
		} else {
			var value User
			e := d.DataTo(&value)
			log.Println("todata")
			lib.CheckError(e)
			result = append(result, value)
			log.Printf("len result: %d\n", len(result))
		}
	}
	return result
}
