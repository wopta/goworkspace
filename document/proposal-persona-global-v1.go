package document

import (
	"github.com/go-pdf/fpdf"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func personaGlobalProposalV1(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) {
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

	woptaHeader(pdf, true)

	pdf.AddPage()

	woptaFooter(pdf)

	woptaPrivacySection(pdf)

	personalDataHandlingSection(pdf, policy, true)

	globalHeader(pdf, true)

	pdf.AddPage()

	globalFooter(pdf)

	globalPrivacySection(pdf, (*policy.Surveys)[len(*policy.Surveys)-1])
}
