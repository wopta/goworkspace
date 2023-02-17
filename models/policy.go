package models

import (
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

func UnmarshalPolicy(data []byte) (Policy, error) {
	var r Policy
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Policy) Marshal() ([]byte, error) {

	return json.Marshal(r)
}

type Policy struct {
	ID              string                 `firestore:"id,omitempty" json:"id,omitempty" bigquery:"id"`
	IdSign          string                 `firestore:"idSign,omitempty" json:"idSign,omitempty" bigquery:"idSign"`
	IdPay           string                 `firestore:"idPay,omitempty" json:"idPay,omitempty" bigquery:"idPay"`
	QuoteQuestions  map[string]interface{} `firestore:"quoteQuestions,omitempty" json:"quoteQuestions,omitempty" bigquery:"-"`
	ContractFileId  string                 `firestore:"contractFileId,omitempty" json:"contractFileId,omitempty" bigquery:"contractFileId"`
	Uid             string                 `firestore:"uid,omitempty" json:"uid,omitempty" bigquery:"uid"`
	ProductUid      string                 `firestore:"productUid,omitempty" json:"productUid,omitempty" bigquery:"productUid"`
	ProductVersion  int                    `firestore:"productVersion,omitempty" json:"productVersion,omitempty" bigquery:"productVersion"`
	ProposalNumber  int                    `firestore:"proposalNumber,omitempty" json:"proposalNumber,omitempty" bigquery:"proposalNumber"`
	Number          int                    `firestore:"number,omitempty" json:"number,omitempty" bigquery:"number"`
	NumberCompany   string                 `firestore:"numberCompany,omitempty" json:"numberCompany,omitempty" bigquery:"numberCompany"`
	Status          string                 `firestore:"status,omitempty" json:"status,omitempty" bigquery:"status"`
	StatusHistory   []string               `firestore:"statusHistory,omitempty" json:"statusHistory ,omitempty" bigquery:"-"`
	RenewHistory    *[]RenewHistory        `firestore:"renewHistory,omitempty" json:"renewHistory,omitempty" bigquery:"-"`
	Transactions    *[]Transaction         `firestore:"transactions,omitempty" json:"transactions,omitempty" bigquery:"-"`
	TransactionsUid *[]string              `firestore:"transactionsUid,omitempty" json:"transactionsUid ,omitempty" bigquery:"-"`
	Company         string                 `firestore:"company,omitempty" json:"company,omitempty" bigquery:"company"`
	Name            string                 `firestore:"name,omitempty" json:"name,omitempty" bigquery:"name"`
	NameDesc        string                 `firestore:"nameDesc,omitempty" json:"nameDesc,omitempty" bigquery:"nameDesc"`
	BigStartDate    civil.DateTime         `bigquery:"startDate"`
	StartDate       time.Time              `firestore:"startDate,omitempty" json:"startDate,omitempty" bigquery:"-"`
	EndDate         time.Time              `firestore:"endDate,omitempty" json:"endDate,omitempty" bigquery:"-"`
	CreationDate    time.Time              `firestore:"creationDate,omitempty" json:"creationDate,omitempty" bigquery:"-"`
	Updated         time.Time              `firestore:"updated,omitempty" json:"updated,omitempty" bigquery:"-"`
	NextPay         time.Time              `firestore:"nextPay,omitempty" json:"nextPay,omitempty" bigquery:"-"`
	NextPayString   string                 `firestore:"nextPayString,omitempty" json:"nextPayString,omitempty"  bigquery:"nextPayString"`
	Payment         string                 `firestore:"payment,omitempty" json:"payment,omitempty" bigquery:"payment"`
	PaymentType     string                 `firestore:"paymentType,omitempty" json:"paymentType,omitempty" bigquery:"paymentType"`
	PaymentSplit    string                 `firestore:"paymentSplit,omitempty" json:"paymentSplit,omitempty" bigquery:"paymentSplit"`
	IsPay           bool                   `firestore:"isPay" json:"isPay,omitempty" bigquery:"isPay"`
	IsAutoRenew     bool                   `firestore:"isAutoRenew,omitempty" json:"isAutoRenew,omitempty" bigquery:"isAutoRenew"`
	IsSign          bool                   `firestore:"isSign" json:"isSign,omitempty" bigquery:"isSign"`
	CoverageType    string                 `firestore:"coverageType,omitempty" json:"coverageType,omitempty" bigquery:"coverageType"`
	Voucher         string                 `firestore:"voucher,omitempty" json:"voucher,omitempty" bigquery:"voucher"`
	Channel         string                 `firestore:"channel,omitempty" json:"channel,omitempty" bigquery:"channel"`
	Covenant        string                 `firestore:"covenant,omitempty" json:"covenant,omitempty" bigquery:"covenant"`
	TaxAmount       float64                `firestore:"taxAmount,omitempty" json:"taxAmount,omitempty" bigquery:"taxAmount"`
	PriceNett       float64                `firestore:"priceNett,omitempty" json:"priceNett,omitempty" bigquery:"priceNett"`
	PriceGross      float64                `firestore:"priceGross,omitempty" json:"priceGross,omitempty" bigquery:"priceGross"`
	Agent           *User                  `firestore:"agent,omitempty" json:"agent,omitempty" bigquery:"-"`
	Contractor      User                   `firestore:"contractor,omitempty" json:"contractor,omitempty" bigquery:"-"`
	Contractors     *[]User                `firestore:"contractors,omitempty" json:"contractors,omitempty" bigquery:"-"`
	DocumentName    string                 `firestore:"documentName,omitempty" json:"documentName,omitempty" bigquery:"-"`
	Statements      []Statement            `firestore:"statements,omitempty" json:"statements,omitempty" bigquery:"-"`
	Survay          []Statement            `firestore:"survey,omitempty" json:"survey,omitempty" bigquery:"-"`
	Attachments     *[]Attachment          `firestore:"attachments,omitempty" json:"attachments,omitempty" bigquery:"-"`
	Assets          []Asset                `firestore:"assets,omitempty" json:"assets,omitempty" bigquery:"-"`
	Claim           *[]Claim               `firestore:"claim,omitempty" json:"claim,omitempty" bigquery:"-"`
}

type RenewHistory struct {
	Amount       float64   `firestore:"amount ,omitempty" json:"amount,omitempty"`
	StartDate    time.Time `firestore:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate      time.Time `firestore:"endDate,omitempty" json:"endDate,omitempty"`
	CreationDate time.Time `firestore:"creationDate,omitempty" json:"creationDate,omitempty"`
}
type Statement struct {
	Title    string `firestore:"title ,omitempty" json:"title,omitempty"`
	Question string `firestore:"question ,omitempty" json:"question,omitempty"`
	Answer   bool   `firestore:"answer ,omitempty" json:"answer,omitempty"`
}

func PolicyToListData(query *firestore.DocumentIterator) []Policy {
	var result []Policy
	for {
		d, err := query.Next()

		if err != nil {
			log.Println("error")
		}
		if err != nil {
			if err == iterator.Done {
				log.Println("iterator.Done")
				break
			}

		}
		var value Policy
		e := d.DataTo(&value)
		log.Println("todata")
		lib.CheckError(e)
		result = append(result, value)

		log.Println(len(result))
	}
	return result
}
