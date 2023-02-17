package broker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func GetPolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result map[string]string
	)
	log.Println("GetPolicy")
	log.Println(r.RequestURI)
	log.Println(r.Header.Get("uid"))
	requestPath := strings.Split(r.RequestURI, "/")
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	json.Unmarshal([]byte(request), &result)
	log.Println(requestPath[2])
	var policy models.Policy
	docsnap := lib.GetFirestore("policy", r.Header.Get("uid"))
	log.Println("to data")
	docsnap.DataTo(&policy)
	res, _ := json.Marshal(policy)
	return string(res), policy, nil
}
