package broker

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func PolicyFiscalcode(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	log.Println("GetPolicyByFiscalCode")
	log.Println(r.RequestURI)
	log.Println(strings.Split(r.RequestURI, "/")[2])

	var policies models.Policy

	// get all policies from firestore
	docsnap := lib.GetFirestore("policy", r.Header.Get("fiscalcode"))
	docsnap.DataTo(&policies)

	// get all policies from wise
	// wiseDoc := wiseProxy.WiseProxyObj()

	res, _ := json.Marshal(policies)

	return string(res), policies, nil
}
