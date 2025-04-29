package broker

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

var (
	errMissingParams       = errors.New("transaction uid param is empty")
	errTransactionNotFound = errors.New("transaction not found")
	errTransactionDeleted  = errors.New("transaction is deleted")
	errNodeNotFound        = errors.New("node not found")
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
		log.SetPrefix("")
	}()

	log.Println("Handler start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.Println("error fetching authToken")
		return "", nil, err
	}

	transactionUid := chi.URLParam(r, "uid")
	if transactionUid == "" {
		err = errMissingParams
		return "", "", err
	}

	param := r.URL.Query().Get("isRenew")
	if isRenew, err = strconv.ParseBool(param); err != nil {
		log.Println("error parsing isRenew")
		return "", "", err
	}

	rawDoc, filename, err := paymentReceiptBuilder(transactionUid, authToken, isRenew)
	if err != nil {
		log.Println("error building raw doc")
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
			return "", "", errTransactionNotFound
		}
		if transaction.IsDelete {
			return "", "", errTransactionDeleted
		}
		policy, err = plcRenew.GetRenewPolicyByUid(transaction.PolicyUid)
		if err != nil {
			return "", "", err
		}
	} else {
		transaction = trx.GetTransactionByUid(transactionUID, "")
		if transaction == nil {
			return "", "", errTransactionNotFound
		}
		if transaction.IsDelete {
			return "", "", errTransactionDeleted
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
			return "", "", errNodeNotFound
		}

		if !plcUtils.CanBeAccessedBy(authToken.Role, policy.ProducerUid, node.Uid) {
			return "", "", fmt.Errorf("node %s cannot access policy %s", authToken.UserID, policy.Uid)
		}
	}

	receiptInfo, err := receiptInfoBuilder(policy, *transaction)
	if err != nil {
		return "", "", err
	}

	doc, err := document.PaymentReceipt(receiptInfo)
	if err != nil {
		return "", "", err
	}

	rawDoc := base64.StdEncoding.EncodeToString(doc)
	filename := fmt.Sprintf(filenameFormat, policy.CodeCompany, lib.ExtractLocalMonth(transaction.EffectiveDate),
		transaction.EffectiveDate.Year())

	return rawDoc, filename, nil
}

func receiptInfoBuilder(policy models.Policy, transaction models.Transaction) (document.ReceiptInfo, error) {
	const dateFormat = "02/01/2006"

	receiptInfo := document.NewReceiptInfo()

	if policy.Company != "" {
		receiptInfo.PolicyInfo.Company = document.CompanyMap[policy.Company]
	}
	if policy.NameDesc != "" {
		receiptInfo.PolicyInfo.ProductDescription = strings.ToUpper(policy.NameDesc)
	}
	if policy.CodeCompany != "" {
		receiptInfo.PolicyInfo.Code = policy.CodeCompany
	}

	if policy.Contractor.Name != "" {
		receiptInfo.CustomerInfo.Fullname = strings.TrimSpace(policy.Contractor.Name + " " + policy.Contractor.Surname)
	}

	address := policy.Contractor.Residence
	if policy.Contractor.Type == models.UserLegalEntity {
		address = policy.Contractor.CompanyAddress
	}

	if address != nil {
		if address.StreetName != "" {
			receiptInfo.CustomerInfo.Address = strings.TrimSpace(address.StreetName + " " + address.StreetNumber)
		}
		if address.PostalCode != "" {
			receiptInfo.CustomerInfo.PostalCode = address.PostalCode
		}
		if address.City != "" {
			receiptInfo.CustomerInfo.City = address.City
		}
		if address.CityCode != "" {
			receiptInfo.CustomerInfo.Province = address.CityCode
		}
	}

	if policy.Contractor.Mail != "" {
		receiptInfo.CustomerInfo.Email = policy.Contractor.Mail
	}
	if policy.Contractor.Phone != "" {
		receiptInfo.CustomerInfo.Phone = policy.Contractor.Phone
	}

	expirationDate := policy.EndDate
	effectiveDate := transaction.EffectiveDate
	if effectiveDate.IsZero() {
		tmp, err := time.Parse(time.DateOnly, transaction.ScheduleDate)
		if err != nil {
			return document.ReceiptInfo{}, err
		}
		effectiveDate = tmp
	}

	tmpExpirationDate := lib.AddMonths(effectiveDate, 12)

	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly):
		tmpExpirationDate = lib.AddMonths(effectiveDate, 1)
	case string(models.PaySplitSingleInstallment):
		tmpExpirationDate = policy.EndDate
	}

	if !tmpExpirationDate.After(expirationDate) {
		expirationDate = tmpExpirationDate
	}

	receiptInfo.Transaction.EffectiveDate = effectiveDate.Format(dateFormat)
	receiptInfo.Transaction.ExpirationDate = expirationDate.Format(dateFormat)
	receiptInfo.Transaction.PriceGross = transaction.Amount

	return receiptInfo, nil
}
