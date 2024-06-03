package internal

import (
	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"google.golang.org/api/iterator"
)

func GetPolicyByUidAndCollection(policyUid, collection string) (models.Policy, error) {
	var (
		policy models.Policy
		iter   *firestore.DocumentIterator
		err    error
	)

	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{Field: "uid", Operator: "==", QueryValue: policyUid},
			{Field: "isDeleted", Operator: "==", QueryValue: false},
		},
	}

	if iter, err = q.FirestoreWherefields(collection); err != nil {
		return models.Policy{}, err
	}

	for {
		docSnap, err := iter.Next()
		if err != nil && err == iterator.Done {
			break
		}
		if err != nil {
			return models.Policy{}, err
		}

		if err = docSnap.DataTo(&policy); err != nil {
			return models.Policy{}, err
		}
	}

	return policy, err
}
