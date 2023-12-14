package policy

import (
	"log"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetPoliciesByQueries(origin string, queries []models.Query, limitValue int) ([]models.Policy, error) {
	fireQueries := lib.Firequeries{
		Queries: make([]lib.Firequery, 0),
	}

	for index, q := range queries {
		log.Printf("query %d/%d field: \"%s\" op: \"%s\" value: \"%v\"", index+1, len(queries), q.Field, q.Op, q.Value)
		value := q.Value
		if q.Type == "dateTime" {
			value, _ = time.Parse(time.RFC3339, value.(string))
		}

		fireQueries.Queries = append(fireQueries.Queries, lib.Firequery{
			Field:      q.Field,
			Operator:   getQueryOperator(q.Op),
			QueryValue: value,
		})
	}

	// convert queries to lib.Query
	queriesLib := make([]models.Query, 0)
	for _, q := range queries {
		queriesLib = append(queriesLib, models.Query{
			Field:      q.Field,
			Op:         getQueryOperator(q.Op),
			Value:      q.Value,
			Type:       q.Type,
		})
	}
	
	return GetPoliciesByQueriesBigQuery(models.WoptaDataset, "policiesViewTmp", queriesLib, limitValue)
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
