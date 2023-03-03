package partnership

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
		policy             models.Policy
		person             models.User
		asset              models.Asset
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
		fmt.Printf("%v", claims)
		// create policy object
		person.BirthDate = claims.BirthDate
		person.Name = claims.Name
		person.Surname = claims.Surname
		person.Mail = claims.Mail
		person.FiscalCode = claims.FiscalCode
		policy.Contractor = person
		asset.Person = &person
		policy.Assets = append(policy.Assets, asset)
	} else {
		fmt.Println(err)
		return "", nil, err
	}

	// call vendibility rules

	// call quoter

	p, err := policy.Marshal()
	
	if err != nil {
		return "", nil, err
	}

	return string(p), policy, err
}

type BeProfClaims struct {
	Name       string `json:"nome"`
	Surname    string `json:"cognome"`
	BirthDate  string `json:"dataDiNascita"`
	Mail       string `json:"email"`
	FiscalCode string `json:"codiceFiscale"`
	VatCode    string `json:"piva"`
	BeProfCode string `json:"codiceBeProf"`
	jwt.RegisteredClaims
}
