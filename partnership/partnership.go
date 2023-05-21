package partnership

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/golang-jwt/jwt/v4"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Partnership")
	functions.HTTP("Partnership", Partnership)
}

func Partnership(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/life",
				Handler: LifePartnershipFx,
				Method:  "GET",
			},
		},
	}
	route.Router(w, r)
}

func LifePartnershipFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy   models.Policy
		person   models.User
		asset    models.Asset
		response PartnershipResponse
	)
	resp.Header().Set("Access-Control-Allow-Methods", "GET")

	jwtData := r.URL.Query().Get("jwt")

	// decode JWT
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.ParseWithClaims(jwtData, &BeprofClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		key, e := b64.StdEncoding.DecodeString(os.Getenv("BEPROF_SIGNING_KEY"))

		return []byte(key), e
	})

	if claims, ok := token.Claims.(*BeprofClaims); ok && token.Valid {
		person.Name = claims.UserFirstname
		person.Surname = claims.UserLastname
		person.Mail = claims.UserEmail
		person.FiscalCode = claims.UserFiscalcode
		person.BirthDate = lib.ExtractBirthdateFromItalianFiscalCode(claims.UserFiscalcode).Format(time.RFC3339)
		person.Phone = claims.UserMobile
		person.Address = claims.UserAddress
		person.PostalCode = claims.UserPostalcode
		person.City = claims.UserCity
		person.CityCode = claims.UserMunicipalityCode
		person.Work = claims.UserEmploymentSector
		person.VatCode = claims.UserPiva

		policy.Contractor = person
		asset.Person = &person
		policy.Assets = append(policy.Assets, asset)
		policy.PartnershipName = models.PartnershipBeProf
		policy.PartnershipData = claims.ToMap()

		response.Policy = policy
		response.Step = 1
	} else {
		fmt.Println(err)
		return "{}", nil, nil
	}

	// verify if this user has already a policy from beprof

	// catalogo e servizio
	// proccedi e i miei servizi

	// call vendibility rules

	// call quoter

	p, err := json.Marshal(response)

	if err != nil {
		return "", nil, err
	}

	return string(p), policy, err
}

/*func extractBirthdateFromItalianFiscalCode(fiscalCode string) time.Time {
	year, _ := strconv.Atoi(fiscalCode[6:8])
	month := getMonth(fiscalCode[8:9])
	day, _ := strconv.Atoi(fiscalCode[9:11])

	if day > 40 {
		day -= 40
	}

	if year < time.Now().Year()-2000 {
		year += 2000
	} else {
		year += 1900
	}

	birthdate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return birthdate
}

func getMonth(monthCode string) int {
	monthMap := map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
		"D": 4,
		"E": 5,
		"H": 6,
		"L": 7,
		"M": 8,
		"P": 9,
		"R": 10,
		"S": 11,
		"T": 12,
	}

	return monthMap[strings.ToUpper(monthCode)]
}*/

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

func (bpc BeprofClaims) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	b, err := json.Marshal(bpc)
	lib.CheckError(err)

	err = json.Unmarshal(b, &m)
	lib.CheckError(err)

	return m

}

type PartnershipResponse struct {
	Policy models.Policy `json:"policy"`
	Step   int           `json:"step"`
}
