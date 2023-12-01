package policy

import (
	"log"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetPoliciesByQueries(origin string, queries []models.Query, limitValue int) ([]models.Policy, error) {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

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

	docSnap, err := fireQueries.FirestoreWhereLimitFields(firePolicy, limitValue)
	return models.PolicyToListData(docSnap), err
}

func getQueryOperator(queryOp string) string {
	switch queryOp {
	case "lte":
		return "<="
	case "gte":
		return ">="
	case "neq":
		return "!="
	}
	return queryOp
}
