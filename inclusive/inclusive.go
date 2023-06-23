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
		e   error
		obj BankAccountMovement
	)
	e = CheckApikey(r)

	obj, e = CheckData(r)

	e = lib.InsertRowsBigQuery("wopta", "inclusive-axa-bank-account", obj)

	return "", nil, e
}

type BankAccountMovement struct {
	Uid            string         `firestore:"-" json:"-" bigquery:"uid"`
	Status         string         `firestore:"-" json:"-" bigquery:"status"`
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
	Code    int    `firestore:"-" json:"code,omitempty" bigquery:"name"`
	Type    string `firestore:"-" json:"type,omitempty" bigquery:"surname"`
	Message string `firestore:"-" json:"message,omitempty" bigquery:"fiscalCode"`
}

func GetErrorJson(code int, typeEr string, message string) error {
	var (
		e     error
		eResp ErrorResponse
		b     []byte
	)
	eResp = ErrorResponse{Code: code, Type: typeEr, Message: message}
	b, e = json.Marshal(eResp)
	e = errors.New(string(b))
	return e
}
func CheckApikey(r *http.Request) error {
	apikey := os.Getenv("HYPE_APIKEY")
	apikeyReq := r.Header.Get("api_key")
	if apikey != apikeyReq {
		return GetErrorJson(401, "Unauthorized", "")
	}
	return nil
}
func CheckData(r *http.Request) (BankAccountMovement, error) {
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(r.Header)
	log.Println(string(req))
	var obj BankAccountMovement
	defer r.Body.Close()
	json.Unmarshal([]byte(req), &obj)
	if obj.Name == "" {
		return obj, GetErrorJson(400, "Bad request", "field name miss")
	}
	if obj.FiscalCode == "" {
		return obj, GetErrorJson(400, "Bad request", "field name miss")
	}
	if obj.Surname == "" {
		return obj, GetErrorJson(400, "Bad request", "field name miss")
	}
	if obj.MovementType != "inclusive" || obj.MovementType != "delete" {
		return obj, GetErrorJson(400, "Bad request", "field MovementType out of enum")
	}
	if obj.StartDate.IsZero() {
		return obj, GetErrorJson(400, "Bad request", "field StartDate miss")
	}
	if obj.EndDate.IsZero() {
		return obj, GetErrorJson(400, "Bad request", "field EndDate miss")
	}
	if obj.HypeId == "" {
		return obj, GetErrorJson(400, "Bad request", "field HypeId miss")
	}
	if obj.GuaranteesCode != "next" || obj.GuaranteesCode != "premium" {
		return obj, GetErrorJson(400, "Bad request", "field GuaranteesCode out of enum")
	}
	return obj, nil
}
func SetData(obj *BankAccountMovement) error {

	return nil
}
