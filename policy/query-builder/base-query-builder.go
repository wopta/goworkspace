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

type baseQueryBuilder struct {
	tableName         string
	tableAlias        string
	randomGenerator   func() string
	paramsHierarchy   []map[string][]string
	paramsWhereClause map[string]string
	orClausesKey      []string
	whereClauses      []string
	limit             uint64
}

func newBaseQueryBuilder(tableName, tableAlias string, randomGenerator func() string,
	paramsHierarchy []map[string][]string, paramsWhereClause map[string]string,
	orClausesKey []string) baseQueryBuilder {
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
		tableName:         tableName,
		tableAlias:        tableAlias,
		randomGenerator:   randomGenerator,
		paramsHierarchy:   paramsHierarchy,
		paramsWhereClause: paramsWhereClause,
		orClausesKey:      orClausesKey,
		whereClauses:      make([]string, 0),
	}
}

func (bq *baseQueryBuilder) getAllowedParams(params map[string]string) []string {
	paramsKeys := lib.GetMapKeys(params)
	for _, value := range bq.paramsHierarchy {
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
		if val, ok := bq.paramsWhereClause[lib.TrimSpace(status)]; ok && val != "" {
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

		if lib.SliceContains(bq.orClausesKey, paramKey) {
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

func (bq *baseQueryBuilder) extractLimit(params map[string]string) error {
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
	bq.limit = uint64(limit)

	return nil
}

func (bq *baseQueryBuilder) parseQuery() string {
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
	if len(bq.whereClauses) != 0 {
		rawQuery.WriteString(strings.Join(bq.whereClauses, " AND "))
	}
	rawQuery.WriteString(fmt.Sprintf(" ORDER BY **tableAlias**.updateDate DESC"))
	rawQuery.WriteString(fmt.Sprintf(" LIMIT %d", bq.limit))

	query := strings.ReplaceAll(rawQuery.String(), "**tableName**", bq.tableName)
	query = strings.ReplaceAll(query, "**tableAlias**", bq.tableAlias)

	return query
}

func (bq *baseQueryBuilder) BuildQuery(params map[string]string) (string, map[string]interface{}) {
	var (
		err           error
		query         string
		allowedParams []string
		whereClauses  []string
		queryParams   map[string]interface{}
	)

	err = bq.extractLimit(params)
	if err != nil {
		log.Printf("Error extracting limit: %v", err)
	}

	allowedParams = bq.getAllowedParams(params)
	if allowedParams == nil {
		return "", nil
	}

	filteredParams := bq.filterParams(params, allowedParams)
	if len(filteredParams) == 0 {
		return "", nil
	}

	whereClauses, queryParams = bq.processParams(allowedParams, filteredParams)
	bq.whereClauses = append(whereClauses, bq.whereClauses...)

	query = bq.parseQuery()

	return query, queryParams
}
