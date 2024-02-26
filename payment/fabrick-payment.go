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
		grossAmounts             []float64
		nettAmounts              []float64
	)

	durationInYears := policy.GetDurationInYears()
	grossAmounts = make([]float64, durationInYears)
	nettAmounts = make([]float64, durationInYears)

	isRecurrent := policy.PaymentMode == models.PaymentModeRecurrent

	customerId := uuid.New().String()

	scheduleDates := getTransactionScheduleDates(policy, origin)

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		priceGross = policy.PriceGrossMonthly
		priceNett = policy.PriceNettMonthly
	}

	if policy.PaymentSplit == string(models.PaySplitYearly) {
		for _, guarantee := range policy.Assets[0].Guarantees {
			for rateIndex := 0; rateIndex < guarantee.Value.Duration.Year; rateIndex++ {
				grossAmounts[rateIndex] += guarantee.Value.PremiumGrossYearly
				nettAmounts[rateIndex] += guarantee.Value.PremiumNetYearly
			}
		}
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

		if policy.PaymentSplit == string(models.PaySplitYearly) {
			priceGross = grossAmounts[index]
			priceNett = nettAmounts[index]
		}

		priceGross = lib.RoundFloat(priceGross, 2)
		priceNett = lib.RoundFloat(priceNett, 2)

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

func getTransactionScheduleDates(policy *models.Policy, origin string) []time.Time {
	var (
		currentScheduleDate time.Time
		response            []time.Time = make([]time.Time, 0)
		yearDuration        int         = 1
	)

	activeTransactions := transaction.GetPolicyActiveTransactions(origin, policy.Uid)

	if len(activeTransactions) == 0 {
		if policy.PaymentMode == models.PaymentModeRecurrent && policy.PaymentSplit == string(models.PaySplitYearly) {
			yearDuration = policy.GetDurationInYears()
		}

		numberOfRates := policy.GetNumberOfRates() * yearDuration

		for i := 0; i < numberOfRates; i++ {
			if i > 0 {
				switch policy.PaymentSplit {
				case string(models.PaySplitYearly):
					currentScheduleDate = policy.StartDate.AddDate(i, 0, 0)
				case string(models.PaySplitMonthly):
					currentScheduleDate = policy.StartDate.AddDate(0, i, 0)
				default:
					log.Printf("unhandled recurrent payment split: %s", policy.PaymentSplit)
					return nil
				}
			}
			response = append(response, currentScheduleDate)
		}
	} else {
		// TODO: handle when policy already has created transactions
		// isFirstSchedule := true
		// for _, tr := range activeTransactions {
		// 	if tr.IsPay {
		// 		continue
		// 	}
		// 	if isFirstSchedule {
		// 		currentScheduleDate = time.Time{}
		// 		isFirstSchedule = false
		// 	}
		// 	currentScheduleDate, err := time.Parse(models.TimeDateOnly, tr.ScheduleDate)
		// 	if err != nil {
		// 		log.Printf("error parsing schedule date %s: %s", tr.ScheduleDate, err.Error())
		// 		return nil
		// 	}
		// 	response = append(response, currentScheduleDate)
		// }
	}

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

		res, err := lib.RetryDo(request, 5)
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
