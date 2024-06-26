package user

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

type GetUsersReq struct {
	Queries []models.Query `json:"queries,omitempty"`
	Limit   int            `json:"limit"`
	Page    int            `json:"page"`
}

type GetUsersResp struct {
	Users []models.User `json:"users"`
}

func GetUsersFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req        GetUsersReq
		response   GetUsersResp
		limitValue = 10
	)

	log.SetPrefix("[GetUsersFx] ")
	defer log.SetPrefix("")

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

	fireUser := lib.GetDatasetByEnv(r.Header.Get("origin"), lib.UserCollection)

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

func GetUserByAuthIdFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[GetUserByAuthIdFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	authId := chi.URLParam(r, "authId")
	log.Println(authId)

	user, e := GetUserByAuthId(r.Header.Get("Origin"), authId)
	jsonString, e := user.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(jsonString), user, e
}

func GetUserByAuthId(origin, authId string) (models.User, error) {
	fireUsers := lib.GetDatasetByEnv(origin, lib.UserCollection)
	userFirebase := lib.WhereLimitFirestore(fireUsers, "authId", "==", authId, 1)
	return models.FirestoreDocumentToUser(userFirebase)
}

func GetUserByFiscalCodeFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[GetUserByFiscalCodeFx] ")
	defer log.SetPrefix("")

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

func GetUserByMailFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[GetUserByMailFx] ")
	defer log.SetPrefix("")

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

func GetAuthUserByMail(origin, mail string) (models.User, error) {
	var user models.User

	authId, err := lib.GetAuthUserIdByEmail(mail)
	if err != nil {
		return user, err
	}

	return GetUserByAuthId(origin, authId)
}
