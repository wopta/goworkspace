package callback

import (
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/bpmn"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
)

var (
	origin, providerId, paymentMethod, flowName, trSchedule string
	ccAddress, toAddress, fromAddress                       mail.Address
	networkNode                                             *models.NetworkNode
	mgaProduct                                              *models.Product
	warrant                                                 *models.Warrant
	sendEmail                                               bool
)

const (
	signFlowKey = "sign"
	payFlowKey  = "pay"
)

func runCallbackBpmn(policy *models.Policy, flowKey string) *bpmn.State {
	var (
		flow     []models.Process
		flowFile *models.NodeSetting
	)
	log.AddPrefix("runCallbackBpmn")
	defer log.PopPrefix()
	log.Println("configuring flow --------------------------")

	toAddress = mail.Address{}
	ccAddress = mail.Address{}
	fromAddress = mail.AddressAnna

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	flowName, flowFile = policy.GetFlow(networkNode, warrant)
	if flowFile == nil {
		log.ErrorF("exiting bpmn - flowFile not loaded")
		return nil
	}
	log.Printf("flowName '%s'", flowName)

	mgaProduct = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	// TODO: fix me - maybe get to/from/cc from flowFile.json?
	switch flowKey {
	case signFlowKey:
		flow = flowFile.SignFlow
	case payFlowKey:
		flow = flowFile.PayFlow
	default:
		log.ErrorF("error flow not set")
		return nil
	}

	flowHandlers := lib.SliceMap[models.Process, string](flow, func(h models.Process) string { return h.Name })
	log.Printf("starting %s flow with set handlers: %s", flowKey, strings.Join(flowHandlers, ","))

	state := bpmn.NewBpmn(*policy)

	addHandlers(state)

	state.RunBpmn(flow)
	return state
}

func addHandlers(state *bpmn.State) {
	addSignHandlers(state)
	addPayHandlers(state)
}

func addSignHandlers(state *bpmn.State) {
	state.AddTaskHandler("setSign", setSign)
	state.AddTaskHandler("addContract", addContract)
	state.AddTaskHandler("sendMailContract", sendMailContract)
	state.AddTaskHandler("fillAttachments", fillAttachments)
	state.AddTaskHandler("setToPay", setToPay)
	state.AddTaskHandler("sendMailPay", sendMailPay)
}

func addPayHandlers(state *bpmn.State) {
	state.AddTaskHandler("updatePolicy", updatePolicy)
	state.AddTaskHandler("payTransaction", payTransaction)
}

func setSign(state *bpmn.State) error {
	log.AddPrefix("SetSign")
	defer log.PopPrefix()

	policy := state.Data
	err := plc.Sign(policy, origin)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func addContract(state *bpmn.State) error {
	policy := state.Data
	plc.AddSignedDocumentsInPolicy(policy, origin)

	return nil
}

func sendMailContract(state *bpmn.State) error {
	policy := state.Data
	log.AddPrefix("sendMailContract")
	defer log.PopPrefix()

	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail {
			toAddress = mail.GetContractorEmail(policy)
			ccAddress = mail.GetNetworkNodeEmail(networkNode)
		} else {
			toAddress = mail.GetNetworkNodeEmail(networkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		toAddress = mail.GetContractorEmail(policy)
	}

	log.Printf(
		"from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	return mail.SendMailContract(*policy, policy.Attachments, fromAddress, toAddress, ccAddress, flowName)
}

func fillAttachments(state *bpmn.State) error {
	policy := state.Data
	log.AddPrefix("FillAttachments")
	err := plc.FillAttachments(policy, origin)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func setToPay(state *bpmn.State) error {
	log.AddPrefix("setToPay")
	defer log.PopPrefix()
	policy := state.Data
	err := plc.SetToPay(policy, origin)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func sendMailPay(state *bpmn.State) error {
	policy := state.Data
	log.AddPrefix("sendMailPay")
	defer log.PopPrefix()
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail {
			toAddress = mail.GetContractorEmail(policy)
			ccAddress = mail.GetNetworkNodeEmail(networkNode)
		} else {
			toAddress = mail.GetNetworkNodeEmail(networkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		toAddress = mail.GetContractorEmail(policy)
	}

	log.Printf(
		"from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailPay(*policy, fromAddress, toAddress, ccAddress, flowName)

	return nil
}

func updatePolicy(state *bpmn.State) error {
	var err error
	policy := state.Data
	log.AddPrefix("updatePolicy")
	defer log.PopPrefix()

	if policy.IsPay || policy.Status != models.PolicyStatusToPay {
		log.Printf("policy already updated with isPay %t and status %s", policy.IsPay, policy.Status)
		return nil
	}

	// Add Policy contract
	err = plc.AddSignedDocumentsInPolicy(policy, origin)
	if err != nil {
		log.ErrorF("error AddContract %s", err.Error())
		return err
	}

	// promote documents from temp bucket to user and connect it to policy
	err = plc.SetUserIntoPolicyContractor(policy, origin)
	if err != nil {
		log.ErrorF("error SetUserIntoPolicyContractor %s", err.Error())
		return err
	}

	// Update Policy as paid
	err = plc.Pay(policy, origin)
	if err != nil {
		log.ErrorF("error Policy Pay %s", err.Error())
		return err
	}

	err = network.UpdateNetworkNodePortfolio(origin, policy, networkNode)
	if err != nil {
		log.ErrorF("error updating %s portfolio %s", networkNode.Type, err.Error())
		return err
	}

	policy.BigquerySave(origin)

	switch flowName {
	case models.ProviderMgaFlow:
		toAddress = mail.GetContractorEmail(policy)
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	case models.RemittanceMgaFlow:
		toAddress = mail.GetNetworkNodeEmail(networkNode)
	case models.MgaFlow, models.ECommerceFlow:
		toAddress = mail.GetContractorEmail(policy)
	}

	// Send mail with the contract to the user
	log.Printf(
		"from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	return mail.SendMailContract(*policy, policy.Attachments, fromAddress, toAddress, ccAddress, flowName)

}

func payTransaction(state *bpmn.State) error {
	log.AddPrefix("fabrickPayment")
	defer log.PopPrefix()
	policy := state.Data
	transaction, _ := tr.GetTransactionToBePaid(policy.Uid, providerId, trSchedule, lib.TransactionsCollection)
	err := tr.Pay(&transaction, origin, paymentMethod)
	if err != nil {
		log.ErrorF("error Transaction Pay %s", err.Error())
		return err
	}

	transaction.BigQuerySave(origin)

	return tr.CreateNetworkTransactions(policy, &transaction, networkNode, mgaProduct)
}
