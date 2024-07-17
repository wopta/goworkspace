package query_builder

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
)

var (
	paramsHierarchy = []map[string][]string{
		{"codeCompany": []string{"codeCompany"}},

		{"proposalNumber": []string{"proposalNumber"}},

		{"insuredFiscalCode": []string{"insuredFiscalCode"}},

		{"contractorName": []string{"contractorName", "contractorSurname"}},
		{"contractorSurname": []string{"contractorName", "contractorSurname"}},

		{"startDateFrom": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"startDateTo": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"company": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"product": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"producerUid": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"status": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"payment": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"renewMonth": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
	}

	paramsWhereClause = map[string]string{
		"codeCompany": "(**tableAlias**.codeCompany = @%s)",

		"proposalNumber": "(**tableAlias**.proposalNumber = CAST(@%s AS INTEGER))",

		"insuredFiscalCode": "(JSON_VALUE(**tableAlias**.data, '$.assets[0].person.fiscalCode') = @%s)",

		"contractorName":    "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.name')), LOWER(@%s)))",
		"contractorSurname": "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.surname')), LOWER(@%s)))",

		"startDateFrom": "(**tableAlias**.startDate >= @%s)",
		"startDateTo":   "(**tableAlias**.startDate <= @%s)",
		"company":       "(**tableAlias**.company = LOWER(@%s))",
		"product":       "(**tableAlias**.name = LOWER(@%s))",
		"producerUid":   "(**tableAlias**.producerUid IN (%s))",
		"renewMonth":    "((**tableAlias**.isDeleted = false OR **tableAlias**.isDeleted IS NULL) AND (EXTRACT(MONTH FROM **tableAlias**.startDate) = CAST(@%s AS INTEGER)))",
		"paid":          "((**tableAlias**.isDeleted = false OR **tableAlias**.isDeleted IS NULL) AND (**tableAlias**.isPay = true))",
		"unpaid":        "((**tableAlias**.isDeleted = false OR **tableAlias**.isDeleted IS NULL) AND (**tableAlias**.isPay = false))",
		"recurrent":     "((**tableAlias**.isDeleted = false OR **tableAlias**.isDeleted IS NULL) AND (**tableAlias**.hasMandate = true))",
		"notRecurrent":  "((**tableAlias**.isDeleted = false OR **tableAlias**.isDeleted IS NULL) AND (**tableAlias**.hasMandate = false OR **tableAlias**.hasMandate IS NULL))",
	}

	orClausesKeys = []string{"status", "payment"}
)

type QueryBuilder interface {
	BuildQuery(map[string]string) (string, map[string]interface{})
}

type BigQueryQueryBuilder struct {
	tableName       string
	tableAlias      string
	randomGenerator func() string
}

func NewBigQueryQueryBuilder(tableName, tableAlias string, randomGenerator func() string) BigQueryQueryBuilder {
	if randomGenerator == nil {
		randomGenerator = func() string {
			var (
				letters  = []rune("abcdefghijklmnopqrstuvwxyz")
				alphanum = []rune("123456789abcdefghijklmnopqrstuvwxyz")
			)
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

			s := make([]rune, 12)
			s[0] = letters[rnd.Intn(len(letters))]
			for i := range s[1:] {
				s[i+1] = alphanum[rnd.Intn(len(alphanum))]
			}
			return string(s)
		}
	}
	return BigQueryQueryBuilder{
		tableName:       tableName,
		tableAlias:      tableAlias,
		randomGenerator: randomGenerator,
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

func (qb *BigQueryQueryBuilder) processOrClauseParam(paramValue string) string {
	whereClauses := make([]string, 0)
	paramsValueList := strings.Split(paramValue, ",")
	for _, status := range paramsValueList {
		if val, ok := paramsWhereClause[lib.TrimSpace(status)]; ok && val != "" {
			whereClauses = append(whereClauses, val)
		}
	}
	return "(" + strings.Join(whereClauses, " OR ") + ")"
}

func (qb *BigQueryQueryBuilder) processProducerUidParam(paramValue string, queryParams map[string]interface{}) string {
	tmp := make([]string, 0)
	for _, uid := range strings.Split(paramValue, ",") {
		randomIdentifier := qb.randomGenerator()
		queryParams[randomIdentifier] = uid
		tmp = append(tmp, fmt.Sprintf("'@%s'", randomIdentifier))
	}
	return fmt.Sprintf(paramsWhereClause["producerUid"], strings.Join(tmp, ", "))
}

func (qb *BigQueryQueryBuilder) processParams(allowedParams []string, filteredParams map[string]string) ([]string, map[string]interface{}) {
	whereClauses := make([]string, 0)
	queryParams := make(map[string]interface{})

	for _, paramKey := range allowedParams {
		paramValue, exists := filteredParams[paramKey]
		if !exists || paramValue == "" {
			continue
		}

		if lib.SliceContains(orClausesKeys, paramKey) {
			whereClause := qb.processOrClauseParam(filteredParams[paramKey])
			whereClauses = append(whereClauses, whereClause)
		} else if paramKey == "producerUid" {
			whereClause := qb.processProducerUidParam(paramValue, queryParams)
			whereClauses = append(whereClauses, whereClause)
		} else {
			randomIdentifier := qb.randomGenerator()
			whereClauses = append(whereClauses, fmt.Sprintf(paramsWhereClause[paramKey], randomIdentifier))
			queryParams[randomIdentifier] = paramValue
		}

	}
	return whereClauses, queryParams
}

func (qb *BigQueryQueryBuilder) BuildQuery(params map[string]string) (string, map[string]interface{}) {
	const queryPrefix = "SELECT **tableAlias**.uid, **tableAlias**.name AS productName, " +
		"**tableAlias**.codeCompany, CAST(**tableAlias**.proposalNumber AS INT64) AS proposalNumber, " +
		"**tableAlias**.nameDesc,**tableAlias**.status, RTRIM(COALESCE(JSON_VALUE(**tableAlias**.data, " +
		"'$.contractor.name'), '') || ' ' || " +
		"COALESCE(JSON_VALUE(**tableAlias**.data, '$.contractor.surname'), '')) AS contractor, " +
		"**tableAlias**.priceGross AS price, **tableAlias**.priceGrossMonthly AS priceMonthly, " +
		"COALESCE(nn.name, '') AS producer, COALESCE(**tableAlias**.producerCode, '') AS producerCode, " +
		"**tableAlias**.startDate, **tableAlias**.endDate, **tableAlias**.paymentSplit, " +
		"**tableAlias**.hasMandate AS hasMandate " +
		"FROM `wopta.**tableName**` **tableAlias** " +
		"LEFT JOIN `wopta.networkNodesView` nn ON nn.uid = **tableAlias**.producerUid " +
		"WHERE "
	var (
		err           error
		limit         = 10
		rawQuery      bytes.Buffer
		queryParams   = make(map[string]interface{})
		whereClauses  = make([]string, 0)
		allowedParams = make([]string, 0)
	)

	if val, ok := params["limit"]; ok {
		limit, err = strconv.Atoi(val)
		if err != nil {
			log.Printf("Failed to parse limit: %v", err)
			return "", nil
		}
		if limit > 100 {
			limit = 100
		}
		delete(params, "limit")
	}

	allowedParams = qb.getAllowedParams(params)
	if allowedParams == nil {
		return "", nil
	}

	filteredParams := qb.filterParams(params, allowedParams)
	if len(filteredParams) == 0 {
		return "", nil
	}

	whereClauses, queryParams = qb.processParams(allowedParams, filteredParams)

	rawQuery.WriteString(queryPrefix)
	rawQuery.WriteString(strings.Join(whereClauses, " AND "))
	rawQuery.WriteString(fmt.Sprintf(" ORDER BY **tableAlias**.updateDate DESC"))
	rawQuery.WriteString(fmt.Sprintf(" LIMIT %d", limit))

	query := strings.ReplaceAll(rawQuery.String(), "**tableName**", qb.tableName)
	query = strings.ReplaceAll(query, "**tableAlias**", qb.tableAlias)

	return query, queryParams
}
