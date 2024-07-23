package renew

import (
	"slices"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"google.golang.org/api/iterator"
)

func GetRenewActiveTransactionsByPolicyUid(policyUid string, annuity int) ([]models.Transaction, error) {
	queries := []lib.Firequery{
		{
			Field:      "policyUid",
			Operator:   "==",
			QueryValue: policyUid,
		},
		{
			Field:      "annuity",
			Operator:   "==",
			QueryValue: annuity,
		},
		{
			Field:      "isDelete",
			Operator:   "==",
			QueryValue: false,
		},
	}

	transactions, err := queryExecutor(queries)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(transactions, sortByEffectiveDate)

	return transactions, nil
}

func GetRenewTransactionsByPolicyUid(policyUid string, annuity int) ([]models.Transaction, error) {
	queries := []lib.Firequery{
		{
			Field:      "policyUid",
			Operator:   "==",
			QueryValue: policyUid,
		},
		{
			Field:      "annuity",
			Operator:   "==",
			QueryValue: annuity,
		},
	}

	transactions, err := queryExecutor(queries)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(transactions, sortByEffectiveDate)

	return transactions, nil
}

func queryExecutor(queries []lib.Firequery) ([]models.Transaction, error) {
	q := lib.Firequeries{Queries: queries}
	transactions := make([]models.Transaction, 0)

	docsnap, err := q.FirestoreWherefields(lib.RenewTransactionCollection)
	if err != nil {
		return nil, err
	}

	for {
		d, err := docsnap.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var transaction models.Transaction
		if err = d.DataTo(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, err
}

func sortByEffectiveDate(a, b models.Transaction) int {
	return a.EffectiveDate.Compare(b.EffectiveDate)
}
