package models

import "github.com/wopta/goworkspace/callback_out/internal"

var (
	Proposal        internal.CallbackoutAction = "Proposal"
	RequestApproval internal.CallbackoutAction = "RequestApproval"
	Emit            internal.CallbackoutAction = "Emit"
	Signed          internal.CallbackoutAction = "Signed"
	Paid            internal.CallbackoutAction = "Paid"
	EmitRemittance  internal.CallbackoutAction = "EmitRemittance"
	Approved        internal.CallbackoutAction = "Approved"
	Rejected        internal.CallbackoutAction = "Rejected"
)

func GetAvailableActions() map[string][]string {
	return map[string][]string{
		Proposal:        {internal.Proposal},
		RequestApproval: {internal.RequestApproval},
		Emit:            {internal.Emit},
		Signed:          {internal.Signed},
		Paid:            {internal.Paid},
		EmitRemittance:  {internal.Emit, internal.Paid},
		Approved:        {internal.Approved},
		Rejected:        {internal.Rejected},
	}
}
