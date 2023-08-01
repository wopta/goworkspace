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
	Data                     string                 `firestore:"-"                           json:"-"                           bigquery:"data"`
}

type Consens struct {
	UserUid         string         `json:"useruid,omitempty" firestore:"useruid" bigquery:"-"`
	Title           string         `json:"title,omitempty" firestore:"title,omitempty" bigquery:"-"`
	Consens         string         `json:"consens,omitempty" firestore:"consens,omitempty" bigquery:"-"`
	Key             int64          `json:"key,omitempty" firestore:"key,omitempty" bigquery:"-"`
	Answer          bool           `json:"answer" firestore:"answer" bigquery:"-"`
	CreationDate    time.Time      `json:"creationDate,omitempty" firestore:"creationDate,omitempty" bigquery:"-"`
	BigCreationDate civil.DateTime `json:"-" firestore:"-" bigquery:"-"`
	Mail            string         `json:"-" firestore:"-" bigquery:"-"`
}

type Address struct {
	StreetName   string `json:"streetName,omitempty" firestore:"streetName" bigquery:"-"`
	StreetNumber string `json:"streetNumber,omitempty" firestore:"streetNumber,omitempty" bigquery:"-"`
	City         string `json:"city,omitempty" firestore:"city" bigquery:"-"`
	PostalCode   string `json:"postalCode,omitempty" firestore:"postalCode" bigquery:"-"`
	Locality     string `json:"locality,omitempty" firestore:"locality" bigquery:"-"`
	CityCode     string `json:"cityCode,omitempty" firestore:"cityCode" bigquery:"-"`
}

func (user *User) initBigqueryData() error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	user.Data = string(userJson)

	birthDate, err := time.Parse(time.RFC3339, user.BirthDate)
	if err != nil {
		return err
	}
	user.BigBirthDate = lib.GetBigQueryNullDateTime(birthDate)

	if user.Residence != nil {
		user.BigResidenceStreetName = user.Residence.StreetName
		user.BigResidenceStreetNumber = user.Residence.StreetNumber
		user.BigResidenceCity = user.Residence.City
		user.BigResidencePostalCode = user.Residence.PostalCode
		user.BigResidenceLocality = user.Residence.Locality
		user.BigResidenceCityCode = user.Residence.CityCode
	}

	if user.Domicile != nil {
		user.BigDomicileStreetName = user.Domicile.StreetName
		user.BigDomicileStreetNumber = user.Domicile.StreetNumber
		user.BigDomicileCity = user.Domicile.City
		user.BigDomicilePostalCode = user.Domicile.PostalCode
		user.BigDomicileLocality = user.Domicile.Locality
		user.BigDomicileCityCode = user.Domicile.CityCode
	}

	user.BigLocation = bigquery.NullGeography{
		// TODO: Check if correct: Geography type uses the WKT format for geometry
		GeographyVal: fmt.Sprintf("POINT (%f %f)", user.Location.Lng, user.Location.Lat),
		Valid:        true,
	}
	user.BigCreationDate = lib.GetBigQueryNullDateTime(user.CreationDate)
	user.BigUpdatedDate = lib.GetBigQueryNullDateTime(user.UpdatedDate)

	return nil
}

func (user *User) BigquerySave(origin string) error {
	table := lib.GetDatasetByEnv(origin, UserCollection)

	if err := user.initBigqueryData(); err != nil {
		return err
	}

	log.Println("user save big query: " + user.Uid)

	return lib.InsertRowsBigQuery(WoptaDataset, table, user)
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
