package win_test

import (
	"testing"

	"gitlab.dev.wopta.it/goworkspace/callback_out/internal"
	md "gitlab.dev.wopta.it/goworkspace/callback_out/models"
	"gitlab.dev.wopta.it/goworkspace/callback_out/win"
)

func TestWinDecodeAction(t *testing.T) {
	mockNodeCode := "Test.001"
	mockClient := win.NewClient(mockNodeCode)

	var testCases = []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "Paid",
			input: md.Paid,
			want:  []string{internal.Paid},
		},
		{
			name:  "Proposal",
			input: md.Proposal,
			want:  []string{internal.Proposal},
		},
		{
			name:  "RequestApproval",
			input: md.RequestApproval,
			want:  []string{internal.RequestApproval},
		},
		{
			name:  "Emit",
			input: md.Emit,
			want:  []string{internal.Emit},
		},
		{
			name:  "EmitRemittance",
			input: md.EmitRemittance,
			want:  []string{internal.Emit, internal.Paid},
		},
		{
			name:  "Unhandled action",
			input: "NON_EXISTING",
			want:  []string{},
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
