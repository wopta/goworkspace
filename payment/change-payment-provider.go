package payment

import (
	"encoding/json"
	"errors"
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
)

type ChangePaymentProviderReq struct {
	PolicyUid    string `json:"policyUid"`
	ProviderName string `json:"providerName"`
}

func ChangePaymentProviderFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
		req    ChangePaymentProviderReq
	)

	log.SetPrefix("ChangePaymentProviderFx ")
	defer log.SetPrefix("")
	log.Println("Handler Start -----------------------------------------------")

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

	activeTransactions := transaction.GetPolicyActiveTransactions("", policy.Uid)
	if len(activeTransactions) == 0 {
		log.Printf("no active transactions found for policy %s", policy.Uid)
		return "{}", nil, err
	}

	customerId := uuid.New().String()
	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, nil, nil)
	//mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	paymentMethods := getPaymentMethods(policy, product)

	for index, tr := range activeTransactions {
		b := getFabrickRequestBody(&policy, index == 0, tr.ScheduleDate, tr.ExpirationDate, customerId, tr.Amount,
			"", paymentMethods)
		if b == "" {
			return "{}", nil, errors.New("fabrick error")
		}
		request := getFabrickPaymentRequest(b)
		if request == nil {
			return "{}", nil, errors.New("fabrick error")
		}
		res, err := lib.RetryDo(request, 5, 10)
		if err != nil {
			return "", nil, err
		}
		if res != nil {
			if res.StatusCode == http.StatusOK {
				resBody, err := io.ReadAll(res.Body)
				defer res.Body.Close()
				if err != nil {
					return "", nil, err
				}
				var result FabrickPaymentResponse
				err = json.Unmarshal(resBody, &result)
				if err != nil {
					return "", nil, err
				}

				if index == 0 {
					// TODO: handle nil pointer
					policy.PayUrl = *result.Payload.PaymentPageURL
				}

				activeTransactions[index].ProviderName = req.ProviderName
				// TODO: handle nil pointer
				activeTransactions[index].ProviderId = *result.Payload.PaymentID

				err = lib.SetFirestoreErr(models.TransactionsCollection, tr.Uid, activeTransactions[index])
				if err != nil {
					return "{}", nil, err
				}

				activeTransactions[index].BigQuerySave("")
			}
		}

	}

	policy.Payment = req.ProviderName

	err = lib.SetFirestoreErr(models.PolicyCollection, policy.Uid, policy)
	if err != nil {
		return "{}", nil, err
	}

	policy.BigquerySave("")

	log.Println("Handler End -------------------------------------------------")

	return "{}", nil, err
}
