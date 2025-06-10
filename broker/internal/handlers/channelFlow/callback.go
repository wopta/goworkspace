package channelFlow

import (
	"time"

	"cloud.google.com/go/bigquery"
	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func CallBackEmit(st bpmn.StorageData) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, st),
		bpmn.GetDataRef("networkNode", &networkNode, st),
		bpmn.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Emit(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackEmitRemittance(st bpmn.StorageData) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, st),
		bpmn.GetDataRef("networkNode", &networkNode, st),
		bpmn.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Emit(*policy.Policy)
	if err = saveAudit(networkNode.NetworkNode, info); err != nil {
		return err
	}
	info = client.Paid(*policy.Policy)
	if err = saveAudit(networkNode.NetworkNode, info); err != nil {
		return err
	}

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackProposal(st bpmn.StorageData) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, st),
		bpmn.GetDataRef("networkNode", &networkNode, st),
		bpmn.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Proposal(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackPaid(st bpmn.StorageData) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, st),
		bpmn.GetDataRef("networkNode", &networkNode, st),
		bpmn.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Paid(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackRequestApproval(st bpmn.StorageData) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, st),
		bpmn.GetDataRef("networkNode", &networkNode, st),
		bpmn.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.RequestApproval(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackApproved(st bpmn.StorageData) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, st),
		bpmn.GetDataRef("networkNode", &networkNode, st),
		bpmn.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Approved(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackRejected(st bpmn.StorageData) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, st),
		bpmn.GetDataRef("networkNode", &networkNode, st),
		bpmn.GetDataRef("clientCallback", &client, st),
	)
	if err != nil {
		return err
	}

	info := client.Rejected(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackSigned(st bpmn.StorageData) error {
	var networkNode *flow.Network
	var policy *flow.Policy
	var client *flow.ClientCallback

	var err = bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, st),
		bpmn.GetDataRef("networkNode", &networkNode, st),
		bpmn.GetDataRef("clientCallback", &client, st),
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

func SaveAudit(st bpmn.StorageData) error {

	node, err := bpmn.GetData[*flow.Network]("networkNode", st)
	if err != nil {
		return err
	}
	res, err := bpmn.GetData[*flow.CallbackInfo]("callbackInfo", st)
	if err != nil {
		return err
	}

	return saveAudit(node.NetworkNode, res.CallbackInfo)
}
func saveAudit(node *models.NetworkNode, callbackInfo base.CallbackInfo) error {
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
