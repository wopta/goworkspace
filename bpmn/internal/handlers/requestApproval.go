package handlers

import (
	"os"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func AddRequestApprovaHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		builder.AddHandler("setRequestApprovalData", setRequestApprova),
		builder.AddHandler("sendRequestApprovalMail", sendRequestApprovalMail),
	)
}

func setRequestApprova(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var origin *flow.String
	var mgaProduct *flow.Product
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("mgaProduct", &mgaProduct, state),
	)
	if err != nil {
		return err
	}
	firePolicy := lib.GetDatasetByEnv(origin.String, lib.PolicyCollection)

	utility.SetRequestApprovalData(policy.Policy, networkNode.NetworkNode, mgaProduct.Product, origin.String)

	log.Printf("saving policy with uid %s to Firestore....", policy.Uid)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendRequestApprovalMail(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var addresses *flow.Addresses
	var flowName *flow.String
	var networkNode *flow.Network
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("addresses", &addresses, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
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
