package broker

import (
	"fmt"
	"log"
	"strings"

	"github.com/wopta/goworkspace/bpmn"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	prd "github.com/wopta/goworkspace/product"
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

	log.Println("[runBrokerBpmn] configuring flow ----------------------------")

	toAddress = mail.Address{}
	ccAddress = mail.Address{}
	fromAddress = mail.AddressAnna

	flowName, flowFile = policy.GetFlow(networkNode, warrant)
	if flowFile == nil {
		log.Println("[runBrokerBpmn] exiting bpmn - flowFile not loaded")
		return nil
	}
	log.Printf("[runBrokerBpmn] flowName '%s'", flowName)

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
		log.Println("[runBrokerBpmn] error flow not set")
		return nil
	}

	flowHandlers := lib.SliceMap[models.Process, string](flow, func(h models.Process) string { return h.Name })
	log.Printf("[runBrokerBpmn] starting %s flow with set handlers: %s", flowKey, strings.Join(flowHandlers, ","))

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
	setLeadData(policy)
	return nil
}

func sendLeadMail(state *bpmn.State) error {
	policy := state.Data

	toAddress = mail.GetContractorEmail(policy)
	ccAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	}

	log.Printf(
		"[sendLeadMail] from '%s', to '%s', cc '%s'",
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

	if policy.ProposalNumber != 0 {
		log.Printf("[setProposalData] policy '%s' already has proposal with number '%d'", policy.Uid, policy.ProposalNumber)
		return nil
	}

	setProposalData(policy)

	log.Printf("[setProposalData] saving proposal n. %d to firestore...", policy.ProposalNumber)

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendProposalMail(state *bpmn.State) error {
	policy := state.Data

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
		"[sendProposalMail] from '%s', to '%s', cc '%s'",
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
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	setRequestApprovalData(policy)

	log.Printf("[setRequestApproval] saving policy with uid %s to Firestore....", policy.Uid)
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
}

func emitData(state *bpmn.State) error {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	policy := state.Data
	emitBase(policy, origin)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendMailSign(state *bpmn.State) error {
	policy := state.Data

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
		"[sendMailSign] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailSign(*policy, fromAddress, toAddress, ccAddress, flowName)
	return nil
}

func sign(state *bpmn.State) error {
	policy := state.Data
	emitSign(policy, origin)
	return nil
}

func pay(state *bpmn.State) error {
	policy := state.Data
	emitPay(policy, origin)
	if policy.PayUrl == "" {
		return fmt.Errorf("missing payment url")
	}
	return nil
}

func setAdvanceBpm(state *bpmn.State) error {
	policy := state.Data
	setAdvance(policy, origin)
	return nil
}

func updateUserAndNetworkNode(state *bpmn.State) error {
	policy := state.Data
	// promote documents from temp bucket to user and connect it to policy
	err := plc.SetUserIntoPolicyContractor(policy, origin)
	if err != nil {
		log.Printf("[putUser] ERROR SetUserIntoPolicyContractor %s", err.Error())
		return err
	}
	return network.UpdateNetworkNodePortfolio(origin, policy, networkNode)
}
