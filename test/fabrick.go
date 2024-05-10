package test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func TestFabrickFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err     error
		rawResp string
		resp    interface{}
	)

	log.SetPrefix("[TestFabrickFx] ")
	defer func() {
		log.Printf("Handler end ----------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	operation := chi.URLParam(r, "operation")

	switch operation {
	case "delete":
		return fabrickDelete(r)
	case "payment-instrument":
		rawResp, resp, err = fabrickPaymentInstrument(r)
	default:
		return "", nil, fmt.Errorf("unhandled operation")
	}

	return rawResp, resp, err
}
