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
	bpmn "github.com/wopta/goworkspace/broker/draftBpnm"
	"github.com/wopta/goworkspace/broker/draftBpnm/flow"
	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	prd "github.com/wopta/goworkspace/product"
	tr "github.com/wopta/goworkspace/transaction"
)

func getFlow(policy *models.Policy, networkNode *models.NetworkNode, storage bpmn.StorageData) (*bpmn.FlowBpnm, error) {
	builder, err := bpmn.NewBpnmBuilder("broker/draftBpnm/flow/channel_flows.json")
	if err != nil {
		return nil, err
	}
	err = addHandlersDraft(builder)
	if err != nil {
		return nil, err
	}
	product = prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, networkNode, warrant)
	flowName, _ = policy.GetFlow(networkNode, warrant)

	//	mgaProduct = product
	//	toAddress = mail.Address{}
	//	ccAddress = mail.Address{}
	//	fromAddress = mail.AddressAnna

	productDraft := flow.ProductDraft{
		product,
	}
	policyDraft := flow.PolicyDraft{
		policy,
	}
	networkDraft := flow.NetworkDraft{
		networkNode,
	}
	if product.Flow == "" {
		product.Flow = policy.Channel
	}
	storage.AddGlobal("policy", &policyDraft)
	storage.AddGlobal("product", &productDraft)
	storage.AddGlobal("node", &networkDraft)
	builder.SetStorage(storage)

	if networkNode == nil || networkNode.CallbackConfig == nil {
		log.Println("no node or callback config available")
		return builder.Build()
	}
	injected, err := getNodeFlow()
	if err != nil {
		return nil, err
	}
	err = builder.Inject(injected)
	if err != nil {
		return nil, err
	}
	return builder.Build()
}

func getNodeFlow() (*bpmn.BpnmBuilder, error) {
	store := bpmn.NewStorageBpnm()
	builder, e := bpmn.NewBpnmBuilder("broker/draftBpnm/flow/node_flows.json")
	if e != nil {
		return nil, e
	}
	callback := flow.CallbackConfig{
		Proposal:        true,
		RequestApproval: true,
		Emit:            true,
		Pay:             true,
		Sign:            true,
	}
	if e := store.AddLocal("config", &callback); e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	err := bpmn.IsError(
		builder.AddHandler("baseCallback", baseRequest),
		builder.AddHandler("winEmit", func(st bpmn.StorageData) error {
			return baseRequest(st)
		}),
		builder.AddHandler("winLead", func(st bpmn.StorageData) error {
			return baseRequest(st)
		}),

		builder.AddHandler("winPay", func(st bpmn.StorageData) error {
			return baseRequest(st)
		}), builder.AddHandler("winProposal", func(st bpmn.StorageData) error {
			return baseRequest(st)
		}),
		builder.AddHandler("winRequestApproval", func(st bpmn.StorageData) error {
			return baseRequest(st)
		}), builder.AddHandler("winSign", func(st bpmn.StorageData) error {
			return baseRequest(st)
		}),
		builder.AddHandler("saveAudit", saveAudit),
	)
	if err != nil {
		return nil, err
	}
	return builder, nil
}

//func inspratica(policy models.Policy, state, producer string) ([]byte, error) {
//	winPolicy := policyDto(policy, producer)
//	payload := inspraticaReq{
//		Anagrafica:   winPolicy.Anagrafica,
//		Garanzie:     winPolicy.Garanzie,
//		IdPratica:    winPolicy.IdPratica,
//		PerAss:       winPolicy.PerAss,
//		Prodotto:     winPolicy.Prodotto,
//		TotaleAnnuo:  winPolicy.TotaleAnnuo,
//		TotaleFirma:  winPolicy.TotaleFirma,
//		TotaleFutura: winPolicy.TotaleFutura,
//		Utente:       winPolicy.Utente,
//		StatoPratica: state,
//	}
//
//	return json.Marshal(payload)
//}
//func winEmit(st bpmn.StorageData) flow.CallbackInfo {
//	path := "restba/extquote/inspratica"
//	node, err := bpmn.GetData[*flow.NetworkDraft]("node", st)
//	body, err := inspratica(policy, "QUOTAZIONE_ACCETTATA", node.ExternalNetworkCode)
//	if err != nil {
//		return internal.CallbackInfo{
//			Request:     nil,
//			RequestBody: nil,
//			Response:    nil,
//			Error:       err,
//		}
//	}
//
//	req, res, err := c.post(bytes.NewReader(body))
//	return internal.CallbackInfo{
//		Request:     req,
//		RequestBody: body,
//		Response:    res,
//		Error:       err,
//	}
//}

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

