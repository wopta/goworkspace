package win_test

import (
	"testing"

	"github.com/wopta/goworkspace/callback_out/internal"
	"github.com/wopta/goworkspace/callback_out/types"
	"github.com/wopta/goworkspace/callback_out/win"
)

func TestWinDecodeAction(t *testing.T) {
	mockNodeCode := "Test.001"
	mockClient := win.NewClient(mockNodeCode)

	res := mockClient.DecodeAction(types.Paid)
	if len(res) != 1 {
		t.Fatalf("expected 1 action for paid but got: %d (%v)", len(res), res)
	}
	if res[0] != internal.Paid {
		t.Fatalf("expected %s action but got: %s", internal.Paid, res[0])
	}

	res = mockClient.DecodeAction(types.Proposal)
	if len(res) != 1 {
		t.Fatalf("expected 1 action for proposal but got: %d (%v)", len(res), res)
	}
	if res[0] != internal.Proposal {
		t.Fatalf("expected %s action but got: %s", internal.Proposal, res[0])
	}

	res = mockClient.DecodeAction(types.RequestApproval)
	if len(res) != 1 {
		t.Fatalf("expected 1 action for request approval but got: %d (%v)", len(res), res)
	}
	if res[0] != internal.RequestApproval {
		t.Fatalf("expected %s action but got: %s", internal.RequestApproval, res[0])
	}

	res = mockClient.DecodeAction(types.Emit)
	if len(res) != 1 {
		t.Fatalf("expected 1 action for emit but got: %d (%v)", len(res), res)
	}
	if res[0] != internal.Emit {
		t.Fatalf("expected %s action but got: %s", internal.Emit, res[0])
	}

	res = mockClient.DecodeAction(types.EmitRemittance)
	if len(res) != 2 {
		t.Fatalf("expected 2 actions for emit remittance but got: %d (%v)", len(res), res)
	}
	if res[0] != internal.Emit {
		t.Fatalf("expected %s action but got: %s", internal.Emit, res[0])
	}
	if res[1] != internal.Paid {
		t.Fatalf("expected %s action but got: %s", internal.Paid, res[1])
	}

	res = mockClient.DecodeAction("WRONG_ACTION")
	if len(res) != 0 {
		t.Fatalf("expected 0 actions for WRONG_ACTION but got: %d (%v)", len(res), res)
	}
}
