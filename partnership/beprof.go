package partnership

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/quote"
	"github.com/wopta/goworkspace/user"
)

func beProfLifePartnership(jwtData string, policy *models.Policy, _ *models.Product, partnershipNode *models.NetworkNode) error {
	var (
		person models.User
		asset  models.Asset
	)

	log.Println("[beProfLifePartnership] decoding jwt")

	token, err := jwt.ParseWithClaims(jwtData, &BeprofClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key, e := b64.StdEncoding.DecodeString(os.Getenv("BEPROF_SIGNING_KEY"))

		return []byte(key), e
	})

	if claims, ok := token.Claims.(*BeprofClaims); ok && token.Valid {
		log.Println("[beProfLifePartnership] setting person info")
		person.Name = claims.UserFirstname
		person.Surname = claims.UserLastname
		person.Mail = claims.UserEmail
		person.FiscalCode = claims.UserFiscalcode
		person.Address = claims.UserAddress
		person.PostalCode = claims.UserPostalcode
		person.City = claims.UserCity
		person.CityCode = claims.UserMunicipalityCode
		person.Work = claims.UserEmploymentSector
		person.VatCode = claims.UserPiva

		person.Normalize()

		if _, personData, err := user.ExtractUserDataFromFiscalCode(person); err == nil {
			person = personData
		}

		policy.Contractor = *person.ToContractor()
		asset.Person = &person
		policy.OfferlName = "default"

		policy.Assets = append(policy.Assets, asset)
		policy.PartnershipData = claims.ToMap()

		quotedPolicy, err := quote.Life(*policy, models.ECommerceChannel, partnershipNode, nil, models.ECommerceFlow)
		policy = &quotedPolicy

		return err
	}

	log.Printf("[beProfLifePartnership] could not validate beprof partnership JWT - %s", err.Error())
	return err
}

type BeprofClaims struct {
	UserBeprofid         int    `json:"user.beprofid"`
	UserFirstname        string `json:"user.firstname"`
	UserLastname         string `json:"user.lastname"`
	UserEmail            string `json:"user.email"`
	UserMobile           string `json:"user.mobile"`
	UserFiscalcode       string `json:"user.fiscalcode"`
	UserPiva             string `json:"user.piva"`
	UserProvince         string `json:"user.province"`
	UserCity             string `json:"user.city"`
	UserPostalcode       string `json:"user.postalcode"`
	UserAddress          string `json:"user.address"`
	UserMunicipalityCode string `json:"user.municipality_code"`
	UserEmploymentSector string `json:"user.employment_sector"`
	ProductCode          string `json:"product.code"`
	ProductPurchaseid    string `json:"product.purchaseid"`
	Price                string `json:"price"`
	jwt.RegisteredClaims
}

func (beprofClaims BeprofClaims) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	b, err := json.Marshal(beprofClaims)
	lib.CheckError(err)

	err = json.Unmarshal(b, &m)
	lib.CheckError(err)

	return m

}
