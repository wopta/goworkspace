package payment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	model "github.com/wopta/goworkspace/models"
)

func FabrickPayObj(data model.Policy, firstSchedule bool, scheduleDate string, customerId string, amount float64) <-chan FabrickPaymentResponse {
	r := make(chan FabrickPaymentResponse)

	go func() {
		defer close(r)
		log.Println("FabrickPay")
		//var b bytes.Buffer
		//fileReader := bytes.NewReader([]byte())

		var urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments"
		client := &http.Client{
			Timeout: time.Second * 15,
		}

		marshal := getfabbricPayments(data, firstSchedule, scheduleDate, customerId, amount)
		log.Printf(data.Uid + " " + string(marshal))
		//log.Println(getFabrickPay(data))
		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(string(marshal)))
		req.Header.Set("api-key", os.Getenv("FABRICK_TOKEN_BACK_API"))
		req.Header.Set("Auth-Schema", "S2S")

		//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-auth-token", os.Getenv("FABRICK_TOKEN_BACK_API"))
		//req.Header.Set("User-Agent", "Go-http-client/1.1")
		//req.Header.Set("Content-Length", strconv.Itoa(len(string(marshal))))
		//req.Header.Set("Host", "35.195.35.137")
		req.Header.Set("Accept", "application/json")

		log.Println(req.Header)

		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			log.Println("header:", res.Header)
			body, err := ioutil.ReadAll(res.Body)
			log.Println("body:", string(body))
			lib.CheckError(err)
			var result FabrickPaymentResponse
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()
			tr := models.SetTransactionPolicy(data, amount, scheduleDate)
			ref, _ := lib.PutFirestore("transactions", tr)
			tr.Uid = ref.ID
			tr.BigPayDate = civil.DateTimeOf(time.Now())
			tr.BigCreationDate = civil.DateTimeOf(time.Now())
			tr.BigStatusHistory = strings.Join(tr.StatusHistory, ",")
			err = lib.InsertRowsBigQuery("wopta", "transactions-day", tr)
			lib.CheckError(err)
			r <- result

		}
	}()
	return r
}
func FabbrickMontlyPay(data model.Policy) FabrickPaymentResponse {

	installment := data.PriceGross / 12
	customerId := uuid.New().String()
	log.Println(data.Uid + " FabbrickMontlyPay")
	layout := "2006-01-02"
	firstres := <-FabrickPayObj(data, true, data.StartDate.Format(layout), customerId, installment)
	time.Sleep(500)
	for i := 1; i <= 11; i++ {
		date := data.StartDate.AddDate(0, i, 0)
		res := <-FabrickPayObj(data, false, date.Format(layout), customerId, installment)
		log.Println(data.Uid+" FabbrickMontlyPay res:", res)
		time.Sleep(500)
	}
	return firstres
}
func FabbrickYearPay(data model.Policy) FabrickPaymentResponse {

	customerId := uuid.New().String()
	log.Println(data.Uid + " FabbrickYearPay")
	res := <-FabrickPayObj(data, false, "", customerId, data.PriceGross)

	return res
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

func getfabbricBase(data model.Policy) string {
	now := time.Now()
	externalId := "pay_id_" + strconv.FormatInt(now.Unix(), 10)
	return `{
		"merchantId": "wop134b31-5926-4b26-1411-726bc9f0b111",
		"externalId": "` + externalId + `",
		"paymentConfiguration": {
			"expirationDate": null,
			"allowedPaymentMethods": null,
			"callbackUrl": "https://europe-west1-positive-apex-350507.cloudfunctions.net/callback/v1/payment",
			"paymentPageRedirectUrls": null
		},
		"bill": {
			"externalId": "` + externalId + `",
			"amount": 122.0,
			"currency": "EUR",
			"description": null,
			"xInfo": null,
			"items": null,
			"subjects": null
		}
	}`
}
func getfabbricPayments(data model.Policy, firstSchedule bool, scheduleDate string, customerId string, amount float64) string {
	var mandate string

	if firstSchedule {
		mandate = "true"
	} else {
		mandate = "false"
	}
	if customerId == "" {
		customerId = uuid.New().String()
	}
	now := time.Now()
	next := now.AddDate(0, 0, 4)
	layout := "2006-01-02T15:04:05.000Z"
	layout2 := "2006-01-02"

	paymentMethods := []string{
		"CREDITCARD",
		//"FBKR2P",
		//"SDD",
		//"SMARTPOS",
	}
	var scheduleTransaction ScheduleTransaction

	var bill Bill
	if scheduleDate != "" {
		scheduleTransaction = ScheduleTransaction{DueDate: scheduleDate, PaymentInstrumentResolutionStrategy: "BY_PAYER"}
		bill.ScheduleTransaction = &scheduleTransaction
	} else {
		scheduleDate = now.Format(layout2)
	}
	log.Println(next.Format(layout))
	externalId := data.Uid + "_" + scheduleDate
	var pay FabrickPaymentsRequest
	pay.MerchantID = "wop134b31-5926-4b26-1411-726bc9f0b111"
	pay.ExternalID = externalId

	bill.ExternalID = externalId
	bill.Amount = amount
	bill.Currency = "EUR"
	bill.Description = "Pagamento polizza n° " + data.NumberCompany

	bill.MandateCreation = mandate

	bill.Items = []Item{{ExternalID: externalId, Amount: amount, Currency: "EUR"}}
	bill.Subjects = &[]Subject{{ExternalID: customerId, Role: "customer", Email: data.Contractor.Mail, Name: data.Contractor.Name + ` ` + data.Contractor.Surname}}

	pay.PaymentConfiguration = PaymentConfiguration{

		//ExpirationDate: next.Format(layout),
		PaymentPageRedirectUrls: PaymentPageRedirectUrls{
			OnSuccess: "https://www.wopta.it",
			OnFailure: "https://www.wopta.it",
			//OnInterruption: "https://www.wopta.it",
		},

		AllowedPaymentMethods: &[]AllowedPaymentMethod{{Role: "payer", PaymentMethods: paymentMethods}},
		CallbackURL:           "https://europe-west1-" + os.Getenv("GOOGLE_PROJECT_ID") + ".cloudfunctions.net/callback/v1/payment?uid=" + data.Uid + `&schedule=` + scheduleDate,
		//PayByLink:             []PayByLink{{Type: "EMAIL", Recipients: data.Contractor.Mail, Template: "pay-by-link"}},
	}
	pay.Bill = bill

	res, _ := pay.Marshal()
	log.Println(data.Uid + "Request payment:" + string(res))
	return string(res)
}
