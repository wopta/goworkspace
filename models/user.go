package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
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
	EmailVerified     bool                `firestore:"emailVerified"               json:"emailVerified,omitempty"`
	Uid               string              `firestore:"uid"                         json:"uid,omitempty"`
	Status            string              `firestore:"status,omitempty"            json:"status,omitempty"`
	StatusHistory     []string            `firestore:"statusHistory,omitempty"     json:"statusHistory,omitempty"`
	BirthDate         string              `firestore:"birthDate"                   json:"birthDate,omitempty"`
	BirthCity         string              `firestore:"birthCity"                   json:"birthCity,omitempty"`
	BirthProvince     string              `firestore:"birthProvince"               json:"birthProvince,omitempty"`
	PictureUrl        string              `firestore:"pictureUrl"                  json:"pictureUrl,omitempty"`
	Location          Location            `firestore:"location"                    json:"location,omitempty"`
	Geo               latlng.LatLng       `firestore:"geo"                         json:"-"`
	Name              string              `firestore:"name"                        json:"name,omitempty"`
	Gender            string              `firestore:"gender"                      json:"gender,omitempty"`
	Type              string              `firestore:"type"                        json:"type,omitempty"`
	Cluster           string              `firestore:"cluster"                     json:"cluster,omitempty"`
	Surname           string              `firestore:"surname"                     json:"surname,omitempty"`
	Address           string              `firestore:"address"                     json:"address,omitempty"`
	PostalCode        string              `firestore:"postalCode"                  json:"postalCode,omitempty"`
	City              string              `firestore:"city"                        json:"city,omitempty"`
	Locality          string              `firestore:"locality"                    json:"locality,omitempty"`
	StreetNumber      string              `firestore:"streetNumber,omitempty"      json:"streetNumber,omitempty"`
	CityCode          string              `firestore:"cityCode"                    json:"cityCode,omitempty"`
	Role              string              `firestore:"role"                        json:"role,omitempty"`
	Work              string              `firestore:"work"                        json:"work,omitempty"`
	WorkType          string              `firestore:"workType"                    json:"workType,omitempty"`
	WorkStatus        string              `firestore:"workStatus,omitempty"        json:"workStatus,omitempty"`
	Mail              string              `firestore:"mail"                        json:"mail,omitempty"`
	Phone             string              `firestore:"phone"                       json:"phone,omitempty"`
	FiscalCode        string              `firestore:"fiscalCode"                  json:"fiscalCode,omitempty"`
	VatCode           string              `firestore:"vatCode"                     json:"vatCode"`
	RiskClass         string              `firestore:"riskClass"                   json:"riskClass,omitempty"`
	CreationDate      time.Time           `firestore:"creationDate,omitempty"      json:"creationDate,omitempty"`
	UpdatedDate       time.Time           `firestore:"updatedDate,omitempty"       json:"updatedDate,omitempty"`
	PoliciesUid       []string            `firestore:"policiesUid"                 json:"policiesUid,omitempty"`
	BigPoliciesUid    string              `firestore:"-"                           json:"-"`
	Claims            *[]Claim            `firestore:"claims"                      json:"claims,omitempty"`
	Consens           *[]Consens          `firestore:"consens"                     json:"consens,omitempty"`
	IsAgent           bool                `firestore:"isAgent,omitempty"           json:"isAgent,omitempty"`
	Height            int                 `firestore:"height"                      json:"height"`
	Weight            int                 `firestore:"weight"                      json:"weight"`
	Json              string              `firestore:"-"                           json:"-"`
	Residence         *Address            `firestore:"residence,omitempty"         json:"residence,omitempty"`
	Domicile          *Address            `firestore:"domicile,omitempty"          json:"domicile,omitempty"`
	IdentityDocuments []*IdentityDocument `firestore:"identityDocuments,omitempty" json:"identityDocuments,omitempty"`
	AuthId            string              `firestore:"authId,omitempty"            json:"authId,omitempty"`
	Statements        []*Statement        `firestore:"statements,omitempty"        json:"statements,omitempty"`
	IsLegalPerson     bool                `firestore:"isLegalPerson,omitempty"     json:"isLegalPerson,omitempty"`
}

