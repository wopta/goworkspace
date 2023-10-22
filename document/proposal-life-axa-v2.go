package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func lifeAxaProposalV2(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	lifeMainHeaderV2(pdf, policy, networkNode, true)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	lifeInsuredInfoSectionV2(pdf, policy)

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

	lifeEmitResumeSectionV2(pdf, policy)

	companiesDescriptionSection(pdf, policy.Company)

	axaHeader(pdf)

	pdf.AddPage()

	axaFooter(pdf)

	axaDeclarationsConsentSection(pdf, policy, true)

	pdf.AddPage()

	axaTableSection(pdf, policy)

	pdf.AddPage()

	axaTablePart2Section(pdf, policy, true)

	pdf.Ln(15)

	axaTablePart3Section(pdf)

	woptaHeader(pdf)

	pdf.AddPage()

	woptaFooter(pdf)

	producerInfo := loadProducerInfo(origin, networkNode)

	allegato3Section(pdf, producerInfo)

	pdf.AddPage()

	allegato4Section(pdf, producerInfo)

	pdf.AddPage()

	allegato4TerSection(pdf, producerInfo)

	pdf.AddPage()

	woptaPrivacySection(pdf)

	personalDataHandlingSection(pdf, policy, true)

	filename, out := saveContract(pdf, policy)
	return filename, out
}
