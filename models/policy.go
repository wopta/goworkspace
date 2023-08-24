package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
	ID                string                       `firestore:"id,omitempty" json:"id,omitempty" bigquery:"id"`
	IdSign            string                       `firestore:"idSign,omitempty" json:"idSign,omitempty" bigquery:"idSign"`
	IdPay             string                       `firestore:"idPay,omitempty" json:"idPay,omitempty" bigquery:"idPay"`
	QuoteQuestions    map[string]interface{}       `firestore:"quoteQuestions,omitempty" json:"quoteQuestions,omitempty" bigquery:"-"`
	ContractFileId    string                       `firestore:"contractFileId,omitempty" json:"contractFileId,omitempty" bigquery:"contractFileId"`
	Uid               string                       `firestore:"uid,omitempty" json:"uid,omitempty" bigquery:"uid"`
	ProductUid        string                       `firestore:"productUid,omitempty" json:"productUid,omitempty" bigquery:"productUid"`
	PayUrl            string                       `firestore:"payUrl,omitempty" json:"payUrl,omitempty" bigquery:"-"`
	SignUrl           string                       `firestore:"signUrl,omitempty" json:"signUrl,omitempty" bigquery:"-"`
	ProductVersion    string                       `firestore:"productVersion,omitempty" json:"productVersion,omitempty" bigquery:"productVersion"`
	ProposalNumber    int                          `firestore:"proposalNumber,omitempty" json:"proposalNumber,omitempty" bigquery:"proposalNumber"`
	OfferlName        string                       `firestore:"offerName,omitempty" json:"offerName,omitempty" bigquery:"offerName"`
	Number            int                          `firestore:"number,omitempty" json:"number,omitempty" bigquery:"number"`
	NumberCompany     int                          `firestore:"numberCompany,omitempty" json:"numberCompany,omitempty" bigquery:"numberCompany"`
	CodeCompany       string                       `firestore:"codeCompany,omitempty" json:"codeCompany,omitempty" bigquery:"codeCompany"`
	Status            string                       `firestore:"status,omitempty" json:"status,omitempty" bigquery:"status"`
	StatusHistory     []string                     `firestore:"statusHistory,omitempty" json:"statusHistory,omitempty" bigquery:"-"`
	BigStatusHistory  string                       `firestore:"-" json:"-" bigquery:"statusHistory"`
	RenewHistory      *[]RenewHistory              `firestore:"renewHistory,omitempty" json:"renewHistory,omitempty" bigquery:"-"`
	Transactions      *[]Transaction               `firestore:"transactions,omitempty" json:"transactions,omitempty" bigquery:"-"`
	TransactionsUid   *[]string                    `firestore:"transactionsUid,omitempty" json:"transactionsUid,omitempty" bigquery:"-"`
	Company           string                       `firestore:"company,omitempty" json:"company,omitempty" bigquery:"company"`
	Name              string                       `firestore:"name,omitempty" json:"name,omitempty" bigquery:"name"`
	NameDesc          string                       `firestore:"nameDesc,omitempty" json:"nameDesc,omitempty" bigquery:"nameDesc"`
	BigStartDate      civil.DateTime               `bigquery:"startDate" firestore:"-" json:"-"`
	BigRenewDate      civil.DateTime               `json:"-" firestore:"-" bigquery:"renewDate"`
	BigEndDate        civil.DateTime               `bigquery:"endDate" firestore:"-" json:"-"`
	BigEmitDate       civil.DateTime               `bigquery:"emitDate" firestore:"-" json:"-"`
	EmitDate          time.Time                    `firestore:"emitDate,omitempty" json:"emitDate,omitempty" bigquery:"-"`
	StartDate         time.Time                    `firestore:"startDate,omitempty" json:"startDate,omitempty" bigquery:"-"`
	RenewDate         time.Time                    `json:"renewDate" firestore:"renewDate" bigquery:"-"`
	EndDate           time.Time                    `firestore:"endDate,omitempty" json:"endDate,omitempty" bigquery:"-"`
	CreationDate      time.Time                    `firestore:"creationDate,omitempty" json:"creationDate,omitempty" bigquery:"-"`
	Updated           time.Time                    `firestore:"updated,omitempty" json:"updated,omitempty" bigquery:"-"`
	NextPay           time.Time                    `firestore:"nextPay,omitempty" json:"nextPay,omitempty" bigquery:"-"`
	NextPayString     string                       `firestore:"nextPayString,omitempty" json:"nextPayString,omitempty" bigquery:"nextPayString"`
	Payment           string                       `firestore:"payment,omitempty" json:"payment,omitempty" bigquery:"payment"`
	PaymentType       string                       `firestore:"paymentType,omitempty" json:"paymentType,omitempty" bigquery:"paymentType"`
	PaymentSplit      string                       `firestore:"paymentSplit,omitempty" json:"paymentSplit,omitempty" bigquery:"paymentSplit"`
	DeleteCode        string                       `json:"deleteCode,omitempty" firestore:"deleteCode,omitempty" bigquery:"-"`
	DeleteDesc        string                       `firestore:"deleteDesc,omitempty" json:"deleteDesc,omitempty" bigquery:"-"`
	DeleteDate        time.Time                    `json:"deleteDate,omitempty" firestore:"deleteDate,omitempty" bigquery:"-"`
	RefundType        string                       `json:"refundType,omitempty" firestore:"refundType,omitempty" bigquery:"-"`
	IsPay             bool                         `firestore:"isPay" json:"isPay,omitempty" bigquery:"isPay"`
	IsAutoRenew       bool                         `firestore:"isAutoRenew,omitempty" json:"isAutoRenew,omitempty" bigquery:"isAutoRenew"`
	IsRenew           bool                         `firestore:"isRenew" json:"isRenew,omitempty" bigquery:"isRenew"`
	IsSign            bool                         `firestore:"isSign" json:"isSign,omitempty" bigquery:"isSign"`
	IsDeleted         bool                         `firestore:"isDeleted" json:"isDeleted,omitempty" bigquery:"isDeleted"`
	DeleteEmited      bool                         `firestore:"deleteEmited" json:"deleteEmited,omitempty" bigquery:"deleteEmited"`
	CompanyEmit       bool                         `firestore:"companyEmit" json:"companyEmit,omitempty" bigquery:"-"`
	CompanyEmitted    bool                         `firestore:"companyEmitted" json:"companyEmitted,omitempty" bigquery:"-"`
	CoverageType      string                       `firestore:"coverageType,omitempty" json:"coverageType,omitempty" bigquery:"coverageType"`
	Voucher           string                       `firestore:"voucher,omitempty" json:"voucher,omitempty" bigquery:"voucher"`
	Channel           string                       `firestore:"channel,omitempty" json:"channel,omitempty" bigquery:"channel"`
	Covenant          string                       `firestore:"covenant,omitempty" json:"covenant,omitempty" bigquery:"covenant"`
	TaxAmount         float64                      `firestore:"taxAmount,omitempty" json:"taxAmount,omitempty" bigquery:"taxAmount"`
	PriceNett         float64                      `firestore:"priceNett,omitempty" json:"priceNett,omitempty" bigquery:"priceNett"`
	PriceGross        float64                      `firestore:"priceGross,omitempty" json:"priceGross,omitempty" bigquery:"priceGross"`
	TaxAmountMonthly  float64                      `json:"taxAmountMonthly,omitempty" firestore:"taxAmountMonthly,omitempty" bigquery:"taxAmountMonthly"`
	PriceNettMonthly  float64                      `json:"priceNettMonthly,omitempty" firestore:"priceNettMonthly,omitempty" bigquery:"priceNettMonthly"`
	PriceGrossMonthly float64                      `json:"priceGrossMonthly,omitempty" firestore:"priceGrossMonthly,omitempty" bigquery:"priceGrossMonthly"`
	Agent             *User                        `firestore:"agent,omitempty" json:"agent,omitempty" bigquery:"-"`
	Contractor        User                         `firestore:"contractor,omitempty" json:"contractor,omitempty" bigquery:"-"`
	Contractors       *[]User                      `firestore:"contractors,omitempty" json:"contractors,omitempty" bigquery:"-"`
	DocumentName      string                       `firestore:"documentName,omitempty" json:"documentName,omitempty" bigquery:"-"`
	Statements        *[]Statement                 `firestore:"statements,omitempty" json:"statements,omitempty" bigquery:"-"`
	Surveys           *[]Survey                    `firestore:"surveys,omitempty" json:"surveys,omitempty" bigquery:"-"`
	Attachments       *[]Attachment                `firestore:"attachments,omitempty" json:"attachments,omitempty" bigquery:"-"`
	Assets            []Asset                      `firestore:"assets,omitempty" json:"assets,omitempty" bigquery:"-"`
	Claim             *[]Claim                     `firestore:"claim,omitempty" json:"claim,omitempty" bigquery:"-"`
	Data              string                       `bigquery:"data" json:"-" firestore:"-"`
	Json              string                       `bigquery:"json" json:"-" firestore:"-"`
	OffersPrices      map[string]map[string]*Price `firestore:"offersPrices,omitempty" json:"offersPrices,omitempty" bigquery:"-"`
	PartnershipName   string                       `json:"partnershipName" firestore:"partnershipName" bigquery:"partnershipName"`
	PartnershipData   map[string]interface{}       `json:"partnershipData" firestore:"partnershipData" bigquery:"-"`
	IsReserved        bool                         `json:"isReserved" firestore:"isReserved" bigquery:"-"`
	RejectReasons     string                       `json:"rejectReasons,omitempty" firestore:"rejectReasons,omitempty" bigquery:"rejectReasons"`
	FundsOrigin       string                       `json:"fundsOrigin,omitempty" firestore:"fundsOrigin,omitempty" bigquery:"-"`
	AgentUid          string                       `json:"agentUid,omitempty" firestore:"agentUid,omitempty" bigquery:"agentUid"`
	AgencyUid         string                       `json:"agencyUid,omitempty" firestore:"agencyUid,omitempty" bigquery:"agencyUid"`
	ReservedInfo      *ReservedInfo                `json:"reservedInfo,omitempty" firestore:"reservedInfo,omitempty" bigquery:"reservedInfo"`
}

