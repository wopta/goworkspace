package callback

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/callback/internal/fabrick"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type paymentHandler interface {
	AnnuityFirstRateFx(http.ResponseWriter, *http.Request) (string, any, error)
	AnnuitySingleRateFx(http.ResponseWriter, *http.Request) (string, any, error)
}

func payment(w http.ResponseWriter, r *http.Request) (string, any, error) {

	log.AddPrefix("Payment")
	defer func() {
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()

	log.Println("Handler start -----------------------------------------------")
	rate := chi.URLParam(r, "rate")
	provider := chi.URLParam(r, "provider")
	log.Printf("Rate '%v' with provider '%v'", rate, provider)
	var handler paymentHandler
	switch provider {
	case models.FabrickPaymentProvider:
		handler = fabrick.FabrickCallback{}
	default:
		return "", nil, fmt.Errorf("Provider '%s' not supported", provider)
	}

	switch rate {
	case "first-rate":
		return handler.AnnuityFirstRateFx(w, r)
	case "single-rate":
		return handler.AnnuitySingleRateFx(w, r)
	}
	return "", nil, fmt.Errorf("Rate '%s' not supported", rate)
}
