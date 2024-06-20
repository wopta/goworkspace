package payment

import (
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment/fabrick"
	"github.com/wopta/goworkspace/payment/manual"
)

type Client interface {
	NewBusiness() (string, []models.Transaction, error)
	Renew() (string, []models.Transaction, error)
	Update() (string, []models.Transaction, error)
	Validate() error
}

func NewClient(client string, policy models.Policy, product models.Product, transactions []models.Transaction, scheduleFirstRate bool, customerId string) Client {
	switch client {
	case models.ManualPaymentProvider:
		return &manual.Client{
			Policy:       policy,
			Transactions: transactions,
		}
	case models.FabrickPaymentProvider:
		return &fabrick.Client{
			Policy:            policy,
			Transactions:      transactions,
			Product:           product,
			ScheduleFirstRate: scheduleFirstRate,
			CustomerId:        customerId,
		}
	default:
		return nil
	}
}
