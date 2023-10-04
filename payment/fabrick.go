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

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	tr "github.com/wopta/goworkspace/transaction"
)

func getFabrickClient(urlstring string, req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}

	req.Header.Set("api-key", os.Getenv("FABRICK_TOKEN_BACK_API"))
	req.Header.Set("Auth-Schema", "S2S")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-auth-token", os.Getenv("FABRICK_PERSISTENT_KEY"))
	req.Header.Set("Accept", "application/json")
	log.Println("[GetFabrickClient]", req)

	return client.Do(req)
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

func FabrickPayFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(io.ReadAll(r.Body))

	var data models.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	log.Println(data.PriceGross)
	lib.CheckError(err)

	paymentMethods := getPaymentMethods(data)

	resultPay := <-FabrickPayObj(data, false, "", time.Now().UTC().AddDate(10, 0, 0).Format(models.TimeDateOnly), "", data.PriceGross,
		data.PriceNett, getOrigin(r.Header.Get("origin")), paymentMethods)

	log.Println(resultPay)
	return "", nil, err
}

func FabrickPayMonthlyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(io.ReadAll(r.Body))

	var data models.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	log.Println(data.PriceGross)
	lib.CheckError(err)

	paymentMethods := getPaymentMethods(data)

	resultPay := FabrickMonthlyPay(data, getOrigin(r.Header.Get("origin")), paymentMethods)
	b, err := json.Marshal(resultPay)
	log.Println(resultPay)
	return string(b), resultPay, err
}

func FabrickPayObj(data models.Policy, firstSchedule bool, scheduleDate string, expireDate string, customerId string, amount float64, amountNet float64, origin string, paymentMethods []string) <-chan FabrickPaymentResponse {
	r := make(chan FabrickPaymentResponse)

	go func() {
		defer close(r)
		log.Println("[FabrickPayObj]")

		var (
			urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments"
		)

		marshal := getFabrickPayments(data, firstSchedule, scheduleDate, expireDate, customerId, amount, origin, paymentMethods)
		log.Printf("[FabrickPayObj] Policy %s: %s", data.Uid, string(marshal))

		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(string(marshal)))
		req.Header.Set("api-key", os.Getenv("FABRICK_TOKEN_BACK_API"))
		req.Header.Set("Auth-Schema", "S2S")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-auth-token", os.Getenv("FABRICK_TOKEN_BACK_API"))
		req.Header.Set("Accept", "application/json")
		log.Printf("[FabrickPayObj] Policy %s request headers: %s", data.Uid, req.Header)
		log.Printf("[FabrickPayObj] Policy %s request body: %s", data.Uid, req.Body)

		res, err := lib.RetryDo(req, 5)
		lib.CheckError(err)

		if res != nil {
			log.Printf("[FabrickPayObj] Policy %s response headers: %s", data.Uid, res.Header)
			body, err := io.ReadAll(res.Body)
			lib.CheckError(err)
			log.Printf("[FabrickPayObj] Policy %s response body: %s", data.Uid, string(body))

			var result FabrickPaymentResponse
			json.Unmarshal([]byte(body), &result)
			defer res.Body.Close()

			tr.PutByPolicy(data, scheduleDate, origin, expireDate, customerId, amount, amountNet, *result.Payload.PaymentID, "", false)

			r <- result
		}
	}()
	return r
}

func FabrickMonthlyPay(data models.Policy, origin string, paymentMethods []string) FabrickPaymentResponse {
	log.Printf("[FabrickMonthlyPay] Policy %s", data.Uid)

	customerId := uuid.New().String()
	firstres := <-FabrickPayObj(data, true, "", "", customerId, data.PriceGrossMonthly, data.PriceNettMonthly, origin, paymentMethods)
	time.Sleep(100)

	for i := 1; i <= 11; i++ {
		date := data.StartDate.AddDate(0, i, 0)
		expireDate := date.AddDate(0, 1, -1)
		res := <-FabrickPayObj(data, false, date.Format(models.TimeDateOnly), expireDate.Format(models.TimeDateOnly), customerId, data.PriceGrossMonthly, data.PriceNettMonthly, origin, paymentMethods)
		log.Printf("[FabrickMonthlyPay] Policy %s - Index %d - response: %v", data.Uid, i, res)
		time.Sleep(100)
	}

	return firstres
}

