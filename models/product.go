package models

import (
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

func UnmarshalProduct(data []byte) (Product, error) {
	var r Product
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Product) Marshal() ([]byte, error) {

	return json.Marshal(r)
}

type Product struct {
	NameTitle         string           `firestore:"nameTitle,omitempty" json:"nameTitle,omitempty"`
	NameSubtitle      string           `firestore:"nameSubtitle,omitempty" json:"nameSubtitle,omitempty"`
	NameDesc          *string          `firestore:"nameDesc,omitempty" json:"nameDesc,omitempty"`
	Companies         []Company        `firestore:"companies,omitempty" json:"companies,omitempty"`
	ProductUid        string           `firestore:"productUid,omitempty" json:"productUid,omitempty"`
	ProductVersion    int              `firestore:"productVersion,omitempty" json:"productVersion,omitempty"`
	Version           string           `firestore:"version,omitempty" json:"version,omitempty"`
	Number            int              `firestore:"number,omitempty" json:"number,omitempty"`
	Name              string           `firestore:"name,omitempty" json:"name,omitempty"`
	Commission        float64          `firestore:"commission,omitempty" json:"commission,omitempty"`
	CommissionRenew   float64          `firestore:"commissionRenew,omitempty" json:"commissionRenew,omitempty"`
	Steps             []Step           `firestore:"steps,omitempty" json:"steps"`
	Offers            map[string]Offer `firestore:"offers,omitempty" json:"offers,omitempty"`
	IsEcommerceActive bool             `json:"isEcommerceActive" firestore:"isEcommerceActive"`
	IsAgencyActive    bool             `json:"isAgencyActive" firestore:"isAgencyActive"`
	IsAgentActive     bool             `json:"isAgentActive" firestore:"isAgentActive"`
	Logo              string           `json:"logo,omitempty" firestore:"logo,omitempty" bigquery:"-"`
}

type Company struct {
	Name                      string               `firestore:"name,omitempty" json:"name,omitempty"`
	Code                      string               `firestore:"code,omitempty" json:"code,omitempty"`
	SequencePrefix            string               `firestore:"sequencePrefix,omitempty" json:"sequencePrefix,omitempty"`
	SequenceStart             int                  `firestore:"sequenceStart,omitempty" json:"sequenceStart,omitempty"`
	SequenceFormat            string               `firestore:"sequenceFormat,omitempty" json:"sequenceFormat,omitempty"`
	EmitTrack                 Track                `firestore:"emitTrack,omitempty" json:"emitTrack,omitempty"`
	Commission                float64              `firestore:"commission,omitempty" json:"commission,omitempty"`
	CommissionRenew           float64              `firestore:"commissionRenew,omitempty" json:"commissionRenew,omitempty"`
	MinimumMonthlyPrice       float64              `firestore:"minimumMonthlyPrice,omitempty" json:"minimumMonthlyPrice,omitempty"`
	MinimumYearlyPrice        float64              `firestore:"minimumYearlyPrice,omitempty" json:"minimumYearlyPrice,omitempty"`
	Guarantees                *[]Guarante          `firestore:"guarantees,omitempty" json:"guarantees,omitempty"`
	GuaranteesMap             map[string]*Guarante `firestore:"guaranteesMap,omitempty" json:"guaranteesMap,omitempty"`
	InformationSetLink        string               `firestore:"informationSetLink,omitempty" json:"informationSetLink,omitempty"`
	IsMonthlyPaymentAvailable bool                 `firestore:"isMonthlyPaymentAvailable" json:"isMonthlyPaymentAvailable"`
	Mandate                   Mandate              `json:"mandate" firestore:"mandate" bigquery:"-"`
	DiscountLimit             float64              `json:"discountLimit" firestore:"discountLimit" bigquery:"-"`
	AgentCode                 string               `json:"agentCode" firestore:"agentCode" bigquery:"-"`
	IsEcommerceActive         bool                 `json:"isEcommerceActive" firestore:"isEcommerceActive" bigquery:"-"`
	IsAgencyActive            bool                 `json:"isAgencyActive" firestore:"isAgencyActive" bigquery:"-"`
	IsAgentActive             bool                 `json:"isAgentActive" firestore:"isAgentActive" bigquery:"-"`
}

type Mandate struct {
	Commission      float64   `json:"commission" firestore:"commission"`
	CommissionRenew float64   `json:"commissionRenew" firestore:"commissionRenew"`
	StartDate       time.Time `json:"startDate" firestore:"startDate"`
	ExpireDate      time.Time `json:"expireDate" firestore:"expireDate"`
	//Bonus TBD
}

type Step struct {
	Widget     string      `firestore:"widget,omitempty" json:"widget"`
	Attributes interface{} `firestore:"attributes,omitempty" json:"attributes"`
	Children   []Child     `firestore:"children,omitempty" json:"children,omitempty"`
}

type Child struct {
	Widget     string      `firestore:"widget,omitempty" json:"widget"`
	Attributes interface{} `firestore:"attributes,omitempty" json:"attributes"`
}

type Offer struct {
	Name        string `firestore:"name,omitempty" json:"name,omitempty"`
	Description string `firestore:"description,omitempty" json:"description,omitempty"`
	Order       int    `firestore:"order,omitempty" json:"order,omitempty"`
}

type Track struct {
	Columns []Column `firestore:"columns,omitempty" json:"columns"`
	Name    string   `firestore:"name,omitempty" json:"name,omitempty"`
	Type    string   `firestore:"type,omitempty" json:"type"`
	Format  string   `firestore:"format,omitempty" json:"format,omitempty"`
}

type Column struct {
	Value  string `firestore:"value,omitempty" json:"value"`
	Name   string `firestore:"name,omitempty" json:"name,omitempty"`
	Type   string `firestore:"type,omitempty" json:"type"`
	Format string `firestore:"format,omitempty" json:"format,omitempty"`
}

func ProductToListData(query *firestore.DocumentIterator) []Product {
	var result []Product
	for {
		d, err := query.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			log.Println(err)
		}
		var value Product
		e := d.DataTo(&value)
		lib.CheckError(e)
		result = append(result, value)

		log.Println(len(result))
	}
	return result
}
