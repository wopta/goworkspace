package payment

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
)

func fabrickPayment(
	policy *models.Policy,
	origin string,
	paymentMethods []string,
	mgaProduct *models.Product,
) FabrickPaymentResponse {
	// TODO: refactor me - all transactions calculations should be done
	// beforehand and should be independent of the payment provider
	log.Printf("generating fabrick payment...")
	log.Printf("policy '%s' payment configuration", policy.Uid)
	log.Printf("provider: %s", policy.Payment)
	log.Printf("paymentSplit: %s", policy.PaymentSplit)
	log.Printf("paymentMode: %s", policy.PaymentMode)

	var (
		response                 FabrickPaymentResponse
		scheduleDate, expireDate string
		priceGross               float64 = policy.PriceGross
		priceNett                float64 = policy.PriceNett
	)

	isRecurrent := policy.PaymentMode == models.PaymentModeRecurrent

	customerId := uuid.New().String()

	scheduleDates := transaction.GetTransactionScheduleDates(policy)

	grossAmounts, nettAmounts := transaction.GetTransactionsAmounts(policy)
	if len(grossAmounts) == 0 || len(nettAmounts) == 0 {
		log.Println("error creating fabrick transactions: empty amounts list")
		return FabrickPaymentResponse{}
	}
	if len(scheduleDates) != len(grossAmounts) || len(scheduleDates) != len(nettAmounts) {
		log.Println("error creating fabrick transactions: schedule dates length doesn't match amounts lengths")
		return FabrickPaymentResponse{}
	}
	log.Printf("creating %d transaction(s)...", len(scheduleDates))

	for index, sd := range scheduleDates {
		isFirstRate := index == 0
		createMandate := isRecurrent && isFirstRate
		if !sd.IsZero() {
			scheduleDate = sd.Format(models.TimeDateOnly)
			expireDate = sd.AddDate(10, 0, 0).Format(models.TimeDateOnly)
		} else {
			scheduleDate = ""
			expireDate = policy.StartDate.AddDate(10, 0, 0).Format(models.TimeDateOnly)
		}
		effectiveDate := policy.StartDate.AddDate(0, index, 0)

		priceGross = lib.RoundFloat(grossAmounts[index], 2)
		priceNett = lib.RoundFloat(nettAmounts[index], 2)

		log.Printf("creating transaction with index '%d' and schedule date '%s' and amount '%.2f' ...", index, scheduleDate, priceGross)

		res := <-createFabrickTransaction(
			policy,
			createMandate,
			scheduleDate,
			expireDate,
			customerId,
			priceGross,
			priceNett,
			origin,
			paymentMethods,
			mgaProduct,
			effectiveDate,
		)

		if len(res.Errors) > 0 {
			log.Printf("error creating fabrick transaction: %v", res.Errors)
			return FabrickPaymentResponse{}
		}

		if isFirstRate {
			response = res
		}
		log.Printf("response: %v", res)
		time.Sleep(100 * time.Nanosecond)
	}

	log.Printf("payment generated: %v", response)

	return response
}

func createFabrickTransaction(
	policy *models.Policy,
	firstSchedule bool,
	scheduleDate, expireDate, customerId string,
	amount, amountNet float64,
	origin string,
	paymentMethods []string,
	mgaProduct *models.Product,
	effectiveDate time.Time,
) <-chan FabrickPaymentResponse {
	r := make(chan FabrickPaymentResponse)

	go func() {
		defer close(r)

		body := getFabrickRequestBody(policy, firstSchedule, scheduleDate, expireDate, customerId, amount, origin, paymentMethods)
		if body == "" {
			return
		}
		request := getFabrickPaymentRequest(body)
		if request == nil {
			return
		}

		log.Printf("policy '%s' request headers: %s", policy.Uid, request.Header)
		log.Printf("policy '%s' request body: %s", policy.Uid, request.Body)

		res, err := lib.RetryDo(request, 5, 10)
		lib.CheckError(err)

		if res != nil {
			log.Printf("policy '%s' response headers: %s", policy.Uid, res.Header)
			body, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			lib.CheckError(err)
			log.Printf("policy '%s' response body: %s", policy.Uid, string(body))

			var result FabrickPaymentResponse

			if res.StatusCode != 200 {
				log.Printf("exiting with statusCode: %d", res.StatusCode)
				result.Errors = append(result.Errors, res.Status, res.StatusCode)
			} else {
				err = json.Unmarshal([]byte(body), &result)
				lib.CheckError(err)

				// TODO: handle result without payment id
				transaction.PutByPolicy(*policy, scheduleDate, origin, expireDate, customerId, amount, amountNet, *result.Payload.PaymentID, "", false, mgaProduct, effectiveDate)
			}

			r <- result
		}
	}()

	return r
}

const (
	currency         string = "EUR"
	expireDateFormat string = "2006-01-02T15:04:05.999999999Z"
	redirectUrl      string = "https://www.wopta.it"
	woptaMerchantId  string = "wop134b31-5926-4b26-1411-726bc9f0b111"
)

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

	pay := FabrickPaymentsRequest{
		MerchantID: woptaMerchantId,
		ExternalID: externalId,
		PaymentConfiguration: PaymentConfiguration{
			PaymentPageRedirectUrls: PaymentPageRedirectUrls{
				OnSuccess: redirectUrl,
				OnFailure: redirectUrl,
			},
			ExpirationDate: expireDate,
			AllowedPaymentMethods: &[]AllowedPaymentMethod{{
				Role: "payer",
				PaymentMethods: lib.SliceMap(
					paymentMethods,
					func(item string) string { return strings.ToUpper(item) },
				),
			}},
			CallbackURL: callbackUrl,
		},
		Bill: Bill{
			ExternalID:      externalId,
			Amount:          amount,
			Currency:        currency,
			Description:     fmt.Sprintf("Pagamento polizza n° %s", policy.CodeCompany),
			MandateCreation: mandate,
			Items: []Item{{
				ExternalID: externalId,
				Amount:     amount,
				Currency:   currency,
			}},
			Subjects: &[]Subject{{
				ExternalID: customerId,
				Role:       "customer",
				Email:      policy.Contractor.Mail,
				Name:       strings.Join([]string{policy.Contractor.Name, policy.Contractor.Surname}, " "),
			}},
		},
	}

	if scheduleDate != "" {
		pay.Bill.ScheduleTransaction = &ScheduleTransaction{
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
