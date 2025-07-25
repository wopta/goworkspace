package models

import (
	"encoding/json"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"google.golang.org/api/iterator"
)

func (r *Product) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Product struct {
	NameTitle          string            `firestore:"nameTitle,omitempty" json:"nameTitle,omitempty"`
	NameSubtitle       string            `firestore:"nameSubtitle,omitempty" json:"nameSubtitle,omitempty"`
	NameDesc           *string           `firestore:"nameDesc,omitempty" json:"nameDesc,omitempty"`
	Companies          []Company         `firestore:"companies,omitempty" json:"companies,omitempty"`
	ProductUid         string            `firestore:"productUid,omitempty" json:"productUid,omitempty"`
	ProductVersion     int               `firestore:"productVersion,omitempty" json:"productVersion,omitempty"`
	Version            string            `firestore:"version,omitempty" json:"version,omitempty"`
	Number             int               `firestore:"number,omitempty" json:"number,omitempty"`
	Name               string            `firestore:"name,omitempty" json:"name,omitempty"`
	Commission         float64           `firestore:"commission,omitempty" json:"commission,omitempty"`
	CommissionRenew    float64           `firestore:"commissionRenew,omitempty" json:"commissionRenew,omitempty"`
	Steps              []Step            `firestore:"steps,omitempty" json:"steps"`
	Offers             map[string]Offer  `firestore:"offers,omitempty" json:"offers,omitempty"`
	Logo               string            `json:"logo,omitempty" firestore:"logo,omitempty" bigquery:"-"`
	PaymentProviders   []PaymentProvider `json:"paymentProviders,omitempty" firestore:"paymentProviders,omitempty" bigquery:"-"`
	Flow               string            `json:"flow,omitempty" firestore:"flow,omitempty" bigquery:"-"` // the name of the flow file to be used
	IsActive           bool              `json:"isActive" firestore:"isActive" bigquery:"-"`
	RenewOffset        int               `json:"renewOffset" firestore:"renewOffset" bigquery:"-"`
	IsAutoRenew        bool              `json:"isAutoRenew" firestore:"isAutoRenew" bigquery:"-"`
	IsRenewable        bool              `json:"isRenewable" firestore:"isRenewable" bigquery:"-"`
	PolicyType         string            `json:"policyType,omitempty" firestore:"policyType,omitempty" bigquery:"-"`
	QuoteType          string            `json:"quoteType" firestore:"quoteType" bigquery:"-"`
	EmitMaxElapsedDays uint              `json:"emitMaxElapsedDays" firestore:"-" bigquery:"-"`
	Categories         []string          `json:"categories" firestore:"-" bigquery:"-"`
	//Setting for remunerazione
	ConsultancyConfig *ConsultancyConfig `json:"consultancyConfig" firestore:"consultancyConfig" bigquery:"-"`
	IsAIAgentEnabled  bool               `json:"isAIAgentEnabled" firestore:"isAIAgentEnabled" bigquery:"-"`

	// DEPRECATED FIELDS

	IsEcommerceActive bool `json:"isEcommerceActive" firestore:"isEcommerceActive"` // DEPRECATED
	IsAgencyActive    bool `json:"isAgencyActive" firestore:"isAgencyActive"`       // DEPRECATED
	IsAgentActive     bool `json:"isAgentActive" firestore:"isAgentActive"`         // DEPRECATED
}

