package callback

import (
	"log"
	"net/http"

	doc "github.com/wopta/goworkspace/document"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func Sign(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Sign")
	log.Println("GET params were:", r.URL.Query())
	var e error
	uid := r.URL.Query().Get("uid")
	envelope := r.URL.Query().Get("envelope")
	action := r.URL.Query().Get("action")

	log.Println(action)
	log.Println(envelope)
	log.Println(uid)

	if action == "workstepFinished" {
		policyF := lib.GetFirestore("policy", uid)
		var policy models.Policy
		policyF.DataTo(policy)
		policy.IsSign = true
		lib.SetFirestore("policy", uid, policy)
		e = lib.InsertRowsBigQuery("wopta", "policy", policy)
		doc.GetFile(policy.ContractFileId, uid)
	}

	return "", nil, e
}
