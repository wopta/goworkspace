package renew

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/policy/renew"
	"github.com/wopta/goworkspace/policy/utils"
	"google.golang.org/api/iterator"
)

var errUnauthorized = errors.New("unauthorized")

type response struct {
	Transactions []models.Transaction `json:"transactions"`
}

func GetRenewTransactionsByPolicyUidFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err    error
		policy models.Policy
		resp   response
	)

	log.SetPrefix("[GetRenewTransactionsByPolicyUidFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		return "", nil, err
	}

	userUid := authToken.UserID
	policyUid := chi.URLParam(r, "policyUid")

	if policy, err = renew.GetRenewPolicyByUid(policyUid); err != nil {
		log.Printf("error fetching policy '%s'", policyUid)
		return "", nil, err
	}

	if !utils.CanBeAccessedBy(authToken.Role, policy.ProducerUid, authToken.UserID) {
		log.Printf("policy %s is not included in %s %s portfolio", policyUid, authToken.Role, userUid)
		return "", nil, errUnauthorized
	}

	transactions, err := GetRenewTransactionsByPolicyUid(policyUid, policy.Annuity)
	if err != nil {
		return "", nil, err
	}

	resp.Transactions = transactions

	bytes, err := json.Marshal(resp)

	return string(bytes), nil, err
}

func GetRenewTransactionsByPolicyUid(policyUid string, annuity int) ([]models.Transaction, error) {
	result := make([]models.Transaction, 0)
	q := lib.Firequeries{
		Queries: []lib.Firequery{
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
		},
	}
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
		var value models.Transaction
		if err = d.DataTo(&value); err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	slices.SortFunc(result, sortByEffectiveDate)

	return result, nil
}

func sortByEffectiveDate(a, b models.Transaction) int {
	return a.EffectiveDate.Compare(b.EffectiveDate)
}
