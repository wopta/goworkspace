package renew_test

import (
	"strings"
	"testing"

	"github.com/wopta/goworkspace/renew"
)

func TestQueryBuilder(t *testing.T) {
	qb := renew.NewBigQueryQueryBuilder()

	t.Run("codeCompany overcome everything", func(t *testing.T) {
		params := map[string]string{
			"codeCompany":       "100100",
			"insuredFiscalCode": "LLLRRR85E05R94Z330F",
		}

		want := `(JSON_VALUE(p.data, '$.codeCompany') = "100100")`

		got := qb.BuildQuery(params)

		if !strings.EqualFold(got, want) {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})

	t.Run("fiscalCode overcome third-level parameters", func(t *testing.T) {
		params := map[string]string{
			"insuredFiscalCode": "LLLRRR85E05R94Z330F",
			"producerCode":      "a1b2c3d4",
		}

		want := `(JSON_VALUE(p.data, '$.assets[0].person.fiscalCode') = "LLLRRR85E05R94Z330F")`

		got := qb.BuildQuery(params)

		if !strings.EqualFold(got, want) {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})

	t.Run("paid renew policies", func(t *testing.T) {
		params := map[string]string{
			"status": "paid",
		}

		want := "(((isDeleted = false OR IS NULL) AND (isPay = true)))"
		got := qb.BuildQuery(params)

		if !strings.EqualFold(got, want) {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})

	t.Run("not paid renew policies", func(t *testing.T) {
		params := map[string]string{
			"status": "unpaid",
		}

		want := "(((isDeleted = false OR IS NULL) AND (isPay = false)))"
		got := qb.BuildQuery(params)

		if !strings.EqualFold(got, want) {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})

	t.Run("renew policies with mandate active", func(t *testing.T) {
		params := map[string]string{
			"payment": "recurrent",
		}

		want := "(((isDeleted = false OR IS NULL) AND (hasMandate = true)))"
		got := qb.BuildQuery(params)

		if !strings.EqualFold(got, want) {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})

	t.Run("renew policies with mandate non active", func(t *testing.T) {
		params := map[string]string{
			"payment": "notRecurrent",
		}

		want := "(((isDeleted = false OR IS NULL) AND (hasMandate = false)))"
		got := qb.BuildQuery(params)

		if !strings.EqualFold(got, want) {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})

	t.Run("combine third-level parameters", func(t *testing.T) {
		params := map[string]string{
			"producerCode":  "a1b2c3d4",
			"startDateFrom": "2024-07-04",
			"startDateTo":   "2024-07-14",
			"status":        "paid",
			"payment":       "recurrent",
		}

		want := `(JSON_VALUE(p.data, '$.startDate') >= "2024-07-04") AND (JSON_VALUE(p.data, '$.startDate') <= "2024-07-14") AND (JSON_VALUE(p.data, '$.producerCode') = "a1b2c3d4") AND (((isDeleted = false OR IS NULL) AND (isPay = true))) AND (((isDeleted = false OR IS NULL) AND (hasMandate = true)))`

		got := qb.BuildQuery(params)

		if !strings.EqualFold(got, want) {
			t.Errorf("expected: %s,\n got: %s", want, got)
		}
	})

	t.Run("combine parameters from differents level", func(t *testing.T) {
		params := map[string]string{
			"producerCode":  "a1b2c3d4",
			"startDateFrom": "2024-07-04",
			"startDateTo":   "2024-07-14",
			"codeCompany":   "100100",
		}

		want := `(JSON_VALUE(p.data, '$.codeCompany') = "100100")`

		got := qb.BuildQuery(params)

		if !strings.EqualFold(got, want) {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})

}
