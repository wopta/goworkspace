package renew

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wopta/goworkspace/lib"
)

var (
	paramsHierarchy = []map[string][]string{
		{"codeCompany": []string{"codeCompany"}},

		{"proposalNumber": []string{"proposalNumber"}},

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

	paramsWhereClause = map[string]string{
		"codeCompany": "(codeCompany = @codeCompany)",

		"proposalNumber": "(proposalNumber = @proposalNumber)",

		"insuredFiscalCode": "(JSON_VALUE(**tableAlias**.data, '$.assets[0].person.fiscalCode') = @insuredFiscalCode)",

		"contractorName":    "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.name')), LOWER(@contractorName)))",
		"contractorSurname": "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.surname')), LOWER(@contractorSurname)))",

		"startDateFrom": "(startDate >= @startDateFrom)",
		"startDateTo":   "(startDate <= @startDateTo)",
		"company":       "(company = LOWER(@company))",
		"product":       "(product = LOWER(@product))",
		"producerCode":  "(producerCode = @producerCode)",
		"paid":          "((isDeleted = false OR isDeleted IS NULL) AND (isPay = true))",
		"unpaid":        "((isDeleted = false OR isDeleted IS NULL) AND (isPay = false))",
		"recurrent":     "((isDeleted = false OR isDeleted IS NULL) AND (hasMandate = true))",
		"notRecurrent":  "((isDeleted = false OR isDeleted IS NULL) AND (hasMandate = false OR hasMandate IS NULL))",
	}

	orClauses = []string{"status", "payment"}
)

type QueryBuilder interface {
	BuildQuery(map[string]string) (string, map[string]interface{})
}

type BigQueryQueryBuilder struct {
	tableName  string
	tableAlias string
}

func NewBigQueryQueryBuilder(tableName, tableAlias string) BigQueryQueryBuilder {
	return BigQueryQueryBuilder{
		tableName:  tableName,
		tableAlias: tableAlias,
	}
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

func (qb *BigQueryQueryBuilder) BuildQuery(params map[string]string) (string, map[string]interface{}) {
	var (
		err           error
		limit         = 10
		rawQuery      bytes.Buffer
		queryParams   = make(map[string]interface{})
		whereClauses  = make([]string, 0)
		allowedParams = make([]string, 0)
	)

	// TODO: handle table alias
	rawQuery.WriteString("SELECT **tableAlias**.uid, **tableAlias**.name AS productName, **tableAlias**.codeCompany, CAST(**tableAlias**.proposalNumber AS INT64) AS proposalNumber, " +
		"**tableAlias**.nameDesc,**tableAlias**.status, RTRIM(COALESCE(JSON_VALUE(**tableAlias**.data, '$.contractor.name'), '') || ' ' || " +
		"COALESCE(JSON_VALUE(**tableAlias**.data, '$.contractor.surname'), '')) AS contractor, " +
		"**tableAlias**.priceGross AS price, **tableAlias**.priceGrossMonthly AS priceMonthly, COALESCE(nn.name, '') AS producer, COALESCE(**tableAlias**.producerCode, '') AS producerCode, **tableAlias**.startDate, " +
		"**tableAlias**.endDate, **tableAlias**.paymentSplit " +
		"FROM `wopta.**tableName**` **tableAlias** " +
		"LEFT JOIN `wopta.networkNodesView` nn ON nn.uid = **tableAlias**.producerUid " +
		"WHERE ")

	if val, ok := params["limit"]; ok {
		limit, err = strconv.Atoi(val)
		if err != nil {
			log.Printf("Failed to parse limit: %v", err)
			return "", nil
		}
		delete(params, "limit")
	}

	allowedParams = qb.getAllowedParams(params)
	if allowedParams == nil {
		// TODO: handle error
		return "", nil
	}

	filteredParams := qb.filterParams(params, allowedParams)

	for _, paramKey := range allowedParams {
		if val, ok := filteredParams[paramKey]; ok && val != "" {
			if lib.SliceContains(orClauses, paramKey) {
				tmpWhereClauses := make([]string, 0)
				statusList := strings.Split(filteredParams[paramKey], ",")
				for _, status := range statusList {
					tmpWhereClauses = append(tmpWhereClauses, paramsWhereClause[status])
				}
				whereClauses = append(whereClauses, "("+strings.Join(tmpWhereClauses, " OR ")+")")
			} else {
				var value interface{} = val
				whereClauses = append(whereClauses, paramsWhereClause[paramKey])
				if paramKey == "proposalNumber" {
					parsedValue, err := strconv.ParseInt(filteredParams[paramKey], 10, 64)
					if err != nil {
						log.Printf("Failed to parse proposalNumber: %v", err)
					}
					value = parsedValue
				}
				queryParams[paramKey] = value
			}
		}
	}

	rawQuery.WriteString(strings.Join(whereClauses, " AND "))
	rawQuery.WriteString(fmt.Sprintf(" LIMIT %d", limit))

	query := strings.ReplaceAll(rawQuery.String(), "**tableName**", qb.tableName)
	query = strings.ReplaceAll(query, "**tableAlias**", qb.tableAlias)

	return query, queryParams
}
