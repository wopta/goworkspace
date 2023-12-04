package transaction

import (
	"log"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func DeleteNetworkTransaction(nt *models.NetworkTransaction) error {
	nt.Status = models.NetworkTransactionStatusDeleted
	nt.StatusHistory = append(nt.StatusHistory, nt.Status)
	nt.IsDelete = true
	nt.DeletionDate = lib.GetBigQueryNullDateTime(time.Now().UTC())

	return nt.SaveBigQuery()
}

func DeleteNetworkTransactionByUid(uid string) error {
	nt := GetNetworkTransactionByUid(uid)
	if nt == nil {
		log.Printf("[DeleteNetworkTransactionByUid] cannot delete, node %s not found", uid)
		return nil
	}

	return DeleteNetworkTransaction(nt)
}
