package document

import (
	"bytes"
	"fmt"
	"github.com/go-pdf/fpdf"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func generateContractDocument(pdf *fpdf.Fpdf, policy *models.Policy) (DocumentGenerated, error) {
	var (
		res DocumentGenerated
		out bytes.Buffer
	)
	err := pdf.Output(&out)
	if err != nil {
		return res, err
	}
	res.FileName = fmt.Sprintf(models.ContractDocumentFormat, policy.NameDesc, policy.CodeCompany)
	res.ParentPath = fmt.Sprintf("temp/%s", policy.Uid)
	res.Bytes = out.Bytes()
	return res, nil
}

func generateProposalDocument(pdf *fpdf.Fpdf, policy *models.Policy) (DocumentGenerated, error) {
	var (
		res DocumentGenerated
		out bytes.Buffer
	)
	err := pdf.Output(&out)
	if err != nil {
		return res, err
	}
	res.FileName = fmt.Sprintf(models.ProposalDocumentFormat, policy.NameDesc, policy.ProposalNumber)
	res.ParentPath = fmt.Sprintf("temp/%s", policy.Uid)
	res.Bytes = out.Bytes()
	return res, nil
}

func generateReservedDocument(pdf *fpdf.Fpdf, policy *models.Policy) (DocumentGenerated, error) {
	var (
		res DocumentGenerated
		out bytes.Buffer
	)
	err := pdf.Output(&out)
	if err != nil {
		return res, err
	}
	res.FileName = fmt.Sprintf(models.RvmInstructionsDocumentFormat, policy.ProposalNumber)
	res.ParentPath = fmt.Sprintf("temp/%s", policy.Uid)
	res.Bytes = out.Bytes()
	return res, nil
}
