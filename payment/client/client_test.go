package client_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/payment/client"
)

func TestClient(t *testing.T) {
	os.Setenv("env", "local-test")

	t.Run("invalid configurations", testClientInvalid)
	t.Run("new business", testClientNewBusiness)
	t.Run("refresh", testClientRefresh)
	t.Run("renew", testClientRenew)
}

func testClientInvalid(t *testing.T) {
	t.Run("transactions number", invalidNumTransactions)
	t.Run("payment configuration", invalidPaymentConfiguration)
}

func testClientNewBusiness(t *testing.T) {
	t.Run("fabrick yearly single", newBusinessFabrickYearlySingle)
	t.Run("fabrick yearly recurrent", newBusinessFabrickYearlyRecurrent)
	t.Run("fabrick monthly", newBusinessFabrickMonthly)
	t.Run("manual remittance", newBusinessManualRemittance)
}

func testClientRefresh(t *testing.T) {
	t.Run("fabrick monthly", refreshFabrickMonthly)
	t.Run("fabrick yearly single", refreshFabrickYearlySingle)
	t.Run("fabrick yearly recurrent", refreshFabrickYearlyRecurrent)
	t.Run("fabrick monthly renewed", refreshFabrickMonthlyRenewed)
	t.Run("fabrick yearly single renewed", refreshFabrickYearlySingleRenewed)
	t.Run("fabrick yearly recurrent renewed", refreshFabrickYearlyRecurrentRenewed)
}

func testClientRenew(t *testing.T) {
	t.Run("fabrick yearly single", renewFabrickYearlySingle)
	t.Run("fabrick yearly recurrent with mandate", renewFabrickYearlyRecurrentWithMandate)
	t.Run("fabrick yearly recurrent without mandate", renewFabrickYearlyRecurrentWithoutMandate)
	t.Run("fabrick monthly with mandate", renewFabrickMonthlyWithMandate)
	t.Run("fabrick monthly without mandate", renewFabrickMonthlyWithoutMandate)
	t.Run("manual remittance", renewManualRemittance)
}

func invalidNumTransactions(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 0)
	product := getProduct()
	transactions := getTransactions(0, models.FabrickPaymentProvider, 0, time.Time{})

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	_, updatedTransactions, err := c.NewBusiness()
	if err == nil {
		t.Fatalf("expected: %02d transactions got: %02d", 0, len(updatedTransactions))
	}
}

func invalidPaymentConfiguration(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeSingle, string(models.PaySplitMonthly), 0)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 0, time.Time{})

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	_, _, err := c.NewBusiness()
	if err == nil {
		t.Fatalf("expected: non-nil error")
	}
}

func newBusinessFabrickYearlySingle(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeSingle, string(models.PaySplitYearly), 0)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 0, time.Time{})

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	payUrl, updatedTransactions, err := c.NewBusiness()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}

	if len(updatedTransactions) != len(transactions) {
		t.Fatalf("expected: %d transactions got: %d", len(transactions), len(updatedTransactions))
	}

	if payUrl == "" {
		t.Fatalf("expected: non-empty payUrl got: %s", payUrl)
	}

	if payUrl != "www.dev.wopta.it/local-00" {
		t.Fatalf("expected: www.dev.wopta.it/local-00, got: %s", payUrl)
	}

	for index, tr := range updatedTransactions {
		if tr.ScheduleDate == "" {
			t.Fatalf("expected: non-empty ScheduleDate")
		}

		if tr.ProviderId != "local" {
			t.Fatalf("expected: %s ProviderName got: %s", "local", tr.ProviderId)
		}

		if tr.UpdateDate.Equal(transactions[index].UpdateDate) {
			t.Fatalf("expected: %s UpdateDate got: %s", transactions[index].UpdateDate, tr.UpdateDate)
		}

		if tr.UserToken == "" {
			t.Fatalf("expected: non-empty UserToken")
		}

		if tr.ProviderName != policy.Payment {
			t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", tr.ProviderName, policy.Payment)
		}
	}
}

