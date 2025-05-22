package test

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

func logFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("TestLog")
	defer log.PopPrefix()

	secutiry := chi.URLParam(r, "severity")
	message := chi.URLParam(r, "message")
	log.Log().CustomLog(message, log.SeverityType(secutiry))
	return "", nil, nil
}
