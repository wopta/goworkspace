package user

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func GetUsersFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req        GetUsersReq
		response   GetUsersResp
		limitValue = 10
	)
	log.Println("GetUsers")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(body, &req)
	lib.CheckError(err)

	if len(req.Queries) == 0 {
		return "", nil, fmt.Errorf("no query specified")
	}

	if req.Limit != 0 {
		limitValue = req.Limit
	}

	fireUser := lib.GetDatasetByEnv(r.Header.Get("origin"), usersCollection)

	fireQueries := lib.Firequeries{
		Queries: make([]lib.Firequery, 0, 5),
	}

	for index, q := range req.Queries {
		log.Printf("query %d/%d field: \"%s\" op: \"%s\" value: \"%v\"\n", index+1, len(req.Queries), q.Field, q.Op, q.Value)
		value := q.Value
		if q.Type == "dateTime" {
			value, _ = time.Parse(time.RFC3339, value.(string))
		}
		fireQueries.Queries = append(fireQueries.Queries, lib.Firequery{
			Field:      q.Field,
			Operator:   q.Op,
			QueryValue: value,
		})
	}

	docsnap, err := fireQueries.FirestoreWhereLimitFields(fireUser, limitValue)
	fmt.Println("printline after FirestoreWhereLimitFields")
	if err != nil {
		return "", nil, err
	}

	response.Users = models.UsersToListData(docsnap)
	log.Printf("GetUsers: found %d users\n", len(response.Users))

	jsonOut, err := json.Marshal(response)
	return string(jsonOut), response, err
}

type GetUsersResp struct {
	Users []models.User `json:"users"`
}

type GetUsersReq struct {
	Queries []struct {
		Field string      `json:"field"`
		Op    string      `json:"op"`
		Value interface{} `json:"value"`
		Type  string      `json:"type"`
	} `json:"queries,omitempty"`
	Limit int `json:"limit"`
	Page  int `json:"page"`
}
