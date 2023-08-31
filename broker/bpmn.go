package broker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wopta/goworkspace/bpmn"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/user"
)

var origin string

const (
	emitFlowKey     = "emit"
	proposalFlowKey = "proposal"
)

func runBrokerBpmn(policy *models.Policy, flowKey string) *bpmn.State {
	log.Println("[runBrokerBpmn] configuring flow")

	var (
		err           error
		flow          []models.Process
		setting       models.NodeSetting
		settingFormat string = "products/%s/setting.json"
	)

	channel := models.GetChannel(policy)
	settingFile := fmt.Sprintf(settingFormat, channel)

	log.Printf("[runBrokerBpmn] loading file for channel %s", channel)
	settingByte := lib.GetFilesByEnv(settingFile)

	err = json.Unmarshal(settingByte, &setting)
	if err != nil {
		log.Printf("[runBrokerBpmn] error unmarshaling setting file: %s", err.Error())
	}

	state := bpmn.NewBpmn(*policy)

	switch flowKey {
	case proposalFlowKey:
		flow = setting.ProposalFlow
		addProposalHandlers(state)
	case emitFlowKey:
		flow = setting.EmitFlow
		addEmitHandlers(state)
	default:
		log.Println("[runBrokerBpmn] error flow not set")
		return nil
	}
	log.Printf("[runBrokerBpmn] using flow %s", flowKey)

	// TODO: use a map function to print only the name of each step
	flowBytes, _ := json.Marshal(flow)
	log.Printf("[runBrokerBpmn] starting %s flow: %s", flowKey, string(flowBytes))

	state.RunBpmn(flow)
	return state
}

func addEmitHandlers(state *bpmn.State) {
	state.AddTaskHandler("emitData", emitData)
	state.AddTaskHandler("sendMailSign", sendMailSign)
	state.AddTaskHandler("sign", sign)
	state.AddTaskHandler("pay", pay)
	state.AddTaskHandler("setAdvice", setAdvanceBpm)
	state.AddTaskHandler("putUser", updateUserAndAgency)
}

func addProposalHandlers(state *bpmn.State) {
	state.AddTaskHandler("setProposalData", setProposalBpm)
	state.AddTaskHandler("sendProposalMail", sendProposalMail)
}

func emitData(state *bpmn.State) error {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	p := state.Data
	emitBase(p, origin)
	return lib.SetFirestoreErr(firePolicy, p.Uid, p)
}

func setAdvanceBpm(state *bpmn.State) error {
	p := state.Data
	setAdvance(p, origin)
	return nil
}

func sendMailSign(state *bpmn.State) error {
	policy := state.Data
	mail.SendMailSign(*policy)
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

func updateUserAndAgency(state *bpmn.State) error {
	policy := state.Data
	user.SetUserIntoPolicyContractor(policy, origin)
	return models.UpdateAgencyPortfolio(policy, origin)
}

func setProposalBpm(state *bpmn.State) error {
	p := state.Data
	setProposalData(p)
	return nil
}

func sendProposalMail(state *bpmn.State) error {
	policy := state.Data
	mail.SendMailProposal(*policy)
	return nil
}
