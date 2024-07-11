package renew

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/models"
)

func DeleteRenewPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
	)

	log.SetPrefix("[DeleteRenewPolicyFx] ")
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
		log.Printf("error getting renew policy %v", err)
		return "", nil, err
	}

	deleteRenewPolicy(&policy)

	return "", nil, nil
}

func deleteRenewPolicy(p *models.Policy) {
	p.IsDeleted = true
	p.Status = models.PolicyStatusDeleted
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.Updated = time.Now().UTC()
}
