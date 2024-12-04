package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
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
	PaymentMode       string                       `json:"paymentMode,omitempty" firestore:"paymentMode,omitempty" bigquery:"paymentMode"`
	DeleteCode        string                       `json:"deleteCode,omitempty" firestore:"deleteCode,omitempty" bigquery:"-"`
	DeleteDesc        string                       `firestore:"deleteDesc,omitempty" json:"deleteDesc,omitempty" bigquery:"-"`
	DeleteDate        time.Time                    `json:"deleteDate,omitempty" firestore:"deleteDate,omitempty" bigquery:"-"`
	RefundType        string                       `json:"refundType,omitempty" firestore:"refundType,omitempty" bigquery:"-"`
	IsPay             bool                         `firestore:"isPay" json:"isPay,omitempty" bigquery:"isPay"`
	IsRenew           bool                         `firestore:"isRenew" json:"isRenew,omitempty" bigquery:"isRenew"`
	IsSign            bool                         `firestore:"isSign" json:"isSign,omitempty" bigquery:"isSign"`
	IsDeleted         bool                         `firestore:"isDeleted" json:"isDeleted,omitempty" bigquery:"isDeleted"`
	DeleteEmited      bool                         `firestore:"deleteEmited" json:"deleteEmited,omitempty" bigquery:"deleteEmited"`
	CompanyEmit       bool                         `firestore:"companyEmit" json:"companyEmit,omitempty" bigquery:"companyEmit"`
	CompanyEmitted    bool                         `firestore:"companyEmitted" json:"companyEmitted,omitempty" bigquery:"companyEmitted"`
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
	PriceGroup        []Price                      `json:"priceGroup,omitempty" firestore:"priceGroup,omitempty" bigquery:"-"`
	Agent             *User                        `firestore:"agent,omitempty" json:"agent,omitempty" bigquery:"-"`
	Contractor        Contractor                   `firestore:"contractor,omitempty" json:"contractor,omitempty" bigquery:"-"`
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
	FundsOrigin       string                       `json:"fundsOrigin,omitempty" firestore:"fundsOrigin,omitempty" bigquery:"-"`
	ReservedInfo      *ReservedInfo                `json:"reservedInfo,omitempty" firestore:"reservedInfo,omitempty" bigquery:"-"`
	BigReasons        string                       `json:"-" firestore:"-" bigquery:"reasons"`
	BigAcceptanceNote string                       `json:"-" firestore:"-" bigquery:"acceptanceNote"`
	BigAcceptanceDate bigquery.NullDateTime        `json:"-" firestore:"-" bigquery:"acceptanceDate"`
	Step              string                       `json:"step,omitempty" firestore:"step,omitempty" bigquery:"step"`
	ProducerCode      string                       `json:"producerCode,omitempty" firestore:"producerCode,omitempty" bigquery:"producerCode"`
	NetworkUid        string                       `json:"networkUid" firestore:"networkUid" bigquery:"networkUid"`
	ProducerUid       string                       `json:"producerUid" firestore:"producerUid" bigquery:"producerUid"`
	ProducerType      string                       `json:"producerType" firestore:"producerType" bigquery:"producerType"`
	Annuity           int                          `json:"annuity" firestore:"annuity" bigquery:"annuity"`
	IsAutoRenew       bool                         `json:"isAutoRenew" firestore:"isAutoRenew" bigquery:"isAutoRenew"`
	IsRenewable       bool                         `json:"isRenewable" firestore:"isRenewable" bigquery:"isRenewable"`
	PolicyType        string                       `json:"policyType" firestore:"policyType" bigquery:"policyType"`
	QuoteType         string                       `json:"quoteType" firestore:"quoteType" bigquery:"quoteType"`
	HasMandate        bool                         `json:"hasMandate" firestore:"hasMandate" bigquery:"hasMandate"`
	DeclaredClaims    []DeclaredClaim              `json:"declaredClaims,omitempty" firestore:"declaredClaims,omitempty" bigquery:"declaredClaims,omitempty"`
	HasBond           string                       `json:"hasBond,omitempty" firestore:"hasBond,omitempty" bigquery:"hasBond,omitempty"`
	Bond              string                       `json:"bond,omitempty" firestore:"bond,omitempty" bigquery:"bond,omitempty"`
	Clause            string                       `json:"clause,omitempty" firestore:"clause,omitempty" bigquery:"clause,omitempty"`

	// DEPRECATED FIELDS

	RejectReasons string `json:"rejectReasons,omitempty" firestore:"rejectReasons,omitempty" bigquery:"-"` // DEPRECATED
	AgentUid      string `json:"agentUid,omitempty" firestore:"agentUid,omitempty" bigquery:"agentUid"`    // DEPRECATED
	AgencyUid     string `json:"agencyUid,omitempty" firestore:"agencyUid,omitempty" bigquery:"agencyUid"` // DEPRECATED
}

