package base

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

type QueryBuilder struct {
	tableName         string
	tableAlias        string
	randomGenerator   func() string
	paramsHierarchy   []map[string][]string
	paramsWhereClause map[string]string
	orClausesKey      []string
	WhereClauses      []string
	limit             uint64
}

func NewQueryBuilder(tableName, tableAlias string, randomGenerator func() string,
	paramsHierarchy []map[string][]string, paramsWhereClause map[string]string,
	orClausesKey []string) QueryBuilder {
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
	return QueryBuilder{
		tableName:         tableName,
		tableAlias:        tableAlias,
		randomGenerator:   randomGenerator,
		paramsHierarchy:   paramsHierarchy,
		paramsWhereClause: paramsWhereClause,
		orClausesKey:      orClausesKey,
		WhereClauses:      make([]string, 0),
	}
}

func (qb *QueryBuilder) getAllowedParams(params map[string]string) []string {
	paramsKeys := lib.GetMapKeys(params)
	for _, value := range qb.paramsHierarchy {
		for k, v := range value {
			if lib.SliceContains(paramsKeys, k) {
				return v
			}
		}
	}
	return nil
}

func (qb *QueryBuilder) filterParams(params map[string]string, allowedParams []string) map[string]string {
	paramsKeys := lib.GetMapKeys(params)
	for _, key := range paramsKeys {
		if !lib.SliceContains(allowedParams, key) {
			delete(params, key)
		}
	}
	return params
}

func (qb *QueryBuilder) processOrClauseParam(paramValue string) string {
	whereClauses := make([]string, 0)
	paramsValueList := strings.Split(paramValue, ",")
	for _, status := range paramsValueList {
		if val, ok := qb.paramsWhereClause[lib.TrimSpace(status)]; ok && val != "" {
			whereClauses = append(whereClauses, val)
		}
	}
	return "(" + strings.Join(whereClauses, " OR ") + ")"
}

func (qb *QueryBuilder) processProducerUidParam(paramValue string, queryParams map[string]interface{}) string {
	tmp := make([]string, 0)
	for _, uid := range strings.Split(paramValue, ",") {
		randomIdentifier := qb.randomGenerator()
		queryParams[randomIdentifier] = lib.TrimSpace(uid)
		tmp = append(tmp, fmt.Sprintf("@%s", randomIdentifier))
	}
	return fmt.Sprintf(qb.paramsWhereClause["producerUid"], strings.Join(tmp, ", "))
}

func (qb *QueryBuilder) processParams(allowedParams []string, filteredParams map[string]string) ([]string, map[string]interface{}) {
	whereClauses := make([]string, 0)
	queryParams := make(map[string]interface{})

	for _, paramKey := range allowedParams {
		paramValue, exists := filteredParams[paramKey]
		if !exists || paramValue == "" {
			continue
		}

		if lib.SliceContains(qb.orClausesKey, paramKey) {
			whereClause := qb.processOrClauseParam(filteredParams[paramKey])
			whereClauses = append(whereClauses, whereClause)
		} else if paramKey == "producerUid" {
			whereClause := qb.processProducerUidParam(paramValue, queryParams)
			whereClauses = append(whereClauses, whereClause)
		} else {
			randomIdentifier := qb.randomGenerator()
			whereClauses = append(whereClauses, fmt.Sprintf(qb.paramsWhereClause[paramKey], randomIdentifier))
			queryParams[randomIdentifier] = paramValue
		}

	}
	return whereClauses, queryParams
}

func (qb *QueryBuilder) extractLimit(params map[string]string) error {
	var (
		err   error
		limit = 10
	)
	if val, ok := params["limit"]; ok {
		limit, err = strconv.Atoi(val)
		if err != nil {
			return err
		}
		if limit > 100 {
			limit = 100
		}
		delete(params, "limit")
	}
	qb.limit = uint64(limit)

	return nil
}

func (qb *QueryBuilder) parseQuery() string {
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
	if len(qb.WhereClauses) != 0 {
		rawQuery.WriteString(strings.Join(qb.WhereClauses, " AND "))
	}
	rawQuery.WriteString(fmt.Sprintf(" ORDER BY **tableAlias**.updateDate DESC"))
	rawQuery.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))

	query := strings.ReplaceAll(rawQuery.String(), "**tableName**", qb.tableName)
	query = strings.ReplaceAll(query, "**tableAlias**", qb.tableAlias)

	return query
}

func (qb *QueryBuilder) Build(params map[string]string) (string, map[string]interface{}) {
	var (
		err           error
		query         string
		allowedParams []string
		whereClauses  []string
		queryParams   map[string]interface{}
	)

	err = qb.extractLimit(params)
	if err != nil {
		log.Printf("Error extracting limit: %v", err)
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
	qb.WhereClauses = append(whereClauses, qb.WhereClauses...)

	query = qb.parseQuery()

	return query, queryParams
}
