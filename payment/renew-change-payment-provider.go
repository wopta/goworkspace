package payment

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment/common"
	plcRenew "github.com/wopta/goworkspace/policy/renew"
	prd "github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/transaction"
	trsRenew "github.com/wopta/goworkspace/transaction/renew"
)

/*
This should be a temporary handler while the imported policies by an agent that works
with an online warrant are not set to the correct provider at import time
*/
func RenewChangePaymentProviderFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                  error
		rawResp              []byte
		payUrl               string
		policy               models.Policy
		activeTransactions   []models.Transaction
		updatedTransactions  []models.Transaction
		req                  ChangePaymentProviderReq
		resp                 ChangePaymentProviderResp
		responseTransactions = make([]models.Transaction, 0)
		unpaidTransactions   = make([]models.Transaction, 0)
	)

	log.SetPrefix("[RenewChangePaymentProviderFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding body")
		return "", nil, err
	}

	if policy, err = plcRenew.GetRenewPolicyByUid(req.PolicyUid); err != nil {
		log.Println("error getting renew policy")
		return "", nil, err
	}

	if strings.EqualFold(policy.Payment, req.ProviderName) {
		log.Printf("can't change payment method to %s for policy with payment method %s", req.ProviderName, policy.Payment)
		return "", nil, errors.New("unable to change payment method")
	}

	if activeTransactions, err = trsRenew.GetRenewActiveTransactionsByPolicyUid(policy.Uid, policy.Annuity); err != nil {
		log.Println("error getting renew transactions")
		return "", nil, err
	}

	for _, tr := range activeTransactions {
		if tr.IsPay {
			responseTransactions = append(responseTransactions, tr)
			continue
		}
		transaction.ReinitializePaymentInfo(&tr, policy.Payment)
		unpaidTransactions = append(unpaidTransactions, tr)
	}

	if len(unpaidTransactions) == 0 {
		log.Printf("no unpaid transactions found for policy %s", policy.Uid)
		return "", nil, errors.New("no unpaid transactions to update")
	}

	policy.SanitizePaymentData()
	policy.Payment = req.ProviderName

	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, nil, nil)

	client := NewClient(policy.Payment, policy, *product, unpaidTransactions, req.ScheduleFirstRate, "")
	payUrl, updatedTransactions, err = client.Update()
	if err != nil {
		log.Printf("error changing payment provider to %s: %s", req.ProviderName, err.Error())
		return "", nil, err
	}

	responseTransactions = append(responseTransactions, updatedTransactions...)
	policy.PayUrl = payUrl
	policy.Updated = time.Now().UTC()
	policy.BigQueryParse()

	if err = common.SaveTransactionsToDB(updatedTransactions, lib.RenewTransactionCollection); err != nil {
		log.Println("error saving transactions")
		return "", nil, err
	}

	if err = lib.SetFirestoreErr(lib.RenewPolicyCollection, policy.Uid, policy); err != nil {
		log.Println("error saving policy to firestore")
		return "", nil, err
	}

	if err = lib.InsertRowsBigQuery(lib.WoptaDataset, lib.RenewPolicyCollection, policy); err != nil {
		log.Println("error saving policy to bigquery")
		return "", nil, err
	}

	resp.Policy = policy
	resp.Transactions = responseTransactions

	if rawResp, err = json.Marshal(resp); err != nil {
		log.Println("error marshaling response")
		return "", nil, err
	}

	return string(rawResp), resp, err
}
