package payment

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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
	resultPay := <-FabrickPayObj(data, false, "", "", data.PriceGross, r.Header.Get("origin"))

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
	resultPay := FabbrickMontlyPay(data, r.Header.Get("origin"))
	b, err := json.Marshal(resultPay)
	log.Println(resultPay)
	return string(b), resultPay, err
}
func CriptoPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	return "", nil, nil
}
