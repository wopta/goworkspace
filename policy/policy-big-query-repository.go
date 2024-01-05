package policy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"

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
	paramRegexp := regexp.MustCompile("[^a-zA-Z]")
	fieldRegexp := regexp.MustCompile("^[a-zA-Z0-9.]*$")

	for index, q := range queries {
		if index == 0 {
			query.WriteString(fmt.Sprintf("SELECT data FROM `%s.%s` WHERE", datasetID, tableID))
		} else {
			query.WriteString(" AND")
		}

		// value parameter should be a random string of only letters
		valueParameter := paramRegexp.ReplaceAllString(q.Field, "")
		params[valueParameter] = q.Value

		op, err := getWhitelistedOperator(q.Op)
		if err != nil {
			return "", nil, err
		}

		if isFieldSafe := fieldRegexp.Match([]byte(q.Field)); !isFieldSafe {
			return "", nil, fmt.Errorf("field name is not allowed: %s", q.Field)
		}

		addQuery(&query, q, op, valueParameter)
	}

	if limit != 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	}

	return query.String(), params, nil
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
	default:
		return "", fmt.Errorf("unknown query operator: %s", queryOp)
	}
}

func addQuery(query *bytes.Buffer, q models.Query, op, valueParameter string) {
	switch q.Type {
	case "dateTime":
		query.WriteString(fmt.Sprintf(" JSON_VALUE(data.%s) %s @%s", q.Field, op, valueParameter))
	case "bool", "boolean":
		if q.Value == false {
			query.WriteString(fmt.Sprintf(" (BOOL(data.%s) %s @%s", q.Field, op, valueParameter))
			query.WriteString(fmt.Sprintf(" OR BOOL(data.%s) IS NULL)", q.Field))
		} else {
			query.WriteString(fmt.Sprintf(" BOOL(data.%s) %s @%s", q.Field, op, valueParameter))
		}
	case "double":
		query.WriteString(fmt.Sprintf(" FLOAT64(data.%s) %s @%s", q.Field, op, valueParameter))
	case "int":
		query.WriteString(fmt.Sprintf(" INT64(data.%s) %s @%s", q.Field, op, valueParameter))
	default:
		if op == "like" {
			query.WriteString(fmt.Sprintf(" REGEXP_CONTAINS(LOWER(JSON_VALUE(data.%s)), LOWER(@%s))", q.Field, valueParameter))
		} else {
			query.WriteString(fmt.Sprintf(" JSON_VALUE(data.%s) %s @%s", q.Field, op, valueParameter))
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
