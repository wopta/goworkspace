package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	lib "github.com/wopta/goworkspace/lib"
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
func ToListData(query *firestore.DocumentIterator) []Policy {
	var result []Policy
	for {
		d, err := query.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			var value *Policy
			e := d.DataTo(value)
			lib.CheckError(e)
			result = append(result, *value)

		}

	}
	return result
}
func GetSequenceByProduct(name string) (string, int) {
	var numberCompany string
	var number int

	rn, e := lib.OrderWhereLimitFirestoreErr("policy", "", "company", "==", name, firestore.Desc, 1)
	if e == nil {
		numberCompany = "49999999"
	} else {
		policy := ToListData(rn)
		intNumberCompany, e := strconv.Atoi(policy[1].NumberCompany)
		lib.CheckError(e)
		numberCompany = fmt.Sprint(intNumberCompany + 1)
		number = policy[1].Number + 1
	}
	r, e := lib.OrderLimitFirestoreErr("policy", "number", firestore.Desc, 1)
	if e != nil {
		number = 0
	} else {
		policy := ToListData(r)
		number = policy[1].Number + 1
	}
	return numberCompany, number
}

type Policy struct {
	ID            string       `firestore:"id,omitempty" json:"id,omitempty"`
	IdSign        string       `firestore:"idPay,omitempty" json:"idPay,omitempty"`
	IdPay         string       `firestore:"idSign,omitempty" json:"idSign,omitempty"`
	Uid           string       `firestore:"uid,omitempty" json:"uid,omitempty"`
	ProductUid    string       `firestore:"productUid,omitempty" json:"productUid,omitempty"`
	Number        int          `firestore:"number,omitempty" json:"number,omitempty"`
	NumberCompany string       `firestore:"numberCompany,omitempty" json:"numberCompany,omitempty"`
	Status        string       `firestore:"status ,omitempty" json:"status ,omitempty"`
	StatusHistory []string     `firestore:"statusHistory ,omitempty" json:"statusHistory ,omitempty"`
	Transactions  []string     `firestore:"transactions ,omitempty" json:"transactions ,omitempty"`
	Company       string       `firestore:"company,omitempty" json:"company,omitempty"`
	Name          string       `firestore:"name,omitempty" json:"name,omitempty"`
	StartDate     time.Time    `firestore:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate       time.Time    `firestore:"endDate,omitempty" json:"endDate,omitempty"`
	CreationDate  time.Time    `firestore:"creationDate,omitempty" json:"creationDate,omitempty"`
	Updated       time.Time    `firestore:"updated,omitempty" json:"updated,omitempty"`
	Payment       string       `firestore:"payment,omitempty" json:"payment,omitempty"`
	PaymentType   string       `firestore:"paymentType,omitempty" json:"paymentType,omitempty"`
	PaymentSplit  string       `firestore:"paymentSplit,omitempty" json:"paymentSplit,omitempty"`
	IsPay         bool         `firestore:"isPay,omitempty" json:"isPay,omitempty"`
	IsSign        bool         `firestore:"isSign,omitempty" json:"isSign,omitempty"`
	CoverageType  string       `firestore:"coverageType,omitempty" json:"coverageType,omitempty"`
	Voucher       string       `firestore:"voucher,omitempty" json:"voucher,omitempty"`
	Channel       string       `firestore:"channel,omitempty" json:"channel,omitempty"`
	Covenant      string       `firestore:"covenant,omitempty" json:"covenant,omitempty"`
	TaxAmount     int64        `firestore:"taxAmount,omitempty" json:"taxAmount,omitempty"`
	PriceNett     int64        `firestore:"priceNett,omitempty" json:"priceNett,omitempty"`
	PriceGross    int64        `firestore:"priceGross,omitempty" json:"priceGross,omitempty"`
	Contractor    User         `firestore:"contractor,omitempty" json:"contractor,omitempty"`
	DocumentName  string       `firestore:"documentName,omitempty" json:"documentName,omitempty"`
	Statements    []Statement  `firestore:"statements,omitempty" json:"statements,omitempty"`
	Attachments   []Attachment `firestore:"attachments,omitempty" json:"attachments,omitempty"`
	Assets        []Asset      `firestore:"guarantees,omitempty" json:"guarantees,omitempty"`
	Claim         []Claim      `firestore:"claim ,omitempty" json:"claim,omitempty"`
}
type Statement struct {
	Question string
	Answer   string
}

func GetDefaultPolicy() (string, interface{}) {

	policy := Policy{
		ID:            "id",
		IdSign:        "idSign",
		IdPay:         "idpay",
		Uid:           "uid",
		Number:        0,
		NumberCompany: "NumberCompany",
		Status:        "Status",
		StatusHistory: []string{"init"},
		Transactions:  []string{""},
		Company:       "Global",
		Name:          "Pmi",
		StartDate:     time.Now(),
		EndDate:       time.Now(),
		CreationDate:  time.Now(),
		Updated:       time.Now(),
		Payment:       "Payment",
		PaymentType:   "PaymentType",
		PaymentSplit:  "PaymentSplit",
		IsPay:         false,
		IsSign:        false,
		CoverageType:  "CoverageType",
		Voucher:       "Voucher",
		Channel:       "Channel",
		Covenant:      "Covenant",
		TaxAmount:     0,
		PriceNett:     0,
		PriceGross:    0,
		Contractor:    User{},
		DocumentName:  "",
		Statements:    []Statement{{}},
		Attachments:   []Attachment{{}},
		Assets: []Asset{{

			Name:    "test",
			Address: "test",
			Type:    "test",
			Building: Building{
				Name:             "test",
				Address:          "test",
				Type:             "test",
				PostalCode:       "test",
				City:             "test",
				BuildingType:     "test",
				BuildingMaterial: "test",
				BuildingYear:     "test",
				SquareMeters:     340,
				IsAllarm:         true,
				Floor:            4,
				Costruction:      "test",
				IsHolder:         true},
			Person: User{},
			Enterprise: Enterprise{
				Name:       "test",
				Address:    "test",
				Type:       "test",
				PostalCode: "test",
				City:       "test",
				VatCode:    "test",
				Ateco:      "test",
				Revenue:    "test",
				Employer:   4},
			IsContractor: true,
			Guarantees: []Guarantee{{
				Type:                       "test",
				Beneficiary:                User{},
				TypeOfSumInsured:           "test",
				Description:                "test",
				Value:                      GuaranteeValue{},
				Slug:                       "test",
				IsBase:                     true,
				IsYour:                     true,
				IsPremium:                  true,
				Base:                       GuaranteeValue{},
				Your:                       GuaranteeValue{},
				Premium:                    GuaranteeValue{},
				Name:                       "test",
				SumInsuredLimitOfIndemnity: 100000,
				Tax:                        10.3,
				Price:                      100000,
				PriceNett:                  100000,
				PriceGross:                 100000,
			}},
		}},
		Claim: []Claim{{}},
	}

	b, e := json.Marshal(policy)
	log.Println(string(b))
	lib.CheckError(e)

	return string(b), policy
}
