package query_builder

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
)

type baseQueryBuilder struct {
	tableName       string
	tableAlias      string
	randomGenerator func() string
}

func newBaseQueryBuilder(tableName, tableAlias string, randomGenerator func() string) baseQueryBuilder {
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
	return baseQueryBuilder{
		tableName:       tableName,
		tableAlias:      tableAlias,
		randomGenerator: randomGenerator,
	}
}

func (bq *baseQueryBuilder) getAllowedParams(params map[string]string) []string {
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

func (bq *baseQueryBuilder) filterParams(params map[string]string, allowedParams []string) map[string]string {
	paramsKeys := lib.GetMapKeys(params)
	for _, key := range paramsKeys {
		if !lib.SliceContains(allowedParams, key) {
			delete(params, key)
		}
	}
	return params
}

func (bq *baseQueryBuilder) processOrClauseParam(paramValue string) string {
	whereClauses := make([]string, 0)
	paramsValueList := strings.Split(paramValue, ",")
	for _, status := range paramsValueList {
		if val, ok := paramsWhereClause[lib.TrimSpace(status)]; ok && val != "" {
			whereClauses = append(whereClauses, val)
		}
	}
	return "(" + strings.Join(whereClauses, " OR ") + ")"
}

func (bq *baseQueryBuilder) processProducerUidParam(paramValue string, queryParams map[string]interface{}) string {
	tmp := make([]string, 0)
	for _, uid := range strings.Split(paramValue, ",") {
		randomIdentifier := bq.randomGenerator()
		queryParams[randomIdentifier] = lib.TrimSpace(uid)
		tmp = append(tmp, fmt.Sprintf("@%s", randomIdentifier))
	}
	return fmt.Sprintf(paramsWhereClause["producerUid"], strings.Join(tmp, ", "))
}

func (bq *baseQueryBuilder) processParams(allowedParams []string, filteredParams map[string]string) ([]string, map[string]interface{}) {
	whereClauses := make([]string, 0)
	queryParams := make(map[string]interface{})

	for _, paramKey := range allowedParams {
		paramValue, exists := filteredParams[paramKey]
		if !exists || paramValue == "" {
			continue
		}

		if lib.SliceContains(orClausesKeys, paramKey) {
			whereClause := bq.processOrClauseParam(filteredParams[paramKey])
			whereClauses = append(whereClauses, whereClause)
		} else if paramKey == "producerUid" {
			whereClause := bq.processProducerUidParam(paramValue, queryParams)
			whereClauses = append(whereClauses, whereClause)
		} else {
			randomIdentifier := bq.randomGenerator()
			whereClauses = append(whereClauses, fmt.Sprintf(paramsWhereClause[paramKey], randomIdentifier))
			queryParams[randomIdentifier] = paramValue
		}

	}
	return whereClauses, queryParams
}

func (bq *baseQueryBuilder) extractLimit(params map[string]string) (uint64, error) {
	var (
		err   error
		limit = 10
	)
	if val, ok := params["limit"]; ok {
		limit, err = strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		if limit > 100 {
			limit = 100
		}
		delete(params, "limit")
	}
	return uint64(limit), nil
}

func (bq *baseQueryBuilder) parseQuery(whereClauses []string, limit uint64) string {
	const queryPrefix = "SELECT **tableAlias**.uid, **tableAlias**.name AS productName, " +
		"**tableAlias**.codeCompany, CAST(**tableAlias**.proposalNumber AS INT64) AS proposalNumber, " +
		"**tableAlias**.nameDesc,**tableAlias**.status, RTRIM(COALESCE(JSON_VALUE(**tableAlias**.data, " +
		"'$.contractor.name'), '') || ' ' || " +
		"COALESCE(JSON_VALUE(**tableAlias**.data, '$.contractor.surname'), '')) AS contractor, " +
		"**tableAlias**.priceGross AS price, **tableAlias**.priceGrossMonthly AS priceMonthly, " +
		"COALESCE(nn.name, '') AS producer, COALESCE(**tableAlias**.producerCode, '') AS producerCode, " +
		"**tableAlias**.startDate, **tableAlias**.endDate, **tableAlias**.paymentSplit, " +
		"COALESCE(**tableAlias**.hasMandate, false) AS hasMandate " +
		"FROM `wopta.**tableName**` **tableAlias** " +
		"LEFT JOIN `wopta.networkNodesView` nn ON nn.uid = **tableAlias**.producerUid " +
		"WHERE "
	var (
		rawQuery bytes.Buffer
	)

	rawQuery.WriteString(queryPrefix)
	rawQuery.WriteString(strings.Join(whereClauses, " AND "))
	rawQuery.WriteString(fmt.Sprintf(" ORDER BY **tableAlias**.updateDate DESC"))
	rawQuery.WriteString(fmt.Sprintf(" LIMIT %d", limit))

	query := strings.ReplaceAll(rawQuery.String(), "**tableName**", bq.tableName)
	query = strings.ReplaceAll(query, "**tableAlias**", bq.tableAlias)

	return query
}
