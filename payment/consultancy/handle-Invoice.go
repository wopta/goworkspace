package consultancy

import (
	"fmt"
	"slices"

	"gitlab.dev.wopta.it/goworkspace/accounting"
	"gitlab.dev.wopta.it/goworkspace/document"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	ItemRate        = "rate"
	ItemConsultancy = "consultancy"
)

func GenerateInvoice(p models.Policy, t models.Transaction) error {
	if !slices.ContainsFunc(t.Items, isItemConsultancy) {
		return nil
	}

	if lib.GetBoolEnv("GENERATE_INVOICE") {
		invoice := accounting.MapPolicyInvoiceInc(p, t, "Contributo per intermediazione")
		documentName := fmt.Sprintf("assets/users/%s/polizza_%s_%d_invoice.pdf",
			p.Contractor.Uid, p.CodeCompany, p.StartDate.AddDate(p.Annuity, 0, 0).Year())
		_, err := accounting.DoInvoicePaid(invoice, documentName)
		if err != nil {
			return err
		}
	}

	// create proforma document
	if proformaResp, err := document.Proforma(p); err == nil {
		proformatAtt := models.Attachment{
			Name:        fmt.Sprintf("Nota informativa %d", p.StartDate.AddDate(p.Annuity, 0, 0).Year()),
			FileName:    proformaResp.Filename,
			IsPrivate:   false,
			ContentType: lib.GetContentType("pdf"),
			Link:        proformaResp.LinkGcs,
			Section:     models.DocumentSectionOther,
			Note:        "",
		}
		*p.Attachments = append(*p.Attachments, proformatAtt)
	} else {
		return err
	}

	return nil
}

func isItemConsultancy(i models.Item) bool {
	return i.Type == ItemConsultancy
}
