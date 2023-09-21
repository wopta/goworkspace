package broker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/wopta/goworkspace/bpmn"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/user"
)

var (
	origin      string
	ccAddress   mail.Address
	toAddress   mail.Address
	fromAddress mail.Address
)

const (
	emitFlowKey     = "emit"
	leadFlowKey     = "lead"
	proposalFlowKey = "proposal"
	requestApprovalFlowKey = "requestApproval"
)

func runBrokerBpmn(policy *models.Policy, flowKey string) *bpmn.State {
	log.Println("[runBrokerBpmn] configuring flow")

	var (
		err           error
		flow          []models.Process
		setting       models.NodeSetting
		settingFormat string = "products/%s/setting.json"
	)

	fromAddress = mail.AddressAnna
	channel := models.GetChannel(policy)
	settingFile := fmt.Sprintf(settingFormat, channel)

	log.Printf("[runBrokerBpmn] loading file for channel %s", channel)
	settingByte := lib.GetFilesByEnv(settingFile)

	err = json.Unmarshal(settingByte, &setting)
	if err != nil {
		log.Printf("[runBrokerBpmn] error unmarshaling setting file: %s", err.Error())
		return nil
	}

	state := bpmn.NewBpmn(*policy)

	// TODO: fix me - maybe get to/from/cc from setting.json?
	switch flowKey {
	case leadFlowKey:
		flow = setting.LeadFlow
		toAddress = mail.GetContractorEmail(policy)
		switch channel {
		case models.AgentChannel:
			ccAddress = mail.GetAgentEmail(policy)
		default:
			ccAddress = mail.Address{}
		}
	case proposalFlowKey:
		flow = setting.ProposalFlow
		toAddress = mail.GetContractorEmail(policy)
		switch channel {
		case models.AgentChannel:
			ccAddress = mail.GetAgentEmail(policy)
		default:
			ccAddress = mail.Address{}
		}
	case requestApprovalFlowKey:
		flow = setting.RequestApprovalFlow
		switch channel {
		case models.AgentChannel:
			toAddress = mail.GetContractorEmail(policy)
			ccAddress = mail.GetAgentEmail(policy)
		default:
			toAddress = mail.Address{}
			ccAddress = mail.Address{}
		}
	case emitFlowKey:
		flow = setting.EmitFlow
		switch channel {
		case models.AgencyChannel:
			toAddress = mail.GetContractorEmail(policy)
			ccAddress = mail.GetEmailByChannel(policy)
		default:
			toAddress = mail.GetEmailByChannel(policy)
			ccAddress = mail.Address{}
		}
	default:
		log.Println("[runBrokerBpmn] error flow not set")
		return nil
	}

	addHandlers(state)

	flowHandlers := lib.SliceMap[models.Process, string](flow, func(h models.Process) string { return h.Name })
	log.Printf("[runBrokerBpmn] starting %s flow with set handlers: %s", flowKey, strings.Join(flowHandlers, ","))

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
		"[sendProposalMail] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailLead(
		*policy,
		fromAddress,
		toAddress,
		ccAddress,
	)
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
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	setProposalData(policy)

	log.Printf("[setProposalData] saving proposal n. %d to firestore...", policy.ProposalNumber)
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

	for index, doc := range policy.ReservedInfo.Documents {
		
		components := strings.Split(doc.Link, "://")

		// Check if there are two components (scheme and path)
		if len(components) == 2 {
			scheme := components[0]
			path := components[1]
	
			// Split the path into bucket and object name
			pathComponents := strings.SplitN(path, "/", 2)
			if len(pathComponents) == 2 {
				bucketName := pathComponents[0]
				objectName := pathComponents[1]
	
				fmt.Println("Scheme:", scheme)
				fmt.Println("Bucket Name:", bucketName)
				fmt.Println("Object Name:", objectName)

				policy.ReservedInfo.Documents[index].Byte = base64.StdEncoding.EncodeToString(lib.GetFromStorage(bucketName, objectName, ""))
			}
		}
	}

	mail.SendMailReserved(*policy, fromAddress, toAddress, ccAddress)
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
	state.AddTaskHandler("putUser", updateUserAndAgency)
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
	mail.SendMailSign(
		*policy,
		fromAddress,
		toAddress,
		ccAddress,
	)
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

func updateUserAndAgency(state *bpmn.State) error {
	policy := state.Data
	user.SetUserIntoPolicyContractor(policy, origin)
	return models.UpdateAgencyPortfolio(policy, origin)
}