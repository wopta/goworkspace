package _script

import (
	"errors"
	"fmt"

	"github.com/mohae/deepcopy"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
	"google.golang.org/api/iterator"

	"time"
)

func CopyTransactionsToBigQuery() {
	//transactions := transaction.GetPolicyTransactions("", "wT6LRDMwSHViSbTCj5GM") DEV
	transactions := transaction.GetPolicyTransactions("", "6tztJcR7KwqBT7JHwnwI")

	for index, tr := range transactions {
		if tr.IsPay {
			t := deepcopy.Copy(tr).(models.Transaction)
			t.Status = models.TransactionStatusToPay
			t.StatusHistory = []string{models.TransactionStatusToPay}
			t.IsPay = false
			t.PaymentMethod = ""
			t.UpdateDate = t.CreationDate
			t.PayDate = time.Time{}
			t.TransactionDate = time.Time{}
			t.BigQuerySave("")
		}
		transactions[index].BigQuerySave("")
	}
}

func CopyAllPoliciesTransactionToBigQuery() error {
	var policyUids = make([]string, 0)

	queries := lib.Firequeries{
		Queries: []lib.Firequery{
			{Field: "isDeleted", Operator: "==", QueryValue: false},
		},
	}

	iter, err := queries.FirestoreWherefields(lib.PolicyCollection)
	if err != nil {
		return fmt.Errorf("unable to query firestore policies: %w", err)
	}
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to iterate over policies: %w", err)
		}

		var policy models.Policy
		err = doc.DataTo(&policy)
		if err != nil {
			return fmt.Errorf("unable to populate policy: %w", err)
		}

		policyUids = append(policyUids, policy.Uid)
	}

	for _, uid := range policyUids {
		transactions := transaction.GetPolicyTransactions("", uid)

		for index, tr := range transactions {
			t := deepcopy.Copy(tr).(models.Transaction)
			t.BigQuerySave("")
			transactions[index].BigQuerySave("")
		}
	}

	return nil
}
