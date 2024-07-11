package internal

import "net/http"

type CallbackInfo struct {
	Request     *http.Request
	RequestBody []byte
	Response    *http.Response
	Error       error
}

type CallbackoutAction = string

var (
	Proposal        CallbackoutAction = "Proposal"
	RequestApproval CallbackoutAction = "RequestApproval"
	Emit            CallbackoutAction = "Emit"
	Paid            CallbackoutAction = "Paid"
)
