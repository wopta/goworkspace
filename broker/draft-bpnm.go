package broker

import (
	"fmt"
	"os"
	"time"

	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
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
	builder, err := bpmn.NewBpnmBuilder("broker/draftBpmn/flow/channel_flows.json")
	if err != nil {
		return nil, err
	}
	err = addHandlersDraft(builder)
	if err != nil {
		return nil, err
	}
	product = prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, networkNode, warrant)
	flowName, _ = policy.GetFlow(networkNode, warrant)

	productDraft := flow.ProductDraft{
		Product: product,
	}
	policyDraft := flow.PolicyDraft{
		Policy: policy,
	}
	networkDraft := flow.NetworkDraft{
		NetworkNode: networkNode,
	}
	address := flow.Addresses{
		ToAddress:   mail.Address{},
		CcAddress:   mail.Address{},
		FromAddress: mail.AddressAnna,
	}
	if product.Flow == "" {
		product.Flow = policy.Channel
	}

	storage.AddGlobal("policy", &policyDraft)
	storage.AddGlobal("product", &productDraft)
	storage.AddGlobal("node", &networkDraft)
	storage.AddGlobal("addresses", &address)
	builder.SetStorage(storage)

	if networkNode == nil || networkNode.CallbackConfig == nil {
		log.InfoF("no node or callback config available, no callback")
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
		builder.AddHandler("rejected", draftrejectPolicy),
		builder.AddHandler("approved", draftapprovePolicy),
		builder.AddHandler("sendAcceptanceMail", sendAcceptanceMail),

		builder.AddHandler("end_accepance", savePolicy),
		builder.AddHandler("end_emit", savePolicy),
		builder.AddHandler("end_lead", savePolicy),
	)
}

func savePolicy(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}
	policy.Updated = time.Now().UTC()
	log.Println("saving to firestore...")
	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)
	err = lib.SetFirestoreErr(firePolicy, policy.Uid, &policy)
	if err != nil {
		return err
	}
	log.Println("firestore saved!")

	policy.BigquerySave(origin)
	return nil
}

func sendAcceptanceMail(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}
	addresses, err := bpmn.GetData[*flow.Addresses]("addresses", state)
	if err != nil {
		return err
	}

	node, err := bpmn.GetData[*flow.NetworkDraft]("node", state)
	if err != nil {
		return err
	}
	log.Printf("toAddress '%s'", addresses.ToAddress.String())
	//TODO: to remove after test
	return nil
	var warrant *models.Warrant
	if node.NetworkNode != nil {
		warrant = node.GetWarrant()
	}
	flowName, _ := policy.GetFlow(node.NetworkNode, warrant)
	mail.SendMailReservedResult(
		*policy.Policy,
		mail.AddressAssunzione,
		addresses.ToAddress,
		mail.Address{},
		flowName,
	)
	return nil
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
	addresses, e := bpmn.GetData[*flow.Addresses]("addresses", state)
	if e != nil {
		return e
	}

	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail {
			addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
			addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode)
		} else {
			addresses.ToAddress = mail.GetNetworkNodeEmail(networkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	}

	log.Printf(
		"[sendMailPay] from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	mail.SendMailPay(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName)

	return nil
}

func sendMailContractDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}

	addresses, e := bpmn.GetData[*flow.Addresses]("addresses", state)
	if e != nil {
		return e
	}

	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail {
			addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
			addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode)
		} else {
			addresses.ToAddress = mail.GetNetworkNodeEmail(networkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	}

	log.Printf(
		"[sendMailContract] from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	mail.SendMailContract(*policy.Policy, nil, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName)

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
	providerId := paymentInfo.FabrickCallback.PaymentID
	transaction, _ := tr.GetTransactionToBePaid(policy.Uid, *providerId, paymentInfo.Schedule, lib.TransactionsCollection)
	err = tr.Pay(&transaction, origin, paymentInfo.PaymentMethod)
	if err != nil {
		log.Printf("[fabrickPayment] ERROR Transaction Pay %s", err.Error())
		return err
	}

	transaction.BigQuerySave(origin)

	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	return tr.CreateNetworkTransactions(policy.Policy, &transaction, networkNode, mgaProduct)
}

func updatePolicyDraft(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if err != nil {
		return err
	}

	addresses, err := bpmn.GetData[*flow.Addresses]("addresses", state)
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
		addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode)
	case models.RemittanceMgaFlow:
		addresses.ToAddress = mail.GetNetworkNodeEmail(networkNode)
	case models.MgaFlow, models.ECommerceFlow:
		addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	}

	// Send mail with the contract to the user
	log.Printf(
		"from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	mail.SendMailContract(*policy.Policy, nil, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName)

	return nil
}

func setLeadBpmnDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}
	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	setLeadData(policy.Policy, *mgaProduct)
	return nil
}

func sendLeadMailDraft(state bpmn.StorageData) error {
	policy, e := bpmn.GetData[*flow.PolicyDraft]("policy", state)
	if e != nil {
		return e
	}

	addresses, e := bpmn.GetData[*flow.Addresses]("addresses", state)
	if e != nil {
		return e
	}

	addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	addresses.CcAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode)
	}

	log.Printf(
		"[sendLeadMail] from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	mail.SendMailLead(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName, []string{})
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

	addresses, e := bpmn.GetData[*flow.Addresses]("addresses", state)
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

	addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode)
	}

	log.Printf(
		"[sendProposalMail] from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)

	mail.SendMailProposal(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName, []string{models.ProposalAttachmentName})
	return nil
}

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
	addresses, e := bpmn.GetData[*flow.Addresses]("addresses", state)
	if e != nil {
		return e
	}

	if policy.Status == models.PolicyStatusWaitForApprovalMga {
		return nil
	}

	addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	addresses.CcAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode)
	case models.ECommerceChannel:
		addresses.ToAddress = mail.Address{} // fail safe for not sending email on ecommerce reserved
	}

	if policy.Name == "qbe" {
		addresses.ToAddress = mail.Address{
			Address: os.Getenv("QBE_RESERVED_MAIL"),
		}
		addresses.CcAddress = mail.Address{}
	}

	mail.SendMailReserved(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName,
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

	addresses, e := bpmn.GetData[*flow.Addresses]("addresses", state)
	if e != nil {
		return e
	}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail {
			addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
			addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode)
		} else {
			addresses.ToAddress = mail.GetNetworkNodeEmail(networkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	}

	log.Printf(
		"[sendMailSign] from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	mail.SendMailSign(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName)
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
	product, e := bpmn.GetData[*flow.ProductDraft]("product", state)
	if e != nil {
		return e
	}

	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	EmitPayCustom(policy.Policy, origin, product.Product, mgaProduct)
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

	addresses, e := bpmn.GetData[*flow.Addresses]("addresses", state)
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

	addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	addresses.CcAddress = mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode)
	}

	log.Printf(
		"[sendEmitProposalMail] from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)

	mail.SendMailProposal(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName, []string{models.ProposalAttachmentName})
	return nil
}
