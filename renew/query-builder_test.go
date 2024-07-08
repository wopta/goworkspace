package renew_test

import (
	"strings"
	"testing"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/renew"
)

func areMapsEqual(m1, m2 map[string]interface{}) bool {
	if len(m1) != len(m2) {
		return false
	}

	for key, value1 := range m1 {
		if value2, exists := m2[key]; !exists || value1 != value2 {
			return false
		}
	}

	return true
}

func TestQueryBuilder(t *testing.T) {
	qb := renew.NewBigQueryQueryBuilder(lib.RenewPolicyViewCollection, "rp")
	var testCases = []struct {
		name   string
		params map[string]string
		want   struct {
			whereClause string
			params      map[string]interface{}
		}
	}{
		{
			"codeCompany overcome everything",
			map[string]string{
				"codeCompany":       "100100",
				"insuredFiscalCode": "LLLRRR85E05R94Z330F",
			},
			struct {
				whereClause string
				params      map[string]interface{}
			}{
				whereClause: "(codeCompany = @codeCompany) LIMIT 10",
				params: map[string]interface{}{
					"codeCompany": "100100",
				}},
		},
		{
			"fiscalCode overcome third-level parameters",
			map[string]string{
				"insuredFiscalCode": "LLLRRR85E05R94Z330F",
				"producerCode":      "a1b2c3d4",
			},
			struct {
				whereClause string
				params      map[string]interface{}
			}{
				whereClause: "(JSON_VALUE(rp.data, '$.assets[0].person.fiscalCode') = @insuredFiscalCode) LIMIT 10",
				params:      map[string]interface{}{"insuredFiscalCode": "LLLRRR85E05R94Z330F"}},
		},
		{
			"paid renew policies",
			map[string]string{
				"status": "paid",
			},
			struct {
				whereClause string
				params      map[string]interface{}
			}{
				whereClause: "(((isDeleted = false OR isDeleted IS NULL) AND " +
					"(isPay = true))) LIMIT 10",
				params: map[string]interface{}{}},
		},
		{
			"not paid renew policies",
			map[string]string{
				"status": "unpaid",
			},
			struct {
				whereClause string
				params      map[string]interface{}
			}{
				whereClause: "(((isDeleted = false OR isDeleted IS NULL) AND " +
					"(isPay = false))) LIMIT 10",
				params: map[string]interface{}{}},
		},
		{
			"renew policies with mandate active",
			map[string]string{
				"payment": "recurrent",
			},
			struct {
				whereClause string
				params      map[string]interface{}
			}{
				whereClause: "(((isDeleted = false OR isDeleted IS NULL) AND " +
					"(hasMandate = true))) LIMIT 10",
				params: map[string]interface{}{}},
		},
		{
			"renew policies with mandate non active",
			map[string]string{
				"payment": "notRecurrent",
			},
			struct {
				whereClause string
				params      map[string]interface{}
			}{
				whereClause: "(((isDeleted = false OR isDeleted IS NULL) AND " +
					"(hasMandate = false OR hasMandate IS NULL))) LIMIT 10",
				params: map[string]interface{}{}},
		},
		{
			"combine third-level parameters",
			map[string]string{
				"producerCode":  "a1b2c3d4",
				"startDateFrom": "2024-07-04",
				"startDateTo":   "2024-07-14",
				"status":        "paid",
				"payment":       "recurrent",
			}, struct {
				whereClause string
				params      map[string]interface{}
			}{
				whereClause: "(startDate >= @startDateFrom) AND (startDate <= @startDateTo) AND " +
					"(producerCode = @producerCode) AND (((isDeleted = false OR isDeleted IS NULL) AND " +
					"(isPay = true))) AND (((isDeleted = false OR isDeleted IS NULL) AND " +
					"(hasMandate = true))) LIMIT 10",
				params: map[string]interface{}{
					"startDateFrom": "2024-07-04",
					"startDateTo":   "2024-07-14",
					"producerCode":  "a1b2c3d4",
				},
			},
		},
		{
			"combine parameters from differents level",
			map[string]string{
				"producerCode":  "a1b2c3d4",
				"startDateFrom": "2024-07-04",
				"startDateTo":   "2024-07-14",
				"codeCompany":   "100100",
			}, struct {
				whereClause string
				params      map[string]interface{}
			}{
				whereClause: "(codeCompany = @codeCompany) LIMIT 10",
				params:      map[string]interface{}{"codeCompany": "100100"}},
		},
	}

	//less := func(a, b string) bool { return a < b }

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotQuery, gotParams := qb.BuildQuery(tc.params)

			whereClauses := strings.TrimSpace(strings.Split(gotQuery, "WHERE ")[1])

			if !strings.EqualFold(whereClauses, tc.want.whereClause) {
				t.Errorf("expected: %s, got: %s", tc.want.whereClause, whereClauses)
			}

			if !areMapsEqual(gotParams, tc.want.params) {
				t.Errorf("expected: %+v, got: %+v", tc.want.params, gotParams)
			}
		})
	}
}
