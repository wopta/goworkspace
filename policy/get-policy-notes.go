package policy

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func getPolicyNotesFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	policyUid := chi.URLParam(r, "policyUid")
	p, e := GetPolicy(policyUid)
	if e != nil {
		return "", nil, e
	}

	notes, err := models.GetPolicyNotes(policyUid)
	if err != nil {
		return "", "", err
	}
	idToken := r.Header.Get("Authorization")
	user, err := lib.GetAuthTokenFromIdToken(idToken)
	if user.Role == lib.UserRoleAgency || user.Role == lib.UserRoleAgent {
		notes.Notes = slices.DeleteFunc(notes.Notes, func(element models.PolicyNote) bool {
			return (element.ReadableByProducer == false) && user.UserID == p.ProducerUid
		})
	}
	res, err := json.Marshal(notes)
	if err != nil {
		return "", "", err
	}
	return string(res), "", nil
}
