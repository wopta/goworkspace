package transaction_test

import (
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/transaction"
	"os"
	"testing"
	"time"
)

type dateInfo struct {
	ScheduleDate   string
	ExpirationDate string
	EffectiveDate  time.Time
}

func getPolicy(paymentSplit string, startDate, endDate time.Time) models.Policy {
	return models.Policy{
		Uid:  "uuid",
		Name: "productName",
		Contractor: models.Contractor{
			Name:    "Test",
			Surname: "Test",
		},
		Company:           "company",
		CodeCompany:       "1234567",
		Payment:           "paymentProvider",
		PaymentSplit:      paymentSplit,
		PriceGross:        100,
		PriceNett:         89.2,
		PriceGrossMonthly: 8.33,
		PriceNettMonthly:  7.43,
		StartDate:         startDate,
		EndDate:           endDate,
	}
}

func outputGenerator(numOutput int, startDate time.Time) []dateInfo {
	output := make([]dateInfo, 0)

	i := 0
	for i < numOutput {
		effectiveDate := startDate.AddDate(0, i, 0)
		output = append(output, dateInfo{
			ScheduleDate:   effectiveDate.Format(time.DateOnly),
			ExpirationDate: effectiveDate.AddDate(10, 0, 0).Format(time.DateOnly),
			EffectiveDate:  effectiveDate,
		})
		i++
	}
	return output
}

func TestCreateTransactionsMonthly(t *testing.T) {
	startDate := time.Date(2023, 03, 14, 0, 0, 0, 0, time.UTC)
	policy := getPolicy(string(models.PaySplitMonthly), startDate, startDate.AddDate(20, 0, 0))

	output := outputGenerator(12, startDate)

	os.Setenv("env", "local-test")
	mgaProduct := product.GetProductV2(models.LifeProduct, models.ProductV2, models.MgaChannel, nil, nil)

	transactions := transaction.CreateTransactions(policy, *mgaProduct, func() string { return "aaaaa" })

	if len(transactions) < 12 {
		t.Fatalf("expected: %02d transactions got: %02d", 12, len(transactions))
	}

	for index, tr := range transactions {
		if tr.PolicyName != "productName" {
			t.Fatalf("expected: %s product got: %s", "productName", tr.PolicyName)
		}

		if tr.Name != "Test Test" {
			t.Fatalf("expected: %s contractor name got: %s", "Test Test", tr.Name)
		}

		if tr.Company != "company" {
			t.Fatalf("expected: %s product got: %s", "company", tr.Company)
		}

		if tr.NumberCompany != "1234567" {
			t.Fatalf("expected: %s codeCompany got: %s", "1234567", tr.NumberCompany)
		}

		if tr.ProviderName != "paymentProvider" {
			t.Fatalf("expected: %s provider name got: %s", "paymentProvider", tr.ProviderName)
		}

		if tr.Status != models.TransactionStatusToPay {
			t.Fatalf("expected: %s status got: %s", models.TransactionStatusToPay, tr.Status)
		}

		if tr.Amount != policy.PriceGrossMonthly {
			t.Fatalf("expected: %.2f price gross got: %.2f", policy.PriceGrossMonthly, tr.Amount)
		}

		if tr.AmountNet != policy.PriceNettMonthly {
			t.Fatalf("expected: %.2f price net got: %.2f", policy.PriceGrossMonthly, tr.AmountNet)
		}

		expectedScheduleDate := output[index].ScheduleDate
		if tr.ScheduleDate != expectedScheduleDate {
			t.Fatalf("expected: %s schedule date got: %s", expectedScheduleDate, tr.ScheduleDate)
		}

		expectedExpirationDate := output[index].ExpirationDate
		if tr.ExpirationDate != expectedExpirationDate {
			t.Fatalf("expected: %s expiration date got: %s", expectedExpirationDate, tr.ExpirationDate)
		}

		expectedEffectiveDate := output[index].EffectiveDate
		if tr.EffectiveDate != expectedEffectiveDate {
			t.Fatalf("expected: %s effective date got: %s", expectedEffectiveDate.String(), tr.EffectiveDate.String())
		}
	}
}

