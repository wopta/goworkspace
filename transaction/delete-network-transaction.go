package transaction

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func DeleteNetworkTransaction(nt *models.NetworkTransaction) error {
	nt.Status = models.NetworkTransactionStatusDeleted
	nt.StatusHistory = append(nt.StatusHistory, nt.Status)
	nt.IsDelete = true
	nt.DeletionDate = lib.GetBigQueryNullDateTime(time.Now().UTC())

	return nt.SaveBigQuery()
}
