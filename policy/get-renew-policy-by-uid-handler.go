package policy

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errRenewPolicyNotFoud = errors.New("renew policy not found")

func GetRenewPolicyByUidFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err    error
		policy models.Policy
	)

	log.SetPrefix("[GetRenewPolicyByUidFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	uid := chi.URLParam(r, "uid")

	if policy, err = getRenewPolicyByUid(uid); err != nil {
		log.Printf("error getting policy with uid: '%s'", uid)
		return "", nil, err
	}

	bytes, err := json.Marshal(policy)
	if err != nil {
		log.Printf("error marshaling policy with uid: '%s'", uid)
		return "", nil, err
	}

	return string(bytes), policy, nil
}

func getRenewPolicyByUid(uid string) (models.Policy, error) {
	var policy models.Policy

	snapshot, err := lib.GetFirestoreErr(lib.RenewPolicyCollection, uid)
	if status.Code(err) == codes.NotFound {
		return models.Policy{}, errRenewPolicyNotFoud
	}
	if err != nil {
		return models.Policy{}, err
	}

	if err = snapshot.DataTo(&policy); err != nil {
		return models.Policy{}, err
	}

	return policy, nil
}
