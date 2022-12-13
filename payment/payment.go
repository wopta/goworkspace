package payment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
				Hendler: FabrickPay,
			},
			{
				Route:   "/v1/cripto",
				Hendler: CriptoPay,
			},
		},
	}
	route.Router(w, r)

}
func FabrickPay(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var data model.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	lib.CheckError(err)
	FabrickPayObj(data)
	return "", nil
}
func CriptoPay(w http.ResponseWriter, r *http.Request) (string, interface{}) {

	return "", nil
}
func FabrickPayObj(data model.Policy) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)
		log.Println("FabrickPay")
		//var b bytes.Buffer
		//fileReader := bytes.NewReader([]byte())
		var urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments"
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		log.Print(getfabbricBase())
		//log.Println(getFabrickPay(data))
		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(getfabbricBase()))
		req.Header.Set("api-key", os.Getenv("FABRICK_TOKEN_BACK_API"))
		req.Header.Set("Auth-Schema", "S2S")
		req.Header.Set("Content-Type", "application/json")
		//header('Content-Length: ' . filesize($pdf));

		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			var result map[string]string
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()
			log.Println("body:", string(body))
			r <- string(body)

		}
	}()
	return r
}

func CriptoPayObj(id string) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)
		log.Println("FabrickPay")
		//var b bytes.Buffer
		//fileReader := bytes.NewReader([]byte())
		var urlstring = os.Getenv("FABRICK_BASEURL") + "v1.0/payments"
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(getCoinqvestPay(id)))
		req.Header.Set("api-key", os.Getenv("FABRICK_TOKEN_API"))
		req.Header.Set("Auth-Schema", "S2S")
		req.Header.Set("Content-Type", "application/json")
		//header('Content-Length: ' . filesize($pdf));

		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			var result map[string]string
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()

			log.Println("body:", string(body))
			r <- result["SspFileId"]

		}
	}()
	return r
}

func getFabrickPay(data model.Policy) string {
	//2022-12-12T10:05:10.000Z
	now := time.Now()
	next := now.AddDate(0, 0, 1)
	layout := "2006-01-02T15:04:05.000Z"
	layout2 := "2006-01-02"
	log.Println(next.Format(layout))
	//"expirationDate": "` + next.Format(layout) + `",
	return `{
		"merchantId": "wop134b31-5926-4b26-1411-726bc9f0b111",
		"externalId": "TST",
		"paymentConfiguration": {
		
			"allowedPaymentMethods": [
				{
					"role": "payer",
					"paymentMethods": [
						"CREDITCARD",
						"SDD"
						
					]
				}
			],
			"payByLink": [
				{
				
					"type": "EMAIL",
					"recipients": "` + data.Contractor.Mail + `",
					"template": "pay-by-link"
				}
			],
			"callbackUrl": "https://europe-west1-positive-apex-350507.cloudfunctions.net/callback/v1/payment",
			"paymentPageRedirectUrls": {
				"onFailure": "https://www.wopta.it",
				"onSuccess": "https://www.wopta.it"
			}
		},
		"bill": {
			"externalId": "TST",
			"amount": ` + fmt.Sprintf("%.2f", data.PriceGross) + `,
			"currency": "EUR",
			"description": "Checkout pagamento",
			"items": [
				{
					"externalId": "TST",
					"amount": ` + fmt.Sprintf("%.2f", data.PriceGross) + `,
					"currency": "EUR",
					"description": "Item 1 Description",
					"xInfo": "{\"cod_azienda\": \"AZ45\",\"divisione\": \" 45\"}"
				}
			],
			"scheduleTransaction": {
				"dueDate": "` + now.Format(layout2) + `",
				"paymentInstrumentResolutionStrategy": "BY_PAYER"
			},
			"mandateCreation": "false",
			"subjects": [
				{
					"role": "customer",
					"externalId": "customer_75052100",
					"email": "` + data.Contractor.Mail + `",
					"name": "` + data.Contractor.Name + ` ` + data.Contractor.Surname + `",
					"xInfo": "{\"key2\": \"value2\"}"
				}
			]
		}
	}`
}

func getCoinqvestPay(id string) string {
	return `{
		"charge":{
		   "customerId":"716dad4c5e5f",
		   "billingCurrency":"USD",
		   "lineItems":[
			  {
				 "description":"T-Shirt",
				 "netAmount":10,
				 "quantity":1,
				 "productId":"P1234"
			  }
		   ],
		   "discountItems":[
			  {
				 "description":"Loyalty Discount",
				 "netAmount":0.5
			  }
		   ],
		   "shippingCostItems":[
			  {
				 "description":"Shipping and Handling",
				 "netAmount":3.99,
				 "taxable":false
			  }
		   ],
		   "taxItems":[
			  {
				 "name":"CA Sales Tax",
				 "percent":0.0825
			  }
		   ]
		},
		"settlementAsset":"USDC:GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN",
		"checkoutLanguage":"en",
		"webhook":"https://www.your-server.com/path/to/webhook",
		"pageSettings":{
		   "returnUrl":"https://www.merchant.com/path/to/complete/checkout",
		   "cancelUrl":"https://www.merchant.com/path/to/cancel/checkout",
		   "shopName":"The T-Shirt Store Ltd.",
		   "displayBuyerInfo":true,
		   "displaySellerInfo":true
		},
		"meta":{
		   "customAttribute":"customValue"
		},
		"anchors":{
		   "BITCOIN":"BTC:GAUTUYY2THLF7SGITDFMXJVYH3LHDSMGEAKSBU267M2K7A3W543CKUEF",
		   "ETHEREUM":"ETH:GBDEVU63Y6NTHJQQZIKVTC23NWLQVP3WJ2RI2OTSJTNYOIGICST6DUXR"
		}
	 }`
}
func getfabbricBase() string {
	return `{
		"merchantId": "wop134b31-5926-4b26-1411-726bc9f0b111",
		"externalId": "paymentXid_20221206",
		"paymentConfiguration": {
			"paymentPageRedirectUrls": {
				"onSuccess": "https://www.google.com",
				"onFailure": "https://www.amazon.it",
				"onInterruption": ""
			}
		},
		"bill": {
			"externalId": "billXid_20221206",
			"amount": 5,
			"currency": "EUR",
			"subjects": [
				{
					"role": "customer",
					"externalId": "wopta_subject_test_20221206",
					"name": "Emanuela Larovere",
					"email": "emanuela.larovere@fabrick.com"
				}
			]
		}
	}`
}
