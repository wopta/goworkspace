package test

import (
	"net/http"

	"github.com/wopta/goworkspace/test/log"
)

func logFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	logger := log.NewLog()
	logger.AddPrefix("titleeee1")
	logger.AddPrefix("titleeee2")
	logger.Error("scoppiooo")
	logger.Warning("waaarr")
	return "", nil, nil
}
