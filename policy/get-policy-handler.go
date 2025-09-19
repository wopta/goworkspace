package policy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func getPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result map[string]string
	)

	log.AddPrefix("GetPolicyFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	uid := chi.URLParam(r, "uid")
	log.Println(uid)

	request := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	json.Unmarshal(request, &result)

	policy, _ := GetPolicy(uid)
	res, _ := json.Marshal(policy)

	log.Println("Handler end -------------------------------------------------")

	return string(res), policy, nil
}

func GetPolicy(uid string) (models.Policy, error) {
	var (
		policy models.Policy
		err    error
	)
	firePolicy := lib.PolicyCollection
	docsnap, err := lib.GetFirestoreErr(firePolicy, uid)
	if err != nil {
		return policy, err
	}
	err = docsnap.DataTo(&policy)
	return policy, err
}

func GetPolicyByCodeCompany(codeCompany string) (res models.Policy, err error) {
	policyQuery := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "codeCompany",
				Operator:   "==",
				QueryValue: codeCompany,
			},
		},
	}
	docsnap, e := policyQuery.FirestoreWherefields(lib.PolicyCollection)
	if e != nil {
		return res, e
	}
	p := models.PolicyToListData(docsnap)
	if len(p) == 0 {
		return res, fmt.Errorf("Policy with codeCompany %v not found", codeCompany)
	}
	return p[0], nil
}
