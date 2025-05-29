package document

import (
	"github.com/go-pdf/fpdf"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func lifeAxaProposalV1(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (DocumentGenerated, error) {
	mainHeader(pdf, policy, true)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	insuredInfoSection(pdf, policy, true)

	guaranteesMap, slugs := loadLifeGuarantees(policy, product)

	lifeGuaranteesTable(pdf, guaranteesMap, slugs)

	avvertenzeBeneficiariSection(pdf)

	beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice := loadLifeBeneficiariesInfo(policy)

	beneficiariesSection(pdf, beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice)

	beneficiaryReferenceSection(pdf, policy)

	surveysSection(pdf, policy, true)

	pdf.AddPage()

	statementsSection(pdf, policy, true)

	pdf.AddPage()

	offerResumeSection(pdf, policy)

	paymentResumeSection(pdf, policy)

	contractWithdrawlSection(pdf, true)

	pdf.AddPage()

	paymentMethodSection(pdf)

	companiesDescriptionSection(pdf, policy.Company)

	axaHeader(pdf, true)

	pdf.AddPage()

	axaFooter(pdf)

	axaDeclarationsConsentSection(pdf, policy, true)

	pdf.AddPage()

	axaTableSection(pdf, policy)

	pdf.AddPage()

	axaTablePart2Section(pdf, policy, true)

	pdf.Ln(15)

	axaTablePart3Section(pdf)

	generatePolicyAnnex(pdf, origin, networkNode, policy, setAnnexHeaderFooter(pdf, networkNode, true))

	woptaHeader(pdf, true)

	pdf.AddPage()

	woptaFooter(pdf)

	woptaPrivacySection(pdf)

	personalDataHandlingSection(pdf, policy, true)

	return generateProposalDocument(pdf, policy)
}
