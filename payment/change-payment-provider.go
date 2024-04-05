package payment

import (
	"encoding/json"
	"errors"
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

type ChangePaymentProviderResp struct {
	Policy       models.Policy        `json:"policy"`
	Transactions []models.Transaction `json:"transactions"`
}

/*
Now this function is used only to change payment provider to Fabrick (info hardcoded in frontend call) for those
policies that have been imported. When we will have multi providers we should delete transactions schedule from
old provider systems and only then send schedule new transactions to new provider systems.
*/
func ChangePaymentProviderFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                  error
		payUrl               string
		policy               models.Policy
		updatedTransactions  []models.Transaction
		req                  ChangePaymentProviderReq
		resp                 ChangePaymentProviderResp
		responseTransactions = make([]models.Transaction, 0)
		unpaidTransactions   = make([]models.Transaction, 0)
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

	policy.SanitizePaymentData()

	if strings.EqualFold(policy.Payment, req.ProviderName) {
		log.Printf("can't change payment method to %s for policy with payment method %s", req.ProviderName, policy.Payment)
		return "{}", nil, errors.New("unable to change payment method")
	}

	unpaidTransactions = transaction.GetPolicyUnpaidTransactions(policy.Uid)
	if len(unpaidTransactions) == 0 {
		log.Printf("no unpaid transactions found for policy %s", policy.Uid)
		return "{}", nil, err
	}

	policy.Payment = req.ProviderName
	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, nil, nil)

	payUrl, updatedTransactions, err = PaymentControllerV2(policy, *product, unpaidTransactions)
	if err != nil {
		log.Printf("error changing payment provider to %s: %s", req.ProviderName, err.Error())
		return "{}", nil, err
	}

	policy.PayUrl = payUrl
	responseTransactions = append(responseTransactions, updatedTransactions...)
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

	resp.Policy = policy
	resp.Transactions = responseTransactions
	rawResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshaling response: %s", err.Error())
		return "{}", nil, err
	}

	log.Println("Handler End -------------------------------------------------")

	return string(rawResp), resp, err
}
