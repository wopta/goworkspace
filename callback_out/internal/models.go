package internal

import "net/http"

type CallbackInfo struct {
	Request     *http.Request
	RequestBody []byte
	Response    *http.Response
	Error       error
}

type CallbackoutAction = string

const (
	Proposal        CallbackoutAction = "Proposal"
	RequestApproval CallbackoutAction = "RequestApproval"
	Emit            CallbackoutAction = "Emit"
	Signed          CallbackoutAction = "Signed"
	Paid            CallbackoutAction = "Paid"
	Approved        CallbackoutAction = "Approved"
	Rejected        CallbackoutAction = "Rejected"
)

type CallbackExternalConfig struct {
	Events   map[string]bool `json:"events"`
	AuthType string          `json:"authType"` // basic, api-key
}
