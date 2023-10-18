package callback

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
	tr "github.com/wopta/goworkspace/transaction"
)

var (
	origin, trSchedule, paymentMethod, flowName string
	ccAddress, toAddress, fromAddress           mail.Address
	networkNode                                 *models.NetworkNode
	mgaProduct                                  *models.Product
	warrant                                     *models.Warrant
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

	log.Println("[runCallbackBpmn] configuring flow --------------------------")

	toAddress = mail.Address{}
	ccAddress = mail.Address{}
	fromAddress = mail.AddressAnna

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	flowName, flowFile = policy.GetFlow(networkNode, warrant)
	if flowFile == nil {
		log.Println("[runCallbackBpmn] exiting bpmn - flowFile not loaded")
		return nil
	}
	log.Printf("[runCallbackBpmn] flowName '%s'", flowName)

	mgaProduct = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	// TODO: fix me - maybe get to/from/cc from flowFile.json?
	switch flowKey {
	case signFlowKey:
		flow = flowFile.SignFlow
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
	case payFlowKey:
		flow = flowFile.PayFlow
		switch policy.Channel {
		case models.NetworkChannel:
			switch networkNode.Type {
			case models.AgentNetworkNodeType:
				toAddress = mail.GetContractorEmail(policy)
				ccAddress = mail.GetNetworkNodeEmail(networkNode)
			case models.AgencyNetworkNodeType:
				toAddress = mail.GetNetworkNodeEmail(networkNode)
			}
		case models.MgaChannel, models.ECommerceChannel:
			toAddress = mail.GetContractorEmail(policy)
		}
	default:
		log.Println("[runCallbackBpmn] error flow not set")
		return nil
	}

	flowHandlers := lib.SliceMap[models.Process, string](flow, func(h models.Process) string { return h.Name })
	log.Printf("[runCallbackBpmn] starting %s flow with set handlers: %s", flowKey, strings.Join(flowHandlers, ","))

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
	policy := state.Data
	err := plc.Sign(policy, origin)
	if err != nil {
		log.Printf("[setSign] ERROR: %s", err.Error())
		return err
	}

	return nil
}

func addContract(state *bpmn.State) error {
	policy := state.Data
	plc.AddContract(policy, origin)

	return nil
}

func sendMailContract(state *bpmn.State) error {
	policy := state.Data
	log.Printf(
		"[sendMailContract] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailContract(*policy, nil, fromAddress, toAddress, ccAddress, flowName)

	return nil
}

func fillAttachments(state *bpmn.State) error {
	policy := state.Data
	err := plc.FillAttachments(policy, origin)
	if err != nil {
		log.Printf("[fillAttachments] ERROR: %s", err.Error())
		return err
	}

	return nil
}

func setToPay(state *bpmn.State) error {
	policy := state.Data
	err := plc.SetToPay(policy, origin)
	if err != nil {
		log.Printf("[setToPay] ERROR: %s", err.Error())
		return err
	}

	return nil
}

func sendMailPay(state *bpmn.State) error {
	policy := state.Data
	log.Printf(
		"[sendMailPay] from '%s', to '%s', cc '%s'",
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

	if policy.IsPay || policy.Status != models.PolicyStatusToPay {
		log.Printf("[updatePolicy] policy already updated with isPay %t and status %s", policy.IsPay, policy.Status)
		return nil
	}

	// promote documents from temp bucket to user and connect it to policy
	err = plc.SetUserIntoPolicyContractor(policy, origin)
	if err != nil {
		log.Printf("[updatePolicy] ERROR SetUserIntoPolicyContractor %s", err.Error())
		return err
	}

	// Add Policy contract
	err = plc.AddContract(policy, origin)
	if err != nil {
		log.Printf("[updatePolicy] ERROR AddContract %s", err.Error())
		return err
	}

	// Update Policy as paid
	err = plc.Pay(policy, origin)
	if err != nil {
		log.Printf("[updatePolicy] ERROR Policy Pay %s", err.Error())
		return err
	}

	err = network.UpdateNetworkNodePortfolio(origin, policy, networkNode)
	if err != nil {
		log.Printf("[updatePolicy] error updating %s portfolio %s", networkNode.Type, err.Error())
		return err
	}

	policy.BigquerySave(origin)

	// Send mail with the contract to the user
	log.Printf(
		"[updatePolicy] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailContract(*policy, nil, fromAddress, toAddress, ccAddress, flowName)

	return nil
}

func payTransaction(state *bpmn.State) error {
	policy := state.Data
	transaction, _ := tr.GetTransactionByPolicyUidAndScheduleDate(policy.Uid, trSchedule, origin)
	err := tr.Pay(&transaction, origin, paymentMethod)
	if err != nil {
		log.Printf("[fabrickPayment] ERROR Transaction Pay %s", err.Error())
		return err
	}

	transaction.BigQuerySave(origin)

	return tr.CreateNetworkTransactions(policy, &transaction, networkNode, mgaProduct)
}