func saveAudit(st bpmn.StorageData) error {
	node, err := bpmn.GetData[*flow.NetworkDraft]("node", st)
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

func addHandlersDraft(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("setProposalData", setProposalDataDraft),
		builder.AddHandler("emitData", emitDataDraft),
		builder.AddHandler("sendMailSign", sendMailSignDraft),
		builder.AddHandler("pay", payDraft),
		builder.AddHandler("setAdvice", setAdvanceDraft),
		builder.AddHandler("putUser", updateUserAndNetworkNodeDraft),
		builder.AddHandler("sendEmitProposalMail", sendEmitProposalMailDraft),
		builder.AddHandler("setLeadData", setLeadBpmnDraft),
		builder.AddHandler("sendLeadMail", sendLeadMailDraft),
		builder.AddHandler("updatePolicy", updatePolicyDraft),
		builder.AddHandler("sign", signDraft),
		builder.AddHandler("payTransaction", payTransactionDraft),
		builder.AddHandler("sendProposalMail", sendProposalMailDraft),
		builder.AddHandler("fillAttachments", fillAttachmentsDraft),
		builder.AddHandler("setToPay", setToPayDraft),
		builder.AddHandler("setSign", setSignDraft),
		builder.AddHandler("sendMailContract", sendMailContractDraft),
		builder.AddHandler("sendMailPay", sendMailPayDraft),
		builder.AddHandler("setRequestApprovalData", setRequestApprovalBpmnDraft),
		builder.AddHandler("sendRequestApprovalMail", sendRequestApprovalMailDraft),
		builder.AddHandler("addContract", addContractDraft),
	)
}

func addContractDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}

	plc.AddContract(policy.Policy, origin)

	return nil
}

func sendMailPayDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}

	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail {
			toAddress = mail.GetContractorEmail(policy.Policy)
			ccAddress = mail.GetNetworkNodeEmail(networkNode)
		} else {
			toAddress = mail.GetNetworkNodeEmail(networkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		toAddress = mail.GetContractorEmail(policy.Policy)
	}

	log.Printf(
		"[sendMailPay] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailPay(*policy.Policy, fromAddress, toAddress, ccAddress, flowName)

	return nil
}

func sendMailContractDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}

	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail {
			toAddress = mail.GetContractorEmail(policy.Policy)
			ccAddress = mail.GetNetworkNodeEmail(networkNode)
		} else {
			toAddress = mail.GetNetworkNodeEmail(networkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		toAddress = mail.GetContractorEmail(policy.Policy)
	}

	log.Printf(
		"[sendMailContract] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailContract(*policy.Policy, nil, fromAddress, toAddress, ccAddress, flowName)

	return nil
}

func setSignDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}

	err = plc.Sign(policy.Policy, origin)
	if err != nil {
		log.Printf("[setSign] ERROR: %s", err.Error())
		return err
	}

	return nil
}

func setToPayDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}
	err = plc.SetToPay(policy.Policy, origin)
	if err != nil {
		log.Printf("[setToPay] ERROR: %s", err.Error())
		return err
	}

	return nil
}

func fillAttachmentsDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}
	err = plc.FillAttachments(policy.Policy, origin)
	if err != nil {
		log.Printf("[fillAttachments] ERROR: %s", err.Error())
		return err
	}

	return nil
}

func payTransactionDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}
	paymentInfo, err := bpmn.GetData[*flow.PaymentInfoBpmn]("paymentInfo", state)
	if err != nil {
		return err
	}
	providerId := *paymentInfo.FabrickCallback.PaymentID
	transaction, _ := tr.GetTransactionToBePaid(policy.Uid, providerId, paymentInfo.Schedule, lib.TransactionsCollection)
	err = tr.Pay(&transaction, origin, paymentInfo.PaymentMethod)
	if err != nil {
		log.Printf("[fabrickPayment] ERROR Transaction Pay %s", err.Error())
		return err
	}

	transaction.BigQuerySave(origin)

	return tr.CreateNetworkTransactions(policy.Policy, &transaction, networkNode, mgaProduct)
	return nil
}

func updatePolicyDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}

	if policy.IsPay || policy.Status != models.PolicyStatusToPay {
		log.ErrorF("policy already updated with isPay %t and status %s", policy.IsPay, policy.Status)
		return nil
	}

	// promote documents from temp bucket to user and connect it to policy
	err = plc.SetUserIntoPolicyContractor(policy.Policy, origin)
	if err != nil {
		log.ErrorF("ERROR SetUserIntoPolicyContractor %s", err.Error())
		return err
	}

	// Add Policy contract
	err = plc.AddContract(policy.Policy, origin)
	if err != nil {
		log.ErrorF("ERROR AddContract %s", err.Error())
		return err
	}

	// Update Policy as paid
	err = plc.Pay(policy.Policy, origin)
	if err != nil {
		log.ErrorF("ERROR Policy Pay %s", err.Error())
		return err
	}

	err = network.UpdateNetworkNodePortfolio(origin, policy.Policy, networkNode)
	if err != nil {
		log.ErrorF("error updating %s portfolio %s", networkNode.Type, err.Error())
		return err
	}

	policy.BigquerySave(origin)

	switch flowName {
	case models.ProviderMgaFlow:
		toAddress = mail.GetContractorEmail(policy.Policy)
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	case models.RemittanceMgaFlow:
		toAddress = mail.GetNetworkNodeEmail(networkNode)
	case models.MgaFlow, models.ECommerceFlow:
		toAddress = mail.GetContractorEmail(policy.Policy)
	}

	// Send mail with the contract to the user
	log.Printf(
		"from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailContract(*policy.Policy, nil, fromAddress, toAddress, ccAddress, flowName)

	return nil
}

func setLeadBpmnDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}
	mgaProduct = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	setLeadData(policy.Policy, *mgaProduct)
	return nil
}

func sendLeadMailDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}

	toAddress = mail.GetContractorEmail(policy.Policy)
	ccAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	}

	log.Printf(
		"[sendLeadMail] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailLead(*policy.Policy, fromAddress, toAddress, ccAddress, flowName, []string{})
	return nil
}

func setProposalDataDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}

	// TODO: remove when a proper solution to handle PMI is found
	if policy.Name == models.PmiProduct {
		return nil
	}

	if policy.ProposalNumber != 0 {
		log.Printf("[setProposalData] policy '%s' already has proposal with number '%d'", policy.Uid, policy.ProposalNumber)
		return nil
	}

	setProposalData(policy.Policy)

	log.Printf("[setProposalData] saving proposal n. %d to firestore...", policy.ProposalNumber)

	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendProposalMailDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}

	// TODO: remove when a proper solution to handle PMI is found
	if policy.Name == models.PmiProduct {
		return nil
	}

	if !sendEmail || policy.IsReserved {
		return nil
	}

	toAddress = mail.GetContractorEmail(policy.Policy)
	ccAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	}

	log.Printf(
		"[sendProposalMail] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)

	mail.SendMailProposal(*policy.Policy, fromAddress, toAddress, ccAddress, flowName, []string{models.ProposalAttachmentName})
	return nil
}

