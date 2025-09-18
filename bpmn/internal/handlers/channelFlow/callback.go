package channelFlow

import (
	"time"

	"cloud.google.com/go/bigquery"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func CallBackEmit(st *bpmnEngine.StorageBpnm) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, st),
		bpmnEngine.GetDataRef("networkNode", &networkNode, st),
		bpmnEngine.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Emit(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackProposal(st *bpmnEngine.StorageBpnm) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, st),
		bpmnEngine.GetDataRef("networkNode", &networkNode, st),
		bpmnEngine.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Proposal(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackPaid(st *bpmnEngine.StorageBpnm) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, st),
		bpmnEngine.GetDataRef("networkNode", &networkNode, st),
		bpmnEngine.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Paid(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackRequestApproval(st *bpmnEngine.StorageBpnm) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, st),
		bpmnEngine.GetDataRef("networkNode", &networkNode, st),
		bpmnEngine.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.RequestApproval(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackApproved(st *bpmnEngine.StorageBpnm) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, st),
		bpmnEngine.GetDataRef("networkNode", &networkNode, st),
		bpmnEngine.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Approved(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackRejected(st *bpmnEngine.StorageBpnm) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, st),
		bpmnEngine.GetDataRef("networkNode", &networkNode, st),
		bpmnEngine.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Rejected(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackSigned(st *bpmnEngine.StorageBpnm) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, st),
		bpmnEngine.GetDataRef("networkNode", &networkNode, st),
		bpmnEngine.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Signed(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

type auditSchema struct {
	CreationDate  bigquery.NullDateTime `bigquery:"creationDate"`
	Client        string                `bigquery:"client"`
	NodeUid       string                `bigquery:"nodeUid"`
	Action        string                `bigquery:"action"`
	ReqMethod     string                `bigquery:"reqMethod"`
	ReqPath       string                `bigquery:"reqPath"`
	ReqBody       string                `bigquery:"reqBody"`
	ResStatusCode int                   `bigquery:"resStatusCode"`
	ResBody       string                `bigquery:"resBody"`
	Error         string                `bigquery:"error"`
}

func SaveAudit(st *bpmnEngine.StorageBpnm) error {

	policy, err := bpmnEngine.GetData[*flow.Policy]("policy", st)
	if err != nil {
		return err
	}
	node, err := bpmnEngine.GetData[*flow.Network]("networkNode", st)
	if err != nil {
		return err
	}
	res, err := bpmnEngine.GetData[*flow.CallbackInfo]("callbackInfo", st)
	if err != nil {
		return err
	}

	return saveAudit(node.NetworkNode, res.CallbackInfo, policy.Policy)
}
func saveAudit(node *models.NetworkNode, callbackInfo base.CallbackInfo, policy *models.Policy) error {
	var audit auditSchema
	audit.CreationDate = lib.GetBigQueryNullDateTime(time.Now().UTC())
	audit.Client = node.CallbackConfig.Name
	audit.NodeUid = node.Uid
	audit.Action = string(callbackInfo.ResAction)

	audit.ReqMethod = callbackInfo.ReqMethod
	audit.ReqPath = callbackInfo.ReqPath

	audit.ResStatusCode = callbackInfo.ResStatusCode
	audit.ResBody = string(callbackInfo.ResBody)
	if callbackInfo.Error != nil {
		audit.Error = callbackInfo.Error.Error()
	}

	log.PrintStruct("SaveAudit:", audit)
	const CallbackOutTableId string = "callback-out"
	if err := lib.InsertRowsBigQuery(lib.WoptaDataset, CallbackOutTableId, audit); err != nil {
		return err
	}
	if callbackInfo.ResStatusCode == 200 {
		policy.AddSystemNote(func(p *models.Policy) models.PolicyNote {
			return models.PolicyNote{
				Text: "Callback " + audit.Client + " eseguita correttamente",
			}
		})
	}
	return nil
}
