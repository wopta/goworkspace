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

func CallBackEmit(st bpmnEngine.StorageData) error {
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

func CallBackProposal(st bpmnEngine.StorageData) error {
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
	return saveAudit(networkNode.NetworkNode, info)
}

func CallBackPaid(st bpmnEngine.StorageData) error {
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

	return saveAudit(networkNode.NetworkNode, info)
}

func CallBackRequestApproval(st bpmnEngine.StorageData) error {
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

	return saveAudit(networkNode.NetworkNode, info)
}

func CallBackApproved(st bpmnEngine.StorageData) error {
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

	return saveAudit(networkNode.NetworkNode, info)
}

func CallBackRejected(st bpmnEngine.StorageData) error {
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

	return saveAudit(networkNode.NetworkNode, info)
}

func CallBackSigned(st bpmnEngine.StorageData) error {
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

	return saveAudit(networkNode.NetworkNode, info)
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

func saveAudit(node *models.NetworkNode, callbackInfo base.CallbackInfo) error {
	log.Println("Saving audit")
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

	const CallbackOutTableId string = "callback-out"
	if err := lib.InsertRowsBigQuery(lib.WoptaDataset, CallbackOutTableId, audit); err != nil {
		return err
	}
	return nil
}
