package payment

import (
	"encoding/json"
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
				Handler: FabrickExpireBill,
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

func FabrickExpireBill(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var transaction models.Transaction
	const expirationTimeSuffix = " 00:00:00"
	//layout := "2006-01-02T15:04:05.000Z"
	layout2 := "2006-01-02"

	log.Println(r.Header.Get("uid"))
	uid := r.Header.Get("uid")
	fireTransactions := lib.GetDatasetByEnv(r.Header.Get("origin"), "transactions")
	docsnap, e := lib.GetFirestoreErr(fireTransactions, uid)
	docsnap.DataTo(&transaction)
	expirationDate := time.Now().UTC().AddDate(0, 0, 1).Format(layout2)
	var urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments/change-expiration"

	req, _ := http.NewRequest(http.MethodPut, urlstring, strings.NewReader(`{"id":"`+transaction.ProviderId+`","newExpirationDate":"`+expirationDate+expirationTimeSuffix+`"}`))
	res, e := getFabrickClient(urlstring, req)
	respBody, e := io.ReadAll(res.Body)
	log.Println("Fabrick res body: ", string(respBody))
	if res.StatusCode != http.StatusOK {
		log.Printf("ExpireBill: fabrick error response status code: %s", res.Status)
		return `{"success":false}`, `{"success":false}`, nil
	}
	transaction.ExpirationDate = expirationDate
	transaction.Status = models.PolicyStatusDeleted
	transaction.StatusHistory = append(transaction.StatusHistory, models.PolicyStatusDeleted)
	transaction.IsDelete = true
	transaction.BigCreationDate = civil.DateTimeOf(transaction.CreationDate)
	transaction.BigStatusHistory = strings.Join(transaction.StatusHistory, ",")
	lib.SetFirestore(fireTransactions, uid, transaction)
	e = lib.InsertRowsBigQuery("wopta", fireTransactions, transaction)

	return `{"success":true}`, `{"success":true}`, e
}
