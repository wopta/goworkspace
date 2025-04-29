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
)

var (
	Proposal        string = "Proposal"
	RequestApproval string = "RequestApproval"
	Emit            string = "Emit"
	Signed          string = "Signed"
	Paid            string = "Paid"
	EmitRemittance  string = "EmitRemittance"
	Approved        string = "Approved"
	Rejected        string = "Rejected"
)

func draftEmitFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		request      EmitRequest
		err          error
		policy       models.Policy
		responseEmit EmitResponse
	)

	log.AddPrefix("EmitFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("error getting authToken")
		return "", nil, err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	origin = r.Header.Get("origin")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", nil, err
	}
	defer r.Body.Close()

	err = json.Unmarshal([]byte(body), &request)
	if err != nil {
		log.Printf("error unmarshaling policy: %s", err.Error())
		return "", nil, err
	}

	uid := request.Uid
	log.Printf("Uid: %s", uid)

	paymentSplit = request.PaymentSplit
	log.Printf("paymentSplit: %s", paymentSplit)

	paymentMode = request.PaymentMode
	log.Printf("paymentMode: %s", paymentMode)

	policy, err = plc.GetPolicy(uid, origin)
	lib.CheckError(err)

	//	if policy.Channel == models.NetworkChannel && policy.ProducerUid != authToken.UserID {
	//		log.Printf("user %s cannot emit policy %s because producer not equal to request user", authToken.UserID, policy.Uid)
	//		return "", nil, errors.New("operation not allowed")
	//	}

	policyJsonLog, _ := policy.Marshal()
	log.Printf("Policy %s JSON: %s", uid, string(policyJsonLog))

	//	if policy.IsPay || policy.IsSign || policy.CompanyEmit || policy.CompanyEmitted || policy.IsDeleted {
	//		log.Printf("cannot emit policy %s because state is not correct", policy.Uid)
	//		return "", nil, errors.New("operation not allowed")
	//	}

	if request.SendEmail == nil {
		sendEmail = true
	} else {
		sendEmail = *request.SendEmail
	}

	productConfig := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	if err = policy.CheckStartDateValidity(productConfig.EmitMaxElapsedDays); err != nil {
		return "", "", err
	}

	emitUpdatePolicy(&policy, request)

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	if policy.IsReserved && policy.Status != models.PolicyStatusApproved {
		log.Printf("cannot emit policy uid %s with status %s and isReserved %t", policy.Uid, policy.Status, policy.IsReserved)
		return "", nil, fmt.Errorf("cannot emit policy uid %s with status %s and isReserved %t", policy.Uid, policy.Status, policy.IsReserved)
	}
	responseEmit, err = emitDraft(&policy, request, origin)
	if err != nil {
		return "", nil, err
	}

	b, err := json.Marshal(responseEmit)

	log.Println("Handler end -------------------------------------------------")

	return string(b), responseEmit, err
}

func emitDraft(policy *models.Policy, request EmitRequest, origin string) (EmitResponse, error) {
	log.Println("[Emit] start ------------------------------------------------")
	var responseEmit EmitResponse

	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)
	fireGuarantee := lib.GetDatasetByEnv(origin, lib.GuaranteeCollection)

	log.Printf("[Emit] Emitting - Policy Uid %s", policy.Uid)
	log.Println("[Emit] starting bpmn flow...")
	flow, err := getFlow(policy, nil)
	if err != nil {
		return responseEmit, err
	}
	err = flow.Run("emit")
	if err != nil {
		return responseEmit, err
	}

	responseEmit = EmitResponse{
		UrlPay:       policy.PayUrl,
		UrlSign:      policy.SignUrl,
		ReservedInfo: policy.ReservedInfo,
		Uid:          policy.Uid,
		CodeCompany:  policy.CodeCompany,
	}

	policy.Updated = time.Now().UTC()
	policyJson, _ := policy.Marshal()
	log.Printf("[Emit] Policy %s: %s", request.Uid, string(policyJson))

	log.Println("[Emit] saving policy to firestore...")
	err = lib.SetFirestoreErr(firePolicy, request.Uid, policy)
	if err != nil {
		return responseEmit, err
	}

	log.Println("[Emit] saving policy to bigquery...")
	policy.BigquerySave(origin)

	log.Println("[Emit] saving guarantees to bigquery...")
	models.SetGuaranteBigquery(*policy, "emit", fireGuarantee)

	log.Println("[Emit] end --------------------------------------------------")
	return responseEmit, nil
}

func getFlow(policy *models.Policy, networkNode *models.NetworkNode) (*bpmn.FlowBpnm, error) {
	builder, err := bpmn.NewBpnmBuilder("broker/draftBpnm/flow/channel_flows.json")
	if err != nil {
		return nil, err
	}
	err = addHandlersDraft(builder)
	if err != nil {
		return nil, err
	}
	storage := bpmn.NewStorageBpnm()
	product = prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, networkNode, warrant)
	flowName, _ = policy.GetFlow(networkNode, warrant)
	productDraft := flow.ProductDraft{
		product,
	}

	mgaProduct = product
	toAddress = mail.Address{}
	ccAddress = mail.Address{}
	fromAddress = mail.AddressAnna

	policyDraft := flow.PolicyDraft{
		policy,
	}
	if product.Flow == "" {
		product.Flow = policy.Channel
	}
	networkDraft := flow.NetworkDraft{
		networkNode,
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
	bpmn.IsError(
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
		builder.AddHandler("saveAudit", func(st bpmn.StorageData) error {
			return nil
		}),
	)
	return builder, nil
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
		builder.AddHandler("updatePolicy", updateUserAndNetworkNodeDraft),
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
	//	policy, err := bpmn.GetData[*bpmn.PolicyDraft]("policy", state)
	//	if err != nil {
	//		return err
	//	}
	//	providerId := *fabrickCallback.PaymentID
	//	transaction, _ := tr.GetTransactionToBePaid(policy.Uid, providerId, trSchedule, lib.TransactionsCollection)
	//	err := tr.Pay(&transaction, origin, paymentMethod)
	//	if err != nil {
	//		log.Printf("[fabrickPayment] ERROR Transaction Pay %s", err.Error())
	//		return err
	//	}
	//
	//	transaction.BigQuerySave(origin)

	//return tr.CreateNetworkTransactions(policy, &transaction, networkNode, mgaProduct)
	return nil
}
func updatePolicyDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}

	if policy.IsPay || policy.Status != models.PolicyStatusToPay {
		log.Printf("[updatePolicy] policy already updated with isPay %t and status %s", policy.IsPay, policy.Status)
		return nil
	}

	// promote documents from temp bucket to user and connect it to policy
	err = plc.SetUserIntoPolicyContractor(policy.Policy, origin)
	if err != nil {
		log.Printf("[updatePolicy] ERROR SetUserIntoPolicyContractor %s", err.Error())
		return err
	}

	// Add Policy contract
	err = plc.AddContract(policy.Policy, origin)
	if err != nil {
		log.Printf("[updatePolicy] ERROR AddContract %s", err.Error())
		return err
	}

	// Update Policy as paid
	err = plc.Pay(policy.Policy, origin)
	if err != nil {
		log.Printf("[updatePolicy] ERROR Policy Pay %s", err.Error())
		return err
	}

	err = network.UpdateNetworkNodePortfolio(origin, policy.Policy, networkNode)
	if err != nil {
		log.Printf("[updatePolicy] error updating %s portfolio %s", networkNode.Type, err.Error())
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
		"[updatePolicy] from '%s', to '%s', cc '%s'",
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
