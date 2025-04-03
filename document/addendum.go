package document

import (
	"errors"
	"fmt"
	"github.com/wopta/goworkspace/lib/log"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/document/pkg/addendum"
	"github.com/wopta/goworkspace/models"
)

const (
	addendumDocumentFormat = "%s_Appendice_%s_%s.pdf"
)

var (
	ErrNotImplemented = errors.New("addendum document not implemented for product")
)

func Addendum(policy *models.Policy) (DocumentResp, error) {
	var (
		err      error
		filename string
		gsLink   string
		out      []byte
	)

	switch policy.Name {
	case models.LifeProduct:
		generator := addendum.NewLifeAddendumGenerator(engine.NewFpdf(), policy)
		if out, err = generator.Generate(); err != nil {
			log.ErrorF("error generating addendum: %v", err)
			return DocumentResp{}, err
		}

		filename = strings.ReplaceAll(fmt.Sprintf(addendumDocumentFormat, policy.NameDesc,
			policy.CodeCompany, time.Now().Format("2006-01-02_15:04:05")), " ", "_")

		if gsLink, err = generator.Save(filename, out); err != nil {
			log.ErrorF("error saving addendum: %v", err)
			return DocumentResp{}, err
		}
	default:
		log.Printf("addendum not implemented for product %s", policy.Name)
		return DocumentResp{}, ErrNotImplemented
	}

	res := DocumentResp{
		LinkGcs:  gsLink,
		Filename: filename,
	}

	return res, nil
}
