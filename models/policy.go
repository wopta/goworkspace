package models

import (
	"encoding/json"
	"time"
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
	ID             string        `firestore:"id,omitempty" json:"id,omitempty"`
	IdSign         string        `firestore:"idSign,omitempty" json:"idSign,omitempty"`
	IdPay          string        `firestore:"idPay,omitempty" json:"idPay,omitempty"`
	Uid            string        `firestore:"uid,omitempty" json:"uid,omitempty"`
	ProductUid     string        `firestore:"productUid,omitempty" json:"productUid,omitempty"`
	ProductVersion int           `firestore:"productVersion,omitempty" json:"productVersion,omitempty"`
	ProposalNumber int           `firestore:"proposalNumber,omitempty" json:"proposalNumber,omitempty"`
	Number         int           `firestore:"number,omitempty" json:"number,omitempty"`
	NumberCompany  string        `firestore:"numberCompany,omitempty" json:"numberCompany,omitempty"`
	Status         string        `firestore:"status ,omitempty" json:"status ,omitempty"`
	StatusHistory  []string      `firestore:"statusHistory ,omitempty" json:"statusHistory ,omitempty"`
	Transactions   []Transaction `firestore:"transactions ,omitempty" json:"transactions ,omitempty"`
	Company        string        `firestore:"company,omitempty" json:"company,omitempty"`
	Name           string        `firestore:"name,omitempty" json:"name,omitempty"`
	StartDate      time.Time     `firestore:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate        time.Time     `firestore:"endDate,omitempty" json:"endDate,omitempty"`
	CreationDate   time.Time     `firestore:"creationDate,omitempty" json:"creationDate,omitempty"`
	Updated        time.Time     `firestore:"updated,omitempty" json:"updated,omitempty"`
	Payment        string        `firestore:"payment,omitempty" json:"payment,omitempty"`
	PaymentType    string        `firestore:"paymentType,omitempty" json:"paymentType,omitempty"`
	PaymentSplit   string        `firestore:"paymentSplit,omitempty" json:"paymentSplit,omitempty"`
	IsPay          bool          `firestore:"isPay,omitempty" json:"isPay,omitempty"`
	IsSign         bool          `firestore:"isSign,omitempty" json:"isSign,omitempty"`
	CoverageType   string        `firestore:"coverageType,omitempty" json:"coverageType,omitempty"`
	Voucher        string        `firestore:"voucher,omitempty" json:"voucher,omitempty"`
	Channel        string        `firestore:"channel,omitempty" json:"channel,omitempty"`
	Covenant       string        `firestore:"covenant,omitempty" json:"covenant,omitempty"`
	TaxAmount      float64       `firestore:"taxAmount,omitempty" json:"taxAmount,omitempty"`
	PriceNett      float64       `firestore:"priceNett,omitempty" json:"priceNett,omitempty"`
	PriceGross     float64       `firestore:"priceGross,omitempty" json:"priceGross,omitempty"`
	Contractor     User          `firestore:"contractor,omitempty" json:"contractor,omitempty"`
	DocumentName   string        `firestore:"documentName,omitempty" json:"documentName,omitempty"`
	Statements     []Statement   `firestore:"statements,omitempty" json:"statements,omitempty"`
	Survay         []Statement   `firestore:"survey,omitempty" json:"survey,omitempty"`
	Attachments    []Attachment  `firestore:"attachments,omitempty" json:"attachments,omitempty"`
	Assets         []Asset       `firestore:"assets,omitempty" json:"assets,omitempty"`
	Claim          []Claim       `firestore:"claim ,omitempty" json:"claim,omitempty"`
}

type Transaction struct {
	Amount        float64   `firestore:"amount ,omitempty" json:"amount,omitempty"`
	Status        string    `firestore:"status ,omitempty" json:"status ,omitempty"`
	PolicyName    string    `firestore:"policyName,omitempty" json:"policName,omitempty"`
	StartDate     time.Time `firestore:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate       time.Time `firestore:"endDate,omitempty" json:"endDate,omitempty"`
	CreationDate  time.Time `firestore:"creationDate,omitempty" json:"creationDate,omitempty"`
	Uid           string    `firestore:"uid,omitempty" json:"uid,omitempty"`
	PolicyUid     string    `firestore:"policyUid,omitempty" json:"policyUid,omitempty"`
	Company       string    `firestore:"company,omitempty" json:"company,omitempty"`
	NumberCompany string    `firestore:"numberCompany,omitempty" json:"numberCompany,omitempty"`
}
type Statement struct {
	Question string `firestore:"question ,omitempty" json:"question,omitempty"`
	Answer   bool   `firestore:"answer ,omitempty" json:"answer,omitempty"`
}
