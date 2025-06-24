package broker

import (
	"fmt"
	"os"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/broker/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/bpmn"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

var (
	origin, flowName, paymentSplit, paymentMode string
	ccAddress, toAddress, fromAddress           mail.Address
	networkNode                                 *models.NetworkNode
	product, mgaProduct                         *models.Product
	warrant                                     *models.Warrant
	sendEmail                                   bool
)

const (
	emitFlowKey            = "emit"
	leadFlowKey            = "lead"
	proposalFlowKey        = "proposal"
	requestApprovalFlowKey = "requestApproval"
)

func runBrokerBpmn(policy *models.Policy, flowKey string) *bpmn.State {
	var (
		flow     []models.Process
		flowFile *models.NodeSetting
	)
	log.AddPrefix("RunBrokerBpmn")
	defer log.PopPrefix()
	log.Println("configuring flow ----------------------------")

	toAddress = mail.Address{}
	ccAddress = mail.Address{}
	fromAddress = mail.AddressAnna

	flowName, flowFile = policy.GetFlow(networkNode, warrant)
	if flowFile == nil {
		log.Println("exiting bpmn - flowFile not loaded")
		return nil
	}
	log.Printf("flowName '%s'", flowName)

	product = prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, networkNode, warrant)
	mgaProduct = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	switch flowKey {
	case leadFlowKey:
		flow = flowFile.LeadFlow
	case proposalFlowKey:
		flow = flowFile.ProposalFlow
	case requestApprovalFlowKey:
		flow = flowFile.RequestApprovalFlow
	case emitFlowKey:
		flow = flowFile.EmitFlow
	default:
		log.Println("error flow not set")
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
	addLeadHandlers(state)
	addProposalHandlers(state)
	addRequestApprovalHandlers(state)
	addEmitHandlers(state)
}

//	======================================
//	LEAD FUNCTIONS
//	======================================

func addLeadHandlers(state *bpmn.State) {
	state.AddTaskHandler("setLeadData", setLeadBpmn)
	state.AddTaskHandler("sendLeadMail", sendLeadMail)
}

func setLeadBpmn(state *bpmn.State) error {
	policy := state.Data
	return utility.SetLeadData(policy, *product, networkNode)
}

func sendLeadMail(state *bpmn.State) error {
	policy := state.Data
	log.AddPrefix("SendEmail")
	defer log.PopPrefix()
	toAddress = mail.GetContractorEmail(policy)
	ccAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	}

	log.Printf(
		"from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailLead(*policy, fromAddress, toAddress, ccAddress, flowName, []string{})
	return nil
}

//	======================================
//	PROPOSAL FUNCTIONS
//	======================================

func addProposalHandlers(state *bpmn.State) {
	state.AddTaskHandler("setProposalData", setProposalBpm)
	state.AddTaskHandler("sendProposalMail", sendProposalMail)
}

func setProposalBpm(state *bpmn.State) error {
	policy := state.Data
	log.AddPrefix("SetProposalData")
	defer log.PopPrefix()
	// TODO: remove when a proper solution to handle PMI is found
	if policy.Name == models.PmiProduct {
		return nil
	}

	if policy.ProposalNumber != 0 {
		log.Printf("policy '%s' already has proposal with number '%d'", policy.Uid, policy.ProposalNumber)
		return nil
	}

	utility.SetProposalData(policy, origin, networkNode, mgaProduct)

	log.Printf("saving proposal n. %d to firestore...", policy.ProposalNumber)

	firePolicy := lib.PolicyCollection
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendProposalMail(state *bpmn.State) error {
	policy := state.Data
	log.AddPrefix("SendProposalEmail")
	defer log.PopPrefix()
	// TODO: remove when a proper solution to handle PMI is found
	if policy.Name == models.PmiProduct {
		return nil
	}

	if !sendEmail || policy.IsReserved {
		return nil
	}

	toAddress = mail.GetContractorEmail(policy)
	ccAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	}

	log.Printf(
		"from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)

	mail.SendMailProposal(*policy, fromAddress, toAddress, ccAddress, flowName, []string{models.ProposalAttachmentName})
	return nil
}

//	======================================
//	REQUEST APPROVAL FUNCTIONS
//	======================================

func addRequestApprovalHandlers(state *bpmn.State) {
	state.AddTaskHandler("setRequestApprovalData", setRequestApprovalBpmn)
	state.AddTaskHandler("sendRequestApprovalMail", sendRequestApprovalMail)
}

func setRequestApprovalBpmn(state *bpmn.State) error {
	policy := state.Data
	log.AddPrefix("SendRequestApproval")
	defer log.PopPrefix()
	firePolicy := lib.PolicyCollection

	utility.SetRequestApprovalData(policy, networkNode, mgaProduct, origin)

	log.Printf("saving policy with uid %s to Firestore....", policy.Uid)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendRequestApprovalMail(state *bpmn.State) error {
	policy := state.Data

	if policy.Status == models.PolicyStatusWaitForApprovalMga {
		return nil
	}

	toAddress = mail.GetContractorEmail(policy)
	ccAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	case models.ECommerceChannel:
		toAddress = mail.Address{} // fail safe for not sending email on ecommerce reserved
	}

	if policy.Name == "qbe" {
		toAddress = mail.Address{
			Address: os.Getenv("QBE_RESERVED_MAIL"),
		}
		ccAddress = mail.Address{}
	}

	mail.SendMailReserved(*policy, fromAddress, toAddress, ccAddress, flowName,
		[]string{models.ProposalAttachmentName})
	return nil
}

//	======================================
//	EMIT FUNCTIONS
//	======================================

func addEmitHandlers(state *bpmn.State) {
	state.AddTaskHandler("setProposalData", setProposalBpm)
	state.AddTaskHandler("emitData", emitData)
	state.AddTaskHandler("sendMailSign", sendMailSign)
	state.AddTaskHandler("sign", sign)
	state.AddTaskHandler("pay", pay)
	state.AddTaskHandler("setAdvice", setAdvanceBpm)
	state.AddTaskHandler("putUser", updateUserAndNetworkNode)
	state.AddTaskHandler("sendEmitProposalMail", sendEmitProposalMail)
}

func emitData(state *bpmn.State) error {
	firePolicy := lib.PolicyCollection
	policy := state.Data
	emitBase(policy, origin)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendMailSign(state *bpmn.State) error {
	policy := state.Data
	log.AddPrefix("SendEmailSign")
	defer log.PopPrefix()

	ccAddress = mail.Address{}
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
	mail.SendMailSign(*policy, fromAddress, toAddress, ccAddress, flowName)
	return nil
}

func sign(state *bpmn.State) error {
	policy := state.Data
	err := utility.SignFiles(policy, product, networkNode, true, origin)
	return err
}

func pay(state *bpmn.State) error {
	policy := state.Data
	utility.EmitPay(policy, origin, product, mgaProduct, networkNode)
	if policy.PayUrl == "" {
		return fmt.Errorf("missing payment url")
	}
	return nil
}

func setAdvanceBpm(state *bpmn.State) error {
	policy := state.Data
	utility.SetAdvance(policy, origin, product, mgaProduct, networkNode, paymentSplit, paymentMode)
	return nil
}

func updateUserAndNetworkNode(state *bpmn.State) error {
	policy := state.Data
	log.AddPrefix("PutUser")
	defer log.PopPrefix()
	// promote documents from temp bucket to user and connect it to policy
	err := plc.SetUserIntoPolicyContractor(policy, origin)
	if err != nil {
		log.ErrorF("error SetUserIntoPolicyContractor %s", err.Error())
		return err
	}
	return network.UpdateNetworkNodePortfolio(origin, policy, networkNode)
}

func sendEmitProposalMail(state *bpmn.State) error {
	policy := state.Data
	log.AddPrefix("sendEmitProposalMail")
	defer log.PopPrefix()
	// TODO: remove when a proper solution to handle PMI is found
	if policy.Name == models.PmiProduct {
		return nil
	}

	if policy.IsReserved {
		return nil
	}

	toAddress = mail.GetContractorEmail(policy)
	ccAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	}

	log.Printf(
		fromAddress.String(),
		"from '%s', to '%s', cc '%s'",
		toAddress.String(),
		ccAddress.String(),
	)

	mail.SendMailProposal(*policy, fromAddress, toAddress, ccAddress, flowName, []string{models.ProposalAttachmentName})
	return nil
}
