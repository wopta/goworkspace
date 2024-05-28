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
)

const (
	currency         string = "EUR"
	expireDateFormat string = "2006-01-02T15:04:05.999999999Z"
	redirectUrl      string = "https://www.wopta.it"
	woptaMerchantId  string = "wop134b31-5926-4b26-1411-726bc9f0b111"
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
	log.Println("[getFabrickClient]", req)

	return client.Do(req)
}

func getFabrickRequestBody(
	policy *models.Policy,
	createMandate, scheduleFirstRate, isFirstRate bool,
	scheduleDate, expireDate, customerId string,
	amount float64,
	origin string,
	paymentMethods []string,
) string {
	var (
		callbackFormat      string    = "%scallback/v1/payment/fabrick/%s?uid=%s&schedule=%s&token=%s&origin=%s"
		callbackEndpoint    string    = "single-rate"
		mandate             string    = "false"
		now                 time.Time = time.Now().UTC()
		requestScheduleDate string    = scheduleDate
	)

	if createMandate {
		mandate = "true"
		if !scheduleFirstRate {
			scheduleDate = ""
		}
	}

	if customerId == "" {
		customerId = uuid.New().String()
	}

	if requestScheduleDate == "" {
		requestScheduleDate = now.Format(models.TimeDateOnly)
	}

	if isFirstRate {
		callbackEndpoint = "first-rate"
	}

	externalId := strings.Join([]string{
		policy.Uid,
		requestScheduleDate,
		policy.CodeCompany,
		strings.ReplaceAll(origin, "https://", ""),
		strconv.FormatInt(now.Unix(), 10),
	}, "_")

	callbackUrl := fmt.Sprintf(
		callbackFormat,
		os.Getenv("WOPTA_BASE_URL"),
		callbackEndpoint,
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

func createFabrickTransaction(
	policy *models.Policy,
	transaction models.Transaction,
	createMandate, scheduleFirstRate, isFirstRate bool,
	customerId string,
	paymentMethods []string,
) <-chan FabrickPaymentResponse {
	r := make(chan FabrickPaymentResponse)

	go func() {
		defer close(r)

		body := getFabrickRequestBody(policy, createMandate, scheduleFirstRate, isFirstRate, transaction.ScheduleDate, transaction.ExpirationDate,
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
			url := fmt.Sprintf("www.dev.wopta.it/%s", transaction.Uid)
			r <- FabrickPaymentResponse{
				Status: &status,
				Errors: nil,
				Payload: &Payload{
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

				var result FabrickPaymentResponse

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

func fabrickExpireBill(providerId string) error {
	var urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments/expirationDate"
	const expirationTimeSuffix = "00:00:00"

	log.Println("starting fabrick expire bill request...")

	expirationDate := fmt.Sprintf(
		"%s %s",
		time.Now().UTC().AddDate(0, 0, -1).Format(models.TimeDateOnly),
		expirationTimeSuffix,
	)

	requestBody := fmt.Sprintf(`{"id":"%s","newExpirationDate":"%s"}`, providerId, expirationDate)
	log.Printf("fabrick expire bill request body: %s", requestBody)

	req, err := http.NewRequest(http.MethodPut, urlstring, strings.NewReader(requestBody))
	if err != nil {
		log.Printf("error creating request: %s", err.Error())
		return err
	}
	res, err := getFabrickClient(urlstring, req)
	if err != nil {
		log.Printf("error getting response: %s", err.Error())
		return err
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("fabrick expire bill response error: %s", err.Error())
		return err
	}
	log.Println("fabrick expire bill response body: ", string(respBody))
	if res.StatusCode != http.StatusOK {
		log.Printf("fabrick expire bill error status %s", res.Status)
		return fmt.Errorf("fabrick expire bill error status %s", res.Status)
	}

	log.Println("fabrick expire bill completed!")

	return nil
}

func fabrickHasMandate(userToken string) (bool, error) {
	if userToken == "" {
		return false, nil
	}

	var (
		urlFormat string = "%s/api/fabrick/pace/v4.0/mods/back/v1.0/payment-instruments?merchantId=%s&subjectXId=%s&status=ACTIVE"
		response  fabrickPaymentInstrumentRes
		found     bool
		url       string = fmt.Sprintf(urlFormat, os.Getenv("FABRICK_BASEURL"), os.Getenv("FABRICK_MERCHANT_ID"), userToken)
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	if os.Getenv("env") == "local" || os.Getenv("env") == "local-test" {
		st := "INACTIVE"
		if userToken == "user-has-token" {
			st = "ACTIVE"
		}
		payload := make([]paymentInstrument, 0)
		payload = append(payload, paymentInstrument{
			Status: st,
		})
		response = fabrickPaymentInstrumentRes{
			Status:  "200",
			Errors:  nil,
			Payload: payload,
		}
	} else {
		res, err := getFabrickClient("", req)
		if err != nil {
			log.Printf("error getting response: %s", err.Error())
			return false, err
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("error status %s", res.Status)
			return false, fmt.Errorf("error status %s", res.Status)
		}

		err = json.NewDecoder(res.Body).Decode(&response)
		defer res.Body.Close()
		if err != nil {
			log.Printf("response error: %s", err.Error())
			return false, err
		}
		log.Printf("response: %+v", response)
	}

	for _, p := range response.Payload {
		if p.Status == "ACTIVE" {
			found = true
			break
		}
	}

	return found, nil
}

type fabrickPaymentInstrumentRes struct {
	Status  string              `json:"status"`
	Errors  []any               `json:"errors"`
	Payload []paymentInstrument `json:"payload"`
}

type paymentInstrument struct {
	Type              string    `json:"type"`
	CreationDate      time.Time `json:"creationDate"`
	ExpiryDate        string    `json:"expiryDate"`
	Status            string    `json:"status"`
	Alias             string    `json:"alias"`
	MakeDefault       bool      `json:"makeDefault"`
	SubjectId         string    `json:"subjectId"`
	SubjectXId        string    `json:"subjectXId"`
	MatchedDossierXId []any     `json:"matchedDossierXId"`
	Xid               string    `json:"xid"`
}
