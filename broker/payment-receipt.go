package broker

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	trx "github.com/wopta/goworkspace/transaction"
)

const (
	filenameFormat = "Quietanza Pagamento Polizza %s rata %s %d"
)

type paymentReceiptResp struct {
	Filename string `json:"filename"`
	RawDoc   string `json:"rawDoc"`
}

func PaymentReceiptFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
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

	rawDoc, filename, err := paymentReceiptBuilder(transactionUid, authToken)
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

func paymentReceiptBuilder(transactionUID string, authToken lib.AuthToken) (string, string, error) {
	transaction := trx.GetTransactionByUid(transactionUID, "")
	if transaction == nil {
		return "", "", errors.New("transaction not found")
	}

	if !transaction.IsPay {
		return "", "", errors.New("transaction is not pay")
	}

	// TODO: what if transaction refers to a renewPolicy
	policy, err := plc.GetPolicy(transaction.PolicyUid, "")
	if err != nil {
		return "", "", err
	}

	if authToken.Role != models.UserRoleAdmin {
		if policy.ProducerUid == "" {
			return "", "", fmt.Errorf("node %s cannot access policy %s", authToken.UserID, policy.Uid)
		}

		node := network.GetNetworkNodeByUid(policy.ProducerUid)
		if node == nil {
			return "", "", errors.New("node not found")
		}

		if policy.ProducerUid != node.Uid && !network.IsChildOf(node.Uid, policy.ProducerUid) {
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
	customerInfo := document.CustomerInfo{
		Fullname:   policy.Contractor.Name + " " + policy.Contractor.Surname,
		Address:    policy.Contractor.Residence.StreetName + " " + policy.Contractor.Residence.StreetNumber,
		PostalCode: policy.Contractor.Residence.PostalCode,
		City:       policy.Contractor.Residence.City,
		Province:   policy.Contractor.Residence.CityCode,
		Email:      policy.Contractor.Mail,
		Phone:      policy.Contractor.Phone,
	}

	expirationDate := lib.AddMonths(transaction.EffectiveDate, 12)
	nextPayment := lib.AddMonths(transaction.EffectiveDate, 12)
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		expirationDate = lib.AddMonths(transaction.EffectiveDate, 1).AddDate(0, 0, -1)
		nextPayment = lib.AddMonths(transaction.EffectiveDate, 1)
	}

	transactionInfo := document.TransactionInfo{
		PolicyCode:     policy.CodeCompany,
		EffectiveDate:  transaction.EffectiveDate.Format("02/01/2006"),
		ExpirationDate: expirationDate.Format("02/01/2006"),
		PriceGross:     fmt.Sprintf("%.2f", transaction.Amount),
		NextPayment:    nextPayment.Format("02/01/2006"),
	}

	return document.ReceiptInfo{
		CustomerInfo: customerInfo,
		Transaction:  transactionInfo,
	}
}
