package base

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
)

type UsersQueryBuilder struct {
	tableName          string
	tableAlias         string
	randomGenerator    func() string
	paramsHierarchy    []map[string][]string
	paramsWhereClause  map[string]string
	toBeTranslatedKeys []string
	WhereClauses       []string
	limit              uint64
}

func NewQueryBuilder(tableName, tableAlias string,
	paramsHierarchy []map[string][]string, paramsWhereClause map[string]string) UsersQueryBuilder {

	var randomGenerator = func() string {
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

	return UsersQueryBuilder{
		tableName:          tableName,
		tableAlias:         tableAlias,
		randomGenerator:    randomGenerator,
		paramsHierarchy:    paramsHierarchy,
		paramsWhereClause:  paramsWhereClause,
		toBeTranslatedKeys: make([]string, 0),
		WhereClauses:       make([]string, 0),
	}
}

func (qb *UsersQueryBuilder) getAllowedParams(params map[string]string) ([]string, error) {
	paramsKeys := lib.GetMapKeys(params)
	for _, value := range qb.paramsHierarchy {
		for k, v := range value {
			if lib.SliceContains(paramsKeys, k) {
				return v, nil
			}
		}
	}
	return nil, errors.New("parameters not allowed")
}

func (qb *UsersQueryBuilder) filterParams(params map[string]string, allowedParams []string) (map[string]string, error) {
	paramsKeys := lib.GetMapKeys(params)
	for _, key := range paramsKeys {
		if !lib.SliceContains(allowedParams, key) {
			delete(params, key)
		}
	}
	if len(params) == 0 {
		return nil, errors.New("parameters not allowed")
	}
	return params, nil
}

func (qb *UsersQueryBuilder) processToBeTranslatedParam(paramValue string) (string, error) {
	whereClauses := make([]string, 0)
	paramsValueList := strings.Split(paramValue, ",")
	for _, status := range paramsValueList {
		if val, ok := qb.paramsWhereClause[lib.TrimSpace(status)]; ok && val != "" {
			whereClauses = append(whereClauses, val)
		}
	}
	if len(whereClauses) == 0 {
		return "", errors.New("error processing params")
	}

	return "(" + strings.Join(whereClauses, " OR ") + ")", nil
}

func (qb *UsersQueryBuilder) processProducerUidParam(paramValue string, queryParams map[string]interface{}) string {
	tmp := make([]string, 0)
	for _, uid := range strings.Split(paramValue, ",") {
		randomIdentifier := qb.randomGenerator()
		queryParams[randomIdentifier] = lib.TrimSpace(uid)
		tmp = append(tmp, fmt.Sprintf("@%s", randomIdentifier))
	}
	return fmt.Sprintf(qb.paramsWhereClause["mail"], strings.Join(tmp, ", "))
}

func (qb *UsersQueryBuilder) processParams(allowedParams []string, filteredParams map[string]string) ([]string,
	map[string]interface{}, error) {
	whereClauses := make([]string, 0)
	queryParams := make(map[string]interface{})

	for _, paramKey := range allowedParams {
		paramValue, exists := filteredParams[paramKey]
		if !exists || paramValue == "" {
			continue
		}

		if lib.SliceContains(qb.toBeTranslatedKeys, paramKey) {
			whereClause, err := qb.processToBeTranslatedParam(filteredParams[paramKey])
			if err != nil {
				return nil, nil, err
			}
			whereClauses = append(whereClauses, whereClause)
		} else if paramKey == "DO_NOT_EXIST" {
			whereClause := qb.processProducerUidParam(paramValue, queryParams)
			if whereClause != "" {
				whereClauses = append(whereClauses, whereClause)
			}
		} else {
			randomIdentifier := qb.randomGenerator()
			whereClauses = append(whereClauses, fmt.Sprintf(qb.paramsWhereClause[paramKey], randomIdentifier))
			queryParams[randomIdentifier] = paramValue
		}

	}
	return whereClauses, queryParams, nil
}

func (qb *UsersQueryBuilder) extractLimit(params map[string]string) error {
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

func (qb *UsersQueryBuilder) parseQuery() string {
	const queryPrefix = "SELECT **tableAlias**.uid, **tableAlias**.name, " +
		"**tableAlias**.surname, **tableAlias**.mail, **tableAlias**.role, **tableAlias**.fiscalCode " +
		"FROM `wopta.**tableName**` **tableAlias** " +
		"WHERE "
	var (
		rawQuery bytes.Buffer
	)

	rawQuery.WriteString(queryPrefix)
	if len(qb.WhereClauses) == 1 {
		rawQuery.WriteString(qb.WhereClauses[0])
	} else if len(qb.WhereClauses) > 1 {
		rawQuery.WriteString(strings.Join(qb.WhereClauses, " AND "))
	}
	rawQuery.WriteString(fmt.Sprintf(" ORDER BY **tableAlias**.updatedDate DESC"))
	rawQuery.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))

	query := strings.ReplaceAll(rawQuery.String(), "**tableName**", qb.tableName)
	query = strings.ReplaceAll(query, "**tableAlias**", qb.tableAlias)

	return query
}

func (qb *UsersQueryBuilder) Build(params map[string]string) (string, map[string]interface{}, error) {
	var (
		err           error
		query         string
		allowedParams []string
		whereClauses  []string
		queryParams   map[string]interface{}
	)

	err = qb.extractLimit(params)
	if err != nil {
		return "", nil, fmt.Errorf("error extracting limit: %w", err)
	}

	allowedParams, err = qb.getAllowedParams(params)
	if err != nil {
		return "", nil, err
	}

	filteredParams, err := qb.filterParams(params, allowedParams)
	if err != nil {
		return "", nil, err
	}

	whereClauses, queryParams, err = qb.processParams(allowedParams, filteredParams)
	if err != nil {
		return "", nil, err
	}
	qb.WhereClauses = append(whereClauses, qb.WhereClauses...)

	query = qb.parseQuery()

	return query, queryParams, nil
}
