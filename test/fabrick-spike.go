package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/product"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	currency         string = "EUR"
	expireDateFormat string = "2006-01-02T15:04:05.999999999Z"
	redirectUrl      string = "https://www.wopta.it"
	woptaMerchantId  string = "wop134b31-5926-4b26-1411-726bc9f0b111"
)

func fabrickSpike() (string, interface{}, error) {
	var (
		err       error
		policy    models.Policy
		policyUid = "UB5wzJw1MgHKcZafILSs"
	)

	log.SetPrefix("[fabrickSpike] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	policy, err = plc.GetPolicy(policyUid, "")
	if err != nil {
		log.Println(err)
		return "", nil, err
	}

	mgaProduct := product.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	commissionMga := lib.RoundFloat(product.GetCommissionByProduct(&policy, mgaProduct, false), 2)

	tomorrow := time.Now().UTC().AddDate(0, 0, 1)

	tr := models.Transaction{
		Amount:         policy.PriceGross,
		AmountNet:      policy.PriceNett,
		Commissions:    commissionMga,
		Status:         models.TransactionStatusToPay,
		PolicyName:     policy.Name,
		Name:           policy.Contractor.Name + " " + policy.Contractor.Surname,
		ScheduleDate:   tomorrow.Format(time.DateOnly),
		ExpirationDate: tomorrow.AddDate(10, 0, 0).Format(time.DateOnly),
		CreationDate:   time.Now().UTC(),
		Uid:            lib.NewDoc(models.TransactionsCollection),
		PolicyUid:      policyUid,
		Company:        policy.Company,
		NumberCompany:  policy.CodeCompany,
		StatusHistory:  []string{models.TransactionStatusToPay},
		IsPay:          false,
		IsEmit:         false,
		IsDelete:       false,
		ProviderName:   models.FabrickPaymentProvider,
		UpdateDate:     time.Now().UTC(),
		EffectiveDate:  time.Now().UTC().AddDate(1, 0, 0),
	}

	paymentMethods := getPaymentMethods(policy, *mgaProduct)
	_, updatedTr, err := fabrickIntegration([]models.Transaction{tr}, paymentMethods, policy)

	for _, t := range updatedTr {
		lib.SetFirestoreErr(models.TransactionsCollection, t.Uid, t)

		t.BigQuerySave("")
	}
	log.Println(len(updatedTr))

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}

func getPaymentMethods(policy models.Policy, product models.Product) []string {
	var paymentMethods = make([]string, 0)

	log.Printf("[GetPaymentMethods] loading available payment methods for %s payment provider", policy.Payment)

	for _, provider := range product.PaymentProviders {
		if provider.Name == policy.Payment {
			for _, config := range provider.Configs {
				if config.Mode == policy.PaymentMode && config.Rate == policy.PaymentSplit {
					paymentMethods = append(paymentMethods, config.Methods...)
				}
			}
		}
	}

	log.Printf("[GetPaymentMethods] found %v", paymentMethods)
	return paymentMethods
}

func fabrickIntegration(transactions []models.Transaction, paymentMethods []string, policy models.Policy) (payUrl string, updatedTransactions []models.Transaction, err error) {
	customerId := "b27422ad-8566-4220-87b7-e7569e8a4dd1"
	now := time.Now().UTC()

	for index, tr := range transactions {
		tr.ProviderName = models.FabrickPaymentProvider

		res := <-createFabrickTransaction(&policy, tr, false, false, customerId, paymentMethods)
		if res.Payload == nil || res.Payload.PaymentPageURL == nil {
			return "", nil, errors.New("error creating transaction on Fabrick")
		}
		log.Printf("transaction %02d payUrl: %s", index+1, *res.Payload.PaymentPageURL)

		//tr.PayUrl = *res.Payload.PaymentPageURL
		tr.ProviderId = *res.Payload.PaymentID
		tr.UserToken = customerId
		tr.UpdateDate = now
		updatedTransactions = append(updatedTransactions, tr)
	}

	return payUrl, updatedTransactions, nil
}

func getFabrickClient(urlstring string, req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}

	req.Header.Set("api-key", os.Getenv("FABRICK_TOKEN_BACK_API"))
	req.Header.Set("Auth-Schema", "S2S")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-auth-token", os.Getenv("FABRICK_PERSISTENT_KEY"))
	req.Header.Set("Accept", "application/json")
	log.Println("[getFabrickClient]", req)

	return client.Do(req)
}

