package channelFlow

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/callback_out/win"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func CallBackEmit(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.Network]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.Policy]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.NetworkNode.ExternalNetworkCode)
	info := win.Emit(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackEmitRemittance(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.Network]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.Policy]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.NetworkNode.ExternalNetworkCode)

	info := win.Emit(*policy.Policy)
	if err = saveAudit(node.NetworkNode, info); err != nil {
		return err
	}
	info = win.Paid(*policy.Policy)
	if err = saveAudit(node.NetworkNode, info); err != nil {
		return err
	}

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}
func CallBackProposal(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.Network]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.Policy]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	info := win.Proposal(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackPaid(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.Network]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.Policy]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	info := win.Paid(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackRequestApproval(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.Network]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.Policy]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	info := win.RequestApproval(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackApproved(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.Network]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.Policy]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	info := win.Approved(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackRejected(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.Network]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.Policy]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	info := win.Rejected(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func CallBackSigned(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.Network]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.Policy]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	info := win.Signed(*policy.Policy)

	st.AddLocal("callbackInfo", &flow.CallbackInfo{CallbackInfo: info})
	return nil
}

func fromCurrentProcessToCallbackoutAction(currentProcess string) (base.CallbackoutAction, error) {
	//EmitRemittance is done directly on CallBackEmitRemittance
	var callbackActions = map[string]base.CallbackoutAction{
		"emitCallBack":            base.Emit,
		"payCallBack":             base.Paid,
		"proposalCallback":        base.Proposal,
		"requestApprovalCallBack": base.RequestApproval,
		"signCallback":            base.Signed,
		"approvedCallBack":        base.Approved,
		"rejectedCallback":        base.Rejected,
	}
	callbackAction, ok := callbackActions[currentProcess]
	if !ok {
		return "", fmt.Errorf("No callback action for process '%s'", currentProcess)
	}
	return callbackAction, nil
}

func BaseRequest(store bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.Policy]("policy", store)
	if e != nil {
		return e

	}
	network := "facile_broker"
	basePath := os.Getenv(fmt.Sprintf("%s_CALLBACK_ENDPOINT", lib.ToUpper(network)))
	if basePath == "" {
		return errors.New("no base path for callback founded")

	}

	rawBody, err := json.Marshal(policy.Policy)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, basePath, bytes.NewReader(rawBody))
	if err != nil {
		return err
	}

	req.SetBasicAuth(
		os.Getenv(fmt.Sprintf("%s_CALLBACK_AUTH_USER", network)),
		os.Getenv(fmt.Sprintf("%s_CALLBACK_AUTH_PASS", network)))
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	info := flow.CallbackInfo{}
	status, err := bpmn.GetStatusFlow(store)
	if err != nil {
		return err
	}
	callbackAction, err := fromCurrentProcessToCallbackoutAction(status.CurrentProcess)
	if err != nil {
		return err
	}
	info.FromRequestResponse(callbackAction, res, req)
	store.AddLocal("callbackInfo", &info)
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
