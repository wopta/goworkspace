package handlers

import (
	"os"

	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/broker/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func AddRequestApprovaHandlers(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("setRequestApprovalData", setRequestApprova),
		builder.AddHandler("sendRequestApprovalMail", sendRequestApprovalMail),
	)
}

func setRequestApprova(state bpmn.StorageData) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var mgaProduct *flow.Product
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
		bpmn.GetDataRef("mgaProduct", &mgaProduct, state),
	)
	if err != nil {
		return err
	}

	utility.SetRequestApprovalData(policy.Policy, networkNode.NetworkNode, mgaProduct.Product)

	log.Printf("saving policy with uid %s to Firestore....", policy.Uid)
	return lib.SetFirestoreErr(lib.PolicyCollection, policy.Uid, policy)
}

func sendRequestApprovalMail(state bpmn.StorageData) error {
	var policy *flow.Policy
	var addresses *flow.Addresses
	var flowName *flow.String
	var networkNode *flow.Network
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
