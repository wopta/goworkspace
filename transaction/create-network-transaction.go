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
	accountType, paymentType, name string,
) (*models.NetworkTransaction, error) {
	log.Printf(
		"[createNetworkTransaction] name '%s' accountType '%s' paymentType '%s' commission '%f' amount '%f'",
		name,
		accountType,
		paymentType,
		commission,
		transaction.Amount,
	)

	var amount float64

	switch paymentType {
	case models.PaymentTypeRemittanceCompany, models.PaymentTypeCommission:
		amount = commission
	case models.PaymentTypeRemittanceMga:
		amount = transaction.Amount - commission
	}

	if accountType == models.AccountTypePassive {
		amount = -amount
	}

	netTransaction := models.NetworkTransaction{
		Uid:              uuid.New().String(),
		PolicyUid:        policy.Uid,
		TransactionUid:   transaction.Uid,
		NetworkUid:       node.NetworkUid,
		NetworkNodeUid:   node.Uid,
		NetworkNodeType:  node.Type,
		AccountType:      accountType,
		PaymentType:      paymentType,
		Amount:           amount,
		AmountNet:        amount, // TBD
		Name:             name,
		Status:           models.NetworkTransactionStatusCreated,
		StatusHistory:    []string{models.NetworkTransactionStatusCreated},
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

func createCompanyNetworkTransaction(policy *models.Policy, transaction *models.Transaction, producerNode *models.NetworkNode) (*models.NetworkTransaction, error) {
	log.Println("[createCompanyNetworkTransaction]")

	var code string

	prod, err := product.GetProduct(policy.Name, policy.ProductVersion, models.UserRoleAdmin)
	if err != nil {
		log.Printf("[createCompanyNetworkTransaction] error getting mga product: %s", err.Error())
		return nil, err
	}

	commissionCompany := product.GetCommissionByNode(policy, prod, false)

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

func CreateNetworkTransactions(policy *models.Policy, transaction *models.Transaction, producerNode *models.NetworkNode) error {
	log.Println("[CreateNetworkTransactions]")

	var (
		err error
	)

	_, err = createCompanyNetworkTransaction(policy, transaction, producerNode)
	if err != nil {
		log.Printf("[CreateNetworkTransactions] error creating company network-transaction: %s", err.Error())
		return err
	}

	if policy.ProducerUid != "" && policy.ProducerType != "partnership" { // use constant
		network.TraverseWithCallbackNetworkByNodeUid(producerNode, "", func(currentNode *models.NetworkNode, currentName string) string {
			var (
				accountType, paymentType string
				prod                     models.Product
				baseName                 string
			)

			for _, p := range currentNode.Products {
				if p.Name == policy.Name {
					prod = p
					break
				}
			}

			accountType = getAccountType(transaction)
			paymentType = getPaymentType(transaction, policy, currentNode)
			commission := product.GetCommissionByNode(policy, &prod, policy.ProducerUid == currentNode.Uid)

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
