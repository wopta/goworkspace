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

	prod, err := product.GetProduct(policy.Name, policy.ProductVersion, models.UserRoleAdmin)
	if err != nil {
		log.Printf("[createCompanyNetworkTransaction] error getting mga product: %s", err.Error())
		return nil, err
	}

	commissionCompany := product.GetCommissionByNode(policy, prod, false)

	name := strings.ToLower(strings.Join([]string{producerNode.Code, policy.Company}, "_"))

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
		network.TraverseNetworkByNodeUid(producerNode, func(currentNode *models.NetworkNode, currentName string) string {
			var (
				accountType, paymentType string
				prod                     models.Product
			)

			for _, p := range currentNode.Products {
				if p.Name == policy.Name {
					prod = p
					break
				}
			}
			isActive := transaction.ProviderName == models.ManualPaymentProvider
			commission := product.GetCommissionByNode(policy, &prod, isActive)
			if isActive {
				accountType = models.AccountTypeActive
				paymentType = models.PaymentTypeRemittanceMga
			} else {
				accountType = models.AccountTypePassive
				paymentType = models.PaymentTypeCommission
			}

			nodeName := strings.ToLower(strings.Join([]string{currentName, currentNode.Code, currentNode.GetName()}, "_"))

			_, err = createNetworkTransaction(policy, transaction, currentNode, commission, accountType, paymentType, nodeName)
			if err != nil {
				log.Printf("[CreateNetworkTransactions] error creating network-transaction: %s", err.Error())
			} else {
				log.Printf("[CreateNetworkTransactions] created network-transaction for node: %s", currentNode.Uid)
			}

			return nodeName
		})
	}

	return err
}
