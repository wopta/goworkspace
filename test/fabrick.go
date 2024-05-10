package test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func TestFabrickFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
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
	default:
		return "", nil, fmt.Errorf("unhandled operation")
	}
}