type RenewHistory struct {
	Amount       float64   `firestore:"amount,omitempty" json:"amount,omitempty"`
	StartDate    time.Time `firestore:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate      time.Time `firestore:"endDate,omitempty" json:"endDate,omitempty"`
	CreationDate time.Time `firestore:"creationDate,omitempty" json:"creationDate,omitempty"`
}

type Survey struct {
	Id                 int64       `json:"id" firestore:"id"`
	Title              string      `firestore:"title,omitempty" json:"title,omitempty"`
	Subtitle           string      `firestore:"subtitle,omitempty" json:"subtitle,omitempty"`
	HasMultipleAnswers *bool       `firestore:"hasMultipleAnswers,omitempty" json:"hasMultipleAnswers,omitempty"`
	Questions          []*Question `firestore:"questions,omitempty" json:"questions,omitempty"`
	Answer             *bool       `firestore:"answer,omitempty" json:"answer,omitempty"`
	HasAnswer          bool        `firestore:"hasAnswer" json:"hasAnswer"`
	ExpectedAnswer     *bool       `firestore:"expectedAnswer,omitempty" json:"expectedAnswer,omitempty"`
	ContractorSign     bool        `json:"contractorSign" firestore:"contractorSign"`
	CompanySign        bool        `json:"companySign" firestore:"companySign"`
}

type Statement struct {
	Id                 int64       `json:"id" firestore:"id"`
	Title              string      `firestore:"title,omitempty" json:"title,omitempty"`
	Subtitle           string      `firestore:"subtitle,omitempty" json:"subtitle,omitempty"`
	HasMultipleAnswers *bool       `firestore:"hasMultipleAnswers,omitempty" json:"hasMultipleAnswers,omitempty"`
	Questions          []*Question `firestore:"questions,omitempty" json:"questions,omitempty"`
	Answer             *bool       `firestore:"answer,omitempty" json:"answer,omitempty"`
	HasAnswer          bool        `firestore:"hasAnswer" json:"hasAnswer"`
	ExpectedAnswer     *bool       `firestore:"expectedAnswer,omitempty" json:"expectedAnswer,omitempty"`
	ContractorSign     bool        `json:"contractorSign" firestore:"contractorSign"`
	CompanySign        bool        `json:"companySign" firestore:"companySign"`
}

type Question struct {
	Id             int64  `json:"id" firestore:"id"`
	Question       string `firestore:"question,omitempty" json:"question,omitempty"`
	IsBold         bool   `firestore:"isBold,omitempty" json:"isBold,omitempty"`
	Indent         bool   `firestore:"indent,omitempty" json:"indent,omitempty"`
	Answer         *bool  `firestore:"answer,omitempty" json:"answer,omitempty"`
	HasAnswer      bool   `firestore:"hasAnswer" json:"hasAnswer"`
	ExpectedAnswer *bool  `firestore:"expectedAnswer,omitempty" json:"expectedAnswer,omitempty"`
}

type Price struct {
	Net      float64 `firestore:"net" json:"net" bigquery:"-"`
	Tax      float64 `firestore:"tax" json:"tax" bigquery:"-"`
	Gross    float64 `firestore:"gross" json:"gross" bigquery:"-"`
	Delta    float64 `firestore:"delta" json:"delta" bigquery:"-"`
	Discount float64 `firestore:"discount" json:"discount" bigquery:"-"`
}

func (policy *Policy) CalculateContractorAge() (int, error) {
	var startDate time.Time
	if policy.StartDate.IsZero() {
		startDate = time.Now()
	} else {
		startDate = policy.StartDate
	}

	birthdate, e := time.Parse(time.RFC3339, policy.Contractor.BirthDate)
	age := startDate.Year() - birthdate.Year()
	log.Printf("[CalculateContractorAge] startDate.YearDay %d - birthdate.YearDay %d", startDate.YearDay(), birthdate.YearDay())

	if startDate.YearDay() < birthdate.YearDay() && !(startDate.Month() == birthdate.Month() && startDate.Day() == birthdate.Day()) {
		age--
	}
	return age, e
}

func (policy *Policy) HasGuarantee(guaranteeSlug string) bool {
	for _, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Slug == guaranteeSlug {
			return true
		}
	}
	return false
}

func (policy *Policy) ExtractGuarantee(guaranteeSlug string) (Guarante, error) {
	for _, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Slug == guaranteeSlug {
			return guarantee, nil
		}
	}
	return Guarante{}, fmt.Errorf("no %s guarantee found", guaranteeSlug)
}

func (policy *Policy) ExtractConsens(consentKey int64) (Consens, error) {
	for _, consent := range *policy.Contractor.Consens {
		if consent.Key == consentKey {
			return consent, nil
		}
	}
	return Consens{}, fmt.Errorf("no consent found with key %d", consentKey)
}

func (policy *Policy) GuaranteesToMap() map[string]Guarante {
	m := make(map[string]Guarante, 0)
	for _, guarantee := range policy.Assets[0].Guarantees {
		m[guarantee.Slug] = guarantee
	}
	return m
}

func (policy *Policy) BigquerySave(origin string) {
	policyBig := lib.GetDatasetByEnv(origin, PolicyCollection)
	policyJson, e := policy.Marshal()
	log.Println(" policy "+policy.Uid, string(policyJson))
	policy.Data = string(policyJson)
	policy.BigStartDate = civil.DateTimeOf(policy.StartDate)
	policy.BigRenewDate = civil.DateTimeOf(policy.RenewDate)
	policy.BigEndDate = civil.DateTimeOf(policy.EndDate)
	policy.BigEmitDate = civil.DateTimeOf(policy.EmitDate)
	policy.BigStatusHistory = strings.Join(policy.StatusHistory, ",")
	if policy.ReservedInfo != nil {
		policy.ReservedInfo.BigReasons = strings.Join(policy.ReservedInfo.Reasons, ",")
	}
	log.Println(" policy save big query: " + policy.Uid)
	e = lib.InsertRowsBigQuery(WoptaDataset, policyBig, policy)
	log.Println(" policy save big query error: ", e)
}

func PolicyToListData(query *firestore.DocumentIterator) []Policy {
	result := make([]Policy, 0)
	for {
		d, err := query.Next()
		if err != nil {
			log.Println("error")
			if err == iterator.Done {
				log.Println("iterator.Done")
				break
			}
			break
		} else {
			var value Policy
			e := d.DataTo(&value)
			log.Println("todata")
			lib.CheckError(e)
			result = append(result, value)
			log.Println(len(result))
		}
	}
	return result
}
