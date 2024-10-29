package broker

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	plcRenew "github.com/wopta/goworkspace/policy/renew"
	plcUtils "github.com/wopta/goworkspace/policy/utils"
	trx "github.com/wopta/goworkspace/transaction"
	trxRenew "github.com/wopta/goworkspace/transaction/renew"
)

const (
	filenameFormat = "Quietanza Pagamento Polizza %s rata %s %d.pdf"
)

type paymentReceiptResp struct {
	Filename string `json:"filename"`
	RawDoc   string `json:"rawDoc"`
}

func PaymentReceiptFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err     error
		isRenew bool
	)

	log.SetPrefix("[PaymentReceiptFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end -------------------------------------------------")
	}()

	log.Println("Handler start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.Printf("error fetching authToken: %s", err.Error())
		return "", nil, err
	}

	transactionUid := chi.URLParam(r, "uid")
	if transactionUid == "" {
		return "", "", errors.New("transaction uid is empty")
	}

	param := r.URL.Query().Get("isRenew")
	if isRenew, err = strconv.ParseBool(param); err != nil {
		log.Printf("error parsing isRenew: %s", err.Error())
		return "", "", err
	}

	rawDoc, filename, err := paymentReceiptBuilder(transactionUid, authToken, isRenew)
	if err != nil {
		log.Printf("error building raw doc: %s", err.Error())
		return "", "", err
	}

	resp := paymentReceiptResp{
		Filename: filename,
		RawDoc:   rawDoc,
	}

	rawResp, err := json.Marshal(resp)

	return string(rawResp), resp, err
}

func paymentReceiptBuilder(transactionUID string, authToken lib.AuthToken, isRenew bool) (string, string, error) {
	var (
		err         error
		policy      models.Policy
		transaction *models.Transaction
	)

	if isRenew {
		transaction = trxRenew.GetRenewTransactionByUid(transactionUID)
		if transaction == nil {
			return "", "", errors.New("transaction not found")
		}
		if !transaction.IsPay {
			return "", "", errors.New("transaction is not paid")
		}
		policy, err = plcRenew.GetRenewPolicyByUid(transaction.PolicyUid)
		if err != nil {
			return "", "", err
		}
	} else {
		transaction = trx.GetTransactionByUid(transactionUID, "")
		if transaction == nil {
			return "", "", errors.New("transaction not found")
		}
		if !transaction.IsPay {
			return "", "", errors.New("transaction is not paid")
		}
		policy, err = plc.GetPolicy(transaction.PolicyUid, "")
		if err != nil {
			return "", "", err
		}
	}

	if authToken.Role != models.UserRoleAdmin {
		if policy.ProducerUid == "" {
			return "", "", fmt.Errorf("node %s cannot access policy %s", authToken.UserID, policy.Uid)
		}

		node := network.GetNetworkNodeByUid(policy.ProducerUid)
		if node == nil {
			return "", "", errors.New("node not found")
		}

		if !plcUtils.CanBeAccessedBy(authToken.Role, policy.ProducerUid, node.Uid) {
			return "", "", fmt.Errorf("node %s cannot access policy %s", authToken.UserID, policy.Uid)
		}
	}

	receiptInfo := receiptInfoBuilder(policy, *transaction)

	doc, err := document.PaymentReceipt(receiptInfo)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return "", "", err
	}

	rawDoc := base64.StdEncoding.EncodeToString(doc)
	filename := fmt.Sprintf(filenameFormat, policy.CodeCompany, lib.ExtractLocalMonth(transaction.EffectiveDate),
		transaction.EffectiveDate.Year())

	return rawDoc, filename, nil
}

func receiptInfoBuilder(policy models.Policy, transaction models.Transaction) document.ReceiptInfo {
	const dateFormat = "02/01/2006"

	customerInfo := document.CustomerInfo{
		Fullname:   policy.Contractor.Name + " " + policy.Contractor.Surname,
		Address:    policy.Contractor.Residence.StreetName + " " + policy.Contractor.Residence.StreetNumber,
		PostalCode: policy.Contractor.Residence.PostalCode,
		City:       policy.Contractor.Residence.City,
		Province:   policy.Contractor.Residence.CityCode,
		Email:      policy.Contractor.Mail,
		Phone:      policy.Contractor.Phone,
	}

	expirationDate := policy.EndDate.AddDate(0, 0, -1)
	nextPayment := "====="

	tmpExpirationDate := lib.AddMonths(transaction.EffectiveDate, 12)
	tmpNextPayment := lib.AddMonths(transaction.EffectiveDate, 12)

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		tmpExpirationDate = lib.AddMonths(transaction.EffectiveDate, 1).AddDate(0, 0, -1)
		tmpNextPayment = lib.AddMonths(transaction.EffectiveDate, 1)
	}

	if !tmpExpirationDate.After(expirationDate) {
		expirationDate = tmpExpirationDate
	}
	if !tmpNextPayment.After(policy.EndDate) {
		nextPayment = tmpNextPayment.Format(dateFormat)
	}

	transactionInfo := document.TransactionInfo{
		PolicyCode:     policy.CodeCompany,
		EffectiveDate:  transaction.EffectiveDate.Format(dateFormat),
		ExpirationDate: expirationDate.Format(dateFormat),
		PriceGross:     humanize.FormatFloat("#.###,##", transaction.Amount) + " €",
		NextPayment:    nextPayment,
	}

	return document.ReceiptInfo{
		CustomerInfo: customerInfo,
		Transaction:  transactionInfo,
	}
}