func TestCreateTransactionsYearly(t *testing.T) {
	startDate := time.Date(2023, 03, 14, 0, 0, 0, 0, time.UTC)
	policy := getPolicy(string(models.PaySplitYearly), startDate, startDate.AddDate(20, 0, 0))

	output := outputGenerator(1, startDate)

	os.Setenv("env", "local-test")
	mgaProduct := product.GetProductV2(models.LifeProduct, models.ProductV2, models.MgaChannel, nil, nil)

	transactions := transaction.CreateTransactions(policy, *mgaProduct, func() string { return "aaaaa" })

	if len(transactions) != 1 {
		t.Fatalf("expected: %02d transactions got: %02d", 1, len(transactions))
	}

	for index, tr := range transactions {
		if tr.PolicyName != "productName" {
			t.Fatalf("expected: %s product got: %s", "productName", tr.PolicyName)
		}

		if tr.Name != "Test Test" {
			t.Fatalf("expected: %s contractor name got: %s", "Test Test", tr.Name)
		}

		if tr.Company != "company" {
			t.Fatalf("expected: %s product got: %s", "company", tr.Company)
		}

		if tr.NumberCompany != "1234567" {
			t.Fatalf("expected: %s codeCompany got: %s", "1234567", tr.NumberCompany)
		}

		if tr.ProviderName != "paymentProvider" {
			t.Fatalf("expected: %s provider name got: %s", "paymentProvider", tr.ProviderName)
		}

		if tr.Status != models.TransactionStatusToPay {
			t.Fatalf("expected: %s status got: %s", models.TransactionStatusToPay, tr.Status)
		}

		if tr.Amount != policy.PriceGross {
			t.Fatalf("expected: %.2f price gross got: %.2f", policy.PriceGrossMonthly, tr.Amount)
		}

		if tr.AmountNet != policy.PriceNett {
			t.Fatalf("expected: %.2f price net got: %.2f", policy.PriceGrossMonthly, tr.AmountNet)
		}

		expectedScheduleDate := output[index].ScheduleDate
		if tr.ScheduleDate != expectedScheduleDate {
			t.Fatalf("expected: %s schedule date got: %s", expectedScheduleDate, tr.ScheduleDate)
		}

		expectedExpirationDate := output[index].ExpirationDate
		if tr.ExpirationDate != expectedExpirationDate {
			t.Fatalf("expected: %s expiration date got: %s", expectedExpirationDate, tr.ExpirationDate)
		}

		expectedEffectiveDate := output[index].EffectiveDate
		if tr.EffectiveDate != expectedEffectiveDate {
			t.Fatalf("expected: %s effective date got: %s", expectedEffectiveDate.String(), tr.EffectiveDate.String())
		}
	}
}

func TestCreateTransactionsSingleInstallment(t *testing.T) {
	startDate := time.Date(2023, 03, 14, 0, 0, 0, 0, time.UTC)
	policy := getPolicy(string(models.PaySplitSingleInstallment), startDate, startDate.AddDate(5, 0, 0))

	output := outputGenerator(1, startDate)

	os.Setenv("env", "local-test")
	mgaProduct := product.GetProductV2(models.LifeProduct, models.ProductV2, models.MgaChannel, nil, nil)

	transactions := transaction.CreateTransactions(policy, *mgaProduct, func() string { return "aaaaa" })

	if len(transactions) != 1 {
		t.Fatalf("expected: %02d transactions got: %02d", 1, len(transactions))
	}

	for index, tr := range transactions {
		if tr.PolicyName != "productName" {
			t.Fatalf("expected: %s product got: %s", "productName", tr.PolicyName)
		}

		if tr.Name != "Test Test" {
			t.Fatalf("expected: %s contractor name got: %s", "Test Test", tr.Name)
		}

		if tr.Company != "company" {
			t.Fatalf("expected: %s product got: %s", "company", tr.Company)
		}

		if tr.NumberCompany != "1234567" {
			t.Fatalf("expected: %s codeCompany got: %s", "1234567", tr.NumberCompany)
		}

		if tr.ProviderName != "paymentProvider" {
			t.Fatalf("expected: %s provider name got: %s", "paymentProvider", tr.ProviderName)
		}

		if tr.Status != models.TransactionStatusToPay {
			t.Fatalf("expected: %s status got: %s", models.TransactionStatusToPay, tr.Status)
		}

		if tr.Amount != policy.PriceGross {
			t.Fatalf("expected: %.2f price gross got: %.2f", policy.PriceGrossMonthly, tr.Amount)
		}

		if tr.AmountNet != policy.PriceNett {
			t.Fatalf("expected: %.2f price net got: %.2f", policy.PriceGrossMonthly, tr.AmountNet)
		}

		expectedScheduleDate := output[index].ScheduleDate
		if tr.ScheduleDate != expectedScheduleDate {
			t.Fatalf("expected: %s schedule date got: %s", expectedScheduleDate, tr.ScheduleDate)
		}

		expectedExpirationDate := output[index].ExpirationDate
		if tr.ExpirationDate != expectedExpirationDate {
			t.Fatalf("expected: %s expiration date got: %s", expectedExpirationDate, tr.ExpirationDate)
		}

		expectedEffectiveDate := output[index].EffectiveDate
		if tr.EffectiveDate != expectedEffectiveDate {
			t.Fatalf("expected: %s effective date got: %s", expectedEffectiveDate.String(), tr.EffectiveDate.String())
		}
	}
}
