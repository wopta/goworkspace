package enrich

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models/catnat"
)

func CatnatAtecoFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("CatnatAtecoFx")
	defer log.PopPrefix()
	fiscalCode := chi.URLParam(r, "fiscalCode")
	if fiscalCode == "" {
		return "", nil, errors.New("FiscalCode no valid")
	}
	client := catnat.NewNetClient()
	response, err := client.EnrichAteco(fiscalCode)
	if err != nil {
		return "", nil, err
	}
	res, err := json.Marshal(response)
	if err != nil {
		return "", nil, err
	}
	return string(res), nil, err
}
