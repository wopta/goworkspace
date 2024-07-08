package renew

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
)

type bigQueryWhereClauseBuilder func(string, func() string) (string, bigquery.QueryParameter)

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
		"codeCompany": "(codeCompany = @%s)",

		"proposalNumber": "(proposalNumber = @%s)",

		"insuredFiscalCode": "(JSON_VALUE(p.data, '$.assets[0].person.fiscalCode') = @%s)",

		"contractorName":    "(REGEXP_CONTAINS(LOWER(JSON_VALUE(p.data, '$.contractor.name')), LOWER(@%s)))",
		"contractorSurname": "(REGEXP_CONTAINS(LOWER(JSON_VALUE(p.data, '$.contractor.surname')), LOWER(@%s)))",

		"startDateFrom": "(startDate >= @%s)",
		"startDateTo":   "(startDate <= @%s)",
		"company":       "(company = LOWER(@%s))",
		"product":       "(product = LOWER(@%s))",
		"producerCode":  "(producerCode = @%s)",
		"paid":          "((isDeleted = false OR isDeleted IS NULL) AND (isPay = true))",
		"unpaid":        "((isDeleted = false OR isDeleted IS NULL) AND (isPay = false))",
		"recurrent":     "((isDeleted = false OR isDeleted IS NULL) AND (hasMandate = true))",
		"notRecurrent":  "((isDeleted = false OR isDeleted IS NULL) AND (hasMandate = false OR hasMandate IS NULL))",
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
	BuildQuery(map[string]string) (string, map[string]interface{})
}

type BigQueryQueryBuilder struct {
	TableName           string
	TableAlias          string
	IdentifierGenerator func() string
}

func NewBigQueryQueryBuilder(tableName, tableAlias string, identifierGenerator func() string) BigQueryQueryBuilder {
	return BigQueryQueryBuilder{
		TableName:           tableName,
		TableAlias:          tableAlias,
		IdentifierGenerator: identifierGenerator,
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
		query         bytes.Buffer
		queryParams   = make(map[string]interface{})
		whereClauses  = make([]string, 0)
		allowedParams = make([]string, 0)
	)

	// TODO: handle table alias
	query.WriteString(fmt.Sprintf("SELECT p.uid, p.name AS productName, p.codeCompany, CAST(p.proposalNumber AS INT64) AS proposalNumber, "+
		"p.nameDesc,p.status, RTRIM(COALESCE(JSON_VALUE(p.data, '$.contractor.name'), '') || ' ' || "+
		"COALESCE(JSON_VALUE(p.data, '$.contractor.surname'), '')) AS contractor, "+
		"p.priceGross AS price, p.priceGrossMonthly AS priceMonthly, COALESCE(nn.name, '') AS producer, COALESCE(p.producerCode, '') AS producerCode, p.startDate, "+
		"p.endDate, p.paymentSplit "+
		"FROM `wopta.%s` %s "+
		"LEFT JOIN `wopta.networkNodesView` nn ON nn.uid = p.producerUid "+
		"WHERE ", qb.TableName, qb.TableAlias))

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
				continue
			}

			identifier := qb.IdentifierGenerator()
			format := paramsWhereClause[paramKey]
			whereClause := fmt.Sprintf(format, identifier)
			whereClauses = append(whereClauses, whereClause)
			if paramKey == "proposalNumber" {
				parsedValue, err := strconv.ParseInt(filteredParams[paramKey], 10, 64)
				if err != nil {
					log.Printf("Failed to parse proposalNumber: %v", err)
				}
				queryParams[identifier] = parsedValue
				continue
			}
			queryParams[identifier] = filteredParams[paramKey]
		}
	}

	query.WriteString(strings.Join(whereClauses, " AND "))
	query.WriteString(fmt.Sprintf(" LIMIT %d", limit))

	return query.String(), queryParams
}
