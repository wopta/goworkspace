package channelFlow

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
	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/callback_out/win"
	"gitlab.dev.wopta.it/goworkspace/lib"
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

func SaveAudit(st bpmn.StorageData) error {
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

	node, err := bpmn.GetData[*flow.Network]("networkNode", st)
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
