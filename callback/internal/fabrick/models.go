package fabrick

import (
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
