package payment

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	model "github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Payment")
	functions.HTTP("Payment", Payment)
}

func Payment(w http.ResponseWriter, r *http.Request) {

	log.Println("Callback")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/v1/fabrick",
				Handler: FabrickPay,
				Method:  "POST",
			},
			{
				Route:   "/v1/fabrick/montly",
				Handler: FabrickPayMontly,
				Method:  "POST",
			},
			{
				Route:   "/v1/cripto",
				Handler: CriptoPay,
				Method:  "POST",
			},
			{
				Route:   "/v1/cripto",
				Handler: CriptoPay,
				Method:  "POST",
			},
			{
				Route:   "/v1/:uid",
				Handler: CriptoPay,
				Method:  http.MethodDelete,
			},
		},
	}
	route.Router(w, r)

}
func FabrickPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))

	var data model.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	log.Println(data.PriceGross)
	lib.CheckError(err)
	resultPay := <-FabrickPayObj(data, false, "", "", data.PriceGross, getOrigin(r.Header.Get("origin")))

	log.Println(resultPay)
	return "", nil, err
}
func FabrickPayMontly(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))

	var data model.Policy
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
	var urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/transactions/change-expiration"

	req, _ := http.NewRequest(http.MethodPut, urlstring, strings.NewReader(`{
		"id": "string",
		"newExpirationDate": "string"
	  }`))
	var data model.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	log.Println(data.PriceGross)
	lib.CheckError(err)
	resultPay := <-FabrickPayObj(data, false, "", "", data.PriceGross, getOrigin(r.Header.Get("origin")))

	log.Println(resultPay)
	return "", nil, err
}