type DeclaredClaim struct {
	GuaranteeSlug string          `json:"guaranteeSlug,omitempty" firestore:"guaranteeSlug,omitempty"`
	Year          int16           `json:"year,omitempty" firestore:"year,omitempty"`
	Total         int16           `json:"total,omitempty" firestore:"total,omitempty"`
	History       []DeclaredClaim `json:"history,omitempty" firestore:"history,omitempty"`
}

type RenewHistory struct {
	Amount       float64   `firestore:"amount,omitempty" json:"amount,omitempty"`
	StartDate    time.Time `firestore:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate      time.Time `firestore:"endDate,omitempty" json:"endDate,omitempty"`
	CreationDate time.Time `firestore:"creationDate,omitempty" json:"creationDate,omitempty"`
}
type PriceGroup struct {
	Name              string  `firestore:"name,omitempty" json:"name,omitempty" bigquery:"name"`
	TaxAmount         float64 `firestore:"taxAmount,omitempty" json:"taxAmount,omitempty" bigquery:"taxAmount"`
	PriceNett         float64 `firestore:"priceNett,omitempty" json:"priceNett,omitempty" bigquery:"priceNett"`
	PriceGross        float64 `firestore:"priceGross,omitempty" json:"priceGross,omitempty" bigquery:"priceGross"`
	TaxAmountMonthly  float64 `json:"taxAmountMonthly,omitempty" firestore:"taxAmountMonthly,omitempty" bigquery:"taxAmountMonthly"`
	PriceNettMonthly  float64 `json:"priceNettMonthly,omitempty" firestore:"priceNettMonthly,omitempty" bigquery:"priceNettMonthly"`
	PriceGrossMonthly float64 `json:"priceGrossMonthly,omitempty" firestore:"priceGrossMonthly,omitempty" bigquery:"priceGrossMonthly"`
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
	Name         string  `firestore:"name,omitempty" json:"name,omitempty" bigquery:"-"`
	Description  string  `firestore:"description,omitempty" json:"description,omitempty" bigquery:"-"`
	Net          float64 `firestore:"net" json:"net" bigquery:"-"`
	Tax          float64 `firestore:"tax" json:"tax" bigquery:"-"`
	Gross        float64 `firestore:"gross" json:"gross" bigquery:"-"`
	Delta        float64 `firestore:"delta" json:"delta" bigquery:"-"`
	Discount     float64 `firestore:"discount" json:"discount" bigquery:"-"`
	NettMonthly  float64 `json:"nettMonthly,omitempty" firestore:"nettMonthly,omitempty" bigquery:"-"`
	GrossMonthly float64 `json:"grossMonthly,omitempty" firestore:"grossMonthly,omitempty" bigquery:"-"`
}

func (p *Policy) Normalize() {
	p.Contractor.Normalize()
	if p.Contractors != nil {
		for index, _ := range *p.Contractors {
			(*p.Contractors)[index].Normalize()
		}
	}
	for index, _ := range p.Assets {
		p.Assets[index].Normalize()
	}
}

func isLeapYear(year int) bool {
	// A leap year is either divisible by 400 or divisible by 4 but not by 100.
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
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

	startDateYearDay := startDate.YearDay()
	if isLeapYear(startDate.Year()) {
		startDateYearDay -= 1
	}

	log.Printf("[CalculateContractorAge] startDate.YearDay %d - birthdate.YearDay %d", startDateYearDay, birthdate.YearDay())

	if startDateYearDay < birthdate.YearDay() && !(startDate.Month() == birthdate.Month() && startDate.Day() == birthdate.Day()) {
		age--
	}

	log.Printf("[CalculateContractorAge] age: %d", age)

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

func (policy *Policy) BigQueryParse() {
	var (
		data []byte
		err  error
	)

	if data, err = policy.Marshal(); err != nil {
		return
	}

	policy.Data = string(data)
	policy.BigStartDate = civil.DateTimeOf(policy.StartDate)
	policy.BigRenewDate = civil.DateTimeOf(policy.RenewDate)
	policy.BigEndDate = civil.DateTimeOf(policy.EndDate)
	policy.BigEmitDate = civil.DateTimeOf(policy.EmitDate)
	policy.BigStatusHistory = strings.Join(policy.StatusHistory, ",")
	if policy.ReservedInfo != nil {
		policy.BigReasons = strings.Join(policy.ReservedInfo.Reasons, ",")
		policy.BigAcceptanceNote = policy.ReservedInfo.AcceptanceNote
		policy.BigAcceptanceDate = lib.GetBigQueryNullDateTime(policy.ReservedInfo.AcceptanceDate)
	}
}

func (policy *Policy) BigquerySave(origin string) {
	log.Printf("[policy.BigquerySave] parsing data for policy %s", policy.Uid)

	policyBig := lib.GetDatasetByEnv(origin, PolicyCollection)

	policy.BigQueryParse()

	log.Println("[policy.BigquerySave] saving to bigquery...")
	if err := lib.InsertRowsBigQuery(WoptaDataset, policyBig, policy); err != nil {
		log.Println("[policy.BigquerySave] error saving policy to bigquery: ", err.Error())
		return
	}
	log.Println("[policy.BigquerySave] bigquery saved!")
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

func GetChannel(policy *Policy) string {
	var channel string

	switch policy.ProducerType {
	case AgentNetworkNodeType:
		channel = AgentChannel
	case AgencyNetworkNodeType:
		channel = AgencyChannel
	default:
		channel = ECommerceChannel
	}

	return channel
}

func (policy *Policy) GetFlow(networkNode *NetworkNode, warrant *Warrant) (string, *NodeSetting) {
	var (
		channel  = policy.Channel
		flowByte []byte
		flowName string
		flowFile NodeSetting
		err      error
	)

	log.Printf("[Policy.GetFlow] loading file for channel %s", channel)

	// Retrocompatibility with old emitted policies without channel when there was only e-commerce
	if channel == "" {
		channel = ECommerceChannel
		policy.Channel = ECommerceChannel
		log.Println("[Policy.GetFlow] overriding unset channel as e-commerce")
	}

	switch channel {
	case NetworkChannel:
		flowName, flowByte = networkNode.GetNetworkNodeFlow(policy.Name, warrant)
	case ECommerceChannel, MgaChannel:
		flowName = channel
		flowByte = lib.GetFilesByEnv(fmt.Sprintf(FlowFileFormat, channel))
	default:
		log.Printf("[Policy.GetFlow] error unavailable channel: '%s'", channel)
		return flowName, nil
	}

	if len(flowByte) == 0 {
		log.Printf("[Policy.GetFlow] error flowFile '%s' empty", flowName)
		return flowName, nil
	}

	err = json.Unmarshal(flowByte, &flowFile)
	if err != nil {
		log.Printf("[Policy.GetFlow] error unmarshaling flow '%s' file: %s", flowName, err.Error())
		return flowName, nil
	}

	return flowName, &flowFile
}

func (policy *Policy) GetNumberOfRates() int {
	switch policy.PaymentSplit {
	case string(PaySplitYear), string(PaySplitYearly), string(PaySplitSingleInstallment):
		return 1
	case string(PaySplitMonthly):
		return 12
	default:
		return 0
	}
}

func (policy *Policy) GetDurationInYears() int {
	return policy.EndDate.Year() - policy.StartDate.Year()
}

func (policy *Policy) SanitizePaymentData() {
	if policy.Payment == "" || policy.Payment == "fabrik" {
		policy.Payment = FabrickPaymentProvider
	}

	if policy.PaymentSplit == string(PaySplitYear) {
		policy.PaymentSplit = string(PaySplitYearly)
	}

	if policy.PaymentMode == "" {
		if policy.PaymentSplit == string(PaySplitYearly) {
			policy.PaymentMode = PaymentModeSingle
		} else {
			policy.PaymentMode = PaymentModeRecurrent
		}
	}
}

func (policy *Policy) CheckStartDateValidity(maxElapsedDays uint) error {
	if policy.CompanyEmit {
		return nil
	}

	truncateDuration := 24 * time.Hour

	now := time.Now().UTC().Truncate(truncateDuration)
	lastValidDate := policy.StartDate.Truncate(truncateDuration).AddDate(0, 0, int(maxElapsedDays))

	if lastValidDate.Before(now) {
		return fmt.Errorf("policy start date expired")
	}

	return nil
}

func (policy *Policy) HasPrivacyConsens() bool {
	if policy.Contractor.Consens != nil {
		for _, v := range *policy.Contractor.Consens {
			if v.Key == 2 {
				return v.Answer
			}
		}
	}

	return false
}
