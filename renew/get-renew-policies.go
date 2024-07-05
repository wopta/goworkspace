package renew

import (
	"log"
	"net/http"
)

func GetRenewPoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
	)

	query := r.URL.Query()
	params := make(map[string]string)
	for key, values := range query {
		params[key] = values[0]
	}

	for key, value := range params {
		log.Printf("key: %s, value: %s", key, value)
	}

	queryBuilder := NewBigQueryQueryBuilder()
	q := queryBuilder.BuildQuery(params)

	log.Printf("resulting query: %s", q)

	return "{}", nil, err

}
