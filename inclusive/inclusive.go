package inclusive

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/civil"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Inclusive")
	functions.HTTP("Inclusive", InclusiveFx)
}

func InclusiveFx(w http.ResponseWriter, r *http.Request) {

	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	log.Println("mail")
	log.Println(r.RequestURI)
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/bankaccount/v1/hype",
				Handler: BankAccountFx,
				Method:  "POST",
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)

}

// TO DO security,payload,error,fasature
func BankAccountFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		e     error
		eResp ErrorResponse
		b     []byte
	)
	apikey := os.Getenv("HYPE_APIKEY")
	apikeyReq := r.Header.Get("api_key")
	if apikey != apikeyReq {
		eResp = ErrorResponse{Code: 1, Type: "bad request", Message: "Name miss"}

	}
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	eResp = ErrorResponse{Code: 1, Type: "bad request", Message: "Name miss"}
	log.Println(string(r.Header))
	log.Println(string(req))
	var obj BankAccountMovement
	defer r.Body.Close()
	json.Unmarshal([]byte(req), &obj)
	if obj.Name == "" {
		eResp.Message = "field name miss"
		b, e = json.Marshal(eResp)
		e = errors.New(string(b))
		return "", nil, e
	}
	e = lib.InsertRowsBigQuery("wopta", "inclusive-axa-bank-account", obj)

	return "", nil, e
}

type BankAccountMovement struct {
	Name           string         `firestore:"-" json:"name,omitempty" bigquery:"name"`             //h-Nome
	Surname        string         `firestore:"-" json:"surname,omitempty" bigquery:"surname"`       //Cognome
	FiscalCode     string         `firestore:"-" json:"fiscalCode,omitempty" bigquery:"fiscalCode"` //Codice fiscale
	HypeId         string         `firestore:"-" json:"hypeId,omitempty" bigquery:"hypeId"`         //h-Ultime 3 / 5 cifre conto corrente
	StartDate      time.Time      `bigquery:"-" firestore:"-" json:"startDate"`                     //h-Data ingresso (inizio validità copertura)
	EndDate        time.Time      `bigquery:"-" firestore:"-" json:"endDate"`
	BigStartDate   civil.DateTime `bigquery:"startDate" firestore:"-" json:"-"`                             //Data ingresso (inizio validità copertura)
	BigEndDate     civil.DateTime `bigquery:"endDate" firestore:"-" json:"-"`                               //Data uscita ()
	MovementType   string         `firestore:"-" json:"movementType,omitempty" bigquery:"movementType"`     //Movimento (ingresso o uscita)
	PolicyNumber   string         `firestore:"-" json:"policyNumber,omitempty" bigquery:"policyNumber"`     //NUMERO POLIZZA
	PolicyType     string         `firestore:"-" json:"policyType,omitempty" bigquery:"policyType"`         //TIPOLOGIA POLIZZA
	GuaranteesCode string         `firestore:"-" json:"guaranteesCode,omitempty" bigquery:"guaranteesCode"` //CODICE CONFIGURAZIONE pacchetti
	AssetType      string         `firestore:"-" json:"assetType,omitempty" bigquery:"assetType"`           //TIPO OGGETTO ASSICURATO
}
type ErrorResponse struct {
	Code    int    `firestore:"-" json:"code,omitempty" bigquery:"name"`          //h-Nome
	Type    string `firestore:"-" json:"type,omitempty" bigquery:"surname"`       //Cognome
	Message string `firestore:"-" json:"message,omitempty" bigquery:"fiscalCode"` //Codice fiscale

}
