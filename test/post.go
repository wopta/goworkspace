package test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/lib"
)

func TestPostFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var request interface{}

	log.SetPrefix("[TestPostFx] ")
	defer log.SetPrefix("")
	log.Println("Handler start -----------------------------------------------")

	operation := chi.URLParam(r, "operation")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	json.Unmarshal(body, &request)
	log.Printf("payload %v", request)

	if operation == "error" {
		return "", nil, GetErrorJson(400, "Bad Request", "Testing error POST")
	} else if operation == "policy-update" {
		policyTransactionsUpdate(int(request.(float64)))
		return "{}", nil, nil
	}

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}
