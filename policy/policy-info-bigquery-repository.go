package policy

import (
	"bytes"
	"cloud.google.com/go/civil"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type PolicyInfo struct {
	Uid            string         `json:"uid" bigquery:"uid"`
	ProductName    string         `json:"productName" bigquery:"productName"`
	CodeCompany    string         `json:"codeCompany" bigquery:"codeCompany"`
	ProposalNumber int            `json:"proposalNumber" bigquery:"proposalNumber"`
	NameDesc       string         `json:"nameDesc" bigquery:"nameDesc"`
	Status         string         `json:"status" bigquery:"status"`
	Contractor     string         `json:"contractor" bigquery:"contractor"`
	Price          float64        `json:"price" bigquery:"price"`
	PriceMonthly   float64        `json:"priceMonthly" bigquery:"priceMonthly"`
	Producer       string         `json:"producer" bigquery:"producer"`
	ProducerCode   string         `json:"producerCode" bigquery:"producerCode"`
	StartDate      civil.DateTime `json:"startDate" bigquery:"startDate"`
	EndDate        civil.DateTime `json:"endDate" bigquery:"endDate"`
	PaymentSplit   string         `json:"paymentSplit" bigquery:"paymentSplit"`
}

func getPoliciesInfoQueriesBigQuery(queries []models.Query, limit int) ([]PolicyInfo, error) {
	query, params, err := buildPolicyInfoQuery(queries, limit)
	if err != nil {
		return nil, err
	}

	policies, err := lib.QueryParametrizedRowsBigQuery[PolicyInfo](query, params)
	if err != nil {
		return nil, err
	}

	return policies, err

}

func buildPolicyInfoQuery(queries []models.Query, limit int) (string, map[string]interface{}, error) {
	var query bytes.Buffer
	params := make(map[string]interface{})
	fieldRegexp := regexp.MustCompile("^[a-zA-Z0-9.]*$")

	query.WriteString(fmt.Sprintf("SELECT p.uid, p.name AS productName, p.codeCompany, CAST(p.proposalNumber AS INT64) AS proposalNumber, " +
		"p.nameDesc,p.status, RTRIM(COALESCE(JSON_VALUE(p.data, '$.contractor.name'), '') || ' ' || " +
		"COALESCE(JSON_VALUE(p.data, '$.contractor.surname'), '')) AS contractor, " +
		"p.priceGross AS price, p.priceGrossMonthly AS priceMonthly, nn.name AS producer, p.producerCode, p.startDate, " +
		"p.endDate, p.paymentSplit " +
		"FROM `wopta.policiesView` p " +
		"INNER JOIN `wopta.networkNodesView` nn ON nn.uid = p.producerUid " +
		"WHERE"))

	for index, q := range queries {
		if index != 0 {
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

		buildWhereClause(&query, q, op, valuesParameters, "p")
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

func buildWhereClause(query *bytes.Buffer, q models.Query, op string, valuesParameter []string, tableAlias string) {
	columnPrefix := "data"
	if tableAlias != "" {
		columnPrefix = fmt.Sprintf("%s.%s", tableAlias, columnPrefix)
	}

	switch q.Type {
	case "dateTime":
		query.WriteString(fmt.Sprintf(" JSON_VALUE(%s.%s) %s @%s", columnPrefix, q.Field, op, valuesParameter[0]))
	case "bool", "boolean":
		if q.Value == false {
			query.WriteString(fmt.Sprintf(" (BOOL(%s.%s) %s @%s", columnPrefix, q.Field, op, valuesParameter[0]))
			query.WriteString(fmt.Sprintf(" OR BOOL(%s.%s) IS NULL)", columnPrefix, q.Field))
		} else {
			query.WriteString(fmt.Sprintf(" BOOL(%s.%s) %s @%s", columnPrefix, q.Field, op, valuesParameter[0]))
		}
	case "double":
		query.WriteString(fmt.Sprintf(" FLOAT64(%s.%s) %s @%s", columnPrefix, q.Field, op, valuesParameter[0]))
	case "int":
		query.WriteString(fmt.Sprintf(" INT64(%s.%s) %s @%s", columnPrefix, q.Field, op, valuesParameter[0]))
	case "array":
		// Assuming q.Field is in format 'arrayFieldName.elementFieldName'
		fields := strings.Split(q.Field, ".")

		arrayField, elementField := fields[0], strings.Join(fields[1:], ".")
		query.WriteString(fmt.Sprintf(" EXISTS(SELECT 1 FROM UNNEST(JSON_EXTRACT_ARRAY(%s.%s)) AS array_element WHERE", columnPrefix, arrayField))
		if op == "like" {
			query.WriteString(fmt.Sprintf(" REGEXP_CONTAINS(LOWER(JSON_EXTRACT_SCALAR(array_element, '$.%s')), LOWER(@%s)))", elementField, valuesParameter[0]))
		} else {
			query.WriteString(fmt.Sprintf(" JSON_EXTRACT_SCALAR(array_element, '$.%s') %s @%s)", elementField, op, valuesParameter[0]))
		}
	default:
		if op == "like" {
			query.WriteString(fmt.Sprintf(" REGEXP_CONTAINS(LOWER(JSON_VALUE(%s.%s)), LOWER(@%s))", columnPrefix, q.Field, valuesParameter[0]))
		} else if op == "in" {
			query.WriteString(fmt.Sprintf(" JSON_VALUE(%s.%s) IN (@", columnPrefix, q.Field))
			query.WriteString(strings.Join(valuesParameter, ", @"))
			query.WriteString(")")
		} else {
			query.WriteString(fmt.Sprintf(" JSON_VALUE(%s.%s) %s @%s", columnPrefix, q.Field, op, valuesParameter[0]))
		}
	}
}
