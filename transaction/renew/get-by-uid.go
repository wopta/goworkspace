package renew

import (
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func GetRenewTransactionByUid(uid string) *models.Transaction {
	var (
		err         error
		transaction *models.Transaction
	)

	docSnap, err := lib.GetFirestoreErr(lib.RenewTransactionCollection, uid)
	if err != nil {
		return nil
	}

	err = docSnap.DataTo(&transaction)
	if err != nil {
		return nil
	}

	return transaction
}