type UserBigquery struct {
	Uid                  string                 `bigquery:"uid"`
	Name                 string                 `bigquery:"name"`
	Surname              string                 `bigquery:"surname"`
	BirthDate            bigquery.NullDateTime  `bigquery:"birthDate"`
	FiscalCode           string                 `bigquery:"fiscalCode"`
	VatCode              string                 `bigquery:"vatCode"`
	ResidenceStretName   string                 `bigquery:"residenceStretName"`
	ResidenceStretNumber string                 `bigquery:"residenceStretNumber"`
	ResidenceCity        string                 `bigquery:"residenceCity"`
	ResidencePostalCode  string                 `bigquery:"residencePostalCode"`
	ResidenceLocality    string                 `bigquery:"residenceLocality"`
	ResidenceCityCode    string                 `bigquery:"residenceCityCode"`
	DomicileStretName    string                 `bigquery:"domicileStretName"`
	DomicileStretNumber  string                 `bigquery:"domicileStretNumber"`
	DomicileCity         string                 `bigquery:"domicileCity"`
	DomicilePostalCode   string                 `bigquery:"domicilePostalCode"`
	DomicileLocality     string                 `bigquery:"domicileLocality"`
	DomicileCityCode     string                 `bigquery:"domicileCityCode"`
	Phone                string                 `bigquery:"phone"`
	Location             bigquery.NullGeography `bigquery:"location"`
	CreationDate         bigquery.NullDateTime  `bigquery:"creationDate"`
	UpdatedDate          bigquery.NullDateTime  `bigquery:"updatedDate"`
	Data                 string                 `bigquery:"data"`
}

func (user User) toBigquery() (UserBigquery, error) {
	userJson, err := json.Marshal(user)
	if err != nil {
		return UserBigquery{}, err
	}
	birthDate, err := time.Parse(time.RFC3339, user.BirthDate)
	if err != nil {
		return UserBigquery{}, err
	}
	geography := bigquery.NullGeography{
		GeographyVal: "",
		Valid:        true,
	}
	return UserBigquery{
		Uid:                  user.Uid,
		Name:                 user.Name,
		Surname:              user.Surname,
		BirthDate:            lib.GetBigQueryNullDateTime(birthDate),
		FiscalCode:           user.FiscalCode,
		VatCode:              user.VatCode,
		ResidenceStretName:   user.Residence.StreetName,
		ResidenceStretNumber: user.Residence.StreetNumber,
		ResidenceCity:        user.Residence.City,
		ResidencePostalCode:  user.Residence.PostalCode,
		ResidenceLocality:    user.Residence.Locality,
		ResidenceCityCode:    user.Residence.CityCode,
		DomicileStretName:    user.Domicile.StreetName,
		DomicileStretNumber:  user.Domicile.StreetNumber,
		DomicileCity:         user.Domicile.City,
		DomicilePostalCode:   user.Domicile.PostalCode,
		DomicileLocality:     user.Domicile.Locality,
		DomicileCityCode:     user.Domicile.CityCode,
		Phone:                user.Phone,
		Location:             geography,
		CreationDate:         lib.GetBigQueryNullDateTime(user.CreationDate),
		UpdatedDate:          lib.GetBigQueryNullDateTime(user.UpdatedDate),
		Data:                 string(userJson),
	}, nil
}

func (user User) BigquerySave(origin string) error {
	table := lib.GetDatasetByEnv(origin, "user")
	agentBigquery, err := user.toBigquery()
	if err != nil {
		return err
	}

	log.Println("user save big query: " + user.Uid)

	return lib.InsertRowsBigQuery("wopta", table, agentBigquery)
}

type Consens struct {
	UserUid         string         `firestore:"useruid"                json:"useruid,omitempty"      bigquery:"useruid"`
	Title           string         `firestore:"title,omitempty"        json:"title,omitempty"        bigquery:"title"`
	Consens         string         `firestore:"consens,omitempty"      json:"consens,omitempty"      bigquery:"consens"`
	Key             int64          `firestore:"key,omitempty"          json:"key,omitempty"          bigquery:"key"`
	Answer          bool           `firestore:"answer"                 json:"answer"                 bigquery:"answer"`
	CreationDate    time.Time      `firestore:"creationDate,omitempty" json:"creationDate,omitempty" bigquery:"-"`
	BigCreationDate civil.DateTime `firestore:"-"                      json:"-"                      bigquery:"creationDate"`
	Mail            string         `firestore:"-"                      json:"-"                      bigquery:"mail"`
}

type Address struct {
	StreetName   string `firestore:"streetName"             json:"streetName,omitempty"   bigquery:"streetName"`
	StreetNumber string `firestore:"streetNumber,omitempty" json:"streetNumber,omitempty" bigquery:"streetNumber"`
	City         string `firestore:"city"                   json:"city,omitempty"         bigquery:"city"`
	PostalCode   string `firestore:"postalCode"             json:"postalCode,omitempty"   bigquery:"postalCode"`
	Locality     string `firestore:"locality"               json:"locality,omitempty"     bigquery:"locality"`
	CityCode     string `firestore:"cityCode"               json:"cityCode,omitempty"     bigquery:"cityCode"`
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
	var err error

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
