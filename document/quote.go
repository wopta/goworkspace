package document

import (
	"errors"

	"github.com/wopta/goworkspace/document/internal/domain"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/document/pkg/quote"
	"github.com/wopta/goworkspace/models"
)

var errProductNotImplemented = errors.New("not implemented")

func Quote(policy *models.Policy, product *models.Product) ([]byte, error) {
	var (
		generator domain.QuoteGenerator
	)

	switch policy.Name {
	case models.LifeProduct:
		generator = quote.NewLifeGenerator(engine.NewFpdf(), policy, product)
	case models.GapProduct:
		return nil, errProductNotImplemented
	case models.PersonaProduct:
		return nil, errProductNotImplemented
	case models.CommercialCombinedProduct:
		return nil, errProductNotImplemented
	default:
		return nil, errProductNotImplemented
	}

	return generator.Exec()
}
