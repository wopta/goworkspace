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
	NameTitle         string            `firestore:"nameTitle,omitempty" json:"nameTitle,omitempty"`
	NameSubtitle      string            `firestore:"nameSubtitle,omitempty" json:"nameSubtitle,omitempty"`
	NameDesc          *string           `firestore:"nameDesc,omitempty" json:"nameDesc,omitempty"`
	Companies         []Company         `firestore:"companies,omitempty" json:"companies,omitempty"`
	ProductUid        string            `firestore:"productUid,omitempty" json:"productUid,omitempty"`
	ProductVersion    int               `firestore:"productVersion,omitempty" json:"productVersion,omitempty"`
	Version           string            `firestore:"version,omitempty" json:"version,omitempty"`
	Number            int               `firestore:"number,omitempty" json:"number,omitempty"`
	Name              string            `firestore:"name,omitempty" json:"name,omitempty"`
	Commission        float64           `firestore:"commission,omitempty" json:"commission,omitempty"`
	CommissionRenew   float64           `firestore:"commissionRenew,omitempty" json:"commissionRenew,omitempty"`
	Steps             []Step            `firestore:"steps,omitempty" json:"steps"`
	Offers            map[string]Offer  `firestore:"offers,omitempty" json:"offers,omitempty"`
	IsEcommerceActive bool              `json:"isEcommerceActive" firestore:"isEcommerceActive"` // DEPRECATED
	IsAgencyActive    bool              `json:"isAgencyActive" firestore:"isAgencyActive"`       // DEPRECATED
	IsAgentActive     bool              `json:"isAgentActive" firestore:"isAgentActive"`         // DEPRECATED
	Logo              string            `json:"logo,omitempty" firestore:"logo,omitempty" bigquery:"-"`
	PaymentProviders  []PaymentProvider `json:"paymentProviders,omitempty" firestore:"paymentProviders,omitempty" bigquery:"-"`
	Flow              string            `json:"flow,omitempty" firestore:"flow,omitempty" bigquery:"-"` // the name of the flow file to be used
	IsActive          bool              `json:"isActive" json:"isActive" bigquery:"-"`
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
	Mandate                   Mandate              `json:"mandate" firestore:"mandate" bigquery:"-"` // DEPRECATED
	DiscountLimit             float64              `json:"discountLimit" firestore:"discountLimit" bigquery:"-"`
	AgentCode                 string               `json:"agentCode" firestore:"agentCode" bigquery:"-"`
	IsEcommerceActive         bool                 `json:"isEcommerceActive" firestore:"isEcommerceActive" bigquery:"-"` // DEPRECATED
	IsAgencyActive            bool                 `json:"isAgencyActive" firestore:"isAgencyActive" bigquery:"-"`       // DEPRECATED
	IsAgentActive             bool                 `json:"isAgentActive" firestore:"isAgentActive" bigquery:"-"`         // DEPRECATED
	AnnulmentCodes            []AnnulmentCode      `json:"annulmentCodes,omitempty" firestore:"annulmentCodes,omitempty" bigquery:"-"`
	CommissionSetting         *CommissionsSetting  `json:"commissionsSetting,omitempty" firestore:"commissionsSetting,omitempty" bigquery:"-"`
	// MaxFreeDiscount           float64              `json:"maxFreeDiscount,omitempty" firestore:"maxFreeDiscount,omitempty" bigquery:"-"`
	// MaxReservedDiscount       float64              `json:"maxReservedDiscount,omitempty" firestore:"maxReservedDiscount,omitempty" bigquery:"-"`
}

// DEPRECATED
type Mandate struct {
	Commission      float64   `json:"commission" firestore:"commission"`
	CommissionRenew float64   `json:"commissionRenew" firestore:"commissionRenew"`
	StartDate       time.Time `json:"startDate" firestore:"startDate"`
	ExpireDate      time.Time `json:"expireDate" firestore:"expireDate"`
	// Bonus TBD
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
	Name        string       `firestore:"name,omitempty" json:"name,omitempty"`
	Description string       `firestore:"description,omitempty" json:"description,omitempty"`
	Order       int          `firestore:"order,omitempty" json:"order,omitempty"`
	Commissions *Commissions `json:"commissions" firestore:"commissions"`
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

type PaymentProvider struct {
	Name    string          `json:"name,omitempty" firestore:"name,omitempty" bigquery:"-"`
	Flows   []string        `json:"flows,omitempty" firestore:"flows,omitempty" bigquery:"-"`
	Methods []PaymentMethod `json:"methods,omitempty" firestore:"methods,omitempty" bigquery:"-"`
	Rates   []string        `json:"rates,omitempty" firestore:"rates,omitempty" bigquery:"-"`
}

type PaymentMethod struct {
	Name  string   `json:"name" firestore:"name" bigquery:"-"`
	Rates []string `json:"rates" firestore:"rates" bigquery:"-"`
}

type AnnulmentCode struct {
	Code        string   `json:"code,omitempty" firestore:"code,omitempty" bigquery:"-"`
	Description string   `json:"description,omitempty" firestore:"description,omitempty" bigquery:"-"`
	RefundTypes []string `json:"refundTypes,omitempty" firestore:"refundTypes,omitempty" bigquery:"-"`
}

type ProductInfo struct {
	Name         string `json:"name"`
	NameTitle    string `json:"nameTitle"`
	NameSubtitle string `json:"nameSubtitle"`
	NameDesc     string `json:"nameDesc"`
	Logo         string `json:"logo"`
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
