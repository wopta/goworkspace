package broker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/callback_out/win"
	"github.com/wopta/goworkspace/lib"
)

func getNodeFlow() (*bpmn.BpnmBuilder, error) {
	store := bpmn.NewStorageBpnm()
	builder, e := bpmn.NewBpnmBuilder("flows/draft/node_flows.json")
	if e != nil {
		return nil, e
	}
	//hard coded, need to be on json
	callback := flow.CallbackConfig{
		Proposal:        true,
		RequestApproval: true,
		Emit:            true,
		Pay:             true,
		Sign:            true,
		Approved:        true,
		Rejected:        true,
	}
	if e := store.AddLocal("config", &callback); e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	err := bpmn.IsError(
		builder.AddHandler("baseCallback", baseRequest),
		builder.AddHandler("winEmit", callBackEmit),
		builder.AddHandler("winSign", callBackSigned),
		builder.AddHandler("saveAudit", saveAudit),
		builder.AddHandler("winPay", callBackPaid),
		builder.AddHandler("winProposal", callBackProposal),
		builder.AddHandler("winRequestApproval", callBackRequestApproval),
		builder.AddHandler("winApproved", callBackApproved),
		builder.AddHandler("winRejected", callBackRejected),
	)
	if err != nil {
		return nil, err
	}
	return builder, nil
}

func callBackEmit(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.NetworkDraft]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	_info := win.Emit(*policy.Policy)

	info := flow.CallbackInfo{
		Request:     _info.Request,
		RequestBody: _info.RequestBody,
		Response:    _info.Response,
		Error:       _info.Error,
	}
	st.AddLocal("callbackInfo", &info)
	return nil
}

func callBackProposal(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.NetworkDraft]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	_info := win.Proposal(*policy.Policy)

	info := flow.CallbackInfo{
		Request:     _info.Request,
		RequestBody: _info.RequestBody,
		Response:    _info.Response,
		Error:       _info.Error,
	}
	st.AddLocal("callbackInfo", &info)
	return nil
}

func callBackPaid(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.NetworkDraft]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	_info := win.Paid(*policy.Policy)

	info := flow.CallbackInfo{
		Request:     _info.Request,
		RequestBody: _info.RequestBody,
		Response:    _info.Response,
		Error:       _info.Error,
	}
	st.AddLocal("callbackInfo", &info)
	return nil
}

func callBackRequestApproval(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.NetworkDraft]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	_info := win.RequestApproval(*policy.Policy)

	info := flow.CallbackInfo{
		Request:     _info.Request,
		RequestBody: _info.RequestBody,
		Response:    _info.Response,
		Error:       _info.Error,
	}
	st.AddLocal("callbackInfo", &info)
	return nil
}

func callBackApproved(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.NetworkDraft]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	_info := win.Approved(*policy.Policy)

	info := flow.CallbackInfo{
		Request:     _info.Request,
		RequestBody: _info.RequestBody,
		Response:    _info.Response,
		Error:       _info.Error,
	}
	st.AddLocal("callbackInfo", &info)
	return nil
}

func callBackRejected(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.NetworkDraft]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	_info := win.Rejected(*policy.Policy)

	info := flow.CallbackInfo{
		Request:     _info.Request,
		RequestBody: _info.RequestBody,
		Response:    _info.Response,
		Error:       _info.Error,
	}
	st.AddLocal("callbackInfo", &info)
	return nil
}

func callBackSigned(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.NetworkDraft]("networkNode", st)
	if err != nil {
		return err
	}
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", st)
	if err != nil {
		return err
	}
	win := win.NewClient(node.ExternalNetworkCode)
	_info := win.Signed(*policy.Policy)

	info := flow.CallbackInfo{
		Request:     _info.Request,
		RequestBody: _info.RequestBody,
		Response:    _info.Response,
		Error:       _info.Error,
	}
	st.AddLocal("callbackInfo", &info)
	return nil
}

func baseRequest(store bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", store)
	if e != nil {
		return e

	}
	network := "facile_broker"
	basePath := os.Getenv(fmt.Sprintf("%s_CALLBACK_ENDPOINT", lib.ToUpper(network)))
	if basePath == "" {
		return errors.New("no base path for callback founded")

	}

	rawBody, err := json.Marshal(policy)
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

	info := flow.CallbackInfo{
		Request:     req,
		RequestBody: rawBody,
		Response:    res,
		Error:       err,
	}
	store.AddLocal("callbackInfo", &info)
	return nil
}

func saveAudit(st bpmn.StorageData) error {
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

	node, err := bpmn.GetData[*flow.NetworkDraft]("networkNode", st)
	if err != nil {
		return err
	}
	res, err := bpmn.GetData[*flow.CallbackInfo]("callbackInfo", st)
	if err != nil {
		return err
	}
	var (
		audit   auditSchema
		resBody []byte
	)

	audit.CreationDate = lib.GetBigQueryNullDateTime(time.Now().UTC())
	audit.Client = node.CallbackConfig.Name
	audit.NodeUid = node.Uid
	audit.Action = res.Action

	audit.ReqBody = string(res.RequestBody)
	if res.Request != nil {
		audit.ReqMethod = res.Request.Method
		audit.ReqPath = res.Request.Host + res.Request.URL.RequestURI()
	}

	if res.Response != nil {
		resBody, _ = io.ReadAll(res.Response.Body)
		defer res.Response.Body.Close()
		audit.ResStatusCode = res.Response.StatusCode
		audit.ResBody = string(resBody)
	}

	if res.Error != nil {
		audit.Error = res.Error.Error()
	}

	const CallbackOutTableId string = "callback-out"
	if err := lib.InsertRowsBigQuery(lib.WoptaDataset, CallbackOutTableId, audit); err != nil {
		return err
	}
	return nil
}
