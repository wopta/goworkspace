package policy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type BigQueryPolicyData struct {
	Data string `json:"data"`
}

func GetPoliciesByQueriesBigQuery(datasetID string, tableID string, queries []models.Query, limit int) ([]models.Policy, error) {
	query, params, err := buildQuery(datasetID, tableID, queries, limit)
	if err != nil {
		return nil, err
	}

	bigQueryPolicies, err := lib.QueryParametrizedRowsBigQuery[BigQueryPolicyData](query, params)
	if err != nil {
		return nil, err
	}

	policies, err := parsePolicies(bigQueryPolicies)

	return policies, err
}

func buildQuery(datasetID string, tableID string, queries []models.Query, limit int) (string, map[string]interface{}, error) {
	var query bytes.Buffer
	params := make(map[string]interface{})
	fieldRegexp := regexp.MustCompile("^[a-zA-Z0-9.]*$")

	for index, q := range queries {
		if index == 0 {
			query.WriteString(fmt.Sprintf("SELECT data FROM `%s.%s` WHERE", datasetID, tableID))
		} else {
			query.WriteString(" AND")
		}

		valuesParameters := make([]string, 0)
		if len(q.Values) > 0 {
			for _, v := range q.Values {
				valueParameter := randomString(12)
				params[valueParameter] = v
				valuesParameters = append(valuesParameters, valueParameter)
			}
		} else {
			// value parameter should be a random string of only letters
			valueParameter := randomString(12)
			params[valueParameter] = q.Value
			valuesParameters = append(valuesParameters, valueParameter)
		}

		op, err := getWhitelistedOperator(q.Op)
		if err != nil {
			return "", nil, err
		}

		if isFieldSafe := fieldRegexp.Match([]byte(q.Field)); !isFieldSafe {
			return "", nil, fmt.Errorf("field name is not allowed: %s", q.Field)
		}

		addQuery(&query, q, op, valuesParameters)
	}

	if limit != 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	}

	return query.String(), params, nil
}

func randomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func getWhitelistedOperator(queryOp string) (string, error) {
	switch queryOp {
	case "lte":
		return "<=", nil
	case "gte":
		return ">=", nil
	case "neq":
		return "!=", nil
	case "==":
		return "=", nil
	case "=", ">", "<", "<=", ">=", "!=", "like":
		return queryOp, nil
	case "in":
		return "in", nil
	default:
		return "", fmt.Errorf("unknown query operator: %s", queryOp)
	}
}

func addQuery(query *bytes.Buffer, q models.Query, op string, valuesParameter []string) {
	switch q.Type {
	case "dateTime":
		query.WriteString(fmt.Sprintf(" JSON_VALUE(data.%s) %s (@%s", q.Field, op, valuesParameter[0]))
	case "bool", "boolean":
		if q.Value == false {
			query.WriteString(fmt.Sprintf(" (BOOL(data.%s) %s @%s", q.Field, op, valuesParameter[0]))
			query.WriteString(fmt.Sprintf(" OR BOOL(data.%s) IS NULL)", q.Field))
		} else {
			query.WriteString(fmt.Sprintf(" BOOL(data.%s) %s @%s", q.Field, op, valuesParameter[0]))
		}
	case "double":
		query.WriteString(fmt.Sprintf(" FLOAT64(data.%s) %s @%s", q.Field, op, valuesParameter[0]))
	case "int":
		query.WriteString(fmt.Sprintf(" INT64(data.%s) %s @%s", q.Field, op, valuesParameter[0]))
	case "array":
		// Assuming q.Field is in format 'arrayFieldName.elementFieldName'
		fields := strings.Split(q.Field, ".")

		arrayField, elementField := fields[0], strings.Join(fields[1:], ".")
		query.WriteString(fmt.Sprintf(" EXISTS(SELECT 1 FROM UNNEST(JSON_EXTRACT_ARRAY(data.%s)) AS array_element WHERE", arrayField))
		if op == "like" {
			query.WriteString(fmt.Sprintf(" REGEXP_CONTAINS(LOWER(JSON_EXTRACT_SCALAR(array_element, '$.%s')), LOWER(@%s)))", elementField, valuesParameter[0]))
		} else {
			query.WriteString(fmt.Sprintf(" JSON_EXTRACT_SCALAR(array_element, '$.%s') %s @%s)", elementField, op, valuesParameter[0]))
		}
	default:
		if op == "like" {
			query.WriteString(fmt.Sprintf(" REGEXP_CONTAINS(LOWER(JSON_VALUE(data.%s)), LOWER(@%s))", q.Field, valuesParameter[0]))
		} else if op == "in" {
			query.WriteString(fmt.Sprintf(" JSON_VALUE(data.%s) IN (@", q.Field))
			query.WriteString(strings.Join(valuesParameter, ", @"))
			query.WriteString(")")
		} else {
			query.WriteString(fmt.Sprintf(" JSON_VALUE(data.%s) %s @%s", q.Field, op, valuesParameter[0]))
		}
	}
}

func parsePolicies(bigQueryPolicies []BigQueryPolicyData) ([]models.Policy, error) {
	policies := make([]models.Policy, 0)
	for _, p := range bigQueryPolicies {
		var policy models.Policy
		err := json.Unmarshal([]byte(p.Data), &policy)
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}
	return policies, nil
}
