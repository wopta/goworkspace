package broker

import (
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
	origin, paymentSplit, flowName    string
	ccAddress, toAddress, fromAddress mail.Address
	networkNode                       *models.NetworkNode
	product, mgaProduct               *models.Product
	warrant                           *models.Warrant
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

	// TODO: fix me - maybe get to/from/cc from flowFile.json?
	switch flowKey {
	case leadFlowKey:
		flow = flowFile.LeadFlow
		toAddress = mail.GetContractorEmail(policy)
		switch policy.Channel {
		case models.NetworkChannel:
			ccAddress = mail.GetNetworkNodeEmail(networkNode)
		}
	case proposalFlowKey:
		flow = flowFile.ProposalFlow
	case requestApprovalFlowKey:
		flow = flowFile.RequestApprovalFlow
		switch policy.Channel {
		case models.NetworkChannel:
			toAddress = mail.GetContractorEmail(policy)
			ccAddress = mail.GetNetworkNodeEmail(networkNode)
		case models.MgaChannel:
			toAddress = mail.GetContractorEmail(policy)
		}
	case emitFlowKey:
		flow = flowFile.EmitFlow
		switch policy.Channel {
		case models.NetworkChannel:
			switch networkNode.Type {
			case models.AgencyNetworkNodeType:
				toAddress = mail.GetContractorEmail(policy)
				ccAddress = mail.GetNetworkNodeEmail(networkNode)
			case models.AgentNetworkNodeType:
				toAddress = mail.GetNetworkNodeEmail(networkNode)
			}
		case models.MgaChannel, models.ECommerceChannel:
			toAddress = mail.GetContractorEmail(policy)
		}
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
	log.Printf(
		"[sendLeadMail] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailLead(*policy, fromAddress, toAddress, ccAddress, flowName)
	return nil
}

//	======================================
//	PROPOSAL FUNCTIONS
//	======================================

func addProposalHandlers(state *bpmn.State) {
	state.AddTaskHandler("setProposalData", setProposalBpm)
}

func setProposalBpm(state *bpmn.State) error {
	policy := state.Data

	setProposalData(policy)

	log.Printf("[setProposalData] saving proposal n. %d to firestore...", policy.ProposalNumber)

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
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
	mail.SendMailReserved(*policy, fromAddress, toAddress, ccAddress, flowName)
	return nil
}

//	======================================
//	EMIT FUNCTIONS
//	======================================

func addEmitHandlers(state *bpmn.State) {
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
