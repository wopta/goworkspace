package base

import (
	"io"
	"net/http"
)

type CallbackInfo struct {
	ReqBody       []byte
	ReqMethod     string
	ReqPath       string
	ResBody       []byte
	ResStatusCode int
	ResAction     CallbackoutAction
	Error         error
}

func (c *CallbackInfo) FromRequestResponse(action CallbackoutAction, resp *http.Response, req *http.Request) error {

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	c.ReqMethod = req.Method
	c.ReqPath = resp.Request.Host + resp.Request.URL.RequestURI()
	c.ReqBody = reqBody

	c.ResBody = resBody
	c.ResStatusCode = resp.StatusCode
	c.ResAction = action
	return nil
}

type CallbackoutAction string

const (
	Emit            CallbackoutAction = "Emit"
	Paid            CallbackoutAction = "Paid"
	Proposal        CallbackoutAction = "Proposal"
	RequestApproval CallbackoutAction = "RequestApproval"
	Signed          CallbackoutAction = "Signed"
	Approved        CallbackoutAction = "Approved"
	Rejected        CallbackoutAction = "Rejected"

	EmitRemittance CallbackoutAction = "EmitRemittance"
)

type CallbackExternalConfig struct {
	Events   map[CallbackoutAction]bool `json:"events"`
	AuthType string                     `json:"authType"` // basic, api-key
}

func GetAvailableActions() map[CallbackoutAction][]CallbackoutAction {
	return map[CallbackoutAction][]CallbackoutAction{
		Proposal:        {Proposal},
		RequestApproval: {RequestApproval},
		Emit:            {Emit},
		Signed:          {Signed},
		Paid:            {Paid},
		EmitRemittance:  {Emit, Paid},
		Approved:        {Approved},
		Rejected:        {Rejected},
	}
}