//	======================================
//	REQUEST APPROVAL FUNCTIONS
//	======================================

func setRequestApprovalBpmnDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}
	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)

	setRequestApprovalData(policy.Policy)

	log.Printf("[setRequestApproval] saving policy with uid %s to Firestore....", policy.Uid)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendRequestApprovalMailDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}

	if policy.Status == models.PolicyStatusWaitForApprovalMga {
		return nil
	}

	toAddress = mail.GetContractorEmail(policy.Policy)
	ccAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	case models.ECommerceChannel:
		toAddress = mail.Address{} // fail safe for not sending email on ecommerce reserved
	}

	if policy.Name == "qbe" {
		toAddress = mail.Address{
			Address: os.Getenv("QBE_RESERVED_MAIL"),
		}
		ccAddress = mail.Address{}
	}

	mail.SendMailReserved(*policy.Policy, fromAddress, toAddress, ccAddress, flowName,
		[]string{models.ProposalAttachmentName})
	return nil
}

func emitDataDraft(state bpmn.StorageData) error {
	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}
	emitBase(policy.Policy, origin)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendMailSignDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}

	ccAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail {
			toAddress = mail.GetContractorEmail(policy.Policy)
			ccAddress = mail.GetNetworkNodeEmail(networkNode)
		} else {
			toAddress = mail.GetNetworkNodeEmail(networkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		toAddress = mail.GetContractorEmail(policy.Policy)
	}

	log.Printf(
		"[sendMailSign] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)
	mail.SendMailSign(*policy.Policy, fromAddress, toAddress, ccAddress, flowName)
	return nil
}

func signDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}
	product, e := bpmn.GetData[*flow.ProductDraft]("product", state)
	if e != nil {
		return e
	}
	EmitSignDraft(policy.Policy, product.Product, origin)
	return nil
}

func EmitSignDraft(policy *models.Policy, product *models.Product, origin string) {
	log.Printf("[emitSign] Policy Uid %s", policy.Uid)

	policy.IsSign = false
	policy.Status = models.PolicyStatusToSign
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusContact, models.PolicyStatusToSign)

	p := <-document.ContractObj(origin, *policy, networkNode, product)
	policy.DocumentName = p.LinkGcs
	_, signResponse, _ := document.NamirialOtpV6(*policy, origin, sendEmail)
	policy.ContractFileId = signResponse.FileId
	policy.IdSign = signResponse.EnvelopeId
	policy.SignUrl = signResponse.Url
}

func payDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}
	emitPay(policy.Policy, origin)
	if policy.PayUrl == "" {
		return fmt.Errorf("missing payment url")
	}
	return nil
}

func setAdvanceDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}
	SetAdvance(policy.Policy, origin)
	return nil
}

func updateUserAndNetworkNodeDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}
	// promote documents from temp bucket to user and connect it to policy
	err := plc.SetUserIntoPolicyContractor(policy.Policy, origin)
	if err != nil {
		log.Printf("[putUser] ERROR SetUserIntoPolicyContractor %s", err.Error())
		return err
	}
	return network.UpdateNetworkNodePortfolio(origin, policy.Policy, networkNode)
}

func sendEmitProposalMailDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}

	// TODO: remove when a proper solution to handle PMI is found
	if policy.Name == models.PmiProduct {
		return nil
	}

	if policy.IsReserved {
		return nil
	}

	toAddress = mail.GetContractorEmail(policy.Policy)
	ccAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	}

	log.Printf(
		"[sendEmitProposalMail] from '%s', to '%s', cc '%s'",
		fromAddress.String(),
		toAddress.String(),
		ccAddress.String(),
	)

	mail.SendMailProposal(*policy.Policy, fromAddress, toAddress, ccAddress, flowName, []string{models.ProposalAttachmentName})
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
