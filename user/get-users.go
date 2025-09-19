package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	qb "gitlab.dev.wopta.it/goworkspace/user/query-builder/pkg"
)

type GetUsersReq struct {
	Queries []models.Query `json:"queries,omitempty"`
	Limit   int            `json:"limit"`
	Page    int            `json:"page"`
}

type GetUsersResp struct {
	Users []models.User `json:"users"`
}

func getUsersFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req        GetUsersReq
		response   GetUsersResp
		limitValue = 10
	)

	log.AddPrefix("GetUsersFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(body, &req)
	lib.CheckError(err)

	if len(req.Queries) == 0 {
		return "", nil, fmt.Errorf("no query specified")
	}

	if req.Limit != 0 {
		limitValue = req.Limit
	}

	fireUser := lib.UserCollection

	fireQueries := lib.Firequeries{
		Queries: make([]lib.Firequery, 0),
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
	if err != nil {
		return "", nil, err
	}

	response.Users = models.UsersToListData(docsnap)
	log.Printf("GetUsers: found %d users\n", len(response.Users))

	jsonOut, err := json.Marshal(response)

	log.Println("Handler end -------------------------------------------------")

	return string(jsonOut), response, err
}

func getUserByAuthIdFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("GetUserByAuthIdFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	authId := chi.URLParam(r, "authId")
	log.Println(authId)

	user, e := GetUserByAuthId(authId)
	jsonString, e := user.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(jsonString), user, e
}

func GetUserByAuthId(authId string) (models.User, error) {
	fireUsers := lib.UserCollection
	userFirebase := lib.WhereLimitFirestore(fireUsers, "authId", "==", authId, 1)
	return models.FirestoreDocumentToUser(userFirebase)
}

func GetUserByFiscalCodeFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("[GetUserByFiscalCodeFx] ")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	fiscalCode := chi.URLParam(r, "fiscalCode")
	log.Println(fiscalCode)

	p, e := GetUserByFiscalCode(fiscalCode)
	jsonString, e := p.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(jsonString), p, e
}

func GetUserByFiscalCode(fiscalCode string) (models.User, error) {
	userFirebase := lib.WhereLimitFirestore(lib.UserCollection, "fiscalCode", "==", fiscalCode, 1)
	return models.FirestoreDocumentToUser(userFirebase)
}

func getUserByMailFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("[GetUserByMailFx] ")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	mail := chi.URLParam(r, "mail")
	log.Println(mail)

	p, e := GetUserByMail(mail)
	jsonString, e := p.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(jsonString), p, e
}

func GetUserByMail(mail string) (models.User, error) {
	userFirebase := lib.WhereLimitFirestore(lib.UserCollection, "mail", "==", mail, 1)
	return models.FirestoreDocumentToUser(userFirebase)
}

func GetAuthUserByMail(mail string) (models.User, error) {
	var user models.User

	authId, err := lib.GetAuthUserIdByEmail(mail)
	if err != nil {
		return user, err
	}

	return GetUserByAuthId(authId)
}

type getUserResp struct {
	Users []userInfo `json:"users"`
}

type userInfo struct {
	Uid        string `json:"uid" bigquery:"uid"`
	Name       string `json:"name" bigquery:"name"`
	Surname    string `json:"surname" bigquery:"surname"`
	Mail       string `json:"mail" bigquery:"mail"`
	Role       string `json:"role" bigquery:"role"`
	FiscalCode string `json:"fiscalCode" bigquery:"fiscalCode"`
}

func getUsersV2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var err error

	log.AddPrefix("GetUsersV2Fx")
	defer func() {
		if err != nil {
			log.ErrorF("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start -----------------------------------------------")

	paramsMap := extractQueryParams(r)
	log.Printf("input params: %v", paramsMap)
	if len(paramsMap) == 0 {
		return "", nil, errors.New("no query params")
	}

	queryBuilder := qb.NewQueryBuilder()
	if queryBuilder == nil {
		return "", nil, errors.New("error initializing query builder")
	}
	query, queryParams, err := queryBuilder.Build(paramsMap)
	if err != nil {
		return "", nil, fmt.Errorf("error generating query: %v", err)
	}

	log.Printf("query: %s\nqueryParams: %+v", query, queryParams)

	users, err := lib.QueryParametrizedRowsBigQuery[userInfo](query, queryParams)
	if err != nil {
		log.ErrorF("error executing query: %s", err.Error())
		return "", nil, err
	}

	resp := &getUserResp{
		Users: users,
	}

	rawResp, err := json.Marshal(resp)

	return string(rawResp), users, err
}

func extractQueryParams(r *http.Request) map[string]string {
	inputParams := r.URL.Query()

	paramsMap := make(map[string]string)
	for key, values := range inputParams {
		paramsMap[key] = values[0]
	}
	return paramsMap
}
