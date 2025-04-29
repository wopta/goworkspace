package renew

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/policy/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errRenewPolicyNotFound = errors.New("renew policy not found")
	errUnauthorized        = errors.New("unauthorized")
)

func GetRenewPolicyByUidFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err    error
		policy models.Policy
	)

	log.AddPrefix("GetRenewPolicyByUidFx")
	defer func() {
		if err != nil {
			log.Error(err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		return "", nil, err
	}

	policyUid := chi.URLParam(r, "uid")
	if policy, err = GetRenewPolicyByUid(policyUid); err != nil {
		log.ErrorF("error getting policy with uid: '%s'", policyUid)
		return "", nil, err
	}

	if !utils.CanBeAccessedBy(authToken.Role, policy.ProducerUid, authToken.UserID) {
		log.Printf("policy %s is not included in %s %s portfolio", policyUid, authToken.Role, authToken.UserID)
		return "", nil, errUnauthorized
	}

	bytes, err := json.Marshal(policy)
	if err != nil {
		log.ErrorF("error marshaling policy with uid: '%s'", policyUid)
		return "", nil, err
	}

	return string(bytes), policy, nil
}

func GetRenewPolicyByUid(uid string) (models.Policy, error) {
	var policy models.Policy

	snapshot, err := lib.GetFirestoreErr(lib.RenewPolicyCollection, uid)
	if status.Code(err) == codes.NotFound {
		return models.Policy{}, errRenewPolicyNotFound
	}
	if err != nil {
		return models.Policy{}, err
	}

	if err = snapshot.DataTo(&policy); err != nil {
		return models.Policy{}, err
	}

	return policy, nil
}
