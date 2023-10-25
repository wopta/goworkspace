package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func lifeAxaProposalV1(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
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

	emitResumeSection(pdf, policy)

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

	filename, out := saveProposal(pdf, policy)
	return filename, out
}
