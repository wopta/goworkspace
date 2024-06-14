package payment

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/models"
)

var (
	errInvalidTransactions = errors.New("invalid number of transactions")
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
		return &ManualClient{
			Policy:       policy,
			Transactions: transactions,
		}
	case models.FabrickPaymentProvider:
		return &FabrickClient{
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

// FABRICK CLIENT

type FabrickClient struct {
	Policy            models.Policy
	Product           models.Product
	Transactions      []models.Transaction
	ScheduleFirstRate bool
	CustomerId        string
}

func (c *FabrickClient) Validate() error {
	if len(c.Transactions) == 0 {
		return errInvalidTransactions
	}

	if err := checkPaymentModes(c.Policy); err != nil {
		return err
	}

	return nil
}

func (c *FabrickClient) getPaymentMethods() []string {
	var paymentMethods = make([]string, 0)

	for _, provider := range c.Product.PaymentProviders {
		if provider.Name == c.Policy.Payment {
			for _, config := range provider.Configs {
				if config.Mode == c.Policy.PaymentMode && config.Rate == c.Policy.PaymentSplit {
					paymentMethods = append(paymentMethods, config.Methods...)
				}
			}
		}
	}

	log.Printf("[FabrickClient.getPaymentMethods] found %v", paymentMethods)
	return paymentMethods
}

func (c FabrickClient) NewBusiness() (string, []models.Transaction, error) {
	log.Println("client fabrick: new business integration")

	if err := c.Validate(); err != nil {
		return "", nil, err
	}

	var updatedTransactions = make([]models.Transaction, 0)

	paymentMethods := c.getPaymentMethods()
	now := time.Now().UTC()
	c.CustomerId = uuid.New().String()

	for index, tr := range c.Transactions {
		isFirstOfBatch := index == 0
		isFirstRateOfAnnuity := c.Policy.StartDate.Month() == tr.EffectiveDate.Month()

		createMandate := c.Policy.PaymentMode == models.PaymentModeRecurrent && isFirstOfBatch
		log.Printf("createMandate: %v", createMandate)

		tr.ProviderName = models.FabrickPaymentProvider

		res := <-createFabrickTransaction(&c.Policy, tr, createMandate, false, isFirstRateOfAnnuity, c.CustomerId, paymentMethods)
		if res.Payload == nil || res.Payload.PaymentPageURL == nil {
			return "", nil, errors.New("error creating transaction on Fabrick")
		}
		log.Printf("transaction %02d payUrl: %s", index+1, *res.Payload.PaymentPageURL)

		tr.ProviderId = *res.Payload.PaymentID
		tr.PayUrl = *res.Payload.PaymentPageURL
		tr.UserToken = c.CustomerId
		tr.UpdateDate = now
		updatedTransactions = append(updatedTransactions, tr)
	}

	if len(updatedTransactions) != len(c.Transactions) {
		return "", nil, fmt.Errorf("invalid number of updatedTransactions")
	}

	return updatedTransactions[0].PayUrl, updatedTransactions, nil
}
func (c FabrickClient) Renew() (string, []models.Transaction, error) {
	log.Println("client fabrick: renew integration")

	if err := c.Validate(); err != nil {
		return "", nil, err
	}

	var (
		updatedTransactions = make([]models.Transaction, 0)
		payUrl              string
		err                 error
	)

	paymentMethods := c.getPaymentMethods()
	hasMandate := false
	now := time.Now().UTC()

	if hasMandate, err = fabrickHasMandate(c.CustomerId); err != nil {
		log.Printf("error checking mandate: %s", err.Error())
	}

	// TODO: this might change in the future. It works as following:
	// - if no customerId is previously provided, scheduleFirstRate will have the
	// inputed value as the true value
	// - if there is a provided customerId, scheduleFirstRate will follow the fact
	// that the user has or not an active mandate
	// - currently the second case is used only in renew.
	if c.CustomerId != "" {
		c.ScheduleFirstRate = hasMandate
	}

	if c.CustomerId == "" {
		c.CustomerId = uuid.New().String()
	}

	for index, tr := range c.Transactions {
		isFirstOfBatch := index == 0
		isFirstRateOfAnnuity := c.Policy.StartDate.Month() == tr.EffectiveDate.Month()

		createMandate := c.Policy.PaymentMode == models.PaymentModeRecurrent && isFirstOfBatch && !hasMandate
		log.Printf("createMandate: %v", createMandate)

		tr.ProviderName = models.FabrickPaymentProvider

		scheduleDate, err := time.Parse(time.DateOnly, tr.ScheduleDate)
		if err != nil {
			log.Printf("error parsing scheduleDate: %s", err.Error())
			return "", nil, err
		}
		if c.ScheduleFirstRate && scheduleDate.Before(now) {
			/*
				sets schedule date to today + 1 in order to avoid corner case in which fabrick is not able to
				execute transaction when recreated at the end of the day
			*/
			tr.ScheduleDate = now.AddDate(0, 0, 1).Format(time.DateOnly)
		}

		res := <-createFabrickTransaction(&c.Policy, tr, createMandate, c.ScheduleFirstRate, isFirstRateOfAnnuity, c.CustomerId, paymentMethods)
		if res.Payload == nil || res.Payload.PaymentPageURL == nil {
			return "", nil, errors.New("error creating transaction on Fabrick")
		}
		if isFirstOfBatch && (!hasMandate || createMandate) {
			payUrl = *res.Payload.PaymentPageURL
		}
		log.Printf("transaction %02d payUrl: %s", index+1, *res.Payload.PaymentPageURL)

		tr.ProviderId = *res.Payload.PaymentID
		tr.PayUrl = *res.Payload.PaymentPageURL
		tr.UserToken = c.CustomerId
		tr.UpdateDate = now
		updatedTransactions = append(updatedTransactions, tr)
	}

	return payUrl, updatedTransactions, nil
}
func (c FabrickClient) Update() (string, []models.Transaction, error) {
	log.Println("client fabrick: update integration")

	if err := c.Validate(); err != nil {
		return "", nil, err
	}

	var updatedTransactions = make([]models.Transaction, 0)

	paymentMethods := c.getPaymentMethods()
	now := time.Now().UTC()
	c.CustomerId = uuid.New().String()

	for index, tr := range c.Transactions {
		isFirstOfBatch := index == 0
		isFirstRateOfAnnuity := c.Policy.StartDate.Month() == tr.EffectiveDate.Month()

		createMandate := c.Policy.PaymentMode == models.PaymentModeRecurrent && isFirstOfBatch
		log.Printf("createMandate: %v", createMandate)

		tr.ProviderName = models.FabrickPaymentProvider

		scheduleDate, err := time.Parse(time.DateOnly, tr.ScheduleDate)
		if err != nil {
			log.Printf("error parsing scheduleDate: %s", err.Error())
			return "", nil, err
		}
		if c.ScheduleFirstRate && scheduleDate.Before(now) {
			/*
				sets schedule date to today + 1 in order to avoid corner case in which fabrick is not able to
				execute transaction when recreated at the end of the day
			*/
			tr.ScheduleDate = now.AddDate(0, 0, 1).Format(time.DateOnly)
		}

		res := <-createFabrickTransaction(&c.Policy, tr, createMandate, c.ScheduleFirstRate, isFirstRateOfAnnuity, c.CustomerId, paymentMethods)
		if res.Payload == nil || res.Payload.PaymentPageURL == nil {
			return "", nil, errors.New("error creating transaction on Fabrick")
		}

		log.Printf("transaction %02d payUrl: %s", index+1, *res.Payload.PaymentPageURL)

		tr.ProviderId = *res.Payload.PaymentID
		tr.PayUrl = *res.Payload.PaymentPageURL
		tr.UserToken = c.CustomerId
		tr.UpdateDate = now
		updatedTransactions = append(updatedTransactions, tr)
	}

	if len(updatedTransactions) != len(c.Transactions) {
		return "", nil, fmt.Errorf("invalid number of updatedTransactions")
	}

	return updatedTransactions[0].PayUrl, updatedTransactions, nil
}

// MANUAL CLIENT

type ManualClient struct {
	Policy       models.Policy
	Transactions []models.Transaction
}

func (c *ManualClient) NewBusiness() (string, []models.Transaction, error) {
	log.Println("client manual: new business integration")

	if err := c.Validate(); err != nil {
		return "", nil, err
	}

	return remittanceIntegration(c.Transactions)
}
func (c *ManualClient) Renew() (string, []models.Transaction, error) {
	log.Println("client manual: renew integration")

	if err := c.Validate(); err != nil {
		return "", nil, err
	}

	return remittanceIntegration(c.Transactions)
}
func (c *ManualClient) Update() (string, []models.Transaction, error) {
	return "", nil, fmt.Errorf("manual integration does not have update")
}
func (c *ManualClient) Validate() error {
	if len(c.Transactions) == 0 {
		return errInvalidTransactions
	}

	if err := checkPaymentModes(c.Policy); err != nil {
		return err
	}

	return nil
}
