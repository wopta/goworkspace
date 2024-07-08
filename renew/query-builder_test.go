package renew_test

import (
	"strings"
	"testing"

	"github.com/wopta/goworkspace/renew"
)

func TestQueryBuilder(t *testing.T) {
	qb := renew.NewBigQueryQueryBuilder(func() string {
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
			`(codeCompany = "@test")`,
		},
		{
			"fiscalCode overcome third-level parameters",
			map[string]string{
				"insuredFiscalCode": "LLLRRR85E05R94Z330F",
				"producerCode":      "a1b2c3d4",
			},
			`(JSON_VALUE(p.data, '$.assets[0].person.fiscalCode') = "@test")`,
		},
		{
			"paid renew policies",
			map[string]string{
				"status": "paid",
			},
			"(((isDeleted = false OR IS NULL) AND (isPay = true)))",
		},
		{
			"not paid renew policies",
			map[string]string{
				"status": "unpaid",
			},
			"(((isDeleted = false OR IS NULL) AND (isPay = false)))",
		},
		{
			"renew policies with mandate active",
			map[string]string{
				"payment": "recurrent",
			},
			"(((isDeleted = false OR IS NULL) AND (hasMandate = true)))",
		},
		{
			"renew policies with mandate non active",
			map[string]string{
				"payment": "notRecurrent",
			},
			"(((isDeleted = false OR IS NULL) AND (hasMandate = false)))",
		},
		{
			"combine third-level parameters",
			map[string]string{
				"producerCode":  "a1b2c3d4",
				"startDateFrom": "2024-07-04",
				"startDateTo":   "2024-07-14",
				"status":        "paid",
				"payment":       "recurrent",
			},
			`(startDate >= "@test") AND (startDate <= "@test") AND (producerCode = "@test") AND (((isDeleted = false OR IS NULL) AND (isPay = true))) AND (((isDeleted = false OR IS NULL) AND (hasMandate = true)))`,
		},
		{
			"combine parameters from differents level",
			map[string]string{
				"producerCode":  "a1b2c3d4",
				"startDateFrom": "2024-07-04",
				"startDateTo":   "2024-07-14",
				"codeCompany":   "100100",
			},
			`(codeCompany = "@test")`,
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
