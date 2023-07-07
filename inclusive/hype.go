package inclusive

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	lib "github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

const (
	dataMovement     = "inclusive_bank_account_movement"
	dataBanckAccount = "inclusive_bank_account"
)

// TO DO security,payload,error,fasature
func BankAccountHypeFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	layout := "2006-01-02"
	layout2 := "HH:mm:ss"

	var (
		e   error
		obj BankAccountMovement
	)
	e = CheckApikey(r)
	if e != nil {
		return "", nil, e
	}
	obj, e = CheckData(r)
	if e != nil {
		return "", nil, e
	}
	obj = SetData(obj)
	e = lib.InsertRowsBigQuery("wopta", dataMovement, obj)
	if obj.MovementType == "insert" {

		e = lib.InsertRowsBigQuery("wopta", dataBanckAccount, obj)
	}
	if obj.MovementType == "delete" {
		e = lib.UpdateRowBigQuery("wopta", dataBanckAccount, map[string]string{
			"status":  obj.Status,
			"endDate": obj.EndDate.Format(layout) + " " + obj.EndDate.Format(layout2),
		}, "fiscalCode='"+obj.FiscalCode+"' and guaranteesCode='"+obj.GuaranteesCode+"'")

	}

	return `{"woptaUid":"` + obj.Uid + `"}`, nil, e
}

type BankAccountMovement struct {
	Uid            string         `firestore:"-" json:"-" bigquery:"uid"`
	Status         string         `firestore:"-" json:"-" bigquery:"status"`
	Name           string         `firestore:"-" json:"name,omitempty" bigquery:"name"`             //h-Nome
	Surname        string         `firestore:"-" json:"surname,omitempty" bigquery:"surname"`       //Cognome
	FiscalCode     string         `firestore:"-" json:"fiscalCode,omitempty" bigquery:"fiscalCode"` //Codice fiscale
	HypeId         string         `firestore:"-" json:"hypeId,omitempty" bigquery:"hypeId"`         //h-Ultime 3 / 5 cifre conto corrente
	StartDate      time.Time      `bigquery:"-" firestore:"-" json:"startDate,omitempty"`           //h-Data ingresso (inizio validità copertura)
	EndDate        time.Time      `bigquery:"-" firestore:"-" json:"endDate,omitempty"`
	BigStartDate   civil.DateTime `bigquery:"startDate" firestore:"-" json:"-"`                             //Data ingresso (inizio validità copertura)
	BigEndDate     civil.DateTime `bigquery:"endDate" firestore:"-" json:"-"`                               //Data uscita ()
	MovementType   string         `firestore:"-" json:"movementType,omitempty" bigquery:"movementType"`     //Movimento (ingresso o uscita)
	PolicyNumber   string         `firestore:"-" json:"policyNumber,omitempty" bigquery:"policyNumber"`     //NUMERO POLIZZA
	PolicyType     string         `firestore:"-" json:"policyType,omitempty" bigquery:"policyType"`         //TIPOLOGIA POLIZZA
	GuaranteesCode string         `firestore:"-" json:"guaranteesCode,omitempty" bigquery:"guaranteesCode"` //CODICE CONFIGURAZIONE pacchetti
	AssetType      string         `firestore:"-" json:"assetType,omitempty" bigquery:"assetType"`           //TIPO OGGETTO ASSICURATO
	Customer       string         `firestore:"-" json:"-" bigquery:"customer"`
	Company        string         `firestore:"-" json:"-" bigquery:"company"`   //Hype
	PolicyUid      string         `firestore:"-" json:"-" bigquery:"policyUid"` //NUMERO POLIZZA
	CustomerId     string         `firestore:"-" json:"-" bigquery:"customerId"`
	BanckAccountId string         `firestore:"-" json:"-" bigquery:"banckAccountId"`
	PolicyName     string         `firestore:"-" json:"-" bigquery:"policyName"`
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
	if obj.MovementType != "insert" && obj.MovementType != "delete" {
		return obj, GetErrorJson(400, "Bad request", "field MovementType out of enum")
	}
	if obj.MovementType == "insert" {
		if obj.StartDate.IsZero() {
			return obj, GetErrorJson(400, "Bad request", "field StartDate miss")
		}
	}
	if obj.MovementType == "delete" {
		res, _ := QueryRowsBigQuery[BankAccountMovement]("wopta",
			"inclusive_axa_bank_account",
			"select * from `wopta."+dataMovement+"` where fiscalCode='"+obj.FiscalCode+"' and guaranteesCode ='"+obj.GuaranteesCode+"'")
		log.Println(len(res))
		if len(res) == 0 {
			return obj, GetErrorJson(400, "Bad request", "insert movement miss")
		}
		if obj.StartDate.IsZero() {

			return obj, GetErrorJson(400, "Bad request", "field StartDate miss")
		}
		if obj.EndDate.IsZero() {
			return obj, GetErrorJson(400, "Bad request", "field EndDate miss")
		}
	}

	if obj.HypeId == "" {
		return obj, GetErrorJson(400, "Bad request", "field HypeId miss")
	}
	if obj.GuaranteesCode != "next" && obj.GuaranteesCode != "premium" {
		return obj, GetErrorJson(400, "Bad request", "field GuaranteesCode out of enum")
	}

	return obj, nil
}
func SetData(obj BankAccountMovement) BankAccountMovement {
	obj.BigStartDate = civil.DateTimeOf(obj.StartDate)
	obj.BigEndDate = civil.DateTimeOf(obj.EndDate)
	obj.PolicyNumber = "180623"
	obj.Uid = uuid.New().String()
	obj.Customer = "hype"
	obj.Company = "axa"
	obj.PolicyType = ""
	obj.PolicyUid = ""
	obj.AssetType = ""
	obj.PolicyName = "Hype Next"
	obj.CustomerId = obj.HypeId
	if obj.MovementType == "insert" {

		obj.Status = "active"
	}
	if obj.MovementType == "delete" {
		obj.Status = "delete"

	}

	return obj
}
func QueryRowsBigQuery[T any](datasetID string, tableID string, query string) ([]T, error) {
	var (
		res  []T
		e    error
		iter *bigquery.RowIterator
	)
	log.Println(query)
	client := getBigqueryClient()
	ctx := context.Background()
	defer client.Close()
	queryi := client.Query(query)
	iter, e = queryi.Read(ctx)
	log.Println(e)
	for {
		var row T
		e := iter.Next(&row)
		log.Println(e)
		if e == iterator.Done {
			return res, e
		}
		if e != nil {
			return res, e
		}
		log.Println(e)
		res = append(res, row)

	}

}

func getBigqueryClient() *bigquery.Client {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	lib.CheckError(err)
	return client
}
