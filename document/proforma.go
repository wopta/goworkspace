package document

import (
	"fmt"
	"log"
	"strings"

	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/document/pkg/proforma"
	"github.com/wopta/goworkspace/models"
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
		log.Printf("error generating proforma: %v", err)
		return DocumentResp{}, err
	}

	filename := strings.ReplaceAll(fmt.Sprintf(proformaDocumentFormat,
		policy.CodeCompany, policy.StartDate.AddDate(policy.Annuity, 0, 0).Year()), " ", "_")

	if gsLink, err = generator.Save(filename, out); err != nil {
		log.Printf("error saving proforma: %v", err)
		return DocumentResp{}, err
	}

	res := DocumentResp{
		LinkGcs:  gsLink,
		Filename: filename,
	}

	return res, nil
}
