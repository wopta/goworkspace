package inclusive

import (
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/civil"
)

type BankAccountMovement struct {
	Tenant         string         `firestore:"-" json:"tenant" bigquery:"tenant"`
	Uid            string         `firestore:"-" json:"-" bigquery:"uid"`
	Id             string         `firestore:"-" json:"id" bigquery:"id"`
	Status         string         `firestore:"-" json:"-" bigquery:"status"`
	Name           string         `firestore:"-" json:"name,omitempty" bigquery:"name"`             //h-Nome
	Surname        string         `firestore:"-" json:"surname,omitempty" bigquery:"surname"`       //Cognome
	FiscalCode     string         `firestore:"-" json:"fiscalCode,omitempty" bigquery:"fiscalCode"` //Codice fiscale
	HypeId         string         `firestore:"-" json:"hypeId,omitempty" bigquery:"hypeId"`         //h-Ultime 3 / 5 cifre conto corrente
	StartDate      time.Time      `bigquery:"-" firestore:"-" json:"startDate,omitempty"`           //h-Data ingresso (inizio validità copertura)
	EndDate        time.Time      `bigquery:"-" firestore:"-" json:"endDate,omitempty"`
	BigStartDate   civil.DateTime `bigquery:"startDate" firestore:"-" json:"-"` //Data ingresso (inizio validità copertura)
	BigEndDate     civil.DateTime `bigquery:"endDate" firestore:"-" json:"-"`
	Address        string         `firestore:"-" json:"address,omitempty" bigquery:"address"`
	City           string         `firestore:"-" json:"city,omitempty" bigquery:"city"`
	CityCode       string         `firestore:"-" json:"cityCode,omitempty" bigquery:"cityCode"`
	PostalCode     string         `firestore:"-" json:"postalCode,omitempty" bigquery:"postalCode"`         //Data uscita ()
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
	Daystart       string         `firestore:"-" json:"-" bigquery:"daystart"`
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

type CountResponseModel struct {
	Total     int `json:"total"`
	Insert    int `json:"insert"`
	Delete    int `json:"delete"`
	Suspended int `json:"suspended"`
}