func FabrickYearPay(data models.Policy, origin string, paymentMethods []string) FabrickPaymentResponse {
	log.Printf("[FabrickYearPay] Policy %s", data.Uid)

	customerId := uuid.New().String()
	res := <-FabrickPayObj(data, false, "", time.Now().UTC().AddDate(10, 0, 0).Format(models.TimeDateOnly), customerId, data.PriceGross, data.PriceNett, origin, paymentMethods)

	return res
}

func getFabrickPayments(data models.Policy, firstSchedule bool, scheduleDate string, expireDate string, customerId string, amount float64, origin string, paymentMethods []string) string {
	log.Printf("[getFabrickPayments] Policy %s", data.Uid)

	var (
		mandate             string
		scheduleTransaction ScheduleTransaction
		bill                Bill
		pay                 FabrickPaymentsRequest
	)

	if firstSchedule {
		mandate = "true"
	} else {
		mandate = "false"
	}
	if customerId == "" {
		customerId = uuid.New().String()
	}
	now := time.Now()

	if scheduleDate != "" {
		scheduleTransaction = ScheduleTransaction{DueDate: scheduleDate, PaymentInstrumentResolutionStrategy: "BY_PAYER"}
		bill.ScheduleTransaction = &scheduleTransaction
	} else {
		scheduleDate = now.Format(models.TimeDateOnly)
	}

	externalId := data.Uid + "_" + scheduleDate + "_" + data.CodeCompany + "_" + strings.ReplaceAll(origin, "https://", "")
	pay.MerchantID = "wop134b31-5926-4b26-1411-726bc9f0b111"
	pay.ExternalID = externalId

	bill.ExternalID = externalId
	bill.Amount = amount
	bill.Currency = "EUR"
	bill.Description = "Pagamento polizza n° " + data.CodeCompany

	bill.MandateCreation = mandate

	bill.Items = []Item{{ExternalID: externalId, Amount: amount, Currency: "EUR"}}
	bill.Subjects = &[]Subject{{ExternalID: customerId, Role: "customer", Email: data.Contractor.Mail, Name: data.Contractor.Name + ` ` + data.Contractor.Surname}}
	callbackUrl := "https://europe-west1-" + os.Getenv("GOOGLE_PROJECT_ID") + ".cloudfunctions.net/callback/v1/payment?uid=" + data.Uid + `&schedule=` + scheduleDate + `&token=` + os.Getenv("WOPTA_TOKEN_API") + `&origin=` + origin
	callbackUrl = strings.Replace(callbackUrl, `\u0026`, `&`, 1)

	log.Printf("[getFabrickPayments] Policy %s callbackUrl: %s", data.Uid, callbackUrl)

	if expireDate != "" {
		tmpExpireDate, err := time.Parse(models.TimeDateOnly, expireDate)
		lib.CheckError(err)
		expireDate = tmpExpireDate.UTC().Format("2006-01-02T15:04:05.999999999Z")
	} else {
		expireDate = time.Now().UTC().AddDate(10, 0, 0).Format("2006-01-02T15:04:05.999999999Z")
	}

	pay.PaymentConfiguration = PaymentConfiguration{
		PaymentPageRedirectUrls: PaymentPageRedirectUrls{
			OnSuccess: "https://www.wopta.it",
			OnFailure: "https://www.wopta.it",
		},
		ExpirationDate: expireDate,
		AllowedPaymentMethods: &[]AllowedPaymentMethod{{Role: "payer", PaymentMethods: lib.SliceMap(paymentMethods,
			func(item string) string { return strings.ToUpper(item) })}},
		CallbackURL: callbackUrl,
	}
	pay.Bill = bill

	res, _ := pay.Marshal()
	result := strings.Replace(string(res), `\u0026`, `&`, -1)

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
	expirationDate := time.Now().UTC().AddDate(0, 0, -1).Format(layout2)
	var urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments/expirationDate"

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