func newBusinessFabrickYearlyRecurrent(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitYearly), 0)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 0, time.Time{})

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	payUrl, updatedTransactions, err := c.NewBusiness()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}

	if len(updatedTransactions) != len(transactions) {
		t.Fatalf("expected: %d transactions got: %d", len(transactions), len(updatedTransactions))
	}

	if payUrl == "" {
		t.Fatalf("expected: non-empty payUrl got: %s", payUrl)
	}

	if payUrl != "www.dev.wopta.it/local-00" {
		t.Fatalf("expected: www.dev.wopta.it/local-00, got: %s", payUrl)
	}

	for index, tr := range updatedTransactions {
		if tr.ScheduleDate == "" {
			t.Fatalf("expected: non-empty ScheduleDate")
		}

		if tr.ProviderId != "local" {
			t.Fatalf("expected: %s ProviderName got: %s", "local", tr.ProviderId)
		}

		if tr.UpdateDate.Equal(transactions[index].UpdateDate) {
			t.Fatalf("expected: %s UpdateDate got: %s", transactions[index].UpdateDate, tr.UpdateDate)
		}

		if tr.UserToken == "" {
			t.Fatalf("expected: non-empty UserToken")
		}

		if tr.ProviderName != policy.Payment {
			t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", tr.ProviderName, policy.Payment)
		}
	}
}

func newBusinessFabrickMonthly(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 0)
	product := getProduct()
	transactions := getTransactions(12, models.FabrickPaymentProvider, 0, time.Time{})

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	payUrl, updatedTransactions, err := c.NewBusiness()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}

	if len(updatedTransactions) != len(transactions) {
		t.Fatalf("expected: %d transactions got: %d", len(transactions), len(updatedTransactions))
	}

	if payUrl == "" {
		t.Fatalf("expected: non-empty payUrl got: %s", payUrl)
	}

	if payUrl != "www.dev.wopta.it/local-00" {
		t.Fatalf("expected: www.dev.wopta.it/local-00, got: %s", payUrl)
	}

	for index, tr := range updatedTransactions {
		if tr.ScheduleDate == "" {
			t.Fatalf("expected: non-empty ScheduleDate")
		}

		if tr.ProviderId != "local" {
			t.Fatalf("expected: %s ProviderName got: %s", "local", tr.ProviderId)
		}

		if tr.UpdateDate.Equal(transactions[index].UpdateDate) {
			t.Fatalf("expected: %s UpdateDate got: %s", transactions[index].UpdateDate, tr.UpdateDate)
		}

		if tr.UserToken == "" {
			t.Fatalf("expected: non-empty UserToken")
		}

		if tr.ProviderName != policy.Payment {
			t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", tr.ProviderName, policy.Payment)
		}
	}
}

