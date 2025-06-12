package fabrick

import (
	"encoding/json"
	"errors"
)

var (
	ErrPolicyNotFound   = errors.New("policy not found")
	ErrProviderIdNotSet = errors.New("providerId not set")
)

type FabrickCallback struct {
}

type FabrickResponse struct {
	Result         bool   `json:"result"`
	RequestPayload string `json:"requestPayload"`
	Locale         string `json:"locale"`
}

type FabrickRequestPayload = map[string]any

type FabrickRequest struct {
	ExternalID string `json:"externalId,omitempty"`
	PaymentID  string `json:"paymentId,omitempty"`
	Bill       struct {
		Transactions []struct {
			PaymentMethod string `json:"paymentMethod,omitempty"`
		} `json:"transactions,omitempty"`
	} `json:"bill,omitempty"`
}

func (r *FabrickRequest) FromPayload(payload FabrickRequestPayload) (string, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(bytes), json.Unmarshal(bytes, r)
}
