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

func runEmitBpmn(policy *models.Policy, channel string) *bpmn.State {
	log.Printf("[runEmitBpmn] configuring flow for %s", channel)

	var (
		err           error
		setting       models.NodeSetting
		settingFormat string = "products/%s/setting.json"
	)

	settingFile := fmt.Sprintf(settingFormat, channel)
	settingByte := lib.GetFilesByEnv(settingFile)

	err = json.Unmarshal(settingByte, &setting)
	if err != nil {
		log.Printf("[runEmitBpmn] error unmarshaling setting file: %s", err.Error())
	}

	state := bpmn.NewBpmn(*policy)
	state.AddTaskHandler("emitData", emitData)
	state.AddTaskHandler("sendMailSign", sendMailSign)
	state.AddTaskHandler("sign", sign)
	state.AddTaskHandler("pay", pay)
	state.AddTaskHandler("setAdvice", setAdvanceBpm)
	state.AddTaskHandler("putUser", updateUserAndAgency)

	// TODO: use a map function to print only the name of each step
	flowBytes, _ := json.Marshal(setting.EmitFlow)
	log.Printf("[runEmitBpmn] starting emit flow: %s", string(flowBytes))

	state.RunBpmn(setting.EmitFlow)
	return state
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