func newBusinessManualRemittance(t *testing.T) {
	policy := getPolicy(models.ManualPaymentProvider, models.PaymentModeSingle, string(models.PaySplitYearly), 0)
	product := getProduct()
	transactions := getTransactions(1, models.ManualPaymentProvider, 0, time.Time{})

	c := client.NewClient(models.ManualPaymentProvider, policy, product, transactions, false, "")
	payUrl, updatedTransactions, err := c.NewBusiness()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}

	if len(updatedTransactions) != len(transactions) {
		t.Fatalf("expected: %d transactions got: %d", len(transactions), len(updatedTransactions))
	}

	if payUrl != "" {
		t.Fatalf("expected: empty payUrl got: %s", payUrl)
	}

	for index, tr := range updatedTransactions {
		if !tr.IsPay {
			t.Fatalf("expected: %v IsPay got: %t", false, tr.IsPay)
		}

		if tr.Status != models.TransactionStatusPay {
			t.Fatalf("expected: %s Status got: %s", models.TransactionStatusPay, tr.Status)
		}

		if len(tr.StatusHistory) != 2 && tr.StatusHistory[1] != models.TransactionStatusPay {
			t.Fatalf("expected: %s StatusHistory[1] got: %s", models.TransactionStatusPay, tr.StatusHistory[1])
		}
		if tr.PayDate.IsZero() {
			t.Fatalf("expected: non-zero PayDate got: %s", tr.PayDate)
		}

		if tr.TransactionDate.IsZero() {
			t.Fatalf("expected: non-zero TransactionDate got: %s", tr.TransactionDate)
		}

		if tr.PaymentMethod != models.PayMethodRemittance {
			t.Fatalf("expected: %s PayMethod got: %s", models.PayMethodRemittance, tr.PaymentMethod)
		}

		if tr.UpdateDate.Equal(transactions[index].UpdateDate) {
			t.Fatalf("expected: %s UpdateDate got: %s", transactions[index].UpdateDate, tr.UpdateDate)
		}

		if tr.ProviderName != policy.Payment {
			t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", tr.ProviderName, policy.Payment)
		}
	}
}

