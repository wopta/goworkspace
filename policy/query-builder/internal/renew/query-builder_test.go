package renew_test

import (
	"strings"
	"testing"

	"gitlab.dev.wopta.it/goworkspace/policy/query-builder/internal/renew"
)

func TestQueryBuilder(t *testing.T) {
	qb := renew.NewQueryBuilder(func() string {
		return "test"
	})
	var testCases = []struct {
		name   string
		params map[string]string
		want   string
	}{
		{
			"codeCompany overcome everything",
			map[string]string{
				"codeCompany":       "100100",
				"insuredFiscalCode": "LLLRRR85E05R94Z330F",
			},
			"(rp.codeCompany = @test) AND (rp.isDeleted = false OR rp.isDeleted IS NULL) " +
				"ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"fiscalCode overcome third-level parameters",
			map[string]string{
				"insuredFiscalCode": "LLLRRR85E05R94Z330F",
				"producerCode":      "a1b2c3d4",
			},
			"(LOWER(JSON_VALUE(rp.data, '$.assets[0].person.fiscalCode')) = LOWER(@test)) AND " +
				"(rp.isDeleted = false OR rp.isDeleted IS NULL) ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"paid renew policies",
			map[string]string{
				"status": "paid",
			},
			"((rp.isPay = true)) AND (rp.isDeleted = false OR rp.isDeleted IS NULL) " +
				"ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"not paid renew policies",
			map[string]string{
				"status": "unpaid",
			},
			"((rp.isPay = false)) AND (rp.isDeleted = false OR rp.isDeleted IS NULL) " +
				"ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"renew policies with mandate active",
			map[string]string{
				"payment": "recurrent",
			},
			"((rp.hasMandate = true)) AND (rp.isDeleted = false OR rp.isDeleted IS NULL) " +
				"ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"renew policies with mandate non active",
			map[string]string{
				"payment": "notRecurrent",
			},
			"((rp.hasMandate = false OR rp.hasMandate IS NULL)) AND " +
				"(rp.isDeleted = false OR rp.isDeleted IS NULL) ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"combine third-level parameters",
			map[string]string{
				"producerUid":   "a1b2c3d4",
				"startDateFrom": "2024-07-04",
				"startDateTo":   "2024-07-14",
				"status":        "paid",
				"payment":       "recurrent",
			}, "(rp.startDate >= @test) AND (rp.startDate <= @test) AND " +
				"(rp.producerUid IN (@test)) AND ((rp.isPay = true)) AND ((rp.hasMandate = true)) AND " +
				"(rp.isDeleted = false OR rp.isDeleted IS NULL) " +
				"ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"combine parameters from differents level",
			map[string]string{
				"producerCode":  "a1b2c3d4",
				"startDateFrom": "2024-07-04",
				"startDateTo":   "2024-07-14",
				"codeCompany":   "100100",
			}, "(rp.codeCompany = @test) AND (rp.isDeleted = false OR rp.isDeleted IS NULL) " +
				"ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"invalid status",
			map[string]string{
				"status": "invalidValue,unpaid",
			},
			"((rp.isPay = false)) AND (rp.isDeleted = false OR rp.isDeleted IS NULL) " +
				"ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"single producer uid",
			map[string]string{
				"producerUid": "aaaa",
			},
			"(rp.producerUid IN (@test)) AND (rp.isDeleted = false OR rp.isDeleted IS NULL) " +
				"ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"multiple producer uid",
			map[string]string{
				"producerUid": "aaa,bbb",
			},
			"(rp.producerUid IN (@test, @test)) AND (rp.isDeleted = false OR rp.isDeleted IS NULL) " +
				"ORDER BY rp.updateDate DESC LIMIT 10)",
		},
		{
			"limit different from default",
			map[string]string{
				"producerUid": "aaa,bbb",
				"limit":       "50",
			},
			"(rp.producerUid IN (@test, @test)) AND (rp.isDeleted = false OR rp.isDeleted IS NULL) " +
				"ORDER BY rp.updateDate DESC LIMIT 50)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, _, _ := qb.Build(tc.params)

			whereClauses := strings.TrimSpace(strings.Split(got, "WHERE ")[1])

			if !strings.EqualFold(whereClauses, tc.want) {
				t.Errorf("expected: %s, got: %s", tc.want, whereClauses)
			}
		})
	}
}

func TestQueryBuilderFail(t *testing.T) {
	qb := renew.NewQueryBuilder(func() string { return "test" })
	var testCases = []struct {
		name   string
		params map[string]string
		want   string
	}{
		{
			"empty params map",
			map[string]string{},
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, _, _ := qb.Build(tc.params)

			if !strings.EqualFold(got, tc.want) {
				t.Errorf("expected: %s, got: %s", tc.want, got)
			}
		})
	}
}
