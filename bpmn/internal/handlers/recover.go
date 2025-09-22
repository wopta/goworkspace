package handlers

import (
	"fmt"
	"os"
	"time"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func AddRecoverHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		builder.AddHandler("addNoteError", addNoteError),
		builder.AddHandler("sendEmailError", sendEmailError),
	)
}

var processNameToIta = map[bpmnEngine.BpmnFlow]string{
	bpmnEngine.Emit:            "Emissione",
	bpmnEngine.Proposal:        "Proposta",
	bpmnEngine.RequestApproval: "RequestApproval",
	bpmnEngine.Acceptance:      "Accettazione",
	bpmnEngine.Pay:             "Pagamento",
	bpmnEngine.Sign:            "Firma",
}

func addNoteError(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
	)
	if err != nil {
		return err
	}
	statusFlow, err := bpmnEngine.GetStatusFlow(state)
	if err != nil {
		return err
	}
	policy.AddSystemNote(models.GetErrorNote(processNameToIta[statusFlow.CurrentProcess]))
	return err
}

func sendEmailError(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
	)
	if err != nil {
		return err
	}
	statusFlow, err := bpmnEngine.GetStatusFlow(state)
	if err != nil {
		return err
	}
	var body string
	body += fmt.Sprintf("Nella data del %v la polizza %v ha avuto problemi nel processo di %v<br><br>", time.Now().Format("2006-01-02"), policy.CodeCompany, processNameToIta[statusFlow.CurrentProcess])
	body += fmt.Sprintf("Execution Id: %v", os.Getenv("GOOGLE_CLOUD_WORKFLOW_EXECUTION_ID"))
	mail.SendBaseEmail(body, fmt.Sprint("Errore nel processo di ", statusFlow.CurrentProcess), lib.GetMailProcessi(processNameToIta[statusFlow.CurrentProcess]))
	return nil
}
