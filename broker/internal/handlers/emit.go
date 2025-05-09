package handlers

import (
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/broker/internal/utility"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
)

func AddEmitHandlers(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("emitData", emitData),
		builder.AddHandler("sendMailSign", sendMailSign),
		builder.AddHandler("sign", sign),
		builder.AddHandler("pay", pay),
		builder.AddHandler("setAdvice", setAdvance),
		builder.AddHandler("putUser", updateUserAndNetworkNode),
		builder.AddHandler("sendEmitProposalMail", sendEmitProposalMail),
		builder.AddHandler("end_emit", savePolicy),
	)
}

func emitData(state bpmn.StorageData) error {
	var origin *flow.StringBpmn
	var policy *flow.PolicyDraft
	var err = bpmn.IsError(
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("policy", &policy, state),
	)
	if err != nil {
		return err
	}

	firePolicy := lib.GetDatasetByEnv(origin.String, lib.PolicyCollection)
	emitBase(policy.Policy)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func emitBase(policy *models.Policy) {
	log.AddPrefix("emitBase")
	defer log.PopPrefix()
	log.Printf("Policy Uid %s", policy.Uid)
	firePolicy := lib.PolicyCollection
	now := time.Now().UTC()

	policy.CompanyEmit = true
	policy.CompanyEmitted = false
	policy.EmitDate = now
	policy.BigEmitDate = civil.DateTimeOf(now)
	company, numb, tot := utility.GetSequenceByCompany(strings.ToLower(policy.Company), firePolicy)
	log.Printf("codeCompany: %s", company)
	log.Printf("numberCompany: %d", numb)
	log.Printf("number: %d", tot)
	policy.Number = tot
	policy.NumberCompany = numb
	policy.CodeCompany = company
	policy.RenewDate = policy.StartDate.AddDate(1, 0, 0)
	policy.BigRenewDate = civil.DateTimeOf(policy.RenewDate)
}

func sendMailSign(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var networkNode *flow.NetworkDraft
	var addresses *flow.Addresses
	var flowName *flow.StringBpmn
	var sendEmail *flow.BoolBpmn
	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
		bpmn.GetDataRef("addresses", &addresses, state),
		bpmn.GetDataRef("flowName", &flowName, state),
		bpmn.GetDataRef("sendEmail", &sendEmail, state),
	)
	if err != nil {
		return err
	}

	switch flowName.String {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail.Bool {
			addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
			addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
		} else {
			addresses.ToAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	}

	log.Printf(
		"[sendMailSign] from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	mail.SendMailSign(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String)
	return nil
}

func sign(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var product *flow.ProductDraft
	var networkNode *flow.NetworkDraft
	var addresses *flow.Addresses
	var flowName *flow.StringBpmn
	var sendEmail *flow.BoolBpmn
	var origin *flow.StringBpmn
	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("product", &product, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
		bpmn.GetDataRef("addresses", &addresses, state),
		bpmn.GetDataRef("flowName", &flowName, state),
		bpmn.GetDataRef("sendEmail", &sendEmail, state),
	)
	if err != nil {
		return err
	}
	utility.EmitSign(policy.Policy, product.Product, networkNode.NetworkNode, sendEmail.Bool, origin.String)
	return nil
}

func pay(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var product *flow.ProductDraft
	var mgaProduct *flow.ProductDraft
	var networkNode *flow.NetworkDraft
	var origin *flow.StringBpmn
	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("product", &product, state),
		bpmn.GetDataRef("mgaProduct", &mgaProduct, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}

	utility.EmitPay(policy.Policy, origin.String, product.Product, mgaProduct.Product, networkNode.NetworkNode)
	if policy.PayUrl == "" {
		return fmt.Errorf("missing payment url")
	}
	return nil
}

func setAdvance(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var product *flow.ProductDraft
	var mgaProduct *flow.ProductDraft
	var networkNode *flow.NetworkDraft
	var origin *flow.StringBpmn
	var paymentSplit *flow.StringBpmn
	var paymentMode *flow.StringBpmn
	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("paymentSplit", &paymentSplit, state),
		bpmn.GetDataRef("paymentMode", &paymentMode, state),
		bpmn.GetDataRef("product", &product, state),
		bpmn.GetDataRef("mgaProduct", &mgaProduct, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}
	utility.SetAdvance(policy.Policy, origin.String, product.Product, mgaProduct.Product, networkNode.NetworkNode, paymentSplit.String, paymentMode.String)
	return nil
}

func updateUserAndNetworkNode(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var networkNode *flow.NetworkDraft
	var origin *flow.StringBpmn
	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}
	// promote documents from temp bucket to user and connect it to policy
	err = plc.SetUserIntoPolicyContractor(policy.Policy, origin.String)
	if err != nil {
		log.ErrorF("[putUser] ERROR SetUserIntoPolicyContractor %s", err.Error())
		return err
	}
	return network.UpdateNetworkNodePortfolio(origin.String, policy.Policy, networkNode.NetworkNode)
}

func sendEmitProposalMail(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var networkNode *flow.NetworkDraft
	var origin *flow.StringBpmn
	var addresses *flow.Addresses
	var flowName *flow.StringBpmn
	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
		bpmn.GetDataRef("flowName", &flowName, state),
		bpmn.GetDataRef("addresses", &addresses, state),
	)
	if err != nil {
		return err
	}
	// TODO: remove when a proper solution to handle PMI is found
	if policy.Name == models.PmiProduct {
		return nil
	}

	if policy.IsReserved {
		return nil
	}

	addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	addresses.CcAddress = mail.Address{}
	switch flowName.String {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
	}

	log.Printf(
		"[sendEmitProposalMail] from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)

	mail.SendMailProposal(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String, []string{models.ProposalAttachmentName})
	return nil
}
