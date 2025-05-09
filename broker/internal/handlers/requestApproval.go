package handlers

import (
	"os"

	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/broker/internal/utility"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

func AddRequestApprovaHandlers(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("setRequestApprovalData", setRequestApprova),
		builder.AddHandler("sendRequestApprovalMail", sendRequestApprovalMail),
	)
}

func setRequestApprova(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var networkNode *flow.NetworkDraft
	var origin *flow.StringBpmn
	var mgaProduct *flow.ProductDraft
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("mgaProduct", &mgaProduct, state),
	)
	if err != nil {
		return err
	}
	firePolicy := lib.GetDatasetByEnv(origin.String, lib.PolicyCollection)

	utility.SetRequestApprovalData(policy.Policy, networkNode.NetworkNode, mgaProduct.Product, origin.String)

	log.Printf("[setRequestApproval] saving policy with uid %s to Firestore....", policy.Uid)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendRequestApprovalMail(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var addresses *flow.Addresses
	var flowName *flow.StringBpmn
	var networkNode *flow.NetworkDraft
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("addresses", &addresses, state),
		bpmn.GetDataRef("flowName", &flowName, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}
	if policy.Status == models.PolicyStatusWaitForApprovalMga {
		return nil
	}

	addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	addresses.CcAddress = mail.Address{}
	switch flowName.String {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
	case models.ECommerceChannel:
		addresses.ToAddress = mail.Address{} // fail safe for not sending email on ecommerce reserved
	}

	if policy.Name == "qbe" {
		addresses.ToAddress = mail.Address{
			Address: os.Getenv("QBE_RESERVED_MAIL"),
		}
		addresses.CcAddress = mail.Address{}
	}

	mail.SendMailReserved(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String,
		[]string{models.ProposalAttachmentName})
	return nil
}
