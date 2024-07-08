package renew

import (
	"crypto/rand"
	"encoding/hex"
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

	queryBuilder := NewBigQueryQueryBuilder(func() string {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			log.Fatalf("Failed to generate random string: %v", err)
		}
		return hex.EncodeToString(b)
	})
	q, _ := queryBuilder.BuildQuery(params)

	log.Printf("resulting query: %s", q)

	return "{}", nil, err

}
