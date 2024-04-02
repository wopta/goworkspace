package transaction_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/transaction"
	"math/rand"
	"os"
	"testing"
	"time"
)

func getPolicy(paymentSplit, paymentProvider, contractorName, contractorSurname string, priceGross, priceNet,
	priceGrossMonthly, priceNetMonthly float64, startDate, endDate time.Time) models.Policy {
	return models.Policy{
		Uid:  uuid.New().String(),
		Name: models.LifeProduct,
		Contractor: models.Contractor{
			Name:    contractorName,
			Surname: contractorSurname,
		},
		Company:           models.AxaCompany,
		CodeCompany:       fmt.Sprintf("%07d", rand.Intn(100-1)+1),
		Payment:           paymentProvider,
		PaymentSplit:      paymentSplit,
		PriceGross:        priceGross,
		PriceNett:         priceNet,
		PriceGrossMonthly: priceGrossMonthly,
		PriceNettMonthly:  priceNetMonthly,
		StartDate:         startDate,
		EndDate:           endDate,
	}
}

func TestCreateTransactionsMonthly(t *testing.T) {
	now := time.Now().UTC()
	policy := getPolicy(string(models.PaySplitMonthly), models.FabrickPaymentProvider, "Test", "Test", 100, 89.2,
		8.33, 7.43, now, now.AddDate(20, 0, 0))

	os.Setenv("env", "local-test")
	mgaProduct := product.GetProductV2(models.LifeProduct, models.ProductV2, models.MgaChannel, nil, nil)

	transactions := transaction.CreateTransactions(policy, *mgaProduct, func() string { return "aaaaa" })

	if len(transactions) < 12 {
		t.Fatalf("expected: %02d transactions got: %02d", 12, len(transactions))
	}

	for index, tr := range transactions {
		if tr.PolicyName != models.LifeProduct {
			t.Fatalf("expected: %s product got: %s", models.LifeProduct, tr.PolicyName)
		}

		expectedName := lib.TrimSpace(fmt.Sprintf("%s %s", policy.Contractor.Name, policy.Contractor.Surname))
		if tr.Name != expectedName {
			t.Fatalf("expected: %s contractor name got: %s", expectedName, tr.Name)
		}

		if tr.Company != models.AxaCompany {
			t.Fatalf("expected: %s product got: %s", models.AxaCompany, tr.Company)
		}

		if tr.NumberCompany != policy.CodeCompany {
			t.Fatalf("expected: %s codeCompany got: %s", policy.CodeCompany, tr.NumberCompany)
		}

		if tr.ProviderName != policy.Payment {
			t.Fatalf("expected: %s provider name got: %s", policy.Payment, tr.ProviderName)
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

		expectedScheduleDate := policy.StartDate.AddDate(0, index, 0).Format(time.DateOnly)
		if tr.ScheduleDate != expectedScheduleDate {
			t.Fatalf("expected: %s schedule date got: %s", expectedScheduleDate, tr.ScheduleDate)
		}

		expectedExpirationDate := policy.StartDate.AddDate(10, index, 0).Format(time.DateOnly)
		if tr.ExpirationDate != expectedExpirationDate {
			t.Fatalf("expected: %s schedule date got: %s", expectedExpirationDate, tr.ExpirationDate)
		}

		expectedEffectiveDate := policy.StartDate.AddDate(0, index, 0)
		if tr.EffectiveDate != expectedEffectiveDate {
			t.Fatalf("expected: %s schedule date got: %s", expectedEffectiveDate.String(), tr.EffectiveDate.String())
		}
	}
}
