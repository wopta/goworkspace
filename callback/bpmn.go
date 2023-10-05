package callback

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/wopta/goworkspace/bpmn"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	tr "github.com/wopta/goworkspace/transaction"
)

var (
	origin, trSchedule, paymentMethod string
	ccAddress, toAddress, fromAddress mail.Address
)

const (
	signFlowKey = "sign"
	payFlowKey  = "pay"
)

func runCallbackBpmn(policy *models.Policy, flowKey string) *bpmn.State {
	log.Println("[runCallbackBpmn] configuring flow")

	var (
		err           error
		flow          []models.Process
		setting       models.NodeSetting
		settingFormat string = "products/%s/setting.json"
	)

	fromAddress = mail.AddressAnna
	channel := policy.Channel
	settingFile := fmt.Sprintf(settingFormat, channel)

	log.Printf("[runCallbackBpmn] loading file for channel %s", channel)
	settingByte := lib.GetFilesByEnv(settingFile)

	err = json.Unmarshal(settingByte, &setting)
	if err != nil {
		log.Printf("[runCallbackBpmn] error unmarshaling setting file: %s", err.Error())
		return nil
	}

	state := bpmn.NewBpmn(*policy)

	// TODO: fix me - maybe get to/from/cc from setting.json?
	switch flowKey {
	case signFlowKey:
		flow = setting.SignFlow
		switch channel {
		case models.AgencyChannel:
			toAddress = mail.GetContractorEmail(policy)
			ccAddress = mail.GetEmailByChannel(policy)
		default:
			toAddress = mail.GetEmailByChannel(policy)
			ccAddress = mail.Address{}
		}
	case payFlowKey:
		flow = setting.PayFlow
		switch channel {
		case models.AgentChannel:
			toAddress = mail.GetContractorEmail(policy)
			ccAddress = mail.GetEmailByChannel(policy)
		default:
			toAddress = mail.GetEmailByChannel(policy)
			ccAddress = mail.Address{}
		}
	default:
		log.Println("[runCallbackBpmn] error flow not set")
		return nil
	}

	addHandlers(state)

	flowHandlers := lib.SliceMap[models.Process, string](flow, func(h models.Process) string { return h.Name })
	log.Printf("[runCallbackBpmn] starting %s flow with set handlers: %s", flowKey, strings.Join(flowHandlers, ","))

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
	mail.SendMailContract(
		*policy,
		nil,
		fromAddress,
		toAddress,
		ccAddress,
	)

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
	mail.SendMailPay(
		*policy,
		fromAddress,
		toAddress,
		ccAddress,
	)

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

	// Update agency if present
	err = models.UpdateAgencyPortfolio(policy, origin)
	if err != nil && err.Error() != "agency not set" {
		log.Printf("[updatePolicy] ERROR updateAgencyPortfolio %s", err.Error())
		return err
	}

	// Update agent if present
	err = models.UpdateAgentPortfolio(policy, origin)
	if err != nil && err.Error() != "agent not set" {
		log.Printf("[updatePolicy] ERROR UpdateAgentPortfolio %s", err.Error())
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
	mail.SendMailContract(
		*policy,
		nil,
		fromAddress,
		toAddress,
		ccAddress,
	)

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

	return nil
}
