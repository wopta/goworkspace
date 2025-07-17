package document

import (
	"github.com/go-pdf/fpdf"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func gapSogessurProposalV1(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode) (DocumentGenerated, error) {
	gapHeaderV1(pdf, policy, networkNode, true)

	gapFooterV1(pdf, policy.NameDesc)

	pdf.AddPage()

	vehicle := policy.Assets[0].Vehicle
	contractor := policy.Contractor
	vehicleOwner := policy.Assets[0].Person
	statements := *policy.Statements

	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 4, "La tua assicurazione è operante sui dati sotto riportati, verifica la loro correttezza"+
		" e segnala eventuali inesattezze", "", "", false)
	pdf.Ln(3)

	gapVehicleDataTableV1(pdf, vehicle)

	gapPersonalInfoTableV1(pdf, contractor, *vehicleOwner)

	gapPolicyDataTableV1(pdf, policy)

	gapPriceTableV1(pdf, policy)

	pdf.Ln(3)

	gapStatementsV1(pdf, statements[:len(statements)-1], policy.Company, true)

	companiesDescriptionSection(pdf, policy.Company)

	woptaGapHeader(pdf, *policy, true)

	pdf.AddPage()

	woptaFooter(pdf)

	printStatement(pdf, statements[len(statements)-1], policy.Company, true)

	generatePolicyAnnex(pdf, networkNode, policy, setAnnexHeaderFooter(pdf, networkNode, true))

	woptaHeader(pdf, true)

	pdf.AddPage()

	woptaFooter(pdf)

	woptaPrivacySection(pdf)

	personalDataHandlingSection(pdf, policy, true)

	return generateProposalDocument(pdf, policy)
}
