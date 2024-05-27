package payment

import (
	"os"
	"testing"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

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

	for i := 0; i < numTransactions; i++ {
		transactions = append(transactions, models.Transaction{
			Amount:         amount,
			AmountNet:      amountNet,
			Commissions:    15,
			Status:         models.TransactionStatusToPay,
			PolicyName:     "local",
			Name:           "Test Test",
			ScheduleDate:   startDate.AddDate(0, i, 0).Format(time.DateOnly),
			ExpirationDate: startDate.AddDate(10, i, 0).Format(time.DateOnly),
			Uid:            "local",
			PolicyUid:      "fjn32onw",
			Company:        "local",
			NumberCompany:  "11111",
			StatusHistory:  []string{models.TransactionStatusToPay},
			ProviderName:   providerName,
			EffectiveDate:  startDate.AddDate(0, i, 0),
			CreationDate:   startDate,
			UpdateDate:     startDate,
			Annuity:        annuity,
		})
	}
	return transactions
}

func TestControllerInvalidNumTransactions(t *testing.T) {
	os.Setenv("env", "local-test")

	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 0)
	product := getProduct()
	transactions := getTransactions(0, models.FabrickPaymentProvider, 0, time.Time{})

	_, updatedTransactions, err := Controller(policy, product, transactions, false, "")
	if err == nil {
		t.Fatalf("expected: %02d transactions got: %02d", 0, len(updatedTransactions))
	}
}

func TestControllerInvalidPaymentConfiguration(t *testing.T) {
	os.Setenv("env", "local-test")

	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeSingle, string(models.PaySplitMonthly), 0)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 0, time.Time{})

	_, _, err := Controller(policy, product, transactions, false, "")
	if err == nil {
		t.Fatalf("expected: non-nil error")
	}
}

func TestControllerFabrickYearlySingle(t *testing.T) {
	os.Setenv("env", "local-test")

	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeSingle, string(models.PaySplitYearly), 0)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 0, time.Time{})

	payUrl, updatedTransactions, err := Controller(policy, product, transactions, false, "")
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}

	if len(updatedTransactions) != len(transactions) {
		t.Fatalf("expected: %d transactions got: %d", len(transactions), len(updatedTransactions))
	}

	if payUrl == "" {
		t.Fatalf("expected: non-empty payUrl got: %s", payUrl)
	}

	if payUrl != "www.dev.wopta.it" {
		t.Fatalf("expected: www.dev.wopta.it, got: %s", payUrl)
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
	}
}

func TestControllerFabrickYearlyRecurrent(t *testing.T) {
	os.Setenv("env", "local-test")

	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitYearly), 0)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 0, time.Time{})

	payUrl, updatedTransactions, err := Controller(policy, product, transactions, false, "")
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}

	if len(updatedTransactions) != len(transactions) {
		t.Fatalf("expected: %d transactions got: %d", len(transactions), len(updatedTransactions))
	}

	if payUrl == "" {
		t.Fatalf("expected: non-empty payUrl got: %s", payUrl)
	}

	if payUrl != "www.dev.wopta.it" {
		t.Fatalf("expected: www.dev.wopta.it, got: %s", payUrl)
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
	}
}

func TestControllerFabrickMonthly(t *testing.T) {
	os.Setenv("env", "local-test")

	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 0)
	product := getProduct()
	transactions := getTransactions(12, models.FabrickPaymentProvider, 0, time.Time{})

	payUrl, updatedTransactions, err := Controller(policy, product, transactions, false, "")
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}

	if len(updatedTransactions) != len(transactions) {
		t.Fatalf("expected: %d transactions got: %d", len(transactions), len(updatedTransactions))
	}

	if payUrl == "" {
		t.Fatalf("expected: non-empty payUrl got: %s", payUrl)
	}

	if payUrl != "www.dev.wopta.it" {
		t.Fatalf("expected: www.dev.wopta.it, got: %s", payUrl)
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
	}
}

func TestControllerRemittance(t *testing.T) {
	os.Setenv("env", "local-test")

	policy := getPolicy(models.ManualPaymentProvider, models.PaymentModeSingle, string(models.PaySplitYearly), 0)
	product := getProduct()
	transactions := getTransactions(1, models.ManualPaymentProvider, 0, time.Time{})

	payUrl, updatedTransactions, err := Controller(policy, product, transactions, false, "")
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

	}
}

func TestControllerReuseCustomerId(t *testing.T) {
	os.Setenv("env", "local-test")

	customerId := "a-random-custumer-id"
	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 0)
	product := getProduct()
	transactions := getTransactions(12, models.FabrickPaymentProvider, 0, time.Time{})

	_, updatedTransactions, err := Controller(policy, product, transactions, false, customerId)
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	for _, trn := range updatedTransactions {
		if trn.UserToken != customerId {
			t.Fatalf("mismatched customerID. Expected: %s - got: %s", customerId, trn.UserToken)
		}
	}

}

func TestControllerRenewYearly(t *testing.T) {
	os.Setenv("env", "local-test")

	policy := getPolicy(models.ManualPaymentProvider, models.PaymentModeSingle, string(models.PaySplitYearly), 1)
	product := getProduct()
	transactions := getTransactions(1, models.ManualPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	_, updatedTransactions, err := Controller(policy, product, transactions, false, "")
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if len(updatedTransactions) != 1 {
		t.Fatalf("expected: 1 got: %d", len(updatedTransactions))
	}
	if updatedTransactions[0].IsPay {
		t.Fatalf("isPay error - expected: false got: %v", updatedTransactions[0].IsPay)
	}
}

func TestControllerRenewMonthly(t *testing.T) {
	os.Setenv("env", "local-test")

	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 1)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	payUrl, updatedTransactions, err := Controller(policy, product, transactions, false, "")
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if updatedTransactions[0].IsPay {
		t.Fatalf("isPay error - expected: false got: %v", updatedTransactions[0].IsPay)
	}
	if updatedTransactions[0].PayUrl != payUrl {
		t.Fatalf("payUrl error - expected: \"\" got: %s", updatedTransactions[0].PayUrl)
	}
}

func TestControllerRenewMonthlyWithExistingMandate(t *testing.T) {
	os.Setenv("env", "local-test")

	policy := getPolicy(models.FabrickPaymentProvider, models.PaymentModeRecurrent, string(models.PaySplitMonthly), 1)
	product := getProduct()
	transactions := getTransactions(1, models.FabrickPaymentProvider, 1, globalDate.AddDate(1, 0, 0))

	payUrl, updatedTransactions, err := Controller(policy, product, transactions, true, "an-user-token")
	if err != nil {
		t.Fatalf("expected: nil error got: %s", err.Error())
	}
	if updatedTransactions[0].IsPay {
		t.Fatalf("isPay error - expected: false got: %v", updatedTransactions[0].IsPay)
	}
	if payUrl != "" {
		t.Fatalf("payUrl error - expected: \"\" got: %s", payUrl)
	}
}
