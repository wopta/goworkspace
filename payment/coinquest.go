package payment

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	lib "github.com/wopta/goworkspace/lib"
)

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
