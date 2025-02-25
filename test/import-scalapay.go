package test

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
)

func ImportScalapay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const (
		layout        = "1/2/2006"
		dataset       = "wopta_inclusive"
		movementTable = "bank_account_movement_scalapay"
		usersTable    = "bank_account_users_scalapay"
	)
	var e error
	log.SetPrefix("[ImportScalapay] ")
	defer log.SetPrefix("")

	log.Println(" ----start -----------------------------------------------")

	rawdata := lib.GetFilesByEnv("data/import-scalapay.csv")

	df := lib.CsvToDataframe(rawdata)

	for k, v := range df.Records() {
		log.Println(k)

		var (
			obj BankAccountMovement
		)
		stringdate := strings.Replace(v[13], "25", "2025", -1)
		startdate, e := time.Parse(layout, stringdate)
		log.Println(e)
		obj.PolicyNumber = "051114"
		obj.Uid = uuid.New().String()
		obj.Customer = "Scalapay"
		obj.Company = "axa"
		obj.PolicyType = ""
		obj.PolicyUid = ""
		obj.AssetType = ""
		obj.PolicyName = "Scalapay base"
		obj.Tenant = "Scalapay"
		obj.Name = v[8]
		obj.Surname = v[7]
		obj.FiscalCode = v[6]
		obj.Address = v[9]
		obj.BigStartDate = civil.DateTimeOf(startdate)
		obj.MovementType = "insert"
		obj.City = v[11]
		obj.CityCode = v[12]
		obj.PostalCode = v[10]
		obj.Status = "active"
		obj.Id = v[4]
		obj.Daystart = strconv.Itoa(obj.StartDate.Day())
		e = lib.InsertRowsBigQuery(dataset, usersTable, obj)
		log.Println(e)
		e = lib.InsertRowsBigQuery(dataset, movementTable, obj)
		log.Println(e)
		if e != nil {
			log.Println(" ------ Error -----------------------------------------------")
			return "", nil, e
		}
	}

	log.Println(" -----end -----------------------------------------------")
	return "", nil, e
}

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
