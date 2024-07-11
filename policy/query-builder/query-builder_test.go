package query_builder_test

import (
	"strings"
	"testing"

	query_builder "github.com/wopta/goworkspace/policy/query-builder"
)

func TestQueryBuilder(t *testing.T) {
	qb := query_builder.NewBigQueryQueryBuilder("test-collection", "rp", func() string {
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
			"(codeCompany = @test) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"fiscalCode overcome third-level parameters",
			map[string]string{
				"insuredFiscalCode": "LLLRRR85E05R94Z330F",
				"producerCode":      "a1b2c3d4",
			},
			"(JSON_VALUE(rp.data, '$.assets[0].person.fiscalCode') = @test) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"paid renew policies",
			map[string]string{
				"status": "paid",
			},
			"(((isDeleted = false OR isDeleted IS NULL) AND (isPay = true))) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"not paid renew policies",
			map[string]string{
				"status": "unpaid",
			},
			"(((isDeleted = false OR isDeleted IS NULL) AND (isPay = false))) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"renew policies with mandate active",
			map[string]string{
				"payment": "recurrent",
			},
			"(((isDeleted = false OR isDeleted IS NULL) AND (hasMandate = true))) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"renew policies with mandate non active",
			map[string]string{
				"payment": "notRecurrent",
			},
			"(((isDeleted = false OR isDeleted IS NULL) AND (hasMandate = false OR hasMandate IS NULL))) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"combine third-level parameters",
			map[string]string{
				"producerUid":   "a1b2c3d4",
				"startDateFrom": "2024-07-04",
				"startDateTo":   "2024-07-14",
				"status":        "paid",
				"payment":       "recurrent",
			}, "(startDate >= @test) AND (startDate <= @test) AND " +
				"(producerUid IN ('@test')) AND (((isDeleted = false OR isDeleted IS NULL) AND " +
				"(isPay = true))) AND (((isDeleted = false OR isDeleted IS NULL) AND " +
				"(hasMandate = true))) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"combine parameters from differents level",
			map[string]string{
				"producerCode":  "a1b2c3d4",
				"startDateFrom": "2024-07-04",
				"startDateTo":   "2024-07-14",
				"codeCompany":   "100100",
			}, "(codeCompany = @test) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"invalid status",
			map[string]string{
				"status": "invalidValue,unpaid",
			},
			"(((isDeleted = false OR isDeleted IS NULL) AND (isPay = false))) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"single producer uid",
			map[string]string{
				"producerUid": "aaaa",
			},
			"(producerUid IN ('@test')) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"multiple producer uid",
			map[string]string{
				"producerUid": "aaa,bbb",
			},
			"(producerUid IN ('@test', '@test')) ORDER BY rp.updateDate DESC LIMIT 10",
		},
		{
			"limit different from default",
			map[string]string{
				"producerUid": "aaa,bbb",
				"limit":       "50",
			},
			"(producerUid IN ('@test', '@test')) ORDER BY rp.updateDate DESC LIMIT 50",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := qb.BuildQuery(tc.params)

			whereClauses := strings.TrimSpace(strings.Split(got, "WHERE ")[1])

			if !strings.EqualFold(whereClauses, tc.want) {
				t.Errorf("expected: %s, got: %s", tc.want, whereClauses)
			}
		})
	}
}

func TestQueryBuilderFail(t *testing.T) {
	qb := query_builder.NewBigQueryQueryBuilder("test-collection", "rp", func() string {
		return "test"
	})
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
			got, _ := qb.BuildQuery(tc.params)

			if !strings.EqualFold(got, tc.want) {
				t.Errorf("expected: %s, got: %s", tc.want, got)
			}
		})
	}
}
