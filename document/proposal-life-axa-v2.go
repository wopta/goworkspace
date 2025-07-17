package document

import (
	"github.com/go-pdf/fpdf"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func lifeAxaProposalV2(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) {
	lifeMainHeaderV2(pdf, policy, networkNode, true)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	lifeInsuredInfoSectionV2(pdf, policy, true)

	guaranteesMap, slugs := loadLifeGuarantees(policy, product)

	lifeGuaranteesTableV2(pdf, guaranteesMap, slugs)

	lifeAvvertenzeBeneficiariSectionV2(pdf)

	beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice := loadLifeBeneficiariesInfo(policy)

	lifeBeneficiariesSectionV2(pdf, beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice)

	lifeBeneficiaryReferenceSectionV2(pdf, policy)

	lifeSurveysSectionV2(pdf, policy, true)

	pdf.AddPage()

	lifeStatementsSectionV2(pdf, policy, true)

	pdf.AddPage()

	lifeOfferResumeSectionV2(pdf, policy)

	lifePaymentResumeSectionV2(pdf, policy)

	lifeContractWithdrawlSectionV2(pdf, true)

	lifePaymentMethodSectionV2(pdf)

	companiesDescriptionSection(pdf, policy.Company)

	axaHeader(pdf, true)

	pdf.AddPage()

	axaFooter(pdf)

	axa2DeclarationsConsentSection(pdf, policy, true)

	_, err := policy.ExtractGuarantee("death")
	if err == nil {
		pdf.AddPage()

		axaTableSection(pdf, policy)
	}

	pdf.AddPage()

	axaTablePart2Section(pdf, policy, true)

	pdf.Ln(15)

	axaTablePart3Section(pdf)

	woptaHeader(pdf, true)

	pdf.AddPage()

	woptaFooter(pdf)

	woptaPrivacySection(pdf)

	personalDataHandlingSection(pdf, policy, true)
}
