package document

import (
	"errors"
	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/document/pkg/addendum"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	addendumDocumentFormat = "%s_Appendice_%s_%s.pdf"
)

var (
	ErrNotImplemented = errors.New("addendum document not implemented for product")
)

func Addendum(policy *models.Policy) (DocumentGenerated, error) {

	switch policy.Name {
	case models.LifeProduct:
		pdf := engine.NewFpdf()
		generator := addendum.NewLifeAddendumGenerator(pdf, policy)
		generator.Generate()
		return generateAddendumDocument(pdf.GetPdf(), policy)
	case models.CatNatProduct:
		pdf := engine.NewFpdf()
		generator := addendum.NewCatnatAddendumGenerator(pdf, policy)
		generator.Generate()
		return generateAddendumDocument(pdf.GetPdf(), policy)
	}

	log.WarningF("addendum not implemented for product %s", policy.Name)
	return DocumentGenerated{}, ErrNotImplemented
}
