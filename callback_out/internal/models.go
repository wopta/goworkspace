package internal

import "net/http"

type CallbackInfo struct {
	Request     *http.Request
	RequestBody []byte
	Response    *http.Response
	Error       error
}

type CallbackoutAction = string

// Externally available types
var (
	ExtProposal        CallbackoutAction = "ExtProposal"
	ExtRequestApproval CallbackoutAction = "ExtRequestApproval"
	ExtEmit            CallbackoutAction = "ExtEmit"
	ExtPaid            CallbackoutAction = "ExtPaid"
	ExtEmitRemittance  CallbackoutAction = "ExtEmitRemittance"
)

// Internal use only
var (
	Proposal        CallbackoutAction = "Proposal"
	RequestApproval CallbackoutAction = "RequestApproval"
	Emit            CallbackoutAction = "Emit"
	Paid            CallbackoutAction = "Paid"
)
