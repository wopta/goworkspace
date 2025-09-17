package fabrick

import (
	"errors"
	"fmt"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/google/uuid"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/payment/internal"
)

type Client struct {
	Policy            models.Policy
	Product           models.Product
	Transactions      []models.Transaction
	ScheduleFirstRate bool
	CustomerId        string
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

func (c *Client) getPaymentMethods() []string {
	var paymentMethods = make([]string, 0)
	log.AddPrefix("FabrickClient.getPaymentMethods")
	defer log.PopPrefix()
	for _, provider := range c.Product.PaymentProviders {
		if provider.Name == c.Policy.Payment {
			for _, config := range provider.Configs {
				if config.Mode == c.Policy.PaymentMode && config.Rate == c.Policy.PaymentSplit {
					paymentMethods = append(paymentMethods, config.Methods...)
				}
			}
		}
	}

	log.Printf("found %v", paymentMethods)
	return paymentMethods
}

func (c Client) NewBusiness() (string, []models.Transaction, error) {
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

func (c Client) Renew() (string, bool, []models.Transaction, error) {
	log.Println("client fabrick: renew integration")

	if err := c.Validate(); err != nil {
		return "", false, nil, err
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
		log.ErrorF("error checking mandate: %s", err.Error())
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
			log.ErrorF("error parsing scheduleDate: %s", err.Error())
			return "", false, nil, err
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
			return "", false, nil, errors.New("error creating transaction on Fabrick")
		}
		if (isFirstOfBatch && c.Policy.PaymentMode != models.PaymentModeRecurrent) || createMandate {
			payUrl = *res.Payload.PaymentPageURL
		}
		log.Printf("transaction %02d payUrl: %s", index+1, *res.Payload.PaymentPageURL)

		tr.ProviderId = *res.Payload.PaymentID
		tr.PayUrl = *res.Payload.PaymentPageURL
		tr.UserToken = c.CustomerId
		tr.UpdateDate = now
		updatedTransactions = append(updatedTransactions, tr)
	}

	return payUrl, hasMandate, updatedTransactions, nil
}
func (c Client) Update() (string, []models.Transaction, error) {
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
			log.ErrorF("error parsing scheduleDate: %s", err.Error())
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
