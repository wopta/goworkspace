package renew

import (
	"fmt"
	"strings"

	"github.com/wopta/goworkspace/lib"
)

var (
	paramsHierarchy = []map[string][]string{
		{"codeCompany": []string{"codeCompany"}},

		{"insuredFiscalCode": []string{"insuredFiscalCode"}},

		{"contractorName": []string{"contractorName", "contractorSurname"}},
		{"contractorSurname": []string{"contractorName", "contractorSurname"}},

		{"startDateFrom": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved"}},
		{"startDateTo": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved"}},
		{"company": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved"}},
		{"product": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved"}},
		{"producerCode": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved"}},
		{"reserved": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved"}},
	}

	paramsQuery = map[string]string{
		"codeCompany": "(JSON_VALUE(p.data, '$.codeCompany') = \"%s\")",

		"insuredFiscalCode": "(JSON_VALUE(p.data, '$.assets[0].person.fiscalCode') = \"%s\")",

		"contractorName":    "(REGEXP_CONTAINS(LOWER(JSON_VALUE(p.data, '$.contractor.name')), LOWER(%s))",
		"contractorSurname": "(REGEXP_CONTAINS(LOWER(JSON_VALUE(p.data, '$.contractor.surname')), LOWER(%s))",

		"startDateFrom": "(JSON_VALUE(p.data, '$.startDate') >= \"%s\")",
		"startDateTo":   "(JSON_VALUE(p.data, '$.startDate') <= \"%s\")",
		"company":       "(JSON_VALUE(p.data, '$.company') = LOWER(\"%s\"))",
		"product":       "(JSON_VALUE(p.data, '$.product') = LOWER(\"%s\"))",
		"producerCode":  "(JSON_VALUE(p.data, '$.producerCode') = \"%s\")",
		"reserved":      "(JSON_VALUE(p.data, '$.reserved') = %s)",
	}
)

type QueryBuilder interface {
	BuildQuery(map[string]string) string
}

type BigQueryQueryBuilder struct{}

func NewBigQueryQueryBuilder() BigQueryQueryBuilder {
	return BigQueryQueryBuilder{}
}

func (qb *BigQueryQueryBuilder) getAllowedParams(params map[string]string) []string {
	paramsKeys := lib.GetMapKeys(params)
	for _, value := range paramsHierarchy {
		for k, v := range value {
			if lib.SliceContains(paramsKeys, k) {
				return v
			}
		}
	}
	return nil
}

func (qb *BigQueryQueryBuilder) filterParams(params map[string]string, allowedParams []string) map[string]string {
	paramsKeys := lib.GetMapKeys(params)
	for _, key := range paramsKeys {
		if !lib.SliceContains(allowedParams, key) {
			delete(params, key)
		}
	}
	return params
}

func (qb *BigQueryQueryBuilder) BuildQuery(params map[string]string) string {
	var (
		queries       = make([]string, 0)
		allowedParams = make([]string, 0)
	)

	allowedParams = qb.getAllowedParams(params)
	if allowedParams == nil {
		// TODO: handle error
		return ""
	}

	filteredParams := qb.filterParams(params, allowedParams)

	for k, v := range filteredParams {
		queries = append(queries, fmt.Sprintf(paramsQuery[k], v))
	}

	return strings.Join(queries, " AND ")
}
