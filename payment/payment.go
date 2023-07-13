package payment

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/civil"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Payment")
	functions.HTTP("Payment", Payment)
}

func Payment(w http.ResponseWriter, r *http.Request) {

	log.Println("Payment")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/v1/fabrick",
				Handler: FabrickPay,
				Method:  "POST",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/fabrick/montly",
				Handler: FabrickPayMontly,
				Method:  "POST",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/cripto",
				Handler: CriptoPay,
				Method:  "POST",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/:uid",
				Handler: FabrickExpireBillFx,
				Method:  http.MethodDelete,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/manual/v1/:transactionUid",
				Handler: ManualPaymentFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
		},
	}
	route.Router(w, r)

}
func FabrickPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(io.ReadAll(r.Body))

	var data models.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	log.Println(data.PriceGross)
	lib.CheckError(err)
	resultPay := <-FabrickPayObj(data, false, "", "", "", data.PriceGross, data.PriceNett, getOrigin(r.Header.Get("origin")))

	log.Println(resultPay)
	return "", nil, err
}
func FabrickPayMontly(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(io.ReadAll(r.Body))

	var data models.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	log.Println(data.PriceGross)
	lib.CheckError(err)
	resultPay := FabbrickMontlyPay(data, getOrigin(r.Header.Get("origin")))
	b, err := json.Marshal(resultPay)
	log.Println(resultPay)
	return string(b), resultPay, err
}

func CriptoPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	return "", nil, nil
}

func getOrigin(origin string) string {
	var result string
	if strings.Contains(origin, "uat") || strings.Contains(origin, "dev") {
		result = "uat"
	} else {
		result = ""
	}
	log.Println(" getOrigin: name:", origin)
	log.Println(" getOrigin result: ", result)
	return result
}

func FabrickExpireBillFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		transaction models.Transaction
	)

	uid := r.Header.Get("uid")
	origin := r.Header.Get("origin")

	fireTransactions := lib.GetDatasetByEnv(origin, "transactions")
	docsnap, err := lib.GetFirestoreErr(fireTransactions, uid)
	lib.CheckError(err)
	docsnap.DataTo(&transaction)

	err = FabrickExpireBill(&transaction)

	if err != nil {
		log.Printf("[FabrickExpireBillFx]: ERROR: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	lib.SetFirestore(fireTransactions, transaction.Uid, transaction)
	err = lib.InsertRowsBigQuery("wopta", fireTransactions, transaction)

	return `{"success":true}`, `{"success":true}`, err
}

const (
	layout               string = "2006-01-02T15:04:05.000Z"
	layout2              string = "2006-01-02"
	expirationTimeSuffix string = " 00:00:00"
)

func FabrickExpireBill(transaction *models.Transaction) error {
	var err error

	expirationDate := time.Now().UTC().AddDate(0, 0, 1).Format(layout2)
	urlstring := os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments/change-expiration"

	req, _ := http.NewRequest(http.MethodPut, urlstring, strings.NewReader(`{"id":"`+transaction.ProviderId+`","newExpirationDate":"`+expirationDate+expirationTimeSuffix+`"}`))
	res, err := getFabrickClient(urlstring, req)
	lib.CheckError(err)

	respBody, err := io.ReadAll(res.Body)
	lib.CheckError(err)
	log.Println("Fabrick res body: ", string(respBody))
	if res.StatusCode != http.StatusOK {
		log.Printf("[FabrickExpireBill] ERROR response status code: %s", res.Status)
		return errors.New("status code " + res.Status)
	}

	transaction.ExpirationDate = expirationDate
	transaction.Status = models.PolicyStatusDeleted
	transaction.StatusHistory = append(transaction.StatusHistory, models.PolicyStatusDeleted)
	transaction.IsDelete = true
	transaction.BigCreationDate = civil.DateTimeOf(transaction.CreationDate)
	transaction.BigStatusHistory = strings.Join(transaction.StatusHistory, ",")

	return err
}
