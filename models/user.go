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

func (u *User) Sanitize() {
	u.Name = lib.ToUpper(u.Name)
	u.Surname = lib.ToUpper(u.Surname)
	u.Gender = lib.ToUpper(u.Gender)
	u.FiscalCode = lib.ToUpper(u.FiscalCode)
	u.VatCode = lib.ToUpper(u.VatCode)
	u.Mail = lib.ToUpper(u.Mail)
	u.Phone = lib.TrimSpace(u.Phone)
	if u.Residence != nil {
		u.Residence.Sanitize()
	}
	if u.Domicile != nil {
		u.Domicile.Sanitize()
	}
	u.BirthDate = lib.TrimSpace(u.BirthDate)
	u.BirthCity = lib.ToUpper(u.BirthCity)
	u.BirthProvince = lib.ToUpper(u.BirthProvince)
	u.Work = lib.TrimSpace(u.Work)
	u.WorkType = lib.TrimSpace(u.WorkType)
	u.WorkStatus = lib.TrimSpace(u.WorkStatus)
	for index, _ := range u.IdentityDocuments {
		u.IdentityDocuments[index].Sanitize()
	}
}

func UnmarshalUser(data []byte) (Claim, error) {
	var r Claim
	err := json.Unmarshal(data, &r)
	return r, err
}

func (u *User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) initBigqueryData() error {
	userJson, err := json.Marshal(u)
	if err != nil {
		return err
	}
	u.Data = string(userJson)

	if u.BirthDate != "" {
		birthDate, err := time.Parse(time.RFC3339, u.BirthDate)
		if err != nil {
			return err
		}
		u.BigBirthDate = lib.GetBigQueryNullDateTime(birthDate)
	}

	if u.Residence != nil {
		u.BigResidenceStreetName = u.Residence.StreetName
		u.BigResidenceStreetNumber = u.Residence.StreetNumber
		u.BigResidenceCity = u.Residence.City
		u.BigResidencePostalCode = u.Residence.PostalCode
		u.BigResidenceLocality = u.Residence.Locality
		u.BigResidenceCityCode = u.Residence.CityCode
	}

	if u.Domicile != nil {
		u.BigDomicileStreetName = u.Domicile.StreetName
		u.BigDomicileStreetNumber = u.Domicile.StreetNumber
		u.BigDomicileCity = u.Domicile.City
		u.BigDomicilePostalCode = u.Domicile.PostalCode
		u.BigDomicileLocality = u.Domicile.Locality
		u.BigDomicileCityCode = u.Domicile.CityCode
	}

	u.BigLocation = bigquery.NullGeography{
		// TODO: Check if correct: Geography type uses the WKT format for geometry
		GeographyVal: fmt.Sprintf("POINT (%f %f)", u.Location.Lng, u.Location.Lat),
		Valid:        true,
	}
	u.BigCreationDate = lib.GetBigQueryNullDateTime(u.CreationDate)
	u.BigUpdatedDate = lib.GetBigQueryNullDateTime(u.UpdatedDate)

	return nil
}

func (u *User) BigquerySave(origin string) error {
	table := lib.GetDatasetByEnv(origin, UserCollection)

	if err := u.initBigqueryData(); err != nil {
		return err
	}

	log.Println("user save big query: " + u.Uid)

	return lib.InsertRowsBigQuery(WoptaDataset, table, u)
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
	fireUser := lib.GetDatasetByEnv(origin, UserCollection)
	docSnap := lib.WhereFirestore(fireUser, "fiscalCode", "==", user.FiscalCode)
	retrievedUser, err := FirestoreDocumentToUser(docSnap)

	if retrievedUser.Uid != "" {
		retrievedUser.IdentityDocuments = append(retrievedUser.IdentityDocuments, user.IdentityDocuments...)
		retrievedUser.Consens = updateUserConsens(retrievedUser.Consens, user.Consens)
		retrievedUser.Address = user.Address
		retrievedUser.PostalCode = user.PostalCode
		retrievedUser.City = user.City
		retrievedUser.Locality = user.Locality
		retrievedUser.CityCode = user.CityCode
		retrievedUser.StreetNumber = user.StreetNumber
		retrievedUser.Location = user.Location
		retrievedUser.Residence = user.Residence
		retrievedUser.Domicile = user.Domicile
		retrievedUser.Phone = user.Phone
		retrievedUser.UpdatedDate = time.Now().UTC()
		if user.Height != 0 {
			retrievedUser.Height = user.Height
		}
		if user.Weight != 0 {
			retrievedUser.Weight = user.Weight
		}
		err = lib.SetFirestoreErr(fireUser, retrievedUser.Uid, retrievedUser)
		if err != nil {
			return "", fmt.Errorf("[UpdateUserByFiscalCode] error saving user %s into firestore %s", retrievedUser.Uid,
				err.Error())
		}

		err = retrievedUser.BigquerySave(origin)
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
