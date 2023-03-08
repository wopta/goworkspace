package partnership

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
		policy models.Policy
		person models.User
		asset  models.Asset
	)
	resp.Header().Set("Access-Control-Allow-Methods", "GET")

	jwtData := r.URL.Query().Get("jwt")

	// decode JWT
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.ParseWithClaims(jwtData, &BeProfClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("BEPROF_SIGNING_KEY")), nil
	})

	if claims, ok := token.Claims.(*BeProfClaims); ok && token.Valid {
		person.Name = claims.Name
		person.Surname = claims.Surname
		person.Mail = claims.Mail
		person.FiscalCode = claims.FiscalCode
		person.BirthDate = extractBirthdateFromItalianFiscalCode(claims.FiscalCode).Format(time.RFC3339)
		policy.Contractor = person
		asset.Person = &person
		policy.Assets = append(policy.Assets, asset)
	} else {
		fmt.Println(err)
		return "", nil, err
	}

	// verify if this user has already a policy from beprof

	// catalogo e servizio
	// proccedi e i miei servizi

	// call vendibility rules

	// call quoter

	p, err := policy.Marshal()

	if err != nil {
		return "", nil, err
	}

	return string(p), policy, err
}

func extractBirthdateFromItalianFiscalCode(fiscalCode string) time.Time {
	year, _ := strconv.Atoi(fiscalCode[6:8])
	month := getMonth(fiscalCode[8:9])
	day, _ := strconv.Atoi(fiscalCode[9:11])

	if day > 40 {
		day -= 40
	}

	if year < 40 {
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
}

type BeProfClaims struct {
	Name       string `json:"nome"`
	Surname    string `json:"cognome"`
	Mail       string `json:"email"`
	FiscalCode string `json:"codiceFiscale"`
	VatCode    string `json:"piva"`
	BeProfCode string `json:"codiceBeProf"`
	jwt.RegisteredClaims
}
