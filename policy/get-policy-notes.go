package policy

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func getPolicyNotes(w http.ResponseWriter, r *http.Request) (string, any, error) {
	policyUid := chi.URLParam(r, "policyUid")
	GetPolicy(policyUid)
	notes, err := models.GetPolicyNotes(policyUid)
	if err != nil {
		return "", "", err
	}
	res, err := json.Marshal(notes)
	if err != nil {
		return "", "", err
	}
	return string(res), "", nil
}
