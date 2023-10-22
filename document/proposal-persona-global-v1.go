package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func personaGlobalProposalV1(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	personaMainHeaderV1(pdf, policy, networkNode, true)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	personaInsuredInfoSection(pdf, policy)

	guaranteesMap, slugs := loadPersonaGuarantees(policy, product)

	personaGuaranteesTable(pdf, guaranteesMap, slugs)

	pdf.Ln(5)

	personaSurveySection(pdf, policy, true)

	personaStatementsSection(pdf, policy, true)

	if policy.HasGuarantee("IPM") {
		pdf.AddPage()
	}

	personaOfferResumeSection(pdf, policy)

	paymentMethodSection(pdf)

	emitResumeSection(pdf, policy)

	companiesDescriptionSection(pdf, policy.Company)

	personalDataHandlingSection(pdf, policy, true)

	filename, out := saveProposal(pdf, policy)
	return filename, out
}