func getFabrickRequestBody(
	policy *models.Policy,
	firstSchedule bool,
	scheduleDate, expireDate, customerId string,
	amount float64,
	origin string,
	paymentMethods []string,
) string {
	var (
		mandate             string    = "false"
		now                 time.Time = time.Now() // should we use .UTC()?
		requestScheduleDate string    = scheduleDate
	)

	if firstSchedule {
		mandate = "true"
		scheduleDate = ""
	}

	if customerId == "" {
		customerId = uuid.New().String()
	}

	if requestScheduleDate == "" {
		requestScheduleDate = now.Format(models.TimeDateOnly)
	}

	externalId := strings.Join([]string{
		policy.Uid,
		requestScheduleDate,
		policy.CodeCompany,
		strings.ReplaceAll(origin, "https://", ""),
		strconv.FormatInt(now.Unix(), 10),
	}, "_")

	callbackUrl := fmt.Sprintf(
		"https://europe-west1-%s.cloudfunctions.net/callback/v1/payment?uid=%s&schedule=%s&token=%s&origin=%s",
		os.Getenv("GOOGLE_PROJECT_ID"),
		policy.Uid,
		requestScheduleDate,
		os.Getenv("WOPTA_TOKEN_API"),
		origin,
	)
	callbackUrl = strings.Replace(callbackUrl, `\u0026`, `&`, 1)

	if expireDate != "" {
		tmpExpireDate, err := time.Parse(models.TimeDateOnly, expireDate)
		if err != nil {
			log.Printf("error parsing expireDate: %s", err.Error())
			return ""
		}
		expireDate = time.Date(
			tmpExpireDate.Year(), tmpExpireDate.Month(), tmpExpireDate.Day(), 2, 30, 30, 30, time.UTC,
		).Format(expireDateFormat)
	} else {
		expireDate = time.Now().UTC().AddDate(10, 0, 0).Format(expireDateFormat)
	}

	pay := payment.FabrickPaymentsRequest{
		MerchantID: woptaMerchantId,
		ExternalID: externalId,
		PaymentConfiguration: payment.PaymentConfiguration{
			PaymentPageRedirectUrls: payment.PaymentPageRedirectUrls{
				OnSuccess: redirectUrl,
				OnFailure: redirectUrl,
			},
			ExpirationDate: expireDate,
			AllowedPaymentMethods: &[]payment.AllowedPaymentMethod{{
				Role: "payer",
				PaymentMethods: lib.SliceMap(
					paymentMethods,
					func(item string) string { return strings.ToUpper(item) },
				),
			}},
			CallbackURL: callbackUrl,
		},
		Bill: payment.Bill{
			ExternalID:      externalId,
			Amount:          amount,
			Currency:        currency,
			Description:     fmt.Sprintf("Pagamento polizza n° %s", policy.CodeCompany),
			MandateCreation: mandate,
			Items: []payment.Item{{
				ExternalID: externalId,
				Amount:     amount,
				Currency:   currency,
			}},
			Subjects: &[]payment.Subject{{
				ExternalID: customerId,
				Role:       "customer",
				Email:      policy.Contractor.Mail,
				Name:       strings.Join([]string{policy.Contractor.Name, policy.Contractor.Surname}, " "),
			}},
		},
	}

	if scheduleDate != "" {
		pay.Bill.ScheduleTransaction = &payment.ScheduleTransaction{
			DueDate:                             scheduleDate,
			PaymentInstrumentResolutionStrategy: "BY_PAYER",
		}
	}

	res, err := pay.Marshal()
	if err != nil {
		log.Printf("error marshalling body: %s", err.Error())
		return ""
	}

	return strings.Replace(string(res), `\u0026`, `&`, -1)
}

func getFabrickPaymentRequest(body string) *http.Request {
	var (
		urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments"
		token     = os.Getenv("FABRICK_TOKEN_BACK_API")
	)

	request, err := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(body))
	if err != nil {
		log.Printf("error generating fabrick payment request: %s", err.Error())
		return nil
	}

	request.Header.Set("api-key", token)
	request.Header.Set("Auth-Schema", "S2S")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-auth-token", token)
	request.Header.Set("Accept", "application/json")

	return request
}

func createFabrickTransaction(
	policy *models.Policy,
	transaction models.Transaction,
	firstSchedule, createMandate bool,
	customerId string,
	paymentMethods []string,
) <-chan payment.FabrickPaymentResponse {
	r := make(chan payment.FabrickPaymentResponse)

	go func() {
		defer close(r)

		body := getFabrickRequestBody(policy, firstSchedule, transaction.ScheduleDate, transaction.ExpirationDate,
			customerId, transaction.Amount, "", paymentMethods)
		if body == "" {
			return
		}
		request := getFabrickPaymentRequest(body)
		if request == nil {
			return
		}

		log.Printf("policy '%s' request headers: %s", policy.Uid, request.Header)
		log.Printf("policy '%s' request body: %s", policy.Uid, request.Body)

		if os.Getenv("env") == "local" || os.Getenv("env") == "local-test" {
			status := "200"
			local := "local"
			url := "www.dev.wopta.it"
			r <- payment.FabrickPaymentResponse{
				Status: &status,
				Errors: nil,
				Payload: &payment.Payload{
					ExternalID:        &local,
					PaymentID:         &local,
					MerchantID:        &local,
					PaymentPageURL:    &url,
					PaymentPageURLB2B: &url,
					TokenB2B:          &local,
					Coupon:            &local,
				},
			}
		} else {

			res, err := lib.RetryDo(request, 5, 10)
			lib.CheckError(err)

			if res != nil {
				log.Printf("policy '%s' response headers: %s", policy.Uid, res.Header)
				body, err := io.ReadAll(res.Body)
				defer res.Body.Close()
				lib.CheckError(err)
				log.Printf("policy '%s' response body: %s", policy.Uid, string(body))

				var result payment.FabrickPaymentResponse

				if res.StatusCode != 200 {
					log.Printf("exiting with statusCode: %d", res.StatusCode)
					result.Errors = append(result.Errors, res.Status, res.StatusCode)
				} else {
					err = json.Unmarshal([]byte(body), &result)
					lib.CheckError(err)
				}

				r <- result
			}
		}
	}()

	return r
}
