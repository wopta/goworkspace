package document

import (
	"fmt"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/document/pkg/proforma"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	proformaDocumentFormat = "nota_informativa_polizza_%s_%d.pdf"
)

func Proforma(policy models.Policy) (DocumentResp, error) {
	var (
		err    error
		gsLink string
		out    []byte
	)

	generator := proforma.NewProformaGenerator(engine.NewFpdf(), &policy)
	if out, err = generator.Generate(); err != nil {
		log.ErrorF("error generating proforma: %v", err)
		return DocumentResp{}, err
	}

	filename := strings.ReplaceAll(fmt.Sprintf(proformaDocumentFormat,
		policy.CodeCompany, policy.StartDate.AddDate(policy.Annuity, 0, 0).Year()), " ", "_")

	if gsLink, err = generator.Save(filename, out); err != nil {
		log.ErrorF("error saving proforma: %v", err)
		return DocumentResp{}, err
	}

	res := DocumentResp{
		LinkGcs:  gsLink,
		Filename: filename,
	}

	return res, nil
}