func (p *Product) GetType() string {
	return "product"
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
	AgentCode                 string               `json:"agentCode" firestore:"agentCode" bigquery:"-"`                 // DEPRECATED
	IsEcommerceActive         bool                 `json:"isEcommerceActive" firestore:"isEcommerceActive" bigquery:"-"` // DEPRECATED
	IsAgencyActive            bool                 `json:"isAgencyActive" firestore:"isAgencyActive" bigquery:"-"`       // DEPRECATED
	IsAgentActive             bool                 `json:"isAgentActive" firestore:"isAgentActive" bigquery:"-"`         // DEPRECATED
	AnnulmentCodes            []AnnulmentCode      `json:"annulmentCodes,omitempty" firestore:"annulmentCodes,omitempty" bigquery:"-"`
	//Setting for commisioni
	CommissionSetting *CommissionsSetting `json:"commissionsSetting,omitempty" firestore:"commissionsSetting,omitempty" bigquery:"-"`
	ProducerCode      string              `json:"producerCode,omitempty" firestore:"producerCode,omitempty" bigquery:"producerCode"`
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
	Title      string      `json:"title"`
	Attributes interface{} `firestore:"attributes,omitempty" json:"attributes"`
	Children   []Child     `firestore:"children,omitempty" json:"children,omitempty"`
	Flows      []string    `json:"flows" firestore:"flows,omitempty"`
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
	Emit    []Column `firestore:"emit,omitempty" json:"Emit,omitempty"`
}

type Column struct {
	Value  string `firestore:"value,omitempty" json:"value"`
	Name   string `firestore:"name,omitempty" json:"name,omitempty"`
	Type   string `firestore:"type,omitempty" json:"type"`
	Format string `firestore:"format,omitempty" json:"format,omitempty"`
	MapFx  string `firestore:"mapFx,omitempty" json:"mapFx,omitempty"`
	Frame  string `firestore:"frame,omitempty" json:"frame,omitempty"`
}

type PaymentProvider struct {
	Name    string          `json:"name,omitempty" firestore:"name,omitempty" bigquery:"-"`
	Flows   []string        `json:"flows,omitempty" firestore:"flows,omitempty" bigquery:"-"`
	Configs []PaymentConfig `json:"configs,omitempty" firestore:"configs,omitempty" bigquery:"-"`
	Methods []PaymentMethod `json:"methods,omitempty" firestore:"methods,omitempty" bigquery:"-"`
	Rates   []string        `json:"rates,omitempty" firestore:"rates,omitempty" bigquery:"-"`
}

