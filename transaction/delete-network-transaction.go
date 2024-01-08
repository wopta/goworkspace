package transaction

import (
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
