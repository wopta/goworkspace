package fabrick

import "encoding/json"

func UnmarshalFabrickPaymentResponse(data []byte) (FabrickPaymentResponse, error) {
	var r FabrickPaymentResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FabrickPaymentResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FabrickPaymentResponse struct {
	Status  *string       `json:"status,omitempty"`
	Errors  []interface{} `json:"errors,omitempty"`
	Payload *Payload      `json:"payload,omitempty"`
}

type Payload struct {
	ExternalID        *string     `json:"externalId,omitempty"`
	PaymentID         *string     `json:"paymentId,omitempty"`
	MerchantID        *string     `json:"merchantId,omitempty"`
	PaymentPageURL    *string     `json:"paymentPageUrl,omitempty"`
	PaymentPageURLB2B *string     `json:"paymentPageUrlB2B,omitempty"`
	TokenB2B          *string     `json:"tokenB2B,omitempty"`
	Coupon            interface{} `json:"coupon"`
}

func UnmarshalFabrickPaymentsRequest(data []byte) (FabrickPaymentsRequest, error) {
	var r FabrickPaymentsRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FabrickPaymentsRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FabrickPaymentsRequest struct {
	MerchantID           string               `json:"merchantId,omitempty"`
	ExternalID           string               `json:"externalId,omitempty"`
	PaymentConfiguration PaymentConfiguration `json:"paymentConfiguration,omitempty"`
	Bill                 Bill                 `json:"bill,omitempty"`
}

type Bill struct {
	ExternalID          string               `json:"externalId,omitempty"`
	Amount              float64              `json:"amount,omitempty"`
	Currency            string               `json:"currency,omitempty"`
	Description         string               `json:"description,omitempty"`
	Items               []Item               `json:"items,omitempty"`
	ScheduleTransaction *ScheduleTransaction `json:"scheduleTransaction,omitempty"`
	MandateCreation     string               `json:"mandateCreation,omitempty"`
	Subjects            *[]Subject           `json:"subjects,omitempty"`
}

type Item struct {
	ExternalID       string  `json:"externalId,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	Description      string  `json:"description,omitempty"`
	PagoPADocumentID string  `json:"pagoPADocumentId,omitempty"`
	XInfo            string  `json:"xInfo,omitempty"`
}

type ScheduleTransaction struct {
	DueDate                             string `json:"dueDate,omitempty"`
	PaymentInstrumentResolutionStrategy string `json:"paymentInstrumentResolutionStrategy,omitempty"`
}

type Subject struct {
	Role       string  `json:"role,omitempty"`
	ExternalID string  `json:"externalId,omitempty"`
	Email      string  `json:"email,omitempty"`
	Name       string  `json:"name,omitempty"`
	XInfo      *string `json:"xInfo,omitempty"`
}

type PaymentConfiguration struct {
	ExpirationDate          string                  `json:"expirationDate,omitempty"`
	AllowedPaymentMethods   *[]AllowedPaymentMethod `json:"allowedPaymentMethods,omitempty"`
	PayByLink               *[]PayByLink            `json:"payByLink,omitempty"`
	CallbackURL             string                  `json:"callbackUrl,omitempty"`
	PaymentPageRedirectUrls PaymentPageRedirectUrls `json:"paymentPageRedirectUrls,omitempty"`
}

type AllowedPaymentMethod struct {
	Role           string   `json:"role,omitempty"`
	PaymentMethods []string `json:"paymentMethods,omitempty"`
}

type PayByLink struct {
	Type       string `json:"type,omitempty"`
	Recipients string `json:"recipients,omitempty"`
	Template   string `json:"template,omitempty"`
}

type PaymentPageRedirectUrls struct {
	OnFailure      string `json:"onFailure,omitempty"`
	OnSuccess      string `json:"onSuccess,omitempty"`
	OnInterruption string `json:"onInterruption,omitempty"`
}

// FABRICK MODELS START

type ItemCallback struct {
	ExternalID  *string     `json:"externalId,omitempty"`
	ItemID      *string     `json:"itemId,omitempty"`
	Amount      *float64    `json:"amount,omitempty"`
	Currency    *string     `json:"currency,omitempty"`
	Description *string     `json:"description,omitempty"`
	XInfo       interface{} `json:"xInfo"`
	Status      *string     `json:"status,omitempty"`
	Xinfo       interface{} `json:"xinfo"`
}

type TransactionCallback struct {
	TransactionID       *string     `json:"transactionId,omitempty"`
	TransactionDateTime interface{} `json:"transactionDateTime"`
	Amount              *float64    `json:"amount,omitempty"`
	Currency            *string     `json:"currency,omitempty"`
	GatewayID           interface{} `json:"gatewayId"`
	AcquirerID          interface{} `json:"acquirerId"`
	Status              *string     `json:"status,omitempty"`
	PaymentMethod       *string     `json:"paymentMethod,omitempty"`
}

type BillCallback struct {
	ExternalID     *string               `json:"externalId,omitempty"`
	BillID         *string               `json:"billId,omitempty"`
	Amount         *float64              `json:"amount,omitempty"`
	Currency       *string               `json:"currency,omitempty"`
	Description    *string               `json:"description,omitempty"`
	ReservedAmount *float64              `json:"reservedAmount,omitempty"`
	ResidualAmount *float64              `json:"residualAmount,omitempty"`
	RefundedAmount *float64              `json:"refundedAmount,omitempty"`
	PaidAmout      *float64              `json:"paidAmout,omitempty"`
	Items          []Item                `json:"items,omitempty"`
	Status         string                `json:"status,omitempty"`
	Transactions   []TransactionCallback `json:"transactions,omitempty"`
}

type FabrickCallback struct {
	ExternalID string        `json:"externalId,omitempty"`
	PaymentID  *string       `json:"paymentId,omitempty"`
	Bill       *BillCallback `json:"bill,omitempty"`
}

// FABRICK MODELS END
