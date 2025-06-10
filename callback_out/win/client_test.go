package win_test

import (
	"testing"

	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/callback_out/win"
)

func TestWinDecodeAction(t *testing.T) {
	mockNodeCode := "Test.001"
	mockClient := win.NewClient(mockNodeCode)

	var testCases = []struct {
		name  string
		input base.CallbackoutAction
		want  []base.CallbackoutAction
	}{
		{
			name:  "Paid",
			input: base.Paid,
			want:  []base.CallbackoutAction{base.Paid},
		},
		{
			name:  "Proposal",
			input: base.Proposal,
			want:  []base.CallbackoutAction{base.Proposal},
		},
		{
			name:  "RequestApproval",
			input: base.RequestApproval,
			want:  []base.CallbackoutAction{base.RequestApproval},
		},
		{
			name:  "Emit",
			input: base.Emit,
			want:  []base.CallbackoutAction{base.Emit},
		},
		{
			name:  "EmitRemittance",
			input: base.EmitRemittance,
			want:  []base.CallbackoutAction{base.Emit, base.Paid},
		},
		{
			name:  "Unhandled action",
			input: "NON_EXISTING",
			want:  []base.CallbackoutAction{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := mockClient.DecodeAction(tc.input)
			if len(got) != len(tc.want) {
				t.Fatalf("expected %d action for paid but got: %d (%v)", len(tc.want), len(got), got)
			}
			for idx, action := range got {
				if action != tc.want[idx] {
					t.Fatalf("expected %s action but got: %s", tc.want[idx], action)
				}
			}
		})
	}
}
