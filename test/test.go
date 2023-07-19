package test

/*

 */
import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var (
	signatureID = 0
)

func init() {
	log.Println("INIT Test")
	functions.HTTP("Test", Test)
}

func Test(w http.ResponseWriter, r *http.Request) {
	log.Println("Test")
	lib.EnableCors(&w, r)

	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/",
				Handler: TestFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/:operation",
				Handler: TestPostFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin},
			},
			{
				Route:   "/:operation",
				Handler: TestGetFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)
}

func TestFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[TestFx]")

	creationDateFrom := time.Now().AddDate(0, 0, -9)
	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "creationDate",
				Operator:   "<",
				QueryValue: creationDateFrom,
			},
		},
	}
	query, _ := q.FirestoreWherefields("policy")
	policies := models.PolicyToListData(query)
	for i, policy := range policies {
		log.Println(i)
		policy.BigquerySave("")
	}
	/*
		fireTransactions := "transactions"

		transactions := models.TransactionToListData(query)
		for i, transaction := range transactions {
			transaction.BigPayDate = lib.GetBigQueryNullDateTime(transaction.PayDate)
			transaction.BigTransactionDate = lib.GetBigQueryNullDateTime(transaction.TransactionDate)
			transaction.BigCreationDate = civil.DateTimeOf(transaction.CreationDate)
			transaction.BigStatusHistory = strings.Join(transaction.StatusHistory, ",")
			log.Println(i)
			log.Println(" Transaction save BigQuery: " + transaction.Uid)
			err := lib.InsertRowsBigQuery("wopta", fireTransactions, transaction)
			if err != nil {
				log.Println("ERROR Transaction "+transaction.Uid+" save BigQuery: ", err)

			}
		}
	*/
	return "", nil, nil
}

func TestPostFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var request interface{}
	operation := r.Header.Get("operation")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	json.Unmarshal([]byte(body), &request)
	log.Printf("[TestPotFx] payload %v", request)

	if operation == "error" {
		return "", nil, GetErrorJson(400, "Bad Request", "Testing error POST")
	}

	return `{"success":true}`, `{"success":true}`, nil
}

func TestGetFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	operation := r.Header.Get("operation")

	if operation == "error" {
		return "", nil, GetErrorJson(401, "Bad Request", "Testing error POST")
	}

	return `{"success":true}`, `{"success":true}`, nil
}

func GetErrorJson(code int, typeEr string, message string) error {
	var (
		e     error
		eResp map[string]interface{} = make(map[string]interface{})
		b     []byte
	)
	eResp["code"] = code
	eResp["type"] = typeEr
	eResp["message"] = message
	b, e = json.Marshal(eResp)
	e = errors.New(string(b))
	return e
}
