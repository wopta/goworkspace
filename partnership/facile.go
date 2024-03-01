package partnership

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func facileLifePartnership(jwtData string, policy *models.Policy, product *models.Product) error {
	var (
		person models.User
		asset  models.Asset
		claims FacileClaims
	)

	log.Println("[facileLifePartnership] decoding jwt")

	err := lib.DecryptJwt(jwtData, os.Getenv("FACILE_SIGNING_KEY"), &claims)
	if err != nil {
		log.Printf("[facileLifePartnership] could not validate facile partnership JWT - %s", err.Error())
		return err
	}

	if claims.ExpiresAt.Before(time.Now()) {
		log.Printf("[facileLifePartnership] jwt expired")
		return fmt.Errorf("jwt expired")
	}

	log.Println("[facileLifePartnership] setting person info")
	person.Name = claims.CustomerName
	person.Surname = claims.CustomerFamilyName
	person.Mail = claims.Email
	birthDate, _ := time.Parse(models.TimeDateOnly, claims.CustomerBirthDate)
	person.BirthDate = birthDate.Format(time.RFC3339)
	person.Phone = fmt.Sprintf("+39%s", claims.Mobile)
	person.Gender = claims.Gender

	person.Normalize()

	policy.Contractor = *person.ToContractor()
	asset.Person = &person
	policy.OfferlName = "default"

	log.Println("[facileLifePartnership] setting death guarantee info")

	deathGuarantee := product.Companies[0].GuaranteesMap["death"]
	deathGuarantee.Value = &models.GuaranteValue{
		Duration: &models.Duration{
			Year: claims.Duration,
		},
		SumInsuredLimitOfIndemnity: float64(claims.InsuredCapital),
	}
	asset.Guarantees = make([]models.Guarante, 0)
	asset.Guarantees = append(asset.Guarantees, *deathGuarantee)

	policy.Assets = append(policy.Assets, asset)
	policy.PartnershipData = claims.ToMap()
	return err
}

type FacileClaims struct {
	CustomerName       string `json:"customerName"`
	CustomerFamilyName string `json:"customerFamilyName"`
	CustomerBirthDate  string `json:"customerBirthDate"`
	Gender             string `json:"gender"`
	Email              string `json:"email"`
	Mobile             string `json:"mobile"`
	IsSmoker           bool   `json:"isSmoker"`
	InsuredCapital     int    `json:"insuredCapital"`
	Duration           int    `json:"duration"`
	jwt.RegisteredClaims
}

func (facileClaims FacileClaims) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	b, err := json.Marshal(facileClaims)
	lib.CheckError(err)

	err = json.Unmarshal(b, &m)
	lib.CheckError(err)

	return m
}