type PaymentConfig struct {
	Rate    string   `json:"rate,omitempty" firestore:"rate,omitempty" bigquery:"-"`
	Methods []string `json:"methods,omitempty" firestore:"methods,omitempty" bigquery:"-"`
	Mode    string   `json:"mode,omitempty" firestore:"mode,omitempty" bigquery:"-"`
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

type ConsultancyConfig struct {
	Min            float64 `json:"min" firestore:"min" bigquery:"-"`
	Max            float64 `json:"max" firestore:"max" bigquery:"-"`
	Step           float64 `json:"step" firestore:"step" bigquery:"-"`
	DefaultValue   float64 `json:"defaultValue" firestore:"defaultValue" bigquery:"-"`
	IsActive       bool    `json:"isActive" firestore:"isActive" bigquery:"-"`
	IsConfigurable bool    `json:"isConfigurable" firestore:"isConfigurable" bigquery:"-"`
}

type ProductInfo struct {
	Name         string        `json:"name"`
	NameTitle    string        `json:"nameTitle"`
	NameSubtitle string        `json:"nameSubtitle"`
	NameDesc     string        `json:"nameDesc"`
	Logo         string        `json:"logo"`
	Version      string        `json:"version"`
	Type         string        `json:"type"`
	Company      string        `json:"company"`               // DEPRECATED
	ExternalUrl  string        `json:"externalUrl,omitempty"` // external integration products
	Products     []ProductInfo `json:"products,omitempty"`    // external integration products
	IsActive     bool          `json:"isActive"`
	Categories   []string      `json:"categories"`
}

func (p *Product) ToProductInfo() ProductInfo {
	return ProductInfo{
		Name:         p.Name,
		NameTitle:    p.NameTitle,
		NameSubtitle: p.NameSubtitle,
		NameDesc:     *p.NameDesc,
		Logo:         p.Logo,
		Version:      p.Version,
		Type:         InternalProductType,
		IsActive:     p.IsActive,
		Categories:   p.Categories,
	}
}

type FormProduct struct {
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	NameTitle    string   `json:"nameTitle"`
	NameSubtitle string   `json:"nameSubtitle"`
	NameDesc     string   `json:"nameDesc"`
	Logo         string   `json:"logo"`
	Version      string   `json:"version"`
	ExternalUrl  string   `json:"externalUrl"`
	IsActive     bool     `json:"isActive"`
	Categories   []string `json:"categories"`
}

func (p *FormProduct) ToProductInfo() ProductInfo {
	return ProductInfo{
		Name:         p.Name,
		NameTitle:    p.NameTitle,
		NameSubtitle: p.NameSubtitle,
		NameDesc:     p.NameDesc,
		Logo:         p.Logo,
		Version:      p.Version,
		ExternalUrl:  p.ExternalUrl,
		Type:         FormProductType,
		IsActive:     p.IsActive,
		Categories:   p.Categories,
	}
}

type ExternalProduct struct {
	Type        string        `json:"type"`
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	ExternalUrl string        `json:"externalUrl"`
	IsActive    bool          `json:"isActive"`
	Products    []ProductInfo `json:"products"`
	Categories  []string      `json:"categories"`
}

func (p *ExternalProduct) ToProductInfo() ProductInfo {
	return ProductInfo{
		Name:        p.Name,
		Version:     p.Version,
		ExternalUrl: p.ExternalUrl,
		Products:    p.Products,
		Type:        ExternalProductType,
		IsActive:    p.IsActive,
		Categories:  p.Categories,
	}
}

type DynamicProduct struct {
	Value interface{}
}

type BaseProduct struct {
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Version  string         `json:"version"`
	IsActive bool           `json:"isActive"`
	Product  DynamicProduct `json:"dynamicProduct"`
}

func (b *BaseProduct) UnmarshalJSON(data []byte) error {
	var dynamicType struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &dynamicType); err != nil {
		return err
	}
	switch dynamicType.Type {
	case ExternalProductType:
		b.Product.Value = new(ExternalProduct)
		if err := json.Unmarshal(data, &b.Product.Value); err != nil {
			return err
		}
		b.Name = (b.Product.Value).(*ExternalProduct).Name
		b.Version = (b.Product.Value).(*ExternalProduct).Version
		b.IsActive = (b.Product.Value).(*ExternalProduct).IsActive
	case FormProductType:
		b.Product.Value = new(FormProduct)
		if err := json.Unmarshal(data, &b.Product.Value); err != nil {
			return err
		}
		b.Name = (b.Product.Value).(*FormProduct).Name
		b.Version = (b.Product.Value).(*FormProduct).Version
		b.IsActive = (b.Product.Value).(*FormProduct).IsActive
	default:
		b.Product.Value = new(Product)
		if err := json.Unmarshal(data, &b.Product.Value); err != nil {
			return err
		}
		b.Name = (b.Product.Value).(*Product).Name
		b.Version = (b.Product.Value).(*Product).Version
		b.IsActive = (b.Product.Value).(*Product).IsActive
	}
	return nil
}

func (b *BaseProduct) ToProductInfo() ProductInfo {
	var res ProductInfo
	switch b.Product.Value.(type) {
	case *ExternalProduct:
		res = (b.Product.Value).(*ExternalProduct).ToProductInfo()
	case *FormProduct:
		res = (b.Product.Value).(*FormProduct).ToProductInfo()
	default:
		res = (b.Product.Value).(*Product).ToProductInfo()
	}
	return res
}

func ProductToListData(query *firestore.DocumentIterator) []Product {
	var result []Product
	for {
		d, err := query.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
		}
		var value Product
		e := d.DataTo(&value)
		lib.CheckError(e)
		result = append(result, value)

		log.Println(len(result))
	}
	return result
}
