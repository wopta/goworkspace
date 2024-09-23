package models

import "github.com/wopta/goworkspace/callback_out/internal"

var (
	Proposal        internal.CallbackoutAction = "Proposal"
	RequestApproval internal.CallbackoutAction = "RequestApproval"
	Emit            internal.CallbackoutAction = "Emit"
	Paid            internal.CallbackoutAction = "Paid"
	EmitRemittance  internal.CallbackoutAction = "EmitRemittance"
)

func GetAvailableActions() map[string][]string {
	return map[string][]string{
		Proposal:        {internal.Proposal},
		RequestApproval: {internal.RequestApproval},
		Emit:            {internal.Emit},
		Paid:            {internal.Paid},
		EmitRemittance:  {internal.Emit, internal.Paid},
	}
}
