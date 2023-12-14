package policy

import (
	"github.com/wopta/goworkspace/models"
	"log"
)

func GetPoliciesByQueries(origin string, requestQueries []models.Query, limitValue int) ([]models.Policy, error) {
	queries := make([]models.Query, 0)
	for index, q := range requestQueries {
		log.Printf("query %d/%d field: \"%s\" op: \"%s\" value: \"%v\"", index+1, len(requestQueries), q.Field, q.Op, q.Value)

		queries = append(queries, models.Query{
			Field: q.Field,
			Op:    getQueryOperator(q.Op),
			Value: q.Value,
			Type:  q.Type,
		})
	}

	return GetPoliciesByQueriesBigQuery(models.WoptaDataset, "policiesViewTmp", queries, limitValue)
}

func getQueryOperator(queryOp string) string {
	switch queryOp {
	case "lte":
		return "<="
	case "gte":
		return ">="
	case "neq":
		return "!="
	case "==":
		return "="
	default:
		return queryOp
	}
}
