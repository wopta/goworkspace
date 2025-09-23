package models

import (
	"encoding/json"
)

func UnmarshalFabrickPaymentResponse(data []byte) (FabrickPaymentResponse, error) {
	var r FabrickPaymentResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FabrickPaymentResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FabrickPaymentResponse struct {
	Status  *string         `json:"status,omitempty"`
	Errors  []interface{}   `json:"errors,omitempty"`
	Payload *FabrickPayload `json:"payload,omitempty"`
}

type FabrickPayload struct {
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
	MerchantID           string                      `json:"merchantId,omitempty"`
	PaymentID            string                      `json:"paymentId,omitempty"`
	ExternalID           string                      `json:"externalId,omitempty"`
	PaymentConfiguration FabrickPaymentConfiguration `json:"paymentConfiguration,omitempty"`
	Bill                 FabrickBill                 `json:"bill,omitempty"`
}

type FabrickBill struct {
	ExternalID          string                      `json:"externalId,omitempty"`
	Amount              float64                     `json:"amount,omitempty"`
	Currency            string                      `json:"currency,omitempty"`
	Description         string                      `json:"description,omitempty"`
	Items               []FabrickItem               `json:"items,omitempty"`
	ScheduleTransaction *FabrickScheduleTransaction `json:"scheduleTransaction,omitempty"`
	MandateCreation     string                      `json:"mandateCreation,omitempty"`
	Subjects            *[]FabrickSubject           `json:"subjects,omitempty"`
	Transactions        []TransactionRequest        `json:"transactions,omitempty"`
}

type FabrickItem struct {
	ExternalID       string  `json:"externalId,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	Description      string  `json:"description,omitempty"`
	PagoPADocumentID string  `json:"pagoPADocumentId,omitempty"`
	XInfo            string  `json:"xInfo,omitempty"`
}

type FabrickScheduleTransaction struct {
	DueDate                             string `json:"dueDate,omitempty"`
	PaymentInstrumentResolutionStrategy string `json:"paymentInstrumentResolutionStrategy,omitempty"`
}

type FabrickSubject struct {
	Role       string  `json:"role,omitempty"`
	ExternalID string  `json:"externalId,omitempty"`
	Email      string  `json:"email,omitempty"`
	Name       string  `json:"name,omitempty"`
	Surname    string  `json:"surname,omitempty"`
	XInfo      *string `json:"xInfo,omitempty"`
}

type FabrickPaymentConfiguration struct {
	ExpirationDate          string                         `json:"expirationDate,omitempty"`
	AllowedPaymentMethods   *[]FabrickAllowedPaymentMethod `json:"allowedPaymentMethods,omitempty"`
	PayByLink               *[]PayByLink                   `json:"payByLink,omitempty"`
	CallbackURL             string                         `json:"callbackUrl,omitempty"`
	PaymentPageRedirectUrls FabrickPaymentPageRedirectUrls `json:"paymentPageRedirectUrls,omitempty"`
}

type FabrickAllowedPaymentMethod struct {
	Role           string   `json:"role,omitempty"`
	PaymentMethods []string `json:"paymentMethods,omitempty"`
}

type PayByLink struct {
	Type       string `json:"type,omitempty"`
	Recipients string `json:"recipients,omitempty"`
	Template   string `json:"template,omitempty"`
}

type FabrickPaymentPageRedirectUrls struct {
	OnFailure      string `json:"onFailure,omitempty"`
	OnSuccess      string `json:"onSuccess,omitempty"`
	OnInterruption string `json:"onInterruption,omitempty"`
}

type TransactionRequest struct {
	TransactionID       *string     `json:"transactionId,omitempty"`
	TransactionDateTime interface{} `json:"transactionDateTime"`
	Amount              *float64    `json:"amount,omitempty"`
	Currency            string      `json:"currency,omitempty"`
	GatewayID           interface{} `json:"gatewayId"`
	AcquirerID          interface{} `json:"acquirerId"`
	Status              *string     `json:"status,omitempty"`
	PaymentMethod       string      `json:"paymentMethod,omitempty"`
}
