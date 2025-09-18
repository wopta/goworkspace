package policy

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func postPolicyNote(w http.ResponseWriter, r *http.Request) (string, any, error) {
	body, err := io.ReadAll(r.Body)
	policyUid := chi.URLParam(r, "policyUid")
	if err != nil {
		return "{}", nil, err
	}
	defer r.Body.Close()
	var note models.PolicyNote
	err = json.Unmarshal(body, &note)
	if err != nil {
		return "", nil, err
	}
	err = models.AddNoteToPolicy(policyUid, note)
	return "{}", nil, err
}
