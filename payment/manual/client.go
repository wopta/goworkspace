package manual

import (
	"fmt"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/payment/internal"
)

type Client struct {
	Policy       models.Policy
	Transactions []models.Transaction
}

func (c *Client) NewBusiness() (string, []models.Transaction, error) {
	log.Println("client manual: new business integration")

	if err := c.Validate(); err != nil {
		return "", nil, err
	}

	payUrl, _, transactions, err := remittanceIntegration(c.Transactions)
	return payUrl, transactions, err
}
func (c *Client) Renew() (string, bool, []models.Transaction, error) {
	log.Println("client manual: renew integration")

	if err := c.Validate(); err != nil {
		return "", false, nil, err
	}

	return remittanceIntegration(c.Transactions)
}
func (c *Client) Update() (string, []models.Transaction, error) {
	return "", nil, fmt.Errorf("manual integration does not have update")
}
func (c *Client) Validate() error {
	if len(c.Transactions) == 0 {
		return internal.ErrInvalidTransactions
	}

	if err := internal.CheckPaymentModes(c.Policy); err != nil {
		return err
	}

	return nil
}

func remittanceIntegration(transactions []models.Transaction) (payUrl string, hasMandate bool, updatedTransaction []models.Transaction, err error) {
	updatedTransaction = make([]models.Transaction, 0)

	for index, tr := range transactions {
		now := time.Now().UTC()
		if index == 0 && tr.Annuity == 0 {
			tr.IsPay = true
			tr.Status = models.TransactionStatusPay
			tr.StatusHistory = append(tr.StatusHistory, models.TransactionStatusPay)
			tr.PayDate = now
			tr.TransactionDate = now
			tr.PaymentMethod = models.PayMethodRemittance
		}
		tr.ProviderId = ""
		tr.UserToken = ""
		tr.UpdateDate = now
		updatedTransaction = append(updatedTransaction, tr)
	}
	return "", false, updatedTransaction, nil
}
