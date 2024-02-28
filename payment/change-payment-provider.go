package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	prd "github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/transaction"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type ChangePaymentProviderReq struct {
	PolicyUid    string `json:"policyUid"`
	ProviderName string `json:"providerName"`
}

func ChangePaymentProviderFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                 error
		payUrl              string
		policy              models.Policy
		updatedTransactions []models.Transaction
		req                 ChangePaymentProviderReq
	)

	log.SetPrefix("ChangePaymentProviderFx ")
	defer log.SetPrefix("")
	log.Println("Handler Start -----------------------------------------------")

	origin := r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("req body: %s", string(body))
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshaling request body: %s", string(body))
		return "{}", nil, err
	}

	policy, err = plc.GetPolicy(req.PolicyUid, "")
	if err != nil {
		log.Printf("no policy found with uid %s: %s", req.PolicyUid, err.Error())
		return "{}", nil, err
	}

	if strings.EqualFold(policy.Payment, req.ProviderName) {
		log.Printf("can't change payment method to %s for policy with payment method %s", req.ProviderName, policy.Payment)
		return "{}", nil, errors.New("unable to change payment method")
	}

	unpaidTransactions := transaction.GetPolicyUnpaidTransactions(policy.Uid)
	if len(unpaidTransactions) == 0 {
		log.Printf("no active transactions found for policy %s", policy.Uid)
		return "{}", nil, err
	}

	policy.Payment = req.ProviderName
	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, nil, nil)
	paymentMethods := getPaymentMethods(policy, product)
	if len(paymentMethods) == 0 {
		log.Printf("no payment methods found for provider %s", req.ProviderName)
		return "{}", nil, errors.New("no payment methods found")
	}

	switch req.ProviderName {
	case models.FabrickPaymentProvider:
		payUrl, updatedTransactions, err = changePaymentProviderToFabrick(origin, policy, unpaidTransactions, paymentMethods)
	default:
		log.Printf("payment provider %s not supported", req.ProviderName)
		return "{}", nil, errors.New("payment provider not supported")
	}

	if err != nil {
		log.Printf("error changing payment provider to %s: %s", req.ProviderName, err.Error())
		return "{}", nil, err
	}

	policy.PayUrl = payUrl
	for _, tr := range updatedTransactions {
		err = lib.SetFirestoreErr(models.TransactionsCollection, tr.Uid, tr)
		if err != nil {
			return "{}", nil, err
		}

		tr.BigQuerySave("")
	}

	err = lib.SetFirestoreErr(models.PolicyCollection, policy.Uid, policy)
	if err != nil {
		return "{}", nil, err
	}

	policy.BigquerySave("")

	log.Println("Handler End -------------------------------------------------")

	return "{}", nil, err
}

func changePaymentProviderToFabrick(origin string, policy models.Policy, transactions []models.Transaction, paymentMethods []string) (string, []models.Transaction, error) {
	var (
		err    error
		payUrl string
	)

	now := time.Now().UTC()
	customerId := uuid.New().String()

	for index, tr := range transactions {
		if index == 0 {
			tr.ScheduleDate = now.Format(models.TimeDateOnly)
			tr.ExpirationDate = now.AddDate(10, 0, 0).Format(models.TimeDateOnly)
		}
		b := getFabrickRequestBody(&policy, index == 0, tr.ScheduleDate, tr.ExpirationDate, customerId, tr.Amount,
			origin, paymentMethods)
		if b == "" {
			return "", nil, errors.New("unable to get fabrick request body")
		}
		request := getFabrickPaymentRequest(b)
		if request == nil {
			return "", nil, errors.New("unable to get fabrick request for policy")
		}
		res, err := lib.RetryDo(request, 5, 10)
		if err != nil {
			log.Printf("error retryDo fabrick request: %s", err.Error())
			return "", nil, err
		}
		if res != nil && res.StatusCode == http.StatusOK {
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Printf("error reading fabrick response body policy %s: %s", policy.Uid, err.Error())
				return "", nil, err
			}
			var result FabrickPaymentResponse
			err = json.Unmarshal(resBody, &result)
			if err != nil {
				log.Printf("error unmarshaling fabrick response policy %s: %s", policy.Uid, err.Error())
				return "", nil, err
			}
			res.Body.Close()

			if index == 0 {
				payUrl = *result.Payload.PaymentPageURL
				transactions[index].ScheduleDate = now.Format(models.TimeDateOnly)
			}

			transactions[index].ProviderName = models.FabrickPaymentProvider
			if result.Payload.PaymentID == nil {
				log.Printf("error nil paymentID fabrick transaction %s", tr.Uid)
				return "", nil, errors.New("error fabrick nil paymentID")
			}
			transactions[index].ProviderId = *result.Payload.PaymentID
			transactions[index].UserToken = customerId
			transactions[index].UpdateDate = time.Now().UTC()
		} else {
			return "", nil, fmt.Errorf("fabrick error: %s", res.Status)
		}
	}

	return payUrl, transactions, err
}
