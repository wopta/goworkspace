package transaction

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
)

func createNetworkTransaction(
	policy *models.Policy,
	transaction *models.Transaction,
	node *models.NetworkNode,
	commission float64, // Amount
	mgaAccountType, paymentType, name string,
) (*models.NetworkTransaction, error) {
	log.Printf(
		"[createNetworkTransaction] name '%s' accountType '%s' paymentType '%s' commission '%f' amount '%f'",
		name,
		mgaAccountType,
		paymentType,
		commission,
		transaction.Amount,
	)

	var amount float64

	switch paymentType {
	case models.PaymentTypeRemittanceCompany, models.PaymentTypeCommission:
		amount = lib.RoundFloat(commission, 2)
	case models.PaymentTypeRemittanceMga:
		amount = lib.RoundFloat(transaction.Amount-commission, 2)
	}

	netTransaction := models.NetworkTransaction{
		Uid:              uuid.New().String(),
		PolicyUid:        policy.Uid,
		TransactionUid:   transaction.Uid,
		NetworkUid:       node.NetworkUid,
		NetworkNodeUid:   node.Uid,
		NetworkNodeType:  node.Type,
		AccountType:      mgaAccountType,
		PaymentType:      paymentType,
		Amount:           amount,
		AmountNet:        amount, // TBD
		Name:             name,
		Status:           models.NetworkTransactionStatusToPay,
		StatusHistory:    []string{models.NetworkTransactionStatusCreated, models.NetworkTransactionStatusToPay},
		IsPay:            false,
		IsConfirmed:      false,
		CreationDate:     lib.GetBigQueryNullDateTime(time.Now().UTC()),
		PayDate:          lib.GetBigQueryNullDateTime(time.Time{}),
		TransactionDate:  lib.GetBigQueryNullDateTime(time.Time{}),
		ConfirmationDate: lib.GetBigQueryNullDateTime(time.Time{}),
	}

	jsonLog, _ := json.Marshal(&netTransaction)

	err := netTransaction.SaveBigQuery()
	if err != nil {
		log.Printf("[createNetworkTransaction] error saving network transaction to bigquery: %s", err.Error())
		return nil, err
	}

	log.Printf("[createNetworkTransaction] network transaction created! %s", string(jsonLog))

	return &netTransaction, err
}

func createCompanyNetworkTransaction(
	policy *models.Policy,
	transaction *models.Transaction,
	producerNode *models.NetworkNode,
	mgaProduct *models.Product,
) (*models.NetworkTransaction, error) {
	log.Println("[createCompanyNetworkTransaction]")

	var code string

	if isGapCamper(policy) {
		log.Println("[GetCommissionByNode] overrinding product commissions for Gap camper")
		mgaProduct.Offers["base"].Commissions.NewBusiness = 0.22
		mgaProduct.Offers["base"].Commissions.NewBusinessPassive = 0.22
		mgaProduct.Offers["base"].Commissions.Renew = 0
		mgaProduct.Offers["base"].Commissions.RenewPassive = 0
		mgaProduct.Offers["complete"].Commissions.NewBusiness = 0.37
		mgaProduct.Offers["complete"].Commissions.NewBusinessPassive = 0.37
		mgaProduct.Offers["complete"].Commissions.Renew = 0
		mgaProduct.Offers["complete"].Commissions.RenewPassive = 0
	}

	commissionMga := product.GetCommissionByProduct(policy, mgaProduct, false)
	commissionCompany := lib.RoundFloat(transaction.Amount-commissionMga, 2)

	if producerNode != nil {
		code = producerNode.Code
	} else {
		code = models.ECommerceChannel
	}

	name := strings.ToUpper(strings.Join([]string{code, policy.Company}, "-"))

	return createNetworkTransaction(
		policy,
		transaction,
		&models.NetworkNode{},
		commissionCompany,
		models.AccountTypePassive,
		models.PaymentTypeRemittanceCompany,
		name,
	)
}

func CreateNetworkTransactions(
	policy *models.Policy,
	transaction *models.Transaction,
	producerNode *models.NetworkNode,
	mgaProduct *models.Product,
) error {
	log.Println("[CreateNetworkTransactions]")

	var err error

	_, err = createCompanyNetworkTransaction(policy, transaction, producerNode, mgaProduct)
	if err != nil {
		log.Printf("[CreateNetworkTransactions] error creating company network-transaction: %s", err.Error())
		return err
	}

	if policy.ProducerUid != "" && policy.ProducerType != models.PartnershipNetworkNodeType {
		network.TraverseWithCallbackNetworkByNodeUid(producerNode, "", func(currentNode *models.NetworkNode, currentName string) string {
			var (
				accountType, paymentType string
				baseName                 string
			)

			warrant := currentNode.GetWarrant()
			if warrant == nil {
				log.Printf("[CreateNetworkTransactions] error getting warrant for node: %s", currentNode.Uid)
				return baseName
			}
			prod := warrant.GetProduct(policy.Name)
			if warrant == nil {
				log.Printf("[CreateNetworkTransactions] error getting product for warrant: %s", warrant.Name)
				return baseName
			}

			accountType = getAccountType(transaction)
			paymentType = getPaymentType(transaction, policy, currentNode)
			commission := product.GetCommissionByProduct(policy, prod, policy.ProducerUid == currentNode.Uid)

			if currentName != "" {
				baseName = strings.ToUpper(strings.Join([]string{currentName, currentNode.Code}, "__"))
			} else {
				baseName = strings.ToUpper(currentNode.Code)
			}
			nodeName := strings.ToUpper(strings.Join([]string{baseName, currentNode.GetName()}, "-"))

			_, err = createNetworkTransaction(policy, transaction, currentNode, commission, accountType, paymentType, nodeName)
			if err != nil {
				log.Printf("[CreateNetworkTransactions] error creating network-transaction: %s", err.Error())
			} else {
				log.Printf("[CreateNetworkTransactions] created network-transaction for node: %s", currentNode.Uid)
			}

			return baseName
		})
	}

	return err
}

func getAccountType(transaction *models.Transaction) string {
	if transaction.ProviderName == models.ManualPaymentProvider {
		return models.AccountTypeActive
	}
	return models.AccountTypePassive
}

func getPaymentType(transaction *models.Transaction, policy *models.Policy, producerNode *models.NetworkNode) string {
	if policy.ProducerUid == producerNode.Uid && transaction.ProviderName == models.ManualPaymentProvider {
		return models.PaymentTypeRemittanceMga
	}
	return models.PaymentTypeCommission
}

func isGapCamper(policy *models.Policy) bool {
	return policy.Name == models.GapProduct &&
		len(policy.Assets) > 0 &&
		policy.Assets[0].Vehicle != nil &&
		policy.Assets[0].Vehicle.VehicleType == "camper"
}
