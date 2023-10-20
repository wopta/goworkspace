package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func lifeAxaV2(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	signatureID = 0

	mainHeaderV2(pdf, policy, networkNode)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	insuredInfoSection(pdf, policy)

	guaranteesMap, slugs := loadLifeGuarantees(policy, product)

	lifeGuaranteesTable(pdf, guaranteesMap, slugs)

	avvertenzeBeneficiariSection(pdf)

	beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice := loadLifeBeneficiariesInfo(policy)

	beneficiariesSection(pdf, beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice)

	beneficiaryReferenceSection(pdf, policy)

	surveysSection(pdf, policy)

	pdf.AddPage()

	statementsSection(pdf, policy)

	offerResumeSection(pdf, policy)

	paymentResumeSection(pdf, policy)

	contractWithdrawlSection(pdf)

	pdf.AddPage()

	paymentMethodSection(pdf)

	emitResumeSection(pdf, policy)

	companiesDescriptionSection(pdf, policy.Company)

	axaHeader(pdf)

	pdf.AddPage()

	axaFooter(pdf)

	axaDeclarationsConsentSection(pdf, policy)

	pdf.AddPage()

	axaTableSection(pdf, policy)

	pdf.AddPage()

	axaTablePart2Section(pdf, policy)

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

	personalDataHandlingSection(pdf, policy)

	filename, out := saveContract(pdf, policy)
	return filename, out
}
