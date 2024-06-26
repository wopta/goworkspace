package accounting

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
)

type PutNetworkTransactionRequest struct {
	IsPay            bool      `json:"isPay"`
	IsConfirmed      bool      `json:"isConfirmed"`
	IsDelete         bool      `json:"isDelete"`
	PayDate          time.Time `json:"payDate"`
	TransactionDate  time.Time `json:"transactionDate"`
	ConfirmationDate time.Time `json:"confirmationDate"`
	DeletionDate     time.Time `json:"deletionDate"`
}

func PutNetworkTransactionFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err     error
		request PutNetworkTransactionRequest
	)

	log.SetPrefix("[PutNetworkTransactionFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	uid := chi.URLParam(r, "uid")
	log.Printf("uid %s", uid)

	networkTransaction := transaction.GetNetworkTransactionByUid(uid)
	if networkTransaction == nil {
		log.Println("error network transaction not found")
		return "", "", fmt.Errorf("no network transaction found for uid %s", uid)
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", "", err
	}

	updateNetworkTransaction(networkTransaction, &request)

	err = networkTransaction.SaveBigQuery()

	log.Println("Handler end -------------------------------------------------")

	return "{}", "", err
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
	if !original.IsDelete && update.IsDelete {
		original.Status = models.NetworkTransactionStatusDeleted
		original.StatusHistory = append(original.StatusHistory, original.Status)
	}
	original.IsPay = update.IsPay
	original.IsConfirmed = update.IsConfirmed
	original.IsDelete = update.IsDelete
	original.PayDate = lib.GetBigQueryNullDateTime(update.PayDate)
	original.TransactionDate = lib.GetBigQueryNullDateTime(update.TransactionDate)
	original.ConfirmationDate = lib.GetBigQueryNullDateTime(update.ConfirmationDate)
	original.DeletionDate = lib.GetBigQueryNullDateTime(update.DeletionDate)
}
