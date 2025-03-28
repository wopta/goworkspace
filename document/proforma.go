package document

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/document/pkg/proforma"
	"github.com/wopta/goworkspace/models"
)

const (
	proformaDocumentFormat = "%s_Proforma_%s_%d_%s.pdf"
)

func Proforma(policy models.Policy) (DocumentResp, error) {
	var (
		err      error
		gsLink string
		out      []byte
	)

	generator := proforma.NewProformaGenerator(engine.NewFpdf(), &policy)
	if out, err = generator.Generate(); err != nil {
		log.Printf("error generating proforma: %v", err)
		return DocumentResp{}, err
	}

	filename := strings.ReplaceAll(fmt.Sprintf(proformaDocumentFormat, policy.NameDesc,
		policy.CodeCompany, policy.Annuity+1, time.Now().Format("2006-01-02_15:04:05")), " ", "_")
	
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
