package claim

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
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
	return "", nil
}
func CriptoPay(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	return "", nil
}
func FabrickPayObj(id string) <-chan string {
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
		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(getFabrickPay(id)))
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

func getFabrickPay(id string) string {
	return `{
		"merchantId": "merchant_id",
		"externalId": "TST_{{$timestamp}}",
		"paymentConfiguration": {
			"expirationDate": "2022-12-12T10:05:10.000Z",
			"allowedPaymentMethods": [
				{
					"role": "payer",
					"paymentMethods": [
						"CREDITCARD",
						"FBKR2P",
						"SDD",
						"SMARTPOS"
					]
				}
			],
			"payByLink": [
				{
					"type": "EMAIL",
					"recipients": "nome.cognome@fabrick.com",
					"template": "pay-by-link"
				}
			],
			"callbackUrl": "https://www.merchant.it.placeholder",
			"paymentPageRedirectUrls": {
				"onFailure": "https://www.merchant.it.placeholder",
				"onSuccess": "https://www.merchant.it.placeholder"
			}
		},
		"bill": {
			"externalId": "TST_{{$timestamp}}",
			"amount": 100.00,
			"currency": "EUR",
			"description": "Checkout pagamento",
			"items": [
				{
					"externalId": "TST_{{$timestamp}}",
					"amount": 50.00,
					"currency": "EUR",
					"description": "Item 1 Description",
					"xInfo": "{\"cod_azienda\": \"AZ45\",\"divisione\": \" 45\"}"
				},
				{
					"externalId": "TST_{{$timestamp}}",
					"amount": 50.00,
					"currency": "EUR",
					"description": "Item 2 Description",
					"xInfo": "{\"cod_azienda\": \"AZ54\",\"divisione\": \" 54\"}"
				}
			],
			"subjects": [
				{
					"role": "customer",
					"externalId": "customer_75052100",
					"email": "nome.cognome@fabrick.com",
					"name": "Mario Bianchi",
					"xInfo": "{\"key2\": \"value2\"}"
				},
				{
					"role": "intermediary",
					"externalId": "AGENZIA_45",
					"email": "age45@fabrick.com",
					"name": "Mario Rossi",
					"xInfo": "{\"customKey1\": \"value1\",\"customKey2\": \"value\"}"
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
