package payment

import "encoding/json"

func UnmarshalWelcome(data []byte) (FabrickPymentsRequest, error) {
	var r FabrickPymentsRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FabrickPymentsRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FabrickPymentsRequest struct {
	MerchantID           string               `json:"merchantId"`
	ExternalID           string               `json:"externalId"`
	PaymentConfiguration PaymentConfiguration `json:"paymentConfiguration"`
	Bill                 Bill                 `json:"bill"`
}

type Bill struct {
	ExternalID          string              `json:"externalId"`
	Amount              int64               `json:"amount"`
	Currency            string              `json:"currency"`
	Description         string              `json:"description"`
	Items               []Item              `json:"items"`
	ScheduleTransaction ScheduleTransaction `json:"scheduleTransaction"`
	MandateCreation     string              `json:"mandateCreation"`
	Subjects            []Subject           `json:"subjects"`
}

type Item struct {
	ExternalID       string `json:"externalId"`
	Amount           int64  `json:"amount"`
	Currency         string `json:"currency"`
	Description      string `json:"description"`
	PagoPADocumentID string `json:"pagoPADocumentId"`
	XInfo            string `json:"xInfo"`
}

type ScheduleTransaction struct {
	DueDate                             string `json:"dueDate"`
	PaymentInstrumentResolutionStrategy string `json:"paymentInstrumentResolutionStrategy"`
}

type Subject struct {
	Role       string  `json:"role"`
	ExternalID string  `json:"externalId"`
	Email      string  `json:"email"`
	Name       string  `json:"name"`
	XInfo      *string `json:"xInfo,omitempty"`
}

type PaymentConfiguration struct {
	ExpirationDate          string                  `json:"expirationDate"`
	AllowedPaymentMethods   []AllowedPaymentMethod  `json:"allowedPaymentMethods"`
	PayByLink               []PayByLink             `json:"payByLink"`
	CallbackURL             string                  `json:"callbackUrl"`
	PaymentPageRedirectUrls PaymentPageRedirectUrls `json:"paymentPageRedirectUrls"`
}

type AllowedPaymentMethod struct {
	Role           string   `json:"role"`
	PaymentMethods []string `json:"paymentMethods"`
}

type PayByLink struct {
	Type       string `json:"type"`
	Recipients string `json:"recipients"`
	Template   string `json:"template"`
}

type PaymentPageRedirectUrls struct {
	OnFailure string `json:"onFailure"`
	OnSuccess string `json:"onSuccess"`
}
