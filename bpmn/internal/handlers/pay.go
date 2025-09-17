package handlers

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/payment/consultancy"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	plcR "gitlab.dev.wopta.it/goworkspace/policy/renew"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
)

func AddPayHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		builder.AddHandler("payTransaction", payTransaction),
		builder.AddHandler("updatePolicyAsPaid", updatePolicyAsPaid),
		builder.AddHandler("promotePolicy", promotePolicy),
		builder.AddHandler("generateInvoice", generateInvoice),
		builder.AddHandler("saveTransactionAndPolicy", saveTransactionAndPolicy),
		builder.AddHandler("createNetworkTransaction", createNetworkTransaction),
	)
}

func createNetworkTransaction(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var flowName *flow.String
	var mgaProduct *flow.Product
	var transaction *flow.Transaction
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
		bpmnEngine.GetDataRef("mgaProduct", &mgaProduct, state),
		bpmnEngine.GetDataRef("transaction", &transaction, state),
	)
	if err != nil {
		return err
	}
	return tr.CreateNetworkTransactions(policy.Policy, transaction.Transaction, networkNode.NetworkNode, mgaProduct.Product)
}
func saveTransactionAndPolicy(state *bpmnEngine.StorageBpnm) error {
	var transaction *flow.Transaction
	var policy *flow.Policy
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("transaction", &transaction, state),
		bpmnEngine.GetDataRef("policy", &policy, state),
	)
	if err != nil {
		return err
	}

	firestoreBatch := map[string]map[string]interface{}{
		lib.TransactionsCollection: {
			transaction.Uid: transaction,
		},
		lib.PolicyCollection: {
			policy.Uid: policy,
		},
	}
	policyCollection := lib.PolicyCollection
	transactionsCollection := lib.TransactionsCollection
	if renewPolicy, err := plcR.GetRenewPolicyByUid(policy.Uid); err == nil {
		policy.Policy = &renewPolicy
		policyCollection = lib.RenewPolicyCollection
		transactionsCollection = lib.RenewTransactionCollection
	}
	if err = lib.SetBatchFirestoreErr(firestoreBatch); err != nil {
		return err
	}
	if err = lib.InsertRowsBigQuery(lib.WoptaDataset, policyCollection, policy); err != nil {
		return err
	}
	if err = lib.InsertRowsBigQuery(lib.WoptaDataset, transactionsCollection, transaction); err != nil {
		return err
	}
	return nil
}

func generateInvoice(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var transaction *flow.Transaction
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("transaction", &transaction, state),
	)

	if err := consultancy.GenerateInvoice(*policy.Policy, *transaction.Transaction); err != nil {
		log.Printf("error handling consultancy: %s", err.Error())
	}
	return err
}
func updatePolicyAsPaid(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
	)
	if err != nil {
		return err
	}
	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)

	if policy.Annuity > 0 {
		policy.Status = models.PolicyStatusRenewed
		policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	}

	policy.Updated = time.Now().UTC()

	policy.BigQueryParse()
	return nil
}
func promotePolicy(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var flowName *flow.String
	var addresses *flow.Addresses
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
		bpmnEngine.GetDataRef("addresses", &addresses, state),
	)
	if err != nil {
		return err
	}

	if policy.IsPay || policy.Status != models.PolicyStatusToPay {
		log.WarningF("policy already updated with isPay %t and status %s", policy.IsPay, policy.Status)
		state.AddLocal("sendEmail", &flow.BoolBpmn{Bool: false})
		return nil
	}

	// Add Policy contract
	err = plc.AddSignedDocumentsInPolicy(policy.Policy)
	if err != nil {
		return err
	}

	// promote documents from temp bucket to user and connect it to policy
	err = plc.PromotePolicy(policy.Policy)
	if err != nil {
		return err
	}

	err = network.UpdateNetworkNodePortfolio(policy.Policy, networkNode.NetworkNode)
	if err != nil {
		return err
	}

	if err = plc.RemoveTempPolicy(policy.Policy); err != nil {
		return err
	}

	state.AddLocal("sendEmail", &flow.BoolBpmn{Bool: true})
	return nil
}

func payTransaction(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var paymentInfo *flow.PaymentInfoBpmn
	var networkNode *flow.Network
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("paymentInfo", &paymentInfo, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}

	providerId := paymentInfo.FabrickPaymentsRequest.PaymentID
	transaction, _ := tr.GetTransactionToBePaid(policy.Uid, providerId, paymentInfo.Schedule, lib.TransactionsCollection)
	err = tr.Pay(&transaction, paymentInfo.PaymentMethod)
	if err != nil {
		log.Error(err)
		return err
	}

	transaction.BigQuerySave()
	state.AddGlobal("transaction", &flow.Transaction{Transaction: &transaction})
	state.AddLocal("skipInvoice", &flow.BoolBpmn{Bool: lib.IsEqual(policy.StartDate.AddDate(policy.Annuity, 0, 0), transaction.EffectiveDate)})
	return nil
}
