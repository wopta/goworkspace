package renew

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
)

type bigQueryWhereClauseBuilder func(value string) (string, bigquery.QueryParameter)

var (
	paramsHierarchy = []map[string][]string{
		{"codeCompany": []string{"codeCompany"}},

		{"insuredFiscalCode": []string{"insuredFiscalCode"}},

		{"contractorName": []string{"contractorName", "contractorSurname"}},
		{"contractorSurname": []string{"contractorName", "contractorSurname"}},

		{"startDateFrom": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved", "status", "payment"}},
		{"startDateTo": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved", "status", "payment"}},
		{"company": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved", "status", "payment"}},
		{"product": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved", "status", "payment"}},
		{"producerCode": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved", "status", "payment"}},
		{"status": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved", "status", "payment"}},
		{"payment": []string{"startDateFrom", "startDateTo", "company", "product", "producerCode", "reserved", "status", "payment"}},
	}

	paramsWhereClause = map[string]bigQueryWhereClauseBuilder{
		"codeCompany": func(value string) (string, bigquery.QueryParameter) {
			rnd := generateRandomString()
			return fmt.Sprintf(`(codeCompany = "@%s")`, rnd), bigquery.QueryParameter{
				Name:  rnd,
				Value: value,
			}
		},

		"insuredFiscalCode": func(value string) (string, bigquery.QueryParameter) {
			rnd := generateRandomString()
			return fmt.Sprintf(`(JSON_VALUE(p.data, '$.assets[0].person.fiscalCode' = "@%s")`, rnd), bigquery.QueryParameter{
				Name:  rnd,
				Value: value,
			}
		},

		"contractorName": func(value string) (string, bigquery.QueryParameter) {
			rnd := generateRandomString()
			return fmt.Sprintf(`(REGEXP_CONTAINS(LOWER(JSON_VALUE(p.data, '$.contractor.name')), LOWER(@%s))`, rnd), bigquery.QueryParameter{
				Name:  rnd,
				Value: value,
			}
		},
		"contractorSurname": func(value string) (string, bigquery.QueryParameter) {
			rnd := generateRandomString()
			return fmt.Sprintf(`(REGEXP_CONTAINS(LOWER(JSON_VALUE(p.data, '$.contractor.surname')), LOWER(@%s))`, rnd), bigquery.QueryParameter{
				Name:  rnd,
				Value: value,
			}
		},

		"startDateFrom": func(value string) (string, bigquery.QueryParameter) {
			rnd := generateRandomString()
			return fmt.Sprintf(`(startDate >= "@%s")`, rnd), bigquery.QueryParameter{
				Name:  rnd,
				Value: value,
			}
		},
		"startDateTo": func(value string) (string, bigquery.QueryParameter) {
			rnd := generateRandomString()
			return fmt.Sprintf(`(startDate <= "@%s")`, rnd), bigquery.QueryParameter{
				Name:  rnd,
				Value: value,
			}
		},
		"company": func(value string) (string, bigquery.QueryParameter) {
			rnd := generateRandomString()
			return fmt.Sprintf(`(company = LOWER(@%s))`, rnd), bigquery.QueryParameter{
				Name:  rnd,
				Value: value,
			}
		},
		"product": func(value string) (string, bigquery.QueryParameter) {
			rnd := generateRandomString()
			return fmt.Sprintf(`(product = LOWER(@%s))`, rnd), bigquery.QueryParameter{
				Name:  rnd,
				Value: value,
			}
		},
		"producerCode": func(value string) (string, bigquery.QueryParameter) {
			rnd := generateRandomString()
			return fmt.Sprintf(`(producerCode = "@%s")`, rnd), bigquery.QueryParameter{
				Name:  rnd,
				Value: value,
			}
		},
		"paid": func(value string) (string, bigquery.QueryParameter) {
			return fmt.Sprintf(`((isDeleted = false OR IS NULL) AND (isPay = true))`), bigquery.QueryParameter{}
		},
		"unpaid": func(value string) (string, bigquery.QueryParameter) {
			return fmt.Sprintf(`((isDeleted = false OR IS NULL) AND (isPay = false))`), bigquery.QueryParameter{}
		},
		"recurrent": func(value string) (string, bigquery.QueryParameter) {
			return fmt.Sprintf(`((isDeleted = false OR IS NULL) AND (hasMandate = true))`), bigquery.QueryParameter{}
		},
		"notRecurrent": func(value string) (string, bigquery.QueryParameter) {
			return fmt.Sprintf(`((isDeleted = false OR IS NULL) AND (hasMandate = false))`), bigquery.QueryParameter{}
		},
	}
	orClauses = []string{"status", "payment"}
)

func generateRandomString() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("Failed to generate random string: %v", err)
	}
	return hex.EncodeToString(b)
}

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

func (qb *BigQueryQueryBuilder) BuildQuery(params map[string]string) (string, []bigquery.QueryParameter) {
	var (
		queryParams   = make([]bigquery.QueryParameter, 0)
		whereClauses  = make([]string, 0)
		allowedParams = make([]string, 0)
	)

	allowedParams = qb.getAllowedParams(params)
	if allowedParams == nil {
		// TODO: handle error
		return "", nil
	}

	filteredParams := qb.filterParams(params, allowedParams)

	paramsKeys := lib.GetMapKeys(params)
	for _, paramKey := range allowedParams {
		if !lib.SliceContains(paramsKeys, paramKey) || filteredParams[paramKey] == "" {
			continue
		}
		if lib.SliceContains(orClauses, paramKey) {
			tmpWhereClauses := make([]string, 0)
			statusList := strings.Split(filteredParams[paramKey], ",")
			for _, status := range statusList {
				whereClause, queryParam := paramsWhereClause[status]("")
				tmpWhereClauses = append(tmpWhereClauses, whereClause)
				queryParams = append(queryParams, queryParam)
			}
			whereClauses = append(whereClauses, "("+strings.Join(tmpWhereClauses, " OR ")+")")
		} else {
			whereClause, queryParam := paramsWhereClause[paramKey](filteredParams[paramKey])
			whereClauses = append(whereClauses, whereClause)
			queryParams = append(queryParams, queryParam)
		}
	}
	return strings.Join(whereClauses, " AND "), queryParams
}
