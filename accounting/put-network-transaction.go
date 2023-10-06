package accounting

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type PutNetworkTransactionRequest struct {
	IsPay            bool      `json:"isPay"`
	IsConfirmed      bool      `json:"isConfirmed"`
	PayDate          time.Time `json:"payDate"`
	TransactionDate  time.Time `json:"transactionDate"`
	ConfirmationDate time.Time `json:"confirmationDate"`
}

func PutNetworkTransactionFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	log.Println("[PutNetworkTransactionFx] Handler start ---------------------")

	var (
		err     error
		request PutNetworkTransactionRequest
	)

	uid := r.Header.Get("uid")

	log.Printf("[PutNetworkTransactionFx] uid %s", uid)

	networkTransaction := GetNetworkTransactionByUid(uid)
	if networkTransaction == nil {
		log.Println("[PutNetworkTransactionFx] error network transaction not found")
		return "", "", fmt.Errorf("no network transaction found for uid %s", uid)
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("[PutNetworkTransactionFx] Request: %s", string(body))
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("[PutNetworkTransactionFx] error unmarshaling request: %s", err.Error())
		return "", "", err
	}

	updateNetworkTransaction(networkTransaction, &request)

	err = networkTransaction.SaveBigQuery()

	return "", "", err
}

func GetNetworkTransactionByUid(uid string) *models.NetworkTransaction {
	log.Printf("[GetNetworkTransactionByUid] uid %s", uid)

	query := fmt.Sprintf(
		"SELECT * FROM `%s.%s` WHERE uid='%s'",
		models.WoptaDataset,
		models.NetworkTransactionCollection,
		uid,
	)

	netTransactions, err := lib.QueryRowsBigQuery[models.NetworkTransaction](query)
	if err != nil || len(netTransactions) == 0 {
		log.Printf("[GetNetworkTransactionsByTransactionUid] error getting network transactions: %s", err.Error())
	}

	return &netTransactions[0]
}

func updateNetworkTransaction(original *models.NetworkTransaction, update *PutNetworkTransactionRequest) {
	if !original.IsPay && update.IsPay {
		original.Status = models.NetworkTransactionStatusPaid
		original.StatusHistory = append(original.StatusHistory, original.Status)
	}
	if !original.IsConfirmed && update.IsConfirmed {
		original.Status = models.NetworkTransactionStatusConfirmed
		original.StatusHistory = append(original.StatusHistory, original.Status)
	}
	original.IsPay = update.IsPay
	original.IsConfirmed = update.IsConfirmed
	original.PayDate = lib.GetBigQueryNullDateTime(update.PayDate)
	original.TransactionDate = lib.GetBigQueryNullDateTime(update.TransactionDate)
	original.ConfirmationDate = lib.GetBigQueryNullDateTime(update.ConfirmationDate)
}