func renewManualRemittance(t *testing.T) {
	policy := getPolicy(models.ManualPaymentProvider, models.PaymentModeSingle, string(models.PaySplitYearly), 1)
	product := getProduct()
	transactions := getTransactions(1, models.ManualPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	c := client.NewClient(models.ManualPaymentProvider, policy, product, transactions, false, "")
	_, _, updatedTransactions, err := c.Renew()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if len(updatedTransactions) != 1 {
		t.Fatalf("expected: 1 got: %d", len(updatedTransactions))
	}
	if updatedTransactions[0].IsPay {
		t.Fatalf("isPay error - expected: false got: %v", updatedTransactions[0].IsPay)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func renewFabrickMonthlyWithoutMandate(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 1)
	product := getProduct()
	transactions := getTransactions(12, models.FabrickPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	payUrl, _, updatedTransactions, err := c.Renew()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if updatedTransactions[0].IsPay {
		t.Fatalf("isPay error - expected: false got: %v", updatedTransactions[0].IsPay)
	}
	if updatedTransactions[0].PayUrl != payUrl {
		t.Fatalf("payUrl error - expected: %s  got: %s", updatedTransactions[0].PayUrl, payUrl)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func renewFabrickMonthlyWithMandate(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 1)
	product := getProduct()
	transactions := getTransactions(12, models.FabrickPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, true, "user-has-token")
	payUrl, _, updatedTransactions, err := c.Renew()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if updatedTransactions[0].IsPay {
		t.Fatalf("isPay error - expected: false got: %v", updatedTransactions[0].IsPay)
	}
	if payUrl != "" {
		t.Fatalf("payUrl error - expected: \"\" got: %s", payUrl)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func renewFabrickYearlyRecurrentWithoutMandate(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitYearly), 1)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	payUrl, _, updatedTransactions, err := c.Renew()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if updatedTransactions[0].IsPay {
		t.Fatalf("isPay error - expected: false got: %v", updatedTransactions[0].IsPay)
	}
	if payUrl != updatedTransactions[0].PayUrl {
		t.Fatalf("payUrl error - expected: %s got: %s", updatedTransactions[0].PayUrl, payUrl)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func renewFabrickYearlyRecurrentWithMandate(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitYearly), 1)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, true, "user-has-token")
	payUrl, _, updatedTransactions, err := c.Renew()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if updatedTransactions[0].IsPay {
		t.Fatalf("isPay error - expected: false got: %v", updatedTransactions[0].IsPay)
	}
	if payUrl != "" {
		t.Fatalf("payUrl error - expected: \"\" got: %s", payUrl)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func renewFabrickYearlySingle(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeSingle, string(models.PaySplitYearly), 1)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	payUrl, _, updatedTransactions, err := c.Renew()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if updatedTransactions[0].IsPay {
		t.Fatalf("isPay error - expected: false got: %v", updatedTransactions[0].IsPay)
	}
	if payUrl == "" {
		t.Fatalf("payUrl error - expected: \"\" got: %s", payUrl)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func refreshFabrickMonthly(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 0)
	product := getProduct()
	transactions := getTransactions(12, models.FabrickPaymentProvider, 0, time.Time{})

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions[5:], false, "")
	payUrl, updatedTransactions, err := c.Update()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if payUrl != "www.dev.wopta.it/local-05" {
		t.Fatalf("wrong payUrl - expected: www.dev.wopta.it/local-05, got: %s", payUrl)
	}
	if updatedTransactions[0].ScheduleDate != globalDate.AddDate(0, 5, 0).Format(time.DateOnly) {
		t.Fatalf("wrong schedule date - expected: %s - got: %s", globalDate.AddDate(0, 5, 0).Format(time.DateOnly), updatedTransactions[0].ScheduleDate)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func refreshFabrickYearlySingle(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeSingle, string(models.PaySplitYearly), 0)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 0, time.Time{})

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	payUrl, updatedTransactions, err := c.Update()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if payUrl != "www.dev.wopta.it/local-00" {
		t.Fatalf("wrong payUrl - expected: www.dev.wopta.it/local-00, got: %s", payUrl)
	}
	if updatedTransactions[0].ScheduleDate != globalDate.Format(time.DateOnly) {
		t.Fatalf("wrong schedule date - expected: %s - got: %s", globalDate.Format(time.DateOnly), updatedTransactions[0].ScheduleDate)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func refreshFabrickYearlyRecurrent(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitYearly), 0)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 0, time.Time{})

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	payUrl, updatedTransactions, err := c.Update()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if payUrl != "www.dev.wopta.it/local-00" {
		t.Fatalf("wrong payUrl - expected: www.dev.wopta.it/local-00, got: %s", payUrl)
	}
	if updatedTransactions[0].ScheduleDate != globalDate.Format(time.DateOnly) {
		t.Fatalf("wrong schedule date - expected: %s - got: %s", globalDate.Format(time.DateOnly), updatedTransactions[0].ScheduleDate)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func refreshFabrickMonthlyRenewed(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 1)
	product := getProduct()
	transactions := getTransactions(12, models.FabrickPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions[5:], false, "")
	payUrl, updatedTransactions, err := c.Update()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if payUrl != "www.dev.wopta.it/local-05" {
		t.Fatalf("wrong payUrl - expected: www.dev.wopta.it/local-05, got: %s", payUrl)
	}
	if updatedTransactions[0].ScheduleDate != globalDate.AddDate(1, 5, 0).Format(time.DateOnly) {
		t.Fatalf("wrong schedule date - expected: %s - got: %s", globalDate.AddDate(1, 5, 0).Format(time.DateOnly), updatedTransactions[0].ScheduleDate)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func refreshFabrickYearlyRecurrentRenewed(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitYearly), 1)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	payUrl, updatedTransactions, err := c.Update()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if payUrl != "www.dev.wopta.it/local-00" {
		t.Fatalf("wrong payUrl - expected: www.dev.wopta.it/local-00, got: %s", payUrl)
	}
	if updatedTransactions[0].ScheduleDate != globalDate.AddDate(1, 0, 0).Format(time.DateOnly) {
		t.Fatalf("wrong schedule date - expected: %s - got: %s", globalDate.AddDate(1, 0, 0).Format(time.DateOnly), updatedTransactions[0].ScheduleDate)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

func refreshFabrickYearlySingleRenewed(t *testing.T) {
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeSingle, string(models.PaySplitYearly), 1)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	c := client.NewClient(models.FabrickPaymentProvider, policy, product, transactions, false, "")
	payUrl, updatedTransactions, err := c.Update()
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if payUrl != "www.dev.wopta.it/local-00" {
		t.Fatalf("wrong payUrl - expected: www.dev.wopta.it/local-00, got: %s", payUrl)
	}
	if updatedTransactions[0].ScheduleDate != globalDate.AddDate(1, 0, 0).Format(time.DateOnly) {
		t.Fatalf("wrong schedule date - expected: %s - got: %s", globalDate.AddDate(1, 0, 0).Format(time.DateOnly), updatedTransactions[0].ScheduleDate)
	}
	if updatedTransactions[0].ProviderName != policy.Payment {
		t.Fatalf("expected providers to match, got transaction: '%s' and policy: '%s'", updatedTransactions[0].ProviderName, policy.Payment)
	}
}

var globalDate = time.Date(2023, 03, 14, 0, 0, 0, 0, time.UTC)

func getPolicy(paymentProvider, paymentMode, paymentSplit string, annuity int) models.Policy {
	return models.Policy{
		Payment:           paymentProvider,
		PaymentMode:       paymentMode,
		PaymentSplit:      paymentSplit,
		PriceGross:        100,
		PriceNett:         89.2,
		PriceGrossMonthly: 8.33,
		PriceNettMonthly:  7.43,
		Annuity:           annuity,
		StartDate:         globalDate,
	}
}

func getProduct() models.Product {
	return models.Product{
		PaymentProviders: []models.PaymentProvider{
			{
				Name:  models.FabrickPaymentProvider,
				Flows: []string{models.ProviderMgaFlow},
				Configs: []models.PaymentConfig{
					{
						Rate:    string(models.PaySplitMonthly),
						Methods: []string{models.PayMethodCard, models.PayMethodSdd},
						Mode:    models.PaymentModeRecurrent,
					},
					{
						Rate:    string(models.PaySplitYearly),
						Methods: []string{models.PayMethodCard, models.PayMethodSdd},
						Mode:    models.PaymentModeRecurrent,
					},
					{
						Rate:    string(models.PaySplitYearly),
						Methods: []string{models.PayMethodCard, "fbkr2p"},
						Mode:    models.PaymentModeSingle,
					},
				},
			},
			{
				Name:  models.ManualPaymentProvider,
				Flows: []string{models.RemittanceMgaFlow},
				Configs: []models.PaymentConfig{
					{
						Rate:    string(models.PaySplitYearly),
						Methods: []string{models.PayMethodRemittance},
						Mode:    models.PaymentModeSingle,
					},
				},
			},
		},
	}
}

func getTransactions(numTransactions int, providerName string, annuity int, startDate time.Time) []models.Transaction {
	transactions := make([]models.Transaction, 0)
	if startDate.IsZero() {
		startDate = globalDate
	}

	if numTransactions == 0 {
		return transactions
	}

	amount := lib.RoundFloat(float64(100/numTransactions), 2)
	amountNet := lib.RoundFloat(float64(80/numTransactions), 2)
	now := time.Now().UTC()

	for i := 0; i < numTransactions; i++ {
		transactions = append(transactions, models.Transaction{
			Amount:         amount,
			AmountNet:      amountNet,
			Commissions:    15,
			Status:         models.TransactionStatusToPay,
			PolicyName:     "local",
			Name:           "Test Test",
			ScheduleDate:   startDate.AddDate(0, i, 0).Format(time.DateOnly),
			ExpirationDate: lib.AddMonths(now, 18).Format(time.DateOnly),
			Uid:            fmt.Sprintf("local-%02d", i),
			PolicyUid:      "fjn32onw",
			Company:        "local",
			NumberCompany:  "11111",
			StatusHistory:  []string{models.TransactionStatusToPay},
			ProviderName:   providerName,
			EffectiveDate:  lib.AddMonths(startDate, i),
			CreationDate:   startDate,
			UpdateDate:     startDate,
			Annuity:        annuity,
		})
	}
	return transactions
}
