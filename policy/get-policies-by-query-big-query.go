package policy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetPoliciesByQueriesBigQuery(datasetID string, tableID string, queries []models.Query, limit int) ([]models.Policy, error) {
	var query bytes.Buffer
	params := make(map[string]interface{})
	for index, q := range queries {
		if index == 0 {
			query.WriteString(fmt.Sprintf("SELECT * FROM `%s.%s` WHERE", datasetID, tableID))
		} else {
			query.WriteString(" AND")
		}

		// value parameter should be a random string of only letter, so filter all numbers
		valueParameter := regexp.MustCompile("[^a-zA-Z]").ReplaceAllString(uuid.New().String(), "")
		addQuery(&query, q, valueParameter)
		params[valueParameter] = q.Value
	}

	if limit != 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	}

	bigQueryPolicies, err := lib.QueryParametrizedRowsBigQuery[models.Policy](query.String(), params)

	if err != nil {
		return nil, err
	}
	policies := make([]models.Policy, 0)
	for _, p := range bigQueryPolicies {
		var policy models.Policy
		err = json.Unmarshal([]byte(p.Data), &policy)
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}
	return policies, err
}

func addQuery(query *bytes.Buffer, q models.Query, valueParameter string) {
	switch q.Type {
	case "dateTime":
		query.WriteString(fmt.Sprintf(" JSON_VALUE(data.%s) %s @%s", q.Field, q.Op, valueParameter))
	case "bool", "boolean":
		query.WriteString(fmt.Sprintf(" BOOL(data.%s) %s @%s", q.Field, q.Op, valueParameter))
		if q.Value == false {
			query.WriteString(fmt.Sprintf(" OR BOOL(data.%s) IS NULL", q.Field))
		}
	case "double":
		query.WriteString(fmt.Sprintf(" FLOAT64(data.%s) %s @%s", q.Field, q.Op, valueParameter))
	case "int":
		query.WriteString(fmt.Sprintf(" INT64(data.%s) %s @%s", q.Field, q.Op, valueParameter))
	default:
		if q.Op == "like" {
			query.WriteString(fmt.Sprintf(" REGEXP_CONTAINS(LOWER(JSON_VALUE(data.%s)), LOWER(@%s))", q.Field, valueParameter))
		} else {
			query.WriteString(fmt.Sprintf(" JSON_VALUE(data.%s) %s @%s", q.Field, q.Op, valueParameter))
		}
	}
}
