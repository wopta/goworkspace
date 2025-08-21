package document

import (
	"errors"

	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/document/pkg/quote"
	"gitlab.dev.wopta.it/goworkspace/models"
)

var errProductNotImplemented = errors.New("not implemented")

func Quote(policy *models.Policy, product *models.Product) ([]byte, error) {
	var (
		generator domain.QuoteGenerator
	)

	switch policy.Name {
	case models.LifeProduct:
		generator = quote.NewLifeGenerator(engine.NewFpdf(), policy, product)
	case models.CatNatProduct:
		generator = quote.NewCatnatGenerator(engine.NewFpdf(), policy, *product)
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
