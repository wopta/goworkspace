package payment

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib/log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment/fabrick"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/policy/renew"
	tr "github.com/wopta/goworkspace/transaction"
	trxRenew "github.com/wopta/goworkspace/transaction/renew"
)

func DeleteTransactionFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		isRenew     bool
		policy      models.Policy
		transaction *models.Transaction
		collection  = lib.TransactionsCollection
	)

	defer func() {
		if err != nil {
			log.ErrorF("error: %s", err.Error())
		}
		log.Println("Handler end -------------------------------------------------")
		log.PopPrefix()
	}()

	log.AddPrefix("[DeleteTransactionFx] ")
	log.Println("Handler start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		return "", nil, err
	}

	uid := chi.URLParam(r, "uid")
	rawIsRenew := r.URL.Query().Get("isRenew")
	if isRenew, err = strconv.ParseBool(rawIsRenew); rawIsRenew != "" && err != nil {
		log.ErrorF("error: %s", err.Error())
		return "", nil, err
	}

	if !isRenew {
		transaction = tr.GetTransactionByUid(uid, "")
		if transaction == nil {
			log.Printf("transaction '%s' not found", uid)
			return "", nil, fmt.Errorf("transaction '%s' not found", uid)
		}
		policy, err = plc.GetPolicy(transaction.PolicyUid, "")
	} else {
		collection = lib.RenewTransactionCollection
		transaction = trxRenew.GetRenewTransactionByUid(uid)
		if transaction == nil {
			log.Printf("transaction '%s' not found", uid)
			return "", nil, fmt.Errorf("transaction '%s' not found", uid)
		}
		policy, err = renew.GetRenewPolicyByUid(transaction.PolicyUid)
	}
	if err != nil {
		log.Printf("policy '%s' not found", transaction.PolicyUid)
		return "", nil, err
	}

	bytes, _ := json.Marshal(transaction)
	log.Printf("found transaction: %s", string(bytes))

	if transaction.ProviderName == models.FabrickPaymentProvider && transaction.ProviderId != "" {
		err = fabrick.FabrickExpireBill(transaction.ProviderId)
		if err != nil {
			log.ErrorF("error deleting transaction on fabrick: %s", err.Error())
			return "", nil, err
		}
	}

	tr.DeleteTransaction(transaction, "Cancellata manualmente")

	err = saveTransaction(transaction, collection)
	if err != nil {
		log.Printf("%s", err.Error())
		return "", nil, err
	}

	if transaction.ProviderName == models.FabrickPaymentProvider && transaction.ProviderId == "" {
		log.Printf("sending warning email...")
		sendMail(authToken, policy, *transaction)
		log.Printf("warning email sent")
	}

	return "{}", nil, err
}

func saveTransaction(transaction *models.Transaction, collection string) error {
	var (
		err error
	)

	transaction.BigQueryParse()
	err = lib.SetFirestoreErr(collection, transaction.Uid, transaction)
	if err != nil {
		return fmt.Errorf("error saving transaction %s in Firestore: %v", transaction.Uid, err.Error())
	}

	err = lib.InsertRowsBigQuery(lib.WoptaDataset, collection, transaction)
	if err != nil {
		log.ErrorF("error saving transaction %s in BigQuery: %v", transaction.Uid, err.Error())
		return err
	}
	return nil
}

func sendMail(authToken lib.AuthToken, policy models.Policy, transaction models.Transaction) {
	const standardLineTemplate = `<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:16px;color:#000000;font-size:14px">%s</p>`
	var message string

	transactionData := fmt.Sprintf("%s %d", lib.ExtractLocalMonth(transaction.EffectiveDate),
		transaction.EffectiveDate.Year())

	lines := []string{
		"Annullo transazione polizza " + policy.CodeCompany + " rata " + transactionData,
		"Attenzione, la transazione è stata annullata correttamente su Woptal, ma non su Fabrick.",
		"Verifica la situazione su Fabrick.",
	}

	for _, line := range lines {
		message += fmt.Sprintf(standardLineTemplate, line)
	}
	message += `<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px"><br></p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">A presto,</p>
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#e50075;font-size:14px"><strong>Anna</strong> di Wopta Assicurazioni</p>`

	mailReq := mail.MailRequest{
		FromAddress:  mail.AddressAnna,
		To:           []string{authToken.Email},
		Cc:           mail.AddressOperations.Address,
		Message:      message,
		Subject:      "Annullo transazione polizza",
		IsHtml:       true,
		TemplateName: "",
		Title:        "Annullo transazione polizza",
	}

	mail.SendMail(mailReq)

}
