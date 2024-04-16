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
	"github.com/wopta/goworkspace/quote"
	"github.com/wopta/goworkspace/user"
)

func eLeadsLifePartnership(jwtData string, policy *models.Policy, product *models.Product, partnershipNode *models.NetworkNode) error {
	var (
		person models.User
		asset  models.Asset
		claims ELeadsClaims
	)

	log.Println("[eLeadsLifePartnership] decoding jwt")

	err := lib.DecryptJwt(jwtData, os.Getenv("ELEADS_SIGNING_KEY"), &claims)
	if err != nil {
		log.Printf("[eLeadsLifePartnership] could not validate eLeads partnership JWT - %s", err.Error())
		return err
	}

	if claims.ExpiresAt.Before(time.Now()) {
		log.Printf("[eLeadsLifePartnership] jwt expired")
		return fmt.Errorf("jwt expired")
	}

	log.Println("[eLeadsLifePartnership] setting person info")
	person.Name = claims.Name
	person.Surname = claims.Surname
	person.Mail = claims.Email
	birthDate, _ := time.Parse(models.TimeDateOnly, claims.BirthDate)
	person.BirthDate = birthDate.Format(time.RFC3339)
	person.Phone = fmt.Sprintf("+39%s", claims.Phone)
	person.FiscalCode = claims.FiscalCode

	if _, personData, err := user.ExtractUserDataFromFiscalCode(person); err == nil {
		person = personData
	}

	person.Normalize()

	policy.Contractor = *person.ToContractor()
	asset.Person = &person
	policy.OfferlName = "default"

	log.Println("[eLeadsLifePartnership] setting death guarantee info")

	deathGuarantee := product.Companies[0].GuaranteesMap["death"]
	deathGuarantee.Value = &models.GuaranteValue{
		Duration: &models.Duration{
			Year: claims.Duration,
		},
		SumInsuredLimitOfIndemnity: float64(claims.SumInsuredLimitOfIndemnity),
	}
	asset.Guarantees = make([]models.Guarante, 0)
	asset.Guarantees = append(asset.Guarantees, *deathGuarantee)

	policy.Assets = append(policy.Assets, asset)
	policy.PartnershipData = claims.ToMap()

	quotedPolicy, err := quote.Life(*policy, models.ECommerceChannel, partnershipNode, nil, models.ECommerceFlow)
	*policy = quotedPolicy

	return err
}

type ELeadsClaims struct {
	Name                       string `json:"name"`
	Surname                    string `json:"surname"`
	Email                      string `json:"email"`
	Phone                      string `json:"phone"`
	FiscalCode                 string `json:"fiscalCode"`
	BirthDate                  string `json:"birthDate"`
	SumInsuredLimitOfIndemnity int    `json:"sumInsuredLimitOfIndemnity"`
	Duration                   int    `json:"duration"`
	jwt.RegisteredClaims
}

func (eLeadsClaims ELeadsClaims) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	b, err := json.Marshal(eLeadsClaims)
	lib.CheckError(err)

	err = json.Unmarshal(b, &m)
	lib.CheckError(err)

	return m

}
