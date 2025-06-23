package handlers

import (
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

func AddEmitHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		builder.AddHandler("emitWithSequence", emitBaseWithSequence),
		builder.AddHandler("emitNoSequence", emitBaseNoSequence),
		builder.AddHandler("sendMailSign", sendMailSign),
		builder.AddHandler("sign", sign),
		builder.AddHandler("pay", pay),
		builder.AddHandler("setAdvice", setAdvance),
		builder.AddHandler("putUser", updateUserAndNetworkNode),
		builder.AddHandler("sendEmitProposalMail", sendEmitProposalMail),
		builder.AddHandler("end_emit", savePolicy),
	)
}

func emitBaseWithSequence(state *bpmnEngine.StorageBpnm) error {
	var origin *flow.String
	var policy *flow.Policy
	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("policy", &policy, state),
	)
	if err != nil {
		return err
	}
	log.AddPrefix("emitBase")
	defer log.PopPrefix()

	log.Printf("Policy Uid %s", policy.Uid)
	firePolicy := lib.GetDatasetByEnv(origin.String, lib.PolicyCollection)
	now := time.Now().UTC()

	policy.CompanyEmit = true
	policy.CompanyEmitted = false
	policy.EmitDate = now
	policy.BigEmitDate = civil.DateTimeOf(now)
	policy.RenewDate = policy.StartDate.AddDate(1, 0, 0)
	policy.BigRenewDate = civil.DateTimeOf(policy.RenewDate)
	company, numb, tot := utility.GetSequenceByCompany(strings.ToLower(policy.Company), firePolicy)
	log.Printf("codeCompany: %s", company)
	log.Printf("numberCompany: %d", numb)
	log.Printf("number: %d", tot)
	policy.Number = tot
	policy.NumberCompany = numb
	policy.CodeCompany = company
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func emitBaseNoSequence(state *bpmnEngine.StorageBpnm) error {
	var origin *flow.String
	var policy *flow.Policy
	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("policy", &policy, state),
	)
	if err != nil {
		return err
	}
	log.AddPrefix("emitBase")
	defer log.PopPrefix()

	log.Printf("Policy Uid %s", policy.Uid)
	now := time.Now().UTC()

	policy.CompanyEmit = true
	policy.CompanyEmitted = false
	policy.EmitDate = now
	policy.BigEmitDate = civil.DateTimeOf(now)
	policy.RenewDate = policy.StartDate.AddDate(1, 0, 0)
	policy.BigRenewDate = civil.DateTimeOf(policy.RenewDate)
	return nil
}

func sendMailSign(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var addresses *flow.Addresses
	var flowName *flow.String
	var sendEmail *flow.BoolBpmn
	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("addresses", &addresses, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
		bpmnEngine.GetDataRef("sendEmail", &sendEmail, state),
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

func sign(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var product *flow.Product
	var networkNode *flow.Network
	var addresses *flow.Addresses
	var flowName *flow.String
	var sendEmail *flow.BoolBpmn
	var origin *flow.String
	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("product", &product, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("addresses", &addresses, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
		bpmnEngine.GetDataRef("sendEmail", &sendEmail, state),
	)
	if err != nil {
		return err
	}
	err = utility.SignFiles(policy.Policy, product.Product, networkNode.NetworkNode, true, origin.String)
	if err != nil {
		return err
	}
	return nil
}

func pay(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var product *flow.Product
	var mgaProduct *flow.Product
	var networkNode *flow.Network
	var origin *flow.String
	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("product", &product, state),
		bpmnEngine.GetDataRef("mgaProduct", &mgaProduct, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
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

func setAdvance(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var product *flow.Product
	var mgaProduct *flow.Product
	var networkNode *flow.Network
	var origin *flow.String
	var paymentSplit *flow.String
	var paymentMode *flow.String
	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("paymentSplit", &paymentSplit, state),
		bpmnEngine.GetDataRef("paymentMode", &paymentMode, state),
		bpmnEngine.GetDataRef("product", &product, state),
		bpmnEngine.GetDataRef("mgaProduct", &mgaProduct, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}
	utility.SetAdvance(policy.Policy, origin.String, product.Product, mgaProduct.Product, networkNode.NetworkNode, paymentSplit.String, paymentMode.String)
	return nil
}

func updateUserAndNetworkNode(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var origin *flow.String
	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
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

func sendEmitProposalMail(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var origin *flow.String
	var addresses *flow.Addresses
	var flowName *flow.String
	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
		bpmnEngine.GetDataRef("addresses", &addresses, state),
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
	mail.SendMailProposal(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String, []string{models.ContractDocumentFormat})
	return nil
}
